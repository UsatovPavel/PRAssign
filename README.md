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
4. Run integration tests
   ```bash
   make test-int
   ```
   Run load tests (after make app)
   ```bash
   make test-load
   ```
5. API openapi.yml TODO fix

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
