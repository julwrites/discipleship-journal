# Discipleship Journal

A P2P Discipleship Journal, designed for privacy and focused mentoring.

## Project Structure
- `api/`: Go Backend (Chi, Pgx)
- `web/`: React Frontend (Vite, TypeScript, Tailwind)
- `docs/`: Documentation and Task Management
- `scripts/`: Utility scripts

## Getting Started

### Local Development (Docker)
To run the full stack locally:

```bash
docker-compose up --build
```

This will start:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Postgres DB**: localhost:5432
- **Migrations**: Auto-run on startup

### Manual Setup
Refer to `api/README.md` and `web/README.md` (if available) for individual setup.

## Documentation
See `docs/tasks/` for the task board and status.
See `docs/architecture/` for design docs.
