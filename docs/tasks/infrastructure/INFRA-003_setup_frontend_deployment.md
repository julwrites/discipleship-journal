---
id: INFRA-003
status: completed
title: Setup Frontend Deployment
priority: high
created: 2025-12-11 06:09:10
category: unknown
type: task
---

# Setup Frontend Deployment

### Description
Configure Firebase Hosting and create a script to deploy the frontend application. This allows for manual deployments and integration into CI/CD pipelines.

### Acceptance Criteria
- [x] **Configuration**:
    - [x] Create `firebase.json` in the project root (configured for `web/dist`).
    - [x] Create `.firebaserc` with default project alias.
- [x] **Scripts**:
    - [x] Create `scripts/deploy-web.sh`.
    - [x] Script runs `npm run build` in `web/` directory.
    - [x] Script runs `firebase deploy --only hosting`.
- [x] **Verification**:
    - [x] Configuration verified against Firebase docs.

### Implementation Notes
- Used `web/dist` as the public directory.
- Configured rewrites to `index.html` for SPA routing.
