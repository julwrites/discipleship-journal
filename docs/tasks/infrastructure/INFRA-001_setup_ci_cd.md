# Task: Setup CI/CD Pipeline

## Task Information
- **Task ID**: INFRA-001
- **Status**: pending
- **Priority**: high
- **Phase**: 1.5
- **Estimated Effort**: 1 day
- **Dependencies**: TEST-001

## Task Details

### Description
Implement Continuous Integration (CI) to automatically run linting and testing on every Pull Request. Prepare Continuous Deployment (CD) workflows for future deployment to staging/production.

### Acceptance Criteria
- [ ] **GitHub Actions Workflow**:
    - [ ] Create `.github/workflows/ci.yml`.
    - [ ] Trigger on `push` to `main` and `pull_request`.
- [ ] **Backend Job**:
    - [ ] Install Go.
    - [ ] Run `go mod tidy` check (ensure clean).
    - [ ] Run `golangci-lint` (or `staticcheck`).
    - [ ] Run `go test ./...`.
- [ ] **Frontend Job**:
    - [ ] Install Node.js.
    - [ ] Run `npm install` (cache node_modules).
    - [ ] Run `npm run lint` (eslint).
    - [ ] Run `npm run test` (vitest).
    - [ ] Run `npm run build` (check for build errors).
- [ ] **E2E Job** (Optional for now, or separate workflow):
    - [ ] Install Playwright browsers.
    - [ ] Run Playwright tests.

### Implementation Notes
- Use `actions/setup-go` and `actions/setup-node`.
- Use `golangci/golangci-lint-action` for Go linting.
- Ensure caching is configured to speed up runs.

---

*Created: 2025-05-18*
*Status: pending*
