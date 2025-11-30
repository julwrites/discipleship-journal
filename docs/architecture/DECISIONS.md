# Architectural Decisions

## Database: PostgreSQL for All
**Decision**: Use PostgreSQL for both relational user data and document-based journal notes.
**Rationale**:
- **Simplicity**: Single database instance to manage and pay for.
- **Flexibility**: PostgreSQL's `JSONB` support allows for schema-less storage of journal notes, similar to MongoDB, while retaining relational integrity for users.
- **Portability**: Standard SQL database, easily portable to Azure, AWS, or self-hosted environments.

## Frontend: React + Vite (CSR)
**Decision**: Build the PWA using React with Vite (Client-Side Rendering) instead of Next.js (SSR) or Vue.
**Rationale**:
- **PWA Suitability**: CSR is ideal for PWAs that need to work offline and feel like native apps. The app loads once and communicates via API.
- **Ecosystem**: React has the largest ecosystem and shares concepts with React Native, facilitating a future native mobile app.
- **Deployment**: Generates static files, easily hosted on Firebase Hosting or any CDN.

## Backend: Go
**Decision**: Use Go (Golang) for the backend service.
**Rationale**:
- **Performance**: High performance and low memory footprint (good for Cloud Run cold starts).
- **Simplicity**: Strong standard library, easy to maintain.
- **Portability**: Compiles to a single binary.

## UI Framework: Tailwind CSS + Shadcn/UI
**Decision**: Use Tailwind CSS with Shadcn/UI components.
**Rationale**:
- **Customizability**: Unlike Material UI, Shadcn/UI provides copy-paste components that are fully customizable.
- **Modern Aesthetic**: Clean, modern look suitable for a consumer-facing app.
- **Performance**: Tailwind generates minimal CSS.

## Authentication: Firebase Auth
**Decision**: Use Firebase Auth for identity management.
**Rationale**:
- **Security**: Handles complex flows (Google Sign-in, token refresh) securely.
- **Integration**: Easy integration with Frontend (Client SDK) and Backend (Admin SDK).
- **Cost**: Generous free tier.

## Testing Strategy
**Decision**: Implement a multi-layered testing strategy.
- **Frontend**: Vitest + React Testing Library for unit/component tests. Playwright for End-to-End (E2E) testing.
- **Backend**: Standard Go `testing` package for unit tests. `testcontainers-go` for integration tests with real database instances.
**Rationale**:
- **Confidence**: Comprehensive coverage ensures stability as the application grows.
- **Speed**: Vitest is significantly faster than Jest.
- **Realism**: Playwright and Testcontainers simulate real-world usage and infrastructure.

## CI/CD: GitHub Actions
**Decision**: Use GitHub Actions for Continuous Integration and Continuous Deployment.
**Rationale**:
- **Integration**: Native integration with the repository.
- **Flexibility**: Large marketplace of actions for Go, Node.js, and Google Cloud.
- **Cost**: Free tier for public/standard repositories is sufficient for start.

## Database Migrations: golang-migrate
**Decision**: Use `golang-migrate/migrate` for versioned database schema changes.
**Rationale**:
- **Control**: explicit SQL migration files (up/down).
- **Compatibility**: Widely used in the Go ecosystem.
- **Automation**: Can be run via CLI in CI/CD or embedded in the application binary.

## Observability: Structured Logging & Metrics
**Decision**: Use `log/slog` (Go 1.21+) for structured logging and prepare for OpenTelemetry.
**Rationale**:
- **Debuggability**: Structured logs are machine-readable and easier to query in tools like Google Cloud Logging.
- **Standard**: `slog` is the new standard library solution, reducing external dependencies.

## API Documentation: OpenAPI (Swagger)
**Decision**: Generate OpenAPI v3 specifications from code comments using `swaggo/swag`.
**Rationale**:
- **Documentation**: Keeps documentation in sync with code.
- **Client Generation**: Allows generating TypeScript clients for the frontend automatically.
- **Testing**: Enables API contract testing.
