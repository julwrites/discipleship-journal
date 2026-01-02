# Discipleship Journal

A P2P Discipleship Journal, designed for privacy and focused mentoring.

## Features
- **Journaling**: Private notes with Bible context.
- **AI Assistant**: Bible-aware AI chat and Q&A.
- **Connections**: Connect with friends and mentors.
- **Social Groups**: Create and join groups for shared discipleship.
- **Sharing**: Share notes with groups for accountability.
- **PWA**: Installable Progressive Web App.

## Project Structure
- `api/`: Go Backend (Chi, Pgx, Postgres)
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
- **Frontend**: http://localhost:3000 (Nginx proxy)
- **Backend API**: http://localhost:8080
- **Postgres DB**: localhost:5432
- **Migrations**: Auto-run on startup

For detailed setup instructions, see:
- [Docker Setup](docs/setup/docker.md)
- [Firebase Setup](docs/setup/firebase.md)
- [Secrets Setup](docs/setup/secrets.md)

### Manual Setup
See `web/README.md` for frontend-specific instructions.

## Documentation
- **Tasks & Status**: `docs/tasks/`
- **Features**: `docs/features/`
- **Architecture**: `docs/architecture/`
- **Testing**: `docs/testing/`
- **Security**: `docs/security/practices.md`

## Security
This project includes pre-commit hooks for secret scanning to prevent accidental commits of sensitive information. To set up:

```bash
./scripts/setup-pre-commit.sh
```

**Important**: Never commit API keys, passwords, or tokens. Always use environment variables.

## Secrets Configuration
See [Secrets Setup](docs/setup/secrets.md) for detailed instructions on setting up database and API secrets for both local development and production deployment.
