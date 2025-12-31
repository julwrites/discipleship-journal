---
id: TESTING-20251231-024800-JULES
status: todo
title: Comprehensive Testing Plan & Fixes
priority: high
created: 2025-12-31 02:48:00
category: testing
dependencies:
  - TESTING-20251231-014339-VIF
type: task
---

# Comprehensive Testing Plan & Fixes

## Context
Current backend tests are failing. Specifically `TestGroupHandler_CreateGroup` is failing due to type mismatches (`uuid.UUID` vs `string`) in `pgxmock` expectations. This indicates a drift between the implementation (which uses `google/uuid` and passes `uuid.UUID` to drivers) and the tests (which likely expect strings or haven't been updated to match the driver's strictness).

## Objective
To ensure 100% functionality verification, we must first **stabilize the existing test suite** and then implement the **Contract Testing strategy**.

## Phase 0: Fix Existing Backend Tests (Immediate Priority)

### 0.1 Fix `TestGroupHandler_CreateGroup`
*   **Issue**: `argument 2 expected [uuid.UUID] does not match actual [string]`.
*   **Cause**: The handler is passing a `string` (or `uuid.UUID` converted to string?) but the mock expects the other type, or vice versa. The error log `expected [uuid.UUID ... ] does not match actual [string ...]` suggests the *mock* was set up with `AnyArg()` or a specific UUID value, but the code passed a string.
*   **Action**: Update `api/handlers/group_test.go` to ensure `ExpectQuery` arguments match the actual types passed by `GroupHandler`. If `GroupHandler` uses `uuid.New()`, the test should probably use `pgxmock.AnyArg()` because the UUID is non-deterministic, OR the handler should accept an ID (dependency injection of generator) - but for now, `AnyArg()` is safer for generated IDs.

### 0.2 Audit Other Handler Tests
*   Run all tests in `api/handlers`.
*   Fix similar mismatches in `NoteHandler`, `ConnectionHandler`, etc.

## Phase 1: Infrastructure & Contract Definition

### 1.1 Backend Integration Infrastructure
*   **Goal**: Enable tests to run against a real PostgreSQL instance to avoid mock drift.
*   **Action**: Add `testcontainers-go` and necessary drivers to `api/go.mod`.
*   **Action**: Implement `api/tests/integration/setup.go` to manage container lifecycle and DB migrations.

### 1.2 Define The API Contract
Based on `web/src/services/api.ts` and `web/src/pages/NoteEditor.tsx`, the Backend MUST support:

| Endpoint | Method | Expected Payload (JSON) | Expected Response (JSON Structure) |
| :--- | :--- | :--- | :--- |
| `/api/notes` | POST | `{"title": string, "content": {"markdown": string}}` | `{"id": UUID}` |
| `/api/notes` | GET | Query: `q`, `startDate`, `endDate`, `sortBy`, `sortOrder` | `{ "data": [...], "meta": ... }` |
| `/api/groups` | POST | `{"name": string, "description": string}` | `{"id": UUID}` |

*   **Action**: Create `api/tests/contract/contract_test.go` that strictly enforces these payloads.

## Phase 2: Backend Verification (The "Real" Test)

### 2.1 Service Integration Tests
*   **Scope**:
    *   `NoteService`: CRUD, Search filters, Soft Delete.
    *   `GroupService`: Membership logic (constraints).

### 2.2 API Contract Tests
*   **Scope**: Create a test suite that:
    1.  Spins up `chi` router (with real DB).
    2.  Sends the "Expected Payload" defined in 1.2.
    3.  Asserts 200 OK and valid JSON response.
    4.  **Constraint**: Tests must pass with the *current* backend implementation (which returns 200 OK for creates).

## Phase 3: Frontend Verification

### 3.1 E2E Mock Audit
*   **Goal**: Ensure Frontend tests (`web/e2e/`) are not testing against "fantasy" APIs.
*   **Action**: Review `web/e2e/journal.spec.ts`.
    *   **Fix**: Update `GET /api/notes` mock to return `{ data: [], meta: ... }` structure.
    *   **Fix**: Update `POST /api/notes` mock to return `{ id: ... }`.

## Execution Order
1.  **Fix Existing Tests**: Modify `api/handlers/group_test.go` and others to pass.
2.  **Infra**: Add `testcontainers`.
3.  **Contract**: Write `api/tests/contract` suite.
4.  **Audit**: Fix Frontend E2E mocks.
