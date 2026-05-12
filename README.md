# DEV.ENV

A CLI tool that manages shared development services (databases, caches, search engines) across all your projects — powered by Docker.

## What is this?

Every project needs MySQL, Redis, Elasticsearch, and so on. The typical solution is a `docker-compose.yml` per project, which means a separate container per project per service — wasted memory, slow startup, and duplicate configuration everywhere.

DEV.ENV runs **one shared pool of containers** for all your projects. Services that support multi-tenancy (like MySQL and Redis) isolate projects using their own built-in mechanisms: a dedicated database per project, a separate Redis database index, a scoped Elasticsearch index. DEV.ENV provisions all of that automatically when you run `dev start`.

When sharing is not appropriate — for example, a service with a custom configuration or a project that needs full isolation — you can mark any service as `dedicated`. DEV.ENV then spins up a container exclusively for that project, managed separately from the shared pool.

Any Docker image can be used, not just the ones DEV.ENV natively supports. For unsupported images, DEV.ENV cannot handle provisioning or project isolation, so on first use it will prompt you to acknowledge the implications before starting.

Written in Go. Ships as a single binary with no runtime dependencies.

## How it works

- One global `docker-compose.yml` lives in `~/.dev.env/` and manages all shared service containers.
- Each project declares which services it needs in a `.dev.env.yaml` file at the project root.
- `dev start` ensures the required containers are running, then provisions project-specific resources (database, user, index, etc.). If `env: true` is set, it also writes a `.env` file with the correct connection strings.
- `dev stop` tracks which projects are using each container. A container only stops when no other active project depends on it.
- Project type detection (based on files like `composer.json` or `package.json`) determines which `.env` template to use when `env: true` is set.

## Features

- **Shared containers, isolated projects.** One MySQL container serves all your projects. Each project gets its own database and user, created automatically.
- **Dedicated containers when you need them.** Mark any service as `dedicated: true` to give that project its own private container, fully isolated from the shared pool.
- **Any Docker image, not just supported ones.** Use any image from Docker Hub. DEV.ENV prompts you once per project to acknowledge the limitations of unsupported images before starting them.
- **Declarative per-project config.** Define required services in `.dev.env.yaml`. No docker-compose knowledge needed.
- **Automatic provisioning.** `dev start` sets up project-scoped resources on first run (databases, users, index env_mappings, etc.).
- **Smart stop behavior.** Containers stay running as long as another project needs them. No accidental shutdowns.
- **Optional `.env` generation.** Set `env: true` in `.dev.env.yaml` to have DEV.ENV write a ready-to-use `.env` with correct hostnames, ports, credentials, and database names. Disabled by default.
- **Conflict-free port allocation.** Each `image:tag` pair gets a deterministic host port derived from the service's default port and the version tag, so multiple versions of the same service never collide. Override any port globally in `~/.dev.env/settings.yaml`.
- **Single binary.** Install once, works globally across all your projects.

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) (with the Docker daemon running)

### Installation

#### Linux / macOS

```sh
curl -sSL https://get.devenv.sh | sh
```

#### Windows

```sh
# TODO
```

## Usage

### setup

Initialize the global `~/.dev.env/` directory structure. Run this once after installation. Any other command will also trigger setup automatically if the directory does not exist yet.

```sh
dev setup
```

```
✓ Created ~/.dev.env/
✓ Created ~/.dev.env/docker/
✓ Created ~/.dev.env/projects/
✓ services.yaml initialized
✓ settings.yaml initialized
```

### init

Initialize DEV.ENV for a project. Runs an interactive wizard that detects your project type, lets you select services, and writes a `.dev.env.yaml` in the current directory.

```sh
dev init
```

```
Detected project type: laravel

Select services (space to toggle, enter to confirm):
  [x] mysql
  [x] redis
  [ ] elasticsearch
  [ ] rabbitmq
  [ ] ...

MySQL tag [latest]: 8.0
MySQL mode: shared or dedicated? [shared]:

Redis tag [latest]:
Redis mode: shared or dedicated? [shared]:

Enable .env generation? [y/N]: y

Writing .dev.env.yaml... done
```

Steps:
1. Detect project type from files on disk (`composer.json` with `laravel/framework` → laravel, `package.json` → node). Show the detected type and let you confirm or override.
2. Show all available services. Toggle which ones to include.
3. For each selected service: prompt for tag (default: latest) and shared vs. dedicated mode.
4. Ask whether to enable `.env` generation.
5. Write `.dev.env.yaml`.

### start

Start all services required by this project and prepare them for use.

```sh
dev start
```

```
✓ mysql:8.0              already running (shared)
✓ redis:latest           started (shared)
✓ elasticsearch:8.11     started (dedicated)
✓ Database "my-project" created
✓ .env written            (only when env: true)
```

On first run, DEV.ENV provisions project-specific resources (creates DB, user, index, etc.). On subsequent runs it checks that everything is still in order. If `env: true` is set in `.dev.env.yaml`, a `.env` file is written on every run.

Shared containers already running for another project are reused. Dedicated containers are always project-private and never shared.

### stop

Stop services that this project was using.

```sh
dev stop
```

```
✓ redis:latest           stopped
✓ elasticsearch:8.11     stopped (dedicated)
~ mysql:8.0              kept running (in use by: api-project)
```

Shared containers only stop when no other active project depends on them. Dedicated containers always stop with the project that owns them.

### status

Show all containers managed by DEV.ENV and which projects are using them.

```sh
dev status
```

```
SERVICE                  STATUS     PROJECTS
mysql:8.0                running    my-project, api-project
redis:latest             running    my-project
elasticsearch:8.11       running    api-project (dedicated)
```

### exec

Open an interactive session inside a running container for the current project, with credentials and database pre-filled.

```sh
dev exec <service>   # connect to a specific service
dev exec             # list available services for this project
```

`<service>` is the image name as declared in `.dev.env.yaml` (e.g. `mysql`, `redis`, `postgres`).

**Examples:**

```sh
dev exec mysql       # opens mysql REPL connected to the project database
dev exec redis       # opens redis-cli on the project's assigned DB index
dev exec postgres    # opens psql connected to the project database
```

**What DEV.ENV opens per service:**

| Service | Session |
|---------|---------|
| `mysql` / `mariadb` / `percona` | `mysql` REPL — root user, project database, password pre-filled |
| `postgres` | `psql` — project user, project database, password pre-filled |
| `mongo` | `mongosh` — root user, password pre-filled |
| `redis` | `redis-cli` — connected to the project's assigned DB index |
| All others | `bash` shell inside the container |

The project must be running (`dev start`) before using `dev exec`.

### config

Read and write project configuration without editing `.dev.env.yaml` directly.

```sh
dev config get <key>
dev config set <key> <value>
dev config list
```

`dev config list` shows all current project config:

```
project   my-project-name
type      laravel
env       true
services  mysql:8.0, redis:latest, elasticsearch:8.11 (dedicated)
```

`dev config get` reads a single value:

```sh
dev config get type
# laravel
```

`dev config set` writes a value back to `.dev.env.yaml`:

```sh
dev config set env true
dev config set type node
```

For nested service config, use dot notation:

```sh
dev config set services.mysql.tag 8.4
dev config set services.mysql.env_map.DB_HOST DATABASE_HOST
```

### global

Run commands across all projects and containers, not just the current project.

```sh
dev global <command>
```

#### global start

Start all containers for all known projects.

```sh
dev global start
```

```
✓ mysql:8.0              started (shared)
✓ redis:latest           started (shared)
✓ elasticsearch:8.11     started (dedicated — api-project)
✓ elasticsearch:8.11     started (dedicated — search-project)
```

#### global stop

Stop all running containers managed by DEV.ENV, regardless of which projects are using them.

```sh
dev global stop
```

```
✓ redis:latest           stopped
✓ elasticsearch:8.11     stopped (dedicated — api-project)
✓ elasticsearch:8.11     stopped (dedicated — search-project)
✓ mysql:8.0              stopped
```

Unlike `dev stop`, `dev global stop` does not check reference counts — it stops everything immediately.

### self-update

Update DEV.ENV to the latest version.

```sh
dev self-update
```

## Project Configuration

Place a `.dev.env.yaml` at the root of your project:

```yaml
project: my-project-name   # used to scope databases, indices, etc.
type: laravel              # determines which .env template to generate
env: true                  # write a .env file on dev start (default: false)
services:
  - image: mysql
    tag: 8.0               # resolves to mysql:8.0
    env_map:                   # rename generated variables to custom names
      DB_HOST: DATABASE_HOST
      DB_DATABASE: DATABASE_NAME
      DB_USERNAME: DATABASE_USER
      DB_PASSWORD: DATABASE_PASS
  - redis                  # shorthand; resolves to redis:latest
  - image: elasticsearch
    tag: 8.11
    dedicated: true        # spin up a private container for this project only
  - image: meilisearch     # unsupported image; requires one-time acknowledgement per project
    tag: v1.7
```

### Field reference

| Field | Required | Default | Description |
|---|---|---|---|
| `project` | Yes | — | Unique project identifier. Used as the database name, Redis key prefix, etc. |
| `type` | No | — | Project type. Used by `dev init` for auto-detection; has no effect on `.env` generation. |
| `env` | No | `false` | If `true`, write a `.env` file with connection strings on `dev start`. |
| `services` | Yes | — | List of services. Can be a string (image name) or an object with `image`, `tag`, and `dedicated`. |

**Service object fields:**

| Field | Required | Default | Description |
|---|---|---|---|
| `image` | Yes | — | Docker image name (e.g. `mysql`, `redis`). |
| `tag` | No | `latest` | Docker image tag (e.g. `8.0`). |
| `dedicated` | No | `false` | If `true`, spin up a private container for this project instead of using the shared pool. |
| `env_map` | No | — | Map of canonical variable names to custom names written to `.env`. See [Variable env_mapping](#variable-env_mapping). |

**Local-override-only fields** (`.dev.env.local.yaml` only, not valid in `.dev.env.yaml`):

| Field | Required | Default | Description |
|---|---|---|---|
| `port` | No | — | Override the host port for this service. **Forces `dedicated: true`** — the service runs as a private container for this project. Not recommended; prefer a global port override in `~/.dev.env/settings.yaml`. |

### Variable env_mapping

By default, DEV.ENV writes variables using the canonical names documented per service (e.g. `DB_HOST`, `DB_DATABASE`). If your application uses different names, use `env_map` to rename them:

```yaml
services:
  - image: mysql
    env_map:
      DB_HOST: DATABASE_HOST
      DB_DATABASE: DATABASE_NAME
      DB_USERNAME: DATABASE_USER
      DB_PASSWORD: DATABASE_PASS
```

Only the variables listed in `env_map` are renamed. Unmapped variables are written using their canonical names. `env_map` has no effect when `env: false`.

### `.env` generation

When `env: true` is set, `dev start` writes a `.env` file at the project root.

**First run (no existing `.env`):** the file is written immediately with no prompts.

**Subsequent runs (`.env` already exists):** DEV.ENV compares what it would write against the current file. If nothing has changed it prints `✓ .env up to date` and stops. If any variable is new or has a different value, it prints a diff and asks for confirmation before writing:

```
Changes to .env:
  + REDIS_DB=2                    (new)
  ~ DB_PORT    3300 → 3380        (changed)

Update .env file? [Y/n]
```

Only the variables DEV.ENV manages are updated. Any variables you have added manually are preserved in a separate section of the file and are never touched.

DEV.ENV generates variables per service based on the declared `services` list. Each service emits its own set of canonical variables:

| Service | Variables written |
|---------|-------------------|
| `mysql` / `mariadb` / `percona` | `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD` |
| `postgres` | `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD` |
| `mongo` | `MONGODB_HOST`, `MONGODB_PORT`, `MONGODB_DATABASE`, `MONGODB_USERNAME`, `MONGODB_PASSWORD`, `MONGODB_URI` |
| `redis` | `REDIS_HOST`, `REDIS_PORT`, `REDIS_DB` |
| `memcached` | `MEMCACHED_HOST`, `MEMCACHED_PORT` |
| `elasticsearch` / `opensearch` | `ELASTICSEARCH_HOST`, `ELASTICSEARCH_INDEX` |
| `meilisearch` | `MEILISEARCH_HOST`, `MEILISEARCH_KEY`, `MEILISEARCH_INDEX` |
| `typesense` | `TYPESENSE_HOST`, `TYPESENSE_PORT`, `TYPESENSE_PROTOCOL`, `TYPESENSE_API_KEY`, `TYPESENSE_COLLECTION` |
| `solr` | `SOLR_HOST`, `SOLR_PORT`, `SOLR_CORE` |
| `cassandra` | `CASSANDRA_HOST`, `CASSANDRA_PORT`, `CASSANDRA_KEYSPACE` |
| `rabbitmq` | `RABBITMQ_HOST`, `RABBITMQ_PORT`, `RABBITMQ_VHOST`, `RABBITMQ_USERNAME`, `RABBITMQ_PASSWORD` |
| `kafka` | `KAFKA_BROKERS`, `KAFKA_TOPIC_PREFIX` |
| `minio` | `MINIO_ENDPOINT`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `MINIO_BUCKET`, `MINIO_USE_SSL` |
| `soketi` | `PUSHER_HOST`, `PUSHER_PORT`, `PUSHER_SCHEME`, `PUSHER_APP_ID`, `PUSHER_APP_KEY`, `PUSHER_APP_SECRET` |
| `reverb` | `REVERB_HOST`, `REVERB_PORT`, `REVERB_SCHEME`, `REVERB_APP_ID`, `REVERB_APP_KEY`, `REVERB_APP_SECRET` |
| `mailpit` / `mailhog` | `MAIL_HOST`, `MAIL_PORT`, `MAIL_MAILER` |

Ports in the `.env` file always reflect the resolved host port (see [Port Allocation](#port-allocation)). Database names, indices, vhosts, and buckets all use the project name. Redis DB index is a unique integer (0–15) assigned per project.

Credentials (passwords, API keys, app secrets) are generated once on first start and reused on every subsequent run — connection strings stay stable across restarts.

The written `.env` has two sections:

```
# Generated by dev — do not edit this section
DB_HOST=127.0.0.1
DB_PORT=3380
...

# User-defined variables
APP_KEY=base64:abc...
APP_ENV=local
```

The first section is managed by DEV.ENV. The second section preserves any variables you have added manually.

### Local overrides (`.dev.env.local.yaml`)

Any field in `.dev.env.yaml` can be overridden locally by placing a `.dev.env.local.yaml` file at the project root. DEV.ENV merges both files at runtime; local values take precedence. This file should be gitignored.

```yaml
# .dev.env.local.yaml
project: my-project-local   # override project name
services:
  mysql:
    tag: "8.4"              # use a different tag locally
    port: 13306             # override host port (forces dedicated — not recommended)
```

Add to `.gitignore`:

```
.dev.env.local.yaml
.env
```

**Project name collision.** When `dev start` detects that the `project:` name is already registered by a different project at a different path, it prompts:

```
! Project name "my-project" is already in use by: /home/user/other-project

  Using the same name means sharing the same database, credentials, and state.
  If this is intentional (e.g. a monorepo), continue. Otherwise, choose a new name.

  New project name [my-project]: my-project-local
```

The chosen name is written automatically to `.dev.env.local.yaml` so `.dev.env.yaml` stays unchanged and version-control-safe.

## Port Allocation

DEV.ENV assigns each `image:tag` pair a deterministic host port so that running multiple versions of the same service simultaneously never causes conflicts.

### How ports are computed

The host port is derived from the service's canonical port and the version tag:

1. Parse the major and minor version from the tag (e.g. `8.0` → `80`, `5.7` → `57`, `7` → `70`, `latest` → `00`).
2. Drop the last two digits of the base port and append the version suffix.

Examples:

| Service | Tag | Base port | Suffix | Host port |
|---------|-----|-----------|--------|-----------|
| `mysql` | `8.0` | 3306 | `80` | `3380` |
| `mysql` | `5.7` | 3306 | `57` | `3357` |
| `postgres` | `16` | 5432 | `16` | `5416` |
| `redis` | `7` | 6379 | `70` | `6370` |
| `elasticsearch` | `8.11` | 9200 | `81` | `9281` |
| `mongo` | `latest` | 27017 | `00` | `27000` |

> Port numbers follow the pattern `<base family><version suffix>`. The "family" (33xx for MySQL, 54xx for PostgreSQL, etc.) stays recognisable.

`dev start` validates that no two services in the current run resolve to the same host port and exits with an error if a collision is detected.

### Overriding ports globally

Add a `ports` map to `~/.dev.env/settings.yaml`, keyed by docker-compose service name:

```yaml
# ~/.dev.env/settings.yaml
default_type: laravel
ports:
  mysql-8-0: 13380
  redis-7-0: 16370
```

Global overrides apply to the shared container and affect all projects that use that service. Use this when the computed default conflicts with another application already running on your machine.

### Overriding ports per project

> **Not recommended.** A per-project port override forces the service to run as a dedicated container for that project, bypassing the shared pool that DEV.ENV is optimised for. Prefer a global override in `~/.dev.env/settings.yaml` instead.

Add `port` to a service entry in `.dev.env.local.yaml`:

```yaml
# .dev.env.local.yaml
services:
  mysql:
    port: 13306
```

DEV.ENV automatically treats the service as `dedicated: true` when a project-level port is set. The dedicated container runs on the overridden port and is fully isolated from the shared instance. This option exists for cases where per-developer port customisation is required and a global override is not practical.

## Global Configuration (`~/.dev.env/`)

DEV.ENV stores all global state in `~/.dev.env/`:

```
~/.dev.env/
├── settings.yaml              # Global user preferences
├── services.yaml              # Registry of all known services and their docker-compose names
├── projects/
│   └── my-project/
│       ├── project.yaml       # Project metadata: disk path and location of .dev.env.yaml
│       ├── state.yaml         # Runtime state: running status, active services, acknowledged images
│       └── secrets.yaml       # Generated credentials: passwords, API keys, app secrets
└── docker/
    └── docker-compose.yml     # Managed by DEV.ENV — do not edit manually
```

### `settings.yaml`

Global user preferences. Currently supports:

| Key | Type | Description |
|-----|------|-------------|
| `default_type` | string | Default project type used when `type` is not set in `.dev.env.yaml`. |
| `ports` | map | Override computed host ports by docker-compose service name (e.g. `mysql-8-0: 13380`). See [Port Allocation](#port-allocation). |

### `services.yaml`

Tracks every `image:tag` pair that has ever been started and assigns each a unique docker-compose service name. This prevents naming collisions when multiple tags of the same image are in use simultaneously (e.g. `mysql:8.0` and `mysql:8.4`).

Service names are derived deterministically: image name + tag with dots replaced by dashes.

```yaml
services:
  mysql-8-0:
    image: mysql
    tag: "8.0"
  mysql-8-4:
    image: mysql
    tag: "8.4"
  redis-latest:
    image: redis
    tag: latest
  elasticsearch-8-11:
    image: elasticsearch
    tag: "8.11"
```

When `dev start` encounters an `image:tag` pair not yet registered, it adds an entry to `services.yaml` and regenerates `docker/docker-compose.yml`. Docker-compose service names map directly to the keys in this file.

### `projects/<name>/project.yaml`

Written by `dev init`. Updated automatically if the project is moved.

```yaml
name: my-project
path: /home/user/code/my-project
```

### `projects/<name>/state.yaml`

Updated by `dev start` and `dev stop`. Records which services this project is actively using and whether any unsupported images have been acknowledged.

```yaml
running: true
services:
  mysql-8-0: shared
  redis-latest: shared
  elasticsearch-8-11: dedicated
acknowledged:
  - meilisearch:v1.7
```

DEV.ENV determines whether a shared container can be stopped by scanning the `state.yaml` of all known projects. A container stops only when no project lists it as an active service.

### `projects/<name>/secrets.yaml`

Stores generated credentials written during the first `dev start`. Never overwritten on subsequent runs — this ensures connection strings in `.env` stay stable across restarts.

```yaml
mysql:
  database: my-project
  username: my-project
  password: "xK9mP2wR..."
rabbitmq:
  vhost: my-project
  username: my-project
  password: "aQ7rL1nZ..."
meilisearch:
  master_key: "mK3nV8pT..."
soketi:
  app_id: "84729"
  app_key: "uW2tH5sX..."
  app_secret: "zJ6pN9qY..."
```

### `docker/docker-compose.yml`

Generated and maintained by DEV.ENV. Rebuilt whenever `services.yaml` gains a new entry. Service names in this file match the keys in `services.yaml`. Not intended for manual editing.

## Supported Services

Supported services can run shared (with automatic project isolation) or dedicated. All services support `dedicated: true` to force a private container.

### Databases

#### MySQL

[Docker Hub](https://hub.docker.com/_/mysql) · [Documentation](https://dev.mysql.com/doc/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate database + user per project | Creates DB, user, and grants privileges |

| Variable | Example value | Description |
|---|---|---|
| `DB_HOST` | `127.0.0.1` | Hostname of the MySQL container |
| `DB_PORT` | `3306` | Exposed port |
| `DB_DATABASE` | `my-project` | Created database, matches project name |
| `DB_USERNAME` | `my-project` | Created user, matches project name |
| `DB_PASSWORD` | `<generated>` | Randomly generated password |

#### MariaDB

[Docker Hub](https://hub.docker.com/_/mariadb) · [Documentation](https://mariadb.org/documentation/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate database + user per project | Creates DB, user, and grants privileges |

Same environment variables as [MySQL](#mysql).

#### Percona

[Docker Hub](https://hub.docker.com/_/percona) · [Documentation](https://docs.percona.com/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate database + user per project | Creates DB, user, and grants privileges |

Same environment variables as [MySQL](#mysql).

#### PostgreSQL

[Docker Hub](https://hub.docker.com/_/postgres) · [Documentation](https://www.postgresql.org/docs/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate database + user per project | Creates DB and user |

| Variable | Example value | Description |
|---|---|---|
| `DB_HOST` | `127.0.0.1` | Hostname of the PostgreSQL container |
| `DB_PORT` | `5432` | Exposed port |
| `DB_DATABASE` | `my-project` | Created database, matches project name |
| `DB_USERNAME` | `my-project` | Created user, matches project name |
| `DB_PASSWORD` | `<generated>` | Randomly generated password |

#### MongoDB

[Docker Hub](https://hub.docker.com/_/mongo) · [Documentation](https://www.mongodb.com/docs/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate database per project | Creates DB and user |

| Variable | Example value | Description |
|---|---|---|
| `MONGODB_HOST` | `127.0.0.1` | Hostname of the MongoDB container |
| `MONGODB_PORT` | `27017` | Exposed port |
| `MONGODB_DATABASE` | `my-project` | Created database, matches project name |
| `MONGODB_USERNAME` | `my-project` | Created user, matches project name |
| `MONGODB_PASSWORD` | `<generated>` | Randomly generated password |
| `MONGODB_URI` | `mongodb://my-project:<password>@127.0.0.1:27017/my-project` | Full connection string |

#### Cassandra

[Docker Hub](https://hub.docker.com/_/cassandra) · [Documentation](https://cassandra.apache.org/doc/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate keyspace per project | Creates keyspace with default replication |

| Variable | Example value | Description |
|---|---|---|
| `CASSANDRA_HOST` | `127.0.0.1` | Hostname of the Cassandra container |
| `CASSANDRA_PORT` | `9042` | Exposed port |
| `CASSANDRA_KEYSPACE` | `my_project` | Created keyspace, derived from project name |

---

### Search Engines

#### Elasticsearch

[Docker Hub](https://hub.docker.com/_/elasticsearch) · [Documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate index per project | Creates index with default env_mappings |

| Variable | Example value | Description |
|---|---|---|
| `ELASTICSEARCH_HOST` | `http://127.0.0.1:9200` | Full base URL of the Elasticsearch container |
| `ELASTICSEARCH_INDEX` | `my-project` | Created index, matches project name |

#### OpenSearch

[Docker Hub](https://hub.docker.com/r/opensearchproject/opensearch) · [Documentation](https://opensearch.org/docs/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate index per project | Creates index with default env_mappings |

| Variable | Example value | Description |
|---|---|---|
| `OPENSEARCH_HOST` | `http://127.0.0.1:9200` | Full base URL of the OpenSearch container |
| `OPENSEARCH_INDEX` | `my-project` | Created index, matches project name |

#### Meilisearch

[Docker Hub](https://hub.docker.com/r/getmeili/meilisearch) · [Documentation](https://www.meilisearch.com/docs/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate index per project | Creates index |

| Variable | Example value | Description |
|---|---|---|
| `MEILISEARCH_HOST` | `http://127.0.0.1:7700` | Full base URL of the Meilisearch container |
| `MEILISEARCH_KEY` | `<generated>` | Master key used to authenticate requests |
| `MEILISEARCH_INDEX` | `my-project` | Created index, matches project name |

#### Typesense

[Docker Hub](https://hub.docker.com/r/typesense/typesense) · [Documentation](https://typesense.org/docs/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate collection per project | Creates collection |

| Variable | Example value | Description |
|---|---|---|
| `TYPESENSE_HOST` | `127.0.0.1` | Hostname of the Typesense container |
| `TYPESENSE_PORT` | `8108` | Exposed port |
| `TYPESENSE_PROTOCOL` | `http` | Connection protocol |
| `TYPESENSE_API_KEY` | `<generated>` | API key used to authenticate requests |
| `TYPESENSE_COLLECTION` | `my-project` | Created collection, matches project name |

#### Solr

[Docker Hub](https://hub.docker.com/_/solr) · [Documentation](https://solr.apache.org/resources.html)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate core per project | Creates core |

| Variable | Example value | Description |
|---|---|---|
| `SOLR_HOST` | `127.0.0.1` | Hostname of the Solr container |
| `SOLR_PORT` | `8983` | Exposed port |
| `SOLR_CORE` | `my-project` | Created core, matches project name |

---

### Caches

#### Redis

[Docker Hub](https://hub.docker.com/_/redis) · [Documentation](https://redis.io/docs/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate DB index per project | Assigns and records a DB index |

| Variable | Example value | Description |
|---|---|---|
| `REDIS_HOST` | `127.0.0.1` | Hostname of the Redis container |
| `REDIS_PORT` | `6379` | Exposed port |
| `REDIS_DB` | `3` | Assigned database index, unique per project |

#### Memcached

[Docker Hub](https://hub.docker.com/_/memcached) · [Documentation](https://memcached.org/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| Yes | No native namespacing | Starts container only |

| Variable | Example value | Description |
|---|---|---|
| `MEMCACHED_HOST` | `127.0.0.1` | Hostname of the Memcached container |
| `MEMCACHED_PORT` | `11211` | Exposed port |

---

### Queues

#### RabbitMQ

[Docker Hub](https://hub.docker.com/_/rabbitmq) · [Documentation](https://www.rabbitmq.com/documentation.html)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate vhost per project | Creates vhost and user |

| Variable | Example value | Description |
|---|---|---|
| `RABBITMQ_HOST` | `127.0.0.1` | Hostname of the RabbitMQ container |
| `RABBITMQ_PORT` | `5672` | Exposed AMQP port |
| `RABBITMQ_VHOST` | `my-project` | Created vhost, matches project name |
| `RABBITMQ_USERNAME` | `my-project` | Created user, matches project name |
| `RABBITMQ_PASSWORD` | `<generated>` | Randomly generated password |

#### Kafka

[Docker Hub](https://hub.docker.com/r/apache/kafka) · [Documentation](https://kafka.apache.org/documentation/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate topic prefix per project | Assigns and records a topic prefix |

| Variable | Example value | Description |
|---|---|---|
| `KAFKA_BROKERS` | `127.0.0.1:9092` | Bootstrap broker address |
| `KAFKA_TOPIC_PREFIX` | `my-project` | Topic prefix assigned to this project |

---

### Storage

#### MinIO

[Docker Hub](https://hub.docker.com/r/minio/minio) · [Documentation](https://min.io/docs/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate bucket per project | Creates bucket |

| Variable | Example value | Description |
|---|---|---|
| `MINIO_ENDPOINT` | `http://127.0.0.1:9000` | Full base URL of the MinIO container |
| `MINIO_ACCESS_KEY` | `<generated>` | Access key for authentication |
| `MINIO_SECRET_KEY` | `<generated>` | Secret key for authentication |
| `MINIO_BUCKET` | `my-project` | Created bucket, matches project name |
| `MINIO_USE_SSL` | `false` | SSL disabled for local development |

---

### Broadcasting

#### Soketi

[Docker Hub](https://quay.io/repository/soketi/soketi) · [Documentation](https://docs.soketi.app/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate app ID + secret per project | Creates app credentials |

| Variable | Example value | Description |
|---|---|---|
| `PUSHER_HOST` | `127.0.0.1` | Hostname of the Soketi container |
| `PUSHER_PORT` | `6001` | Exposed port |
| `PUSHER_SCHEME` | `http` | Connection protocol |
| `PUSHER_APP_ID` | `<generated>` | App ID for this project |
| `PUSHER_APP_KEY` | `<generated>` | App key for this project |
| `PUSHER_APP_SECRET` | `<generated>` | App secret for this project |

#### Reverb

[Documentation](https://laravel.com/docs/reverb)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | Separate app ID + secret per project | Creates app credentials |

| Variable | Example value | Description |
|---|---|---|
| `REVERB_HOST` | `127.0.0.1` | Hostname of the Reverb container |
| `REVERB_PORT` | `8080` | Exposed port |
| `REVERB_SCHEME` | `http` | Connection protocol |
| `REVERB_APP_ID` | `<generated>` | App ID for this project |
| `REVERB_APP_KEY` | `<generated>` | App key for this project |
| `REVERB_APP_SECRET` | `<generated>` | App secret for this project |

---

### Mail

#### Mailpit

[Docker Hub](https://hub.docker.com/r/axllent/mailpit) · [Documentation](https://mailpit.axllent.org/docs/)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | None — single shared inbox by design | Starts container only |

| Variable | Example value | Description |
|---|---|---|
| `MAIL_HOST` | `127.0.0.1` | Hostname of the Mailpit container |
| `MAIL_PORT` | `1025` | Exposed SMTP port |
| `MAIL_MAILER` | `smtp` | Mailer driver |

The Mailpit web UI is available at `http://127.0.0.1:8025`.

#### MailHog

[Docker Hub](https://hub.docker.com/r/mailhog/mailhog) · [Documentation](https://github.com/mailhog/MailHog)

| Dedicated only | Shared isolation | What DEV.ENV provisions |
|---|---|---|
| No | None — single shared inbox by design | Starts container only |

| Variable | Example value | Description |
|---|---|---|
| `MAIL_HOST` | `127.0.0.1` | Hostname of the MailHog container |
| `MAIL_PORT` | `1025` | Exposed SMTP port |
| `MAIL_MAILER` | `smtp` | Mailer driver |

The MailHog web UI is available at `http://127.0.0.1:8025`.

---

### Unsupported images

Any Docker Hub image can be used. When `dev start` encounters an unsupported image for the first time in a project, it pauses and prompts:

```
! meilisearch:v1.7 is not a supported service.

  DEV.ENV will start and stop this container, but:
  - No project isolation. All projects using this image share the same container.
  - No provisioning. No databases, users, or indices will be created.
  - No .env generation. You must add connection details to your .env manually.

  Acknowledge and continue? [y/N]
```

Answering yes saves the acknowledgement to the project's state file (`~/.dev.env/projects/<name>/state.json`) so you are not prompted again for that project. The container then runs shared like any other service — no isolation, no provisioning. If you need isolation, use `dedicated: true` explicitly.

## Supported Project Types

Project type is used by `dev init` for auto-detection only — it does not affect which `.env` variables are generated. Variables are always determined by the declared `services` list.

| Type | Detection |
|---|---|
| `laravel` | `composer.json` with `laravel/framework` |
| `node` | `package.json` |

> More types are planned. Detection logic is extensible.
