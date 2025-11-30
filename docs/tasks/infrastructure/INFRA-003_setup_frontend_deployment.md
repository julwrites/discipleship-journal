# Task: Setup Frontend Deployment

## Task Information
- **Task ID**: INFRA-003
- **Status**: pending
- **Priority**: high
- **Phase**: 1.5
- **Estimated Effort**: 0.5 days
- **Dependencies**: FOUNDATION-001

## Task Details

### Description
Configure Firebase Hosting and create a script to deploy the frontend application. This allows for manual deployments and integration into CI/CD pipelines.

### Acceptance Criteria
- [ ] **Configuration**:
    - [ ] Create `firebase.json` in the project root (configured for `web/dist`).
    - [ ] Create `.firebaserc` with default project alias.
- [ ] **Scripts**:
    - [ ] Create `scripts/deploy-web.sh` (or add to Makefile).
    - [ ] Script should run `npm run build` in `web/` directory.
    - [ ] Script should run `firebase deploy --only hosting`.
- [ ] **Verification**:
    - [ ] Verify manual deployment works to a staging/dev project.

### Implementation Notes
- Ensure `firebase-tools` is required (or installed via the script/CI).
- Use `web/dist` as the public directory in `firebase.json`.
- Configure rewrites to `index.html` for SPA routing.

---

*Created: 2025-05-18*
*Status: pending*
