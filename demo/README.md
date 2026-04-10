# Go Demo App

Scaffolded with Fiber, MySQL, sqlc, and Goose.

## Setup

```bash
# Start MySQL
make docker-up

# Run migrations
make migrate-up

# Generate sqlc models
make sqlc-generate

# Run app
make run
```
