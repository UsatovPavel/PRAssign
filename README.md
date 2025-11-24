# PRAssign

## Quick Start
1. Clone repository:
   ```bash
   git clone https://github.com/UsatovPavel/PRAssign.git
   cd PRAssign
   ```
2. Set environment variables (create `.env` file):
   ```ini
   AUTH_SECRET=your_strong_secret_here
   ```
3. Build and run:
   ```bash
   make app
   ```

## Key Features

- **Authentication**: Token-based auth for protected endpoints.
- **Dockerized CI & Testing**: Docker Compose setups for build, integration and load(with k6) tests
- **Code Quality & Linting**: golangci-lint integrated into the build pipeline (configurable, runs in CI/build).
- **Team Management**: Create and retrieve teams with member/admin access checks.
- **User Administration**: Toggle user active state and fetch per-user review data (self or admin).
- **Pull Request Workflow**: Create, merge and reassign PRs with author/admin authorization.
- **Statistics & Reporting:** Aggregated assignment metrics and per-user reports for admins.


## Tech Stack
- **Language**: Go 1.24
- **Framework**: Gin
- **Migrations**: golang-migrate
- **Database**:  PostgreSQL
- **Build Tool**: Docker v2/make
- **Linter**    golangci-lint 2.6.2 
- **Load testing**: k6(js scripts)
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

## Testing

4. Run integration tests
   ```bash
   make test-int
   ```
   Run load tests (after
   ```bash
   make app
   ```
   )
   ```bash
   make test-load
   ```
