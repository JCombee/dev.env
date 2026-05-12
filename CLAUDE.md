# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Feature Development Workflow

All new features follow this mandatory 3-phase gate process. **Do not advance to the next phase without explicit user approval.**

### Phase 1 — Documentation
Plan the feature and write its documentation first. This includes:
- What the feature does and why
- CLI command/flag changes
- Config schema changes
- File layout changes
- Behavior edge cases

**Stop here. Wait for user approval before writing any code or tests.**

Once approved: update `README.md` to reflect the feature (commands, flags, config, behavior), then commit it before proceeding.

### Phase 2 — Tests
Write all tests before implementing:
- Unit tests for each internal package involved
- E2E tests covering the golden path and relevant edge cases (use the `e2e/` harness)
- Tests must be runnable but are expected to fail at this stage

**Stop here. Wait for user approval before implementing.**

### Phase 3 — Implementation
Implement the feature to make the tests pass. No scope creep beyond what the approved docs and tests define.

---

## Build & Test Commands

```bash
# Build (output always goes to build/)
go build -ldflags="-s -w -X github.com/jcombee/devenv/version.Version=<ver>" -o build/dev.exe .

# Unit tests only (exclude e2e)
go test $(go list ./... | grep -v /e2e) -v

# Run a single unit test
go test ./internal/compose/... -run TestParsePS_JSONArray -v

# E2E tests (builds dev + fakedocker binaries automatically)
go test ./e2e/ -v -count=1

# Lint
golangci-lint run ./...

# Vet
go vet ./...
```

**Policy:** Never run compiled binaries against the host system. All testing uses unit tests with temp dirs (`t.TempDir()`) or the e2e harness. Docker is not required for unit or e2e tests.

## Architecture

### Command Layer (`cmd/`)

Each cobra command delegates to a `run<Cmd>(runner compose.Runner)` function accepting the `Runner` interface. This enables unit-testing of command logic without Docker. The `PersistentPreRunE` on the root command calls `global.Setup()` automatically on first run.

Commands: `setup`, `init`, `start`, `stop`, `status`, `config` (list/get/set), `version`.

### Key Interfaces

**`compose.Runner`** (`internal/compose/runner.go`) — abstracts `docker compose` calls:
- `Up(composeFile string, services ...string) error`
- `Stop(composeFile string, services ...string) error`
- `PS(composeFile string) ([]ServiceStatus, error)`

`ExecRunner` (production) shells out to real Docker. E2E tests use a compiled `fakedocker` binary injected via PATH.

### State & Storage (`~/.dev.env/`)

```
~/.dev.env/
  settings.yaml          # global settings
  services.yaml          # map of composeName → {image, tag} for all tracked services
  docker/
    docker-compose.yml   # shared compose file (regenerated on every start)
    secrets.yaml         # admin credentials per container (generated once, never overwritten)
  projects/<name>/
    project.yaml         # metadata (name, path)
    state.yaml           # running: bool, services: map[composeName]"shared"|"dedicated"
    secrets.yaml         # per-project credentials
    docker-compose.yml   # dedicated services only
```

All YAML I/O goes through `internal/store` (`Read`/`Write`). `internal/global` owns the global dir structure. `internal/project` owns per-project files.

### Service Name Convention

`compose.ServiceName(image, tag)` → replaces `.` and `:` with `-` → e.g., `mysql:8.0` → `mysql-8-0`.
Dedicated services get project prefix: `myproject-mysql-8-0`.

### Reference Counting (`internal/refcount`)

Before stopping a shared container, `refcount.ActiveProjectsUsingService(composeName)` scans **all** `state.yaml` files. `stop.go` marks `state.Running = false` and saves it **before** calling refcount, so the stopping project is excluded from the count.

### Config (`internal/config`)

`.dev.env.yaml` (committed) merged with `.dev.env.local.yaml` (gitignored) — local overrides project name/type. `ServiceEntry` supports shorthand YAML (`"redis"` or `"redis:7"`) via custom unmarshaler.

Project type auto-detected (`internal/config/detect.go`): checks for `artisan` (Laravel), `package.json` (Node), etc.

### E2E Harness (`e2e/harness_test.go`)

`TestMain` builds both `dev` and `fakedocker` binaries into a shared temp dir. Each test gets a `Harness` with an isolated `HOME` and `FAKEDOCKER_DIR`. `fakedocker` logs every call as JSON; `h.DockerCalls()` reads that log. Exit codes tested via `h.AssertExitCode(err, code)`.

## Error Handling

Use `errors.As` not type assertions on errors (enforced by `errorlint` linter). Wrap errors with `fmt.Errorf("context: %w", err)`. Commands exit non-zero on all errors — no raw stack traces to users.

## Conventions

- Conventional commits required (`feat:`, `fix:`, `refactor:`, etc.) — used by `git-cliff` for CHANGELOG
- No build artifacts in repo root — all binaries in `build/`
- Tests must set `HOME` and `USERPROFILE` via `t.Setenv` when touching `~/.dev.env/`
- `internal/services/registry.go` is the single source of truth for supported images, ports, and provisioning categories
