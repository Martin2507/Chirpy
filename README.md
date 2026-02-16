
```md
# Chirpy

Chirpy is a small service that exposes a REST API for creating users, posting short messages (“chirps”), and managing authentication.

## Why it exists

This project was built as part of the Boot.dev backend curriculum to practice:
- Building HTTP APIs
- Working with JSON, routing, and status codes
- Authentication (tokens) and basic user management
- Persisting data (depending on your setup)

## Tech stack

- Language: Go
- API style: REST + JSON
- Storage: (fill in: in-memory / file / SQL)
- Auth: (fill in: JWT / opaque tokens)

## Features

- Create and manage users
- Login / token-based auth
- Create, list, and delete chirps
- Basic validation and error handling

## Getting started

### Prerequisites
- Go installed (1.20+ recommended)
- (Optional) Database installed/configured if your version uses one

### Install
```bash
git clone https://github.com/<your-username>/chirpy.git
cd chirpy

```

### Configure

Set any required environment variables (adjust as needed):

```bash
export PORT=8080
export DB_URL="<your-db-connection-string>"
export JWT_SECRET="<your-secret>"

```

### Run

```bash
go run .

```

The server should start on:

-   `http://localhost:8080`

## API overview

Base URL:  `http://localhost:8080`

### Health check

-   `GET /api/healthz`

### Users

-   `POST /api/users`  — create a user

### Auth

-   `POST /api/login`  — authenticate and receive a token

### Chirps

-   `GET /api/chirps`  — list chirps
-   `POST /api/chirps`  — create a chirp (auth required)
-   `DELETE /api/chirps/{id}`  — delete a chirp (auth required)

> Note: Endpoint paths may differ slightly depending on your implementation—edit this section to match your routes.

## Example request

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"body":"hello from chirpy"}'

```

## Project status

Built for learning and iteration. Not intended for production use without additional hardening (rate limiting, observability, etc.).

## License

MIT (or choose another license)