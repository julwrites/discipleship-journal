# AI Agent Guidelines

This repository is maintained with the assistance of an AI Agent. The following guidelines must be followed:

## Code Style

### Frontend (React/Vite)
*   Use TypeScript for all new code.
*   Use Functional Components and Hooks.
*   Use Tailwind CSS for styling.
*   Use Shadcn/UI for UI components.
*   Follow standard directory structure: `components`, `pages`, `hooks`, `services`, `types`.

### Backend (Go)
*   Follow standard Go idioms (Effective Go).
*   Use `chi` or `gin` for routing.
*   Use `pgx` for PostgreSQL interactions.
*   Keep handler logic separate from business logic (services).
*   Use environment variables for configuration.

## Commit Messages
*   Use descriptive commit messages.
*   Format: `type(scope): description`.

## Testing
*   Write unit tests for utility functions and core business logic.
*   Ensure the code builds and runs before submitting.
