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
