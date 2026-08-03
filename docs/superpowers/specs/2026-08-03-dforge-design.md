# dforge — Design Specification

> One command to bootstrap, inspect and manage Docker Compose projects.

**Status:** Approved design — ready for implementation planning
**Date:** 2026-08-03
**Source brief:** `brief.md`

---

## 1. Summary

dforge is a cross-platform, offline-first, rule-based (AI-free) CLI that improves the
Docker Compose developer experience: it bootstraps projects, generates best-practice
Docker artifacts, adds/removes infrastructure services, diagnoses configuration issues,
validates environment files, and wraps common Docker Compose runtime operations.

The guiding philosophy: **make Docker Compose feel effortless without hiding how Docker
actually works.**

### Non-goals

No AI, no cloud, no telemetry, no auth, no project hosting, no Kubernetes (initial version).

---

## 2. Locked design decisions

These were resolved during brainstorming and are binding for implementation:

| Decision | Choice |
|---|---|
| Scope of this spec | All 9 commands (`init`, `doctor`, `add`, `remove`, `env`, `status`, `logs`, `shell`, `clean`). Implementation plan is phased internally. |
| Compose file editing | `yaml.v3` **Node API** — preserve user comments, key ordering, and formatting. Write a `.bak` backup before mutating. |
| Interactivity | **Full Bubble Tea + Lipgloss TUI** in MVP. All logic lives in pure `internal/` packages; the TUI is a thin layer. Every interactive flow has a flag-based non-interactive equivalent and auto-falls back when stdout is not a TTY. |
| Generated Dockerfiles | **Production-grade**: multi-stage, non-root user, pinned base image tags (never `latest`), `HEALTHCHECK`, accompanying `.dockerignore`. dforge's own output must score highly against `doctor`. |
| `doctor` scoring | **Weighted-by-severity** via a check registry. `Score = 100 − Σ penalties`, clamped to `[0,100]`. |
| Distribution | Module `github.com/yunkhngn/dforge`. Cross-platform binaries + Homebrew tap via **GoReleaser** on tagged releases. |

---

## 3. Tech stack

- **Go** (single static binary, `go:embed` for templates)
- **Cobra** — command structure
- **Viper** — configuration (optional user config, e.g. default base images)
- **gopkg.in/yaml.v3** — compose parsing/mutation via the Node API
- **Bubble Tea** — interactive TUI flows
- **Lipgloss** — terminal styling / theme
- **GoReleaser** — release + Homebrew tap automation

---

## 4. Architecture

Guiding rule: **all logic lives in pure `internal/` packages; `cmd/` and `internal/ui/`
are thin shells.** This makes every command testable without a Docker daemon and without
a TTY.

```
dforge/
├── cmd/                 # Cobra commands: parse flags, call internal, render via ui
├── internal/
│   ├── detect/          # framework detection (pure, signal-based)
│   ├── generator/       # render Dockerfile / compose / .dockerignore from embedded templates
│   ├── compose/         # parse + mutate compose.yaml (yaml.Node), backup, diff
│   ├── services/        # catalog of addable services (image, ports, env, volumes, healthcheck)
│   ├── doctor/          # check registry + weighted scoring engine
│   ├── env/             # .env parsing, compare to .env.example, redaction
│   ├── docker/          # runtime wrapper behind a Client interface (exec docker/compose)
│   └── ui/              # Bubble Tea models + Lipgloss theme (prompts, spinners, tables)
├── templates/           # //go:embed  dockerfiles/, compose/, services/
├── configs/             # default configuration
├── docs/
├── testdata/            # fixture projects + golden files
└── main.go
```

### Command → layer mapping

| Command | Primary packages | Needs Docker daemon? |
|---|---|---|
| `init` | detect, generator, ui | No |
| `add` / `remove` | compose, services, ui | No |
| `doctor` | doctor, compose, env | No |
| `env` | env, ui | No |
| `status` / `logs` / `shell` / `clean` | docker, ui | **Yes** |

The `docker` package is the only impure layer. It exposes a `Client` interface with a
real exec-based implementation; runtime commands depend on the interface so their logic
is mockable.

---

## 5. Component specifications

### 5.1 detect

Signal-based framework detection. Inspects marker files and returns ranked candidates
with a confidence signal:

| Framework | Primary signals |
|---|---|
| Next.js | `package.json` dependency `next` |
| React | `package.json` dependency `react` (without `next`) |
| NestJS | `package.json` dependency `@nestjs/core` |
| Express | `package.json` dependency `express` |
| Spring Boot | `pom.xml` / `build.gradle` with spring-boot plugin/dependency |
| Go | `go.mod` |
| Rust | `Cargo.toml` |

Behavior:
- Single confident match → use it.
- Ambiguous (multiple candidates) or none → Bubble Tea picker, or `--framework <name>`
  flag for non-interactive use.
- Monorepo: detect per top-level directory; user selects the target directory/service.

### 5.2 generator

Renders artifacts from `text/template` files embedded via `go:embed`. For each supported
framework it produces:
- a **multi-stage** Dockerfile (build stage → minimal runtime stage), running as a
  **non-root** user, with a **pinned** base image tag and a `HEALTHCHECK`;
- a `compose.yaml` service entry (build context, ports, env, restart policy, healthcheck);
- a `.dockerignore`.

Never overwrites an existing file without confirmation. `--force` / `--yes` overrides the
prompt (and is required when non-interactive).

### 5.3 compose

Loads `compose.yaml` into a `yaml.Node` tree and mutates in place so comments, key
ordering, and formatting survive round-trips.

- `add`: insert service node, add any required named volumes, add env entries if needed.
- `remove`: delete service node; prune named volumes and networks **only if unused by any
  remaining service**.
- Always writes a `<file>.bak` backup before overwriting.
- Exposes a human-readable diff of intended changes for confirmation.

### 5.4 services

Declarative catalog of addable services. Each entry defines: default **pinned** image,
exposed ports, environment variables, named volumes, and a healthcheck.

Supported: `postgres`, `mysql`, `redis`, `mongodb`, `rabbitmq`, `minio`, `mailpit`,
`meilisearch`.

### 5.5 doctor

A registry of `Check`s. Each check has: id, human title, severity
(`error` / `warn` / `info`), weight (penalty points), and a pure evaluation function over
collected project state.

Checks (from brief): Docker installed, Docker Compose installed, Dockerfile exists,
compose.yaml exists, .dockerignore exists, HEALTHCHECK present, restart policy set, image
tag is not `latest`, port conflicts, missing volumes, missing environment file, running as
root (from generated/inspected Dockerfile).

Scoring: `Score = 100 − Σ(penalty of each failing check)`, clamped to `[0,100]`. Output
groups results with ✓ / ⚠ / ✗ symbols followed by the score, e.g.:

```
✓ Docker installed
✓ compose.yaml
⚠ Missing HEALTHCHECK
⚠ Running as root
⚠ Using latest tag

Score: 82/100
```

### 5.6 env

Parses `.env` and `.env.example`. Reports:
- missing variables (in example, absent in `.env`)
- extra variables (in `.env`, absent in example)
- duplicate keys
- empty values
- missing `.env` file entirely

**Safety: only variable names are ever printed. Values are always redacted.**

### 5.7 docker (runtime wrapper)

`Client` interface wrapping `docker` / `docker compose` invocations. Real implementation
shells out; a mock is used in tests.

- `status`: list services with running/health state and exposed ports.
- `logs <service>`: stream/print logs for a service (future: follow, level filter, colors).
- `shell <service>`: exec into a container, auto-detecting `bash` then falling back to `sh`.
- `clean`: remove dangling images, unused volumes, unused networks, build cache.
  **Requires confirmation unless `--force`.**

Degrades gracefully when Docker is not installed or the daemon is not running: clear
message, non-zero exit.

### 5.8 ui

A single Lipgloss theme (colors, ✓/⚠/✗ symbols, spacing) used by all commands for
predictable, human-friendly output. Bubble Tea drives interactive flows: framework picker
(`init`), service picker (`add`), and confirmations (`clean`, overwrite). When stdout is
not a TTY, commands automatically fall back to non-interactive behavior driven by flags.

---

## 6. Cross-cutting: safety rules

- Never overwrite existing files without confirmation.
- Never expose environment variable values.
- Never execute destructive Docker commands automatically; require confirmation before
  cleaning resources (`--force` to bypass intentionally).
- Create a `.bak` backup before modifying compose files.

---

## 7. Error handling

- File-not-found / permission errors surface with actionable messages and non-zero exit.
- Runtime commands detect missing Docker via the `docker` package and report clearly
  rather than emitting raw exec errors.
- Ambiguous detection and overwrite conflicts are recoverable (prompt) or fail loudly with
  guidance when non-interactive.

---

## 8. Testing strategy

- **Pure packages** (`detect`, `generator`, `compose`, `services`, `doctor`, `env`):
  table-driven unit tests. Generated Dockerfiles, compose files, and doctor reports are
  verified with **golden files** in `testdata/`.
- **compose** mutation: fixture-in → assert output preserves comments and ordering.
- **docker** package: command logic tested against a mocked `Client`. Real integration
  tests are gated behind a build tag (`//go:build integration`) and skipped by default.

---

## 9. Success criteria

A developer can:
- Bootstrap a Docker Compose project in under one minute.
- Diagnose common Docker configuration issues (with a score).
- Add common infrastructure services with a single command.
- Validate Docker environment configuration safely (names only, no values).
- Manage Docker Compose projects through a clean, consistent CLI.

dforge's own generated artifacts pass its own `doctor` with a high score.

---

## 10. Out of scope for this spec (future roadmap)

- v1.1: template marketplace, multiple compose profiles, project presets.
- v1.2: image optimization suggestions, compose visualization, health reports, export
  project summary.
- Kubernetes support.
