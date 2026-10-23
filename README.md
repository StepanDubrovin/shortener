# shortener

URL shortening service built with Go and PostgreSQL.

## Features

- Shorten long URLs via REST API
- Redirect to original URL by short code
- Persistent PostgreSQL storage

## Requirements

- Go 1.27.2 or newer
- Docker with Docker Compose v2

## Run

### 1. Configure PostgreSQL

```bash
cp .env.example .env
```

Set `POSTGRES_PASSWORD` in `.env` before starting the database for the first time.
The examples below use the default database and user names, both `shortener`.

### 2. Start the database

```bash
docker compose up -d
docker compose ps
```

Check that PostgreSQL is ready to accept connections:

```bash
docker compose exec postgres pg_isready -U shortener -d shortener
```

If it is still starting, repeat the readiness check after a few seconds.

### 3. Apply the migration

Run this once on a new database, after PostgreSQL is ready:

```bash
docker compose exec -T postgres \
  psql -U shortener -d shortener -v ON_ERROR_STOP=1 \
  < migrations/001_create_urls.sql
```

### 4. Start the application

Set `DATABASE_URL` using the credentials from `.env`. Replace `change_me` below
with your PostgreSQL password; percent-encode any special characters in the URL credentials.

```bash
export DATABASE_URL='postgres://shortener:change_me@localhost:5432/shortener?sslmode=disable'
go run ./cmd/shortener
```

Docker Compose reads `.env` automatically; the Go application does not. It reads
`DATABASE_URL` from its process environment. When running from GoLand, set this
variable in the run configuration's environment variables.

The HTTP server listens on `:8080`. PostgreSQL is exposed on `127.0.0.1:5432`.
`sslmode=disable` is used for the local database connection.

Database files are stored in a Docker volume and survive application and container restarts.
To stop the database while keeping its data:

```bash
docker compose down
```

PostgreSQL initialization settings in `.env` apply only when the data volume is first created.
Changing them later does not update existing database users or passwords.

## API

### Shorten URL

```bash
curl -i -X POST http://localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/very-long-path"}'
```

Returns `201 Created` with a JSON response:

```json
{"code": "1", "short": "http://localhost:8080/1"}
```

### Follow short link

Open `http://localhost:8080/1` in a browser or:

```bash
curl -i http://localhost:8080/1
```

Returns `302 Found` with a `Location` header pointing to the original URL.

Use the code returned by the POST request; `1` is an example, not a guaranteed first code.

### Errors

| Status | Meaning |
| --- | --- |
| `400 Bad Request` | Invalid JSON or an empty URL |
| `404 Not Found` | Unknown or invalid short code |
| `500 Internal Server Error` | Storage operation failed |

## Stack

- Go (`net/http`, `encoding/json`, `database/sql`)
- PostgreSQL 17
- pgx PostgreSQL driver
- Docker Compose
