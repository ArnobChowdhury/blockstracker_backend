# ![logo](./assets/32x32.png) BlocksTracker Backend

This is the backend service for [BlocksTracker](https://blocks-tracker.com/), a privacy-first habit and task tracking app. Built with Go, PostgreSQL, and Redis.

## 🛠 Tech Stack

- **Go** (Gin)
- **PostgreSQL**
- **Redis**
- **Docker Compose** for local development

## 🚀 Getting Started

Before running the services, you must create a configuration file for `air`, our live-reloading tool. An example file is provided for you. Copy it to create your local configuration:

```bash
cp air.example.toml .air.toml
```

Once that's done, you can build and run the containers:

```bash
docker-compose up --build -d
```

The backend will be available at `http://localhost:5000`.

## 🔑 Environment Variables

The application requires the following environment variables to be set. You can create a `.env` file in the root of the project, and `docker-compose` will automatically use it.

- `DB_USER`: The username for the PostgreSQL database.
- `DB_PASSWORD`: The password for the PostgreSQL database.
- `DB_NAME`: The name of the PostgreSQL database.
- `TEST_DB_NAME`: The name of the PostgreSQL database to use for integration tests.
- `JWT_ACCESS_SECRET`: The secret key for signing JWT access tokens.
- `JWT_REFRESH_SECRET`: The secret key for signing JWT refresh tokens.
- `REDIS_PASSWORD`: Password for the Redis server (leave empty if none).

## 📂 Project Structure

- `cmd/` – Main application entry points
- `internal/` – Core application logic
- `pkg/` – Shared packages
- `routes/` – HTTP routes and handlers
- `handlers/` – Handlers

## 🧰 Common Tasks

This project uses `task` for running common development scripts. You can see all available tasks by running `task --list`. Here are the primary commands:

| Task                    | Description                                    |
| ----------------------- | ---------------------------------------------- |
| `task migrate-up`       | Run database migrations (up).                  |
| `task migrate-down`     | Rollback the last database migration.          |
| `task test`             | Run all unit and integration tests.            |
| `task unit-test`        | Run only the unit tests.                       |
| `task integration-test` | Run only the integration tests.                |
| `task generate-swagger` | Generate/update the Swagger API documentation. |

### Creating a New Migration

To create a new database migration file, you first need to get a shell inside the running `go_app` container:

```bash
docker compose exec -it go_app /bin/bash
```

Then, from within the container's shell, run the `goose create` command, replacing `<your_migration_name>` with a descriptive name:

```bash
goose create <your_migration_name> sql

```

## 🏛️ Architectural Decisions

This project uses Architectural Decision Records (ADRs) to document important architectural choices, their context, and their consequences. You can find them in the `docs/adr` directory.

## ⚖️ License

BlocksTracker is free and open-source software, licensed under the [GPLv3](https://www.gnu.org/licenses/gpl-3.0.html).
