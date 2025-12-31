---
id: TESTING-20251231-020000-JULES
status: todo
title: Comprehensive Testing Plan for 100% Functionality
priority: high
created: 2025-12-31 02:00:00
category: testing
dependencies:
  - TESTING-20251231-014339-VIF
type: task
---

# Comprehensive Testing Plan

## Objective
To ensure that all **currently implemented functionality** (as defined in `docs/features/README.md`) works 100%, bridging the gap between Frontend logic (`web/src/services/api.ts`) and Backend implementation (`api/handlers` & `api/services`).

## Phase 1: Infrastructure & Contract Definition

### 1.1 Backend Integration Infrastructure
*   **Goal**: Enable tests to run against a real PostgreSQL instance.
*   **Action**: Add `testcontainers-go` and necessary drivers to `api/go.mod`.
*   **Action**: Implement `api/tests/integration/setup.go` to manage container lifecycle and DB migrations.

### 1.2 Define The API Contract
Based on `web/src/services/api.ts` and `web/src/pages/NoteEditor.tsx`, the Backend MUST support:

| Endpoint | Method | Expected Payload (JSON) | Expected Response (JSON Structure) |
| :--- | :--- | :--- | :--- |
| `/api/notes` | POST | `{"title": string, "content": {"markdown": string}}` | `{"id": UUID, ...}` |
| `/api/notes/{id}` | PUT | `{"title": string, "content": {"markdown": string}}` | `{"success": true}` or Updated Note |
| `/api/notes` | GET | Query: `q`, `startDate`, `endDate`, `sortBy`, `sortOrder` | `[{"id": UUID, "title": string, "content": null, ...}]` |
| `/api/connections/request` | POST | `{"receiver_email": string}` | `{"success": true, ...}` |
| `/api/groups/{id}/shares` | POST | `{"note_id": UUID, "comment": string}` | `{"id": UUID, ...}` |
| `/api/chat` | POST | `{"passage": string, "themes": string[], "prompt": string}` | `{"response": string}` |

*   **Risk**: `NoteEditor.tsx` sends `content: { markdown: markdown }`. Backend *must* parse this as `map[string]interface{}` or a struct with `Markdown string`. If Backend expects raw string or different JSON structure, it fails.
*   **Action**: Create `api/tests/contract/contract_test.go` that strictly enforces these payloads.

## Phase 2: Backend Verification (The "Real" Test)

### 2.1 Service Integration Tests
*   **Goal**: Verify business logic against real DB (constraints, dates, search).
*   **Scope**:
    *   `NoteService`: CRUD, Search filters, Soft Delete.
    *   `GroupService`: Membership logic (constraints).
    *   `ReadingPlanService`: Progress tracking (SQL joins).

### 2.2 API Contract Tests
*   **Goal**: Verify that the Router+Handler stack accepts the exact JSON the frontend sends.
*   **Scope**: Create a test suite that:
    1.  Spins up `chi` router (with real DB).
    2.  Sends the "Expected Payload" defined in 1.2.
    3.  Asserts 200 OK and valid JSON response.

## Phase 3: Frontend Verification

### 3.1 E2E Mock Audit
*   **Goal**: Ensure Frontend tests (`web/e2e/`) are not testing against "fantasy" APIs.
*   **Action**: Review `web/e2e/journal.spec.ts` and others.
    *   Check: Does `journal.spec.ts` send `{content: {markdown: ...}}`?
    *   Check: Does the mock response match the Backend's `Note` struct (e.g., `created_at` format)?
*   **Action**: Update mocks if discrepancies found.

### 3.2 Smoke Test Script
*   **Goal**: Quick manual/automated verification of the full stack.
*   **Action**: Create `scripts/smoke_test.sh` that:
    1.  Uses `curl`.
    2.  Authenticates (using a test token or local bypass).
    3.  Creates a note.
    4.  Reads it back.
    5.  Deletes it.

## Execution Order
1.  **Infra**: Add `testcontainers`.
2.  **Contract**: Write `api/tests/contract` suite (this is the most high-value step).
3.  **Fix**: If Contract tests fail, fix the Backend Handlers immediately.
4.  **Audit**: Check Frontend E2E mocks.
