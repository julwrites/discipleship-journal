# System Architecture

## Overview
The Discipleship Journal is a web application designed to help users document their spiritual journey. It consists of a React-based single-page application (SPA) frontend and a Go-based REST API backend.

## Tech Stack

### Frontend (`web/`)
*   **Framework**: React 19 (Vite)
*   **Language**: TypeScript
*   **Styling**: Tailwind CSS, Shadcn/UI (Radix Primitives)
*   **State Management**: React Hooks
*   **Routing**: React Router DOM v7
*   **Authentication**: Firebase Auth (Client SDK)
*   **API Interaction**: Standard `fetch` or custom services.

### Backend (`api/`)
*   **Language**: Go 1.24+
*   **Web Framework**: Chi (Router)
*   **Database**: PostgreSQL (pgx driver)
*   **Authentication**: Firebase Admin SDK (Middleware verification)
*   **AI Integration**: Google Generative AI (likely Gemini, inferred from dependencies)
*   **Configuration**: Environment variables (`.env`)

## Data Flow
1.  **User Interaction**: Users interact with the React frontend.
2.  **Authentication**: Users sign in via Firebase Auth. The frontend obtains an ID token.
3.  **API Requests**: The frontend sends HTTP requests to the Backend API, including the ID Token in the Authorization header.
4.  **Middleware**: Backend middleware verifies the Firebase ID Token.
5.  **Handlers**: Verified requests are routed to handlers (`api/handlers`).
6.  **Database**: Handlers interact with PostgreSQL via the `database` package.
7.  **AI Services**: Chat and Ask features interact with external AI APIs.

## Directory Structure
*   `api/`: Backend source code.
    *   `handlers/`: HTTP request handlers.
    *   `middleware/`: HTTP middleware (Auth, Logging).
    *   `database/`: Database connection and queries.
*   `web/`: Frontend source code.
    *   `src/components/`: Reusable UI components.
    *   `src/pages/`: Application pages.
    *   `src/services/`: API client services.
