# Discipleship Journal Web

The frontend for the Discipleship Journal, built with React, TypeScript, Vite, and Tailwind CSS.

## Features
- **Journaling**: Create and manage notes with Markdown support.
- **AI Integration**: Chat with Bible-aware AI.
- **Social**: Connect with friends, join groups, and share notes.
- **PWA**: Installable on mobile and desktop.

## Development

### Prerequisites
- Node.js 18+
- npm

### Setup
```bash
npm install
```

### Run
```bash
npm run dev
```
The app will be available at `http://localhost:5173`.

### Test
```bash
npm run test       # Unit tests (Vitest)
npm run test:e2e   # E2E tests (Playwright)
```

### Build
```bash
npm run build
```

## Architecture
- **Components**: Shadcn/UI (Radix UI)
- **State/API**: React Hooks, Context
- **Routing**: React Router
- **Auth**: Firebase Auth
