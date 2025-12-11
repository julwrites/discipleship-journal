---
id: INFRA-004
status: completed
title: Setup Backend Deployment
priority: high
created: 2025-12-11 06:09:10
category: unknown
type: task
---

# Setup Backend Deployment

### Description
Create a script to deploy the Backend API to Google Cloud Run. This ensures consistent deployment parameters and enables CI/CD integration.

### Acceptance Criteria
- [x] **Docker**:
    - [x] Verify `api/Dockerfile` is optimized.
    - [x] Create `.dockerignore` to exclude unnecessary files.
- [x] **Scripts**:
    - [x] Create `scripts/deploy-api.sh`.
    - [x] Script builds the Docker image.
    - [x] Script pushes the image to GCR.
    - [x] Script deploys to Cloud Run using `gcloud run deploy`.
- [x] **Configuration**:
    - [x] Defined default region and project ID.
- [x] **Verification**:
    - [x] Script syntax verified.

### Implementation Notes
- Using `gcloud builds submit` for simplified build/push.
- Configured for `us-central1`.
