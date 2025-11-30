# Task: Setup Backend Deployment

## Task Information
- **Task ID**: INFRA-004
- **Status**: pending
- **Priority**: high
- **Phase**: 1.5
- **Estimated Effort**: 0.5 days
- **Dependencies**: FOUNDATION-001

## Task Details

### Description
Create a script to deploy the Backend API to Google Cloud Run. This ensures consistent deployment parameters and enables CI/CD integration.

### Acceptance Criteria
- [ ] **Docker**:
    - [ ] Verify `api/Dockerfile` is optimized (multi-stage build).
    - [ ] Create `.dockerignore` to exclude unnecessary files.
- [ ] **Scripts**:
    - [ ] Create `scripts/deploy-api.sh` (or add to Makefile).
    - [ ] Script should build the Docker image.
    - [ ] Script should push the image to Google Container Registry (GCR) or Artifact Registry.
    - [ ] Script should deploy to Cloud Run using `gcloud run deploy`.
- [ ] **Configuration**:
    - [ ] Define necessary environment variables (DB_URL, etc.) or Secret Manager references in the deploy command.
- [ ] **Verification**:
    - [ ] Verify manual deployment works to a staging/dev environment.

### Implementation Notes
- Use `gcloud builds submit` if simpler for building/pushing.
- Ensure the service account has necessary permissions (Cloud SQL Client).
- Region should be configurable (default to `us-central1` or similar).

---

*Created: 2025-05-18*
*Status: pending*
