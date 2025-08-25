# Backend: Postgres in Docker

This backend currently uses an in-memory store. To run Postgres locally and seed the schema, use Docker Compose in this folder.

## Quick start

1. Copy the example env file and adjust if needed:

   cp .env.example .env

2. Start Postgres and Adminer:

   docker compose up -d

3. Verify DB is healthy and seeded:

   - Adminer: http://localhost:8081
   - System: PostgreSQL
   - Server: db  (use the Docker service name here; not localhost)
   - User: learnlang (or POSTGRES_USER)
   - Password: learnlang (or POSTGRES_PASSWORD)
   - Database: learnlang (or POSTGRES_DB)

   Notes:
   - In Adminer (inside Docker), "localhost" refers to the Adminer container itself. Use "db" to reach Postgres.
   - From your host (psql, app), use host "localhost" and port 5432 (or POSTGRES_PORT if changed).

The container runs init SQL from `db/migrations/` on first start, creating tables and seeding the `languages` table (Hindi, German).

## Connection string

Use `DATABASE_URL` (defined in `.env`) when you wire up the backend to Postgres, e.g.

postgres://learnlang:learnlang@localhost:5432/learnlang?sslmode=disable

## Notes

- Data persists in a named Docker volume (`pgdata`). To reset:

  docker compose down -v

- Migrations are simple init SQL for now. You can adopt a tool like `goose` or `migrate` later.

## Next steps to switch code from memory to Postgres

- Introduce a store interface with methods used by handlers (packs, vocabs, languages).
- Keep the current in-memory implementation for tests; add a Postgres implementation using `database/sql` + `pgx`.
- Select the implementation via env (e.g., `STORE=postgres` with `DATABASE_URL`).
- Ensure uniqueness via DB constraints already in the schema.

## Public directory for static files
- because backend startup will fail without upload folder start go wtih
UPLOAD_DIR=/Users/tathagat/dev/projects/learnlang.app/backend/public/uploads go run .

## Test database

- On first boot, Docker will also create a `learnlang_test` database and apply the same schema/seeds.
- The backend automatically switches to `learnlang_test` when running `go test`, or when `TEST_MODE=1` is set.
- Override the test DB name via `POSTGRES_TEST_DB` if needed.

Examples:

```
# run tests explicitly against the test DB
TEST_MODE=1 go test ./...

# override the test DB name
POSTGRES_TEST_DB=my_app_test TEST_MODE=1 go test ./...
```