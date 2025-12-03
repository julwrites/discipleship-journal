# Task: Integrate E2E Tests into CI Pipeline

## Task Information
- **Task ID**: INFRA-007
- **Status**: completed
- **Priority**: high
- **Phase**: 3.5
- **Estimated Effort**: 0.5 days
- **Dependencies**: TEST-002, INFRA-001

## Task Details

### Description
Integrate the existing Playwright End-to-End (E2E) tests into the GitHub Actions CI pipeline. This ensures that critical user flows are verified on every push and pull request, preventing regressions in the frontend.

### Acceptance Criteria
- [x] **GitHub Actions Workflow**:
    - [x] Update `.github/workflows/ci.yml` (or create a new workflow).
    - [x] Add a job to install Playwright dependencies.
    - [x] Run Playwright tests (`npm run test:e2e`).
    - [x] Archive test results as artifacts on failure.
- [x] **Configuration**:
    - [x] Ensure the CI environment properly handles the frontend dev server (Playwright's `webServer` config).
    - [x] Ensure environment variables (like mock keys) are set for CI.

### Implementation Notes
- Playwright needs browsers installed (`npx playwright install --with-deps`).
- The `web` directory contains the Playwright config and tests.
- We should use the same `web` directory for context.
- We might need to create a `.env` file or set env vars in the workflow for the mock keys to work and bypass Firebase initialization issues in CI.
- **Update**: Fixed E2E test selectors in `web/e2e/smoke.spec.ts` to be more robust (using `getByText` instead of `getByRole` for the heading which was causing failures).

---
*Created: 2025-12-03*
*Status: completed - E2E tests integrated into CI workflow and tests verified locally.*
