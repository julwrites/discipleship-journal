---
id: INFRA-006
status: completed
title: Setup Continuous Deployment (CD) Pipeline
priority: high
created: 2025-12-11 06:09:10
category: unknown
type: task
---

# Setup Continuous Deployment (CD) Pipeline

### Description
Implement a GitHub Actions workflow to automatically deploy the Backend to Cloud Run and the Frontend to Firebase Hosting when changes are pushed to the `main` branch.

### Acceptance Criteria
- [x] **Deployment Workflow**:
    - [x] Create `.github/workflows/deploy.yml`.
    - [x] Trigger on `push` to `main` (or specific tags/manual dispatch).
- [x] **Backend Job**:
    - [x] Authenticate with GCP.
    - [x] Build and Push Docker image.
    - [x] Deploy to Cloud Run.
- [x] **Frontend Job**:
    - [x] Build React app (`npm run build`).
    - [x] Deploy to Firebase Hosting.
- [x] **Security**:
    - [x] Use GitHub Secrets for credentials (`GCP_SA_KEY`, etc.).

### Implementation Notes
- Use `google-github-actions/auth` for GCP authentication.
- Use `google-github-actions/deploy-cloudrun` for Cloud Run.
- Use `FirebaseExtended/action-hosting-deploy` or `firebase-tools` for Frontend.
- This pipeline connects the scripts created in INFRA-003/004 with the CI setup in INFRA-001.
