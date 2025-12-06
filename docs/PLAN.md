# Discipleship Journal PWA - Project Plan

## Overview
A Personal Web App (PWA) for spiritual journaling, integrated with Bible AI, to be hosted on `app.navteens.org`. The app focuses on three core interactions: Journaling, Search/Dashboard, and AI Chat.

## Architecture

### Tech Stack
- **Frontend**: React, TypeScript, Vite, Tailwind CSS, Shadcn/UI.
- **Backend**: Go (Chi router, pgx driver).
- **Database**: PostgreSQL (Cloud SQL).
  - Uses Relational schema for Users.
  - Uses JSONB for Journal Notes (NoSQL-like flexibility).
- **Authentication**: Firebase Auth (Google Sign-in).
  - Frontend: Firebase Client SDK.
  - Backend: Firebase Admin SDK for token verification.
- **AI Integration**: [BibleAIAPI](https://github.com/julwrites/BibleAIAPI) (Gemini).
- **Hosting**:
  - Frontend: Firebase Hosting.
  - Backend: Google Cloud Run.
  - Database: Google Cloud SQL.

## Roadmap

### Phase 1: Foundation (MVP)
- **Repository Setup**: Monorepo structure (`web/`, `api/`).
- **Authentication**: Google Sign-in via Firebase.
- **Backend Core**: Middleware, Database connection, User/Note CRUD.
- **Frontend Core**: Dashboard, basic routing.

### Phase 1.5: Engineering Excellence (Pre-Core)
- **Testing Infrastructure**: Setup Vitest (Frontend), Go Test (Backend), and Playwright (E2E).
- **CI/CD**: GitHub Actions for automated testing and linting.
- **Observability**: Structured logging and basic error tracking.
- **Database Migrations**: Setup versioned schema migrations (`golang-migrate`).
- **API Documentation**: Setup Swagger/OpenAPI generation (`swaggo`).
- **Deployment Scripts**: Create scripts for deploying Frontend (Firebase) and Backend (Cloud Run).

### Phase 2: Core Features
1.  **User Settings**:
    - Configure Preferred Bible Version.
    - Configure Globally Unique Username.
2.  **Dashboard**:
    - Display Journal Notes (MRU order).
    - Search Bar (bottom).
    - Floating Action Button (+) for new notes.
    - Chat with AI entry point.
3.  **Journaling**:
    - Context: Select Bible passages.
    - Editor: Markdown support (Images, GIFs, Emojis).
    - AI Assistance: Ask questions with context (passages + note content).
4.  **AI Chat Workflow**:
    - Define Bible Passages.
    - Define Themes/Keywords (optional).
    - System Prompt (User input).
    - **Result**: Collated into a new Journal Note.

### Phase 3: Polish & Deployment
- **PWA Features**: Manifest, Service Workers, Offline capabilities.
- **Deployment**: CI/CD pipeline to GCP.
- **Domain**: `app.navteens.org`.

## Future Work (Post-MVP)
- **Social Features**: Groups (Phase 4).
- **Native Apps**: React Native (sharing logic with PWA).

## Infrastructure & Costs
- **Estimated Cost (1000 users)**: ~$10-15 USD/month.
  - Cloud SQL is the primary cost driver.
  - Cloud Run and Firebase Hosting likely fall within free tiers.
