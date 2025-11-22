# Go Microservice Boilerplate (Hexagonal Architecture)

A **ready‑to‑use** starter kit for building production‑grade microservices in Go using the Hexagonal (Ports & Adapters) architecture.

## Features
- Clean separation of **Domain**, **Ports**, and **Adapters**
- HTTP API with **chi** router, health‑check and Prometheus metrics
- PostgreSQL persistence via **pgx** with a repository adapter
- Docker multi‑stage build and **docker‑compose** for local development
- Makefile with common dev commands (`run`, `test`, `lint`, `docker‑build`, `docker‑run`)
- GitHub Actions CI that runs tests, lint and builds the Docker image
- MIT licensed – free for commercial use

## Quick Start
```bash
# Clone the repo
git clone https://github.com/yourorg/go-microservice-boilerplate.git
cd go-microservice-boilerplate

# Fetch Go dependencies
go mod tidy

# Start PostgreSQL and the service (Docker Compose)
docker compose up -d

# Run the server locally (optional)
make run
```
The service will be reachable at `http://localhost:8080`. Try the health endpoint:
```bash
curl http://localhost:8080/healthz
```
You should see `{"status":"ok"}`.

## Project Layout
```
cmd/                # entry points (main.go)
internal/           # private packages
  domain/           # business models & services
  infra/            # adapters (DB, external APIs)
pkg/                # reusable public packages (if any)
Dockerfile          # multi‑stage build
docker-compose.yml  # local dev stack
Makefile            # dev shortcuts
go.mod, go.sum      # module definition
```

## Development
- **Run tests**: `make test`
- **Lint**: `make lint` (requires golangci‑lint installed)
- **Build binary**: `make build`
- **Docker image**: `make docker-build`

## Contributing
Feel free to open issues or PRs. Follow the contribution guidelines in `CONTRIBUTING.md` (to be added).

## License
MIT – see the `LICENSE` file.
