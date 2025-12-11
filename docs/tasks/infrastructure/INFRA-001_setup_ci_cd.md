---
id: INFRA-001
status: Completed
title: Setup CI/CD Pipeline
priority: high
created: 2025-12-11 06:09:10
category: unknown
type: task
---

# Setup CI/CD Pipeline

### Description
Implement Continuous Integration (CI) to automatically run linting and testing on every Pull Request. Prepare Continuous Deployment (CD) workflows for future deployment to staging/production.

### Acceptance Criteria
- [x] **GitHub Actions Workflow**:
    - [x] Create `.github/workflows/ci.yml`.
    - [x] Trigger on `push` to `main` and `pull_request`.
- [x] **Backend Job**:
    - [x] Install Go.
    - [x] Run `go mod tidy` check (ensure clean).
    - [x] Run `golangci-lint` (or `staticcheck`).
    - [x] Run `go test ./...`.
- [x] **Frontend Job**:
    - [x] Install Node.js.
    - [x] Run `npm install` (cache node_modules).
    - [x] Run `npm run lint` (eslint).
    - [x] Run `npm run test` (vitest).
    - [x] Run `npm run build` (check for build errors).
- [ ] **E2E Job** (Optional for now, or separate workflow):
    - [ ] Install Playwright browsers.
    - [ ] Run Playwright tests.

### Implementation Notes
- Use `actions/setup-go` and `actions/setup-node`.
- Use `golangci/golangci-lint-action` for Go linting.
- Ensure caching is configured to speed up runs.
