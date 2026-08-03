# dforge

[![Go Reference](https://pkg.go.dev/badge/github.com/yunkhngn/dforge.svg)](https://pkg.go.dev/github.com/yunkhngn/dforge)
[![Release](https://img.shields.io/github/v/release/yunkhngn/dforge)](https://github.com/yunkhngn/dforge/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/yunkhngn/dforge)](go.mod)

One command to bootstrap, inspect, and manage Docker Compose projects offline.

`dforge` is an offline-first Go CLI tool designed to simplify Docker Compose workflows. It detects your tech stack, generates production-grade Dockerfiles and Compose configurations, mutates `compose.yaml` while preserving comments, audits container security/best practices with a weighted health score, and provides a lightweight runtime wrapper over `docker compose`.

---

## ✨ Features

- **🚀 Smart Stack Detection (`dforge init`)**: Automatically detects Next.js, React, NestJS, Express, Spring Boot, Go, and Rust projects.
- **🛡️ Production-Grade Dockerfiles**: Generates multi-stage Dockerfiles with non-root users (`USER`), explicit `HEALTHCHECK` instructions, and pinned image tags (never `:latest`).
- **💬 Comment-Preserving AST Editor (`dforge add` / `remove`)**: Add or remove services in `compose.yaml` while preserving existing comments, formatting, and top-level volume definitions.
- **🩺 Automated Doctor (`dforge doctor`)**: Runs 12+ security and best-practice checks, outputting a weighted health score out of 100 with actionable feedback.
- **🔒 Privacy-First Env Validator (`dforge env`)**: Compares `.env` against `.env.example` for missing, extra, or empty keys while **never printing variable values**.
- **⚡ Runtime Container Wrapper (`status`, `logs`, `shell`, `clean`)**: Easily check container health, tail service logs, drop into shell sessions, or clean dangling Docker resources with interactive safety confirmation.
- **🌐 100% Offline**: Zero telemetry, zero AI/cloud APIs, zero runtime network dependencies.

---

## 📦 Installation

### Homebrew (macOS / Linux)

```bash
brew tap yunkhngn/tap
brew install dforge
```

### Go Install

```bash
go install github.com/yunkhngn/dforge@latest
```

### Build from Source

```bash
git clone https://github.com/yunkhngn/dforge.git
cd dforge
go build -o dforge main.go
```

---

## ⚡ Quick Start

### 1. Initialize a Project

Navigate to your application directory and run `dforge init`:

```bash
cd my-nextjs-app
dforge init
```

`dforge` will detect Next.js and generate:
- `Dockerfile` (multi-stage, non-root, pinned `node:20-alpine`)
- `.dockerignore` (optimized node ignore rules)
- `compose.yaml` (port mapped `3000:3000`, `restart: unless-stopped`)

### 2. Add Predefined Infrastructure Services

Add a PostgreSQL database and Redis cache with pinned images, health checks, and volume persistence:

```bash
dforge add postgres
dforge add redis
```

### 3. Audit Container Setup

Check your setup against Docker best practices:

```bash
dforge doctor
```

Output:
```
✓ Docker installed
✓ Docker Compose installed
✓ Dockerfile exists
✓ compose.yaml exists
✓ .dockerignore exists
✓ HEALTHCHECK present
✓ restart policy set
✓ no :latest image tag
✓ not running as root

Score: 100/100
```

### 4. Validate Environment Variables

Compare `.env` with `.env.example` without exposing secrets:

```bash
dforge env
```

---

## 🛠️ Command Reference

| Command | Usage | Description |
| :--- | :--- | :--- |
| `init` | `dforge init [--framework <fw>]` | Detect framework and generate `Dockerfile`, `.dockerignore`, and `compose.yaml`. |
| `doctor` | `dforge doctor` | Audit Docker setup and calculate health score (0-100). |
| `add` | `dforge add <service>` | Add a predefined service to `compose.yaml`. |
| `remove` | `dforge remove <service>` | Remove a service from `compose.yaml` and prune orphan volumes. |
| `env` | `dforge env` | Compare `.env` with `.env.example` (shows variable names only). |
| `status` | `dforge status` | Display service status, health state, and published ports. |
| `logs` | `dforge logs [service] [-f]` | Show or tail container logs. |
| `shell` | `dforge shell <service>` | Open interactive shell (`bash` with `sh` fallback) inside container. |
| `clean` | `dforge clean [--force]` | Remove dangling images, unused volumes, networks, and build cache. |

### Global Flags

- `--yes`: Automatically answer "yes" to confirmation prompts.
- `--force`: Bypass safety checks and prompt confirmations.

---

## 🗄️ Service Catalog

`dforge` comes with a built-in catalog of 8 pinned, production-ready infrastructure services:

| Service | Base Image Tag | Default Port(s) | Default Healthcheck |
| :--- | :--- | :--- | :--- |
| `postgres` | `postgres:16-alpine` | `5432:5432` | `pg_isready -U postgres` |
| `mysql` | `mysql:8.4` | `3306:3306` | `mysqladmin ping -h localhost` |
| `redis` | `redis:7.4-alpine` | `6379:6379` | `redis-cli ping` |
| `mongodb` | `mongo:7.0` | `27017:27017` | `mongosh --eval db.adminCommand('ping')` |
| `rabbitmq` | `rabbitmq:3.13-management` | `5672:5672`, `15672:15672` | `rabbitmq-diagnostics ping` |
| `minio` | `minio/minio:RELEASE...` | `9000:9000`, `9001:9001` | `mc ready local` |
| `mailpit` | `axllent/mailpit:v1.20` | `1025:1025`, `8025:8025` | — |
| `meilisearch` | `getmeili/meilisearch:v1.9` | `7700:7700` | — |

---

## 🏗️ Architecture

`dforge` follows a strict clean architecture:

```
main.go                              # Entrypoint -> cmd.Execute()
cmd/                                 # Cobra CLI commands & presentation seams
  root.go init.go doctor.go add.go remove.go env.go status.go logs.go shell.go clean.go
internal/
  detect/                            # Signal-based framework detector
  generator/                         # Embedded Go template engine & renderer
  services/                          # Pinned service catalog registry
  compose/                           # gopkg.in/yaml.v3 AST comment-preserving editor
  env/                               # Redacted .env file comparison engine
  doctor/                            # Health check registry & weighted scoring engine
  docker/                            # Docker client interface, exec driver & test mocks
  fsutil/                            # Atomic write & backup file utilities (.bak)
  ui/                                # Lipgloss styles & Bubble Tea TUI components
```

- **Logic / Shell Separation**: Pure business logic lives in `internal/` packages and is unit tested without Docker or a TTY. `cmd/` and `internal/ui/` act as thin presentation wrappers.
- **Safety First**: Any destructive operation requires user confirmation unless `--force` or `--yes` is passed. Automatic `<file>.bak` backups are created before file mutations.

---

## 🧪 Testing

Run the full unit test suite across all 10 packages:

```bash
go test ./... -v
```

Run static analysis:

```bash
go vet ./...
```

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
