# ![logo](./assets/32x32.png) BlocksTracker Backend

This is the backend service for [BlocksTracker](https://blocks-tracker.com/), a privacy-first habit and task tracking app. Built with Go, PostgreSQL, and Redis.

## 🛠 Tech Stack

- **Go** (Gin)
- **PostgreSQL**
- **Redis**
- **Docker Compose** for local development

## 🚀 Getting Started

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

This project uses [`task`](https://taskfile.dev) for development scripts. Examples:

```bash
task migrate-up
task test
task generate-swagger
```

## ⚖️ License

BlocksTracker is free and open-source software, licensed under the [GPLv3](https://www.gnu.org/licenses/gpl-3.0.html).
