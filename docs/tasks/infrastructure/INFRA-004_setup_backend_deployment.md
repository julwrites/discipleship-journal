# Task: Setup Backend Deployment

## Task Information
- **Task ID**: INFRA-004
- **Status**: completed
- **Priority**: high
- **Phase**: 1.5
- **Estimated Effort**: 0.5 days
- **Dependencies**: FOUNDATION-001

## Task Details

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

---

*Created: 2025-05-18*
*Status: completed*
