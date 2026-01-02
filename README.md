# PRAssign
Go with Gin microservice for Avito TestTask. Provides API for user/pullrequests/team management,assignments statistics.
## Quick Start
1. Clone repository:
   ```bash
   git clone https://github.com/UsatovPavel/PRAssign.git
   cd PRAssign
   ```
2. Set environment variables (create `.env` file):
   ```ini
   SERVER_PORT=8080
   TEST_SERVER_PORT=18080
   AUTH_SECRET=your_strong_secret_here

   POSTGRES_USER=pr-assignment
   POSTGRES_PASSWORD=pr-assignment
   POSTGRES_DB=pr-assignment
   POSTGRES_DB_TEST=pr-assignment-test

   FACTORIAL_KAFKA_BOOTSTRAP=kafka1:9092,kafka2:9092,kafka3:9092
   FACTORIAL_KAFKA_TOPIC_TASKS=factorial.tasks
   FACTORIAL_KAFKA_TOPIC_RESULTS=factorial.results
   FACTORIAL_KAFKA_GROUP_RESULTS=factorial-results-consumer
   FACTORIAL_RESULTS_TTL=1h
   FACTORIAL_RESULTS_TIMEOUT=30s
   FACTORIAL_RESULTS_CONSUMER_ENABLED=1

   APP_REPLICAS=1
   ```
3. Build and run:
   ```bash
   make app
   ```
   Or clone [Scala project](https://github.com/UsatovPavel/AsyncFactorial) into a neighboring folder and run:
   ```bash
   make all
   ```
   logs: ../go.log (proxy), ../scala.log (consumer, filtered).

## Key Features

- **Authentication**: Token-based auth for protected endpoints.
- **Dockerized CI & Testing**: Docker Compose setups for build, integration and load (with k6) tests.
- **Code Quality & Linting**: golangci-lint integrated into the build pipeline (configurable, runs in CI/build).
- **Team Management**: Create and retrieve teams with member/admin access checks.
- **User Administration**: Toggle user active state and fetch per-user review data (self or admin).
- **Pull Request Workflow**: Create, merge and reassign PRs with author/admin authorization.
- **Statistics & Reporting:** Aggregated assignment metrics and per-user reports for admins.
- **Kafka factorial pipeline:** POST → Kafka → Scala → Kafka → Go consumer → Postgres; GET returns accumulated results. Includes e2e/stress/resilience tests.
- **Scaling + proxy:** nginx load-balances multiple Go replicas; Scala consumer is scalable. All settings come from `.env`.


## Tech Stack
- **Language**: Go 1.23
- **Framework**: Gin
- **Kafka client**: Sarama
- **Database**: PostgreSQL + golang-migrate
- **Proxy**: nginx 
- **Build Tools**: Docker Compose v2, make
- **Linter**: golangci-lint 2.6.2
- **Load testing**: k6 (JS scripts)
---

## Prerequisites
- Go 1.24
- Docker & DockerCompose v2
- PostgreSQL 15+

--- 

## API Reference

### Authentication
| Endpoint      | Method | Description                                 | Rights                                |
| ------------- | ------ | ------------------------------------------- | ------------------------------------- |
| `/auth/token` | POST   | Issue a token by username                   | Public            |

### Health
| Endpoint  | Method | Description           | Rights                                |
| --------- | ------ | --------------------- | ------------------------------------- |
| `/health` | GET    | Service health check. | Public            |

### Team Management
| Endpoint    | Method | Description                | Rights                                                                 |
| ----------- | ------ | -------------------------- | ---------------------------------------------------------------------- |
| `/team/add` | POST   | Create/add a team.         | Admin/user listed among team members in the request |
| `/team/get` | GET    | Retrieve team information. | Admin/team member of the requested team |

### Users Management
| Endpoint             | Method | Description                        | Rights                                |
| -------------------- | ------ | ---------------------------------- | ------------------------------------- |
| `/users/setIsActive` | POST   | Set/update user active flag.       | Admin, Requested user |
| `/users/getReview`   | GET    | Get list of reviews by user_id.    | Admin, Requested user |

### Pull request Management
| Endpoint                | Method | Description                 | Rights                                         |
| ----------------------- | ------ | --------------------------- | ---------------------------------------------- |
| `/pullRequest/create`   | POST   | Create a pull request.      | PR author or admin |
| `/pullRequest/merge`    | POST   | Merge a pull request.       | PR author or admin |
| `/pullRequest/reassign` | POST   | Reassign PR reviewer/owner. | PR author or admin |

### Statistics Management
| Endpoint                               | Method | Description                                           | Rights      |
| -------------------------------------- | ------ | ----------------------------------------------------- | ----------- |
| `/statistics/assignments/users`        | GET    | Assignments statistics by users.                      | Admin   |
| `/statistics/assignments/pullrequests` | GET    | Assignments statistics by pull requests.              | Admin   |
| `/statistics/assignments/user/:id`     | GET    | Assignments statistics for a specific user.           | Admin  |

### Factorial
| Endpoint                               | Method | Description                                           | Rights      |
| -------------------------------------- | ------ | ----------------------------------------------------- | ----------- |
| `/factorial`                           | POST   | Enqueue factorial tasks (Job-Id header required).     | Auth (token) |
| `/factorial/:job_id/result`            | GET    | Get aggregated factorial results by job_id.           | Auth (token) |

## Testing

4. Run integration tests
   ```bash
   make test-int
   ```
   Run load/e2e tests (after
   ```bash
   make app
   ```
   )
   ```bash
   make test-load
   ```
   ```bash
   make test-e2e
   ```
   Run resilience tests (after `make all`):
   ```bash
   make resilence TEST=<name>
   ```