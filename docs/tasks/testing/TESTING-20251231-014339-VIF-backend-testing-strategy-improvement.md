---
id: TESTING-20251231-014339-VIF
status: in_progress
title: Backend Testing Strategy Improvement
priority: high
created: 2025-12-31 01:43:39
category: testing
dependencies:
type: task
---

# Backend Testing Strategy Improvement

## Context
The current backend testing suite primarily consists of unit tests with mocks (`pgxmock`, `testify`). While these provide good logic coverage, they fail to catch integration issues and discrepancies between the frontend's expected API contract and the backend's actual implementation. Users have reported that "many backend API's fail when called from the frontend".

## Objective
To achieve "100% certainty" in backend releases, we will implement a multi-layered testing strategy:
1.  **Integration Tests**: Verify Service <-> DB interactions using real (containerized) PostgreSQL.
2.  **Contract/API Tests**: Verify Router <-> Handler <-> Service <-> DB flows using payloads that strictly match the frontend's `api.ts`.
3.  **Enhanced Unit Tests**: Fill coverage gaps in existing handler tests.

## Plan

### 1. Enable Integration Testing Infrastructure
*   [x] Add `testcontainers-go` and `pgx/v5/stdlib` (if needed) to `api/go.mod`.
*   [x] Create a `api/tests/integration` package.
*   [x] Implement a `SetupIntegrationDB` helper that spins up a Postgres container, applies migrations, and returns a connection pool.
*   [x] **Constraint**: Ensure these tests are skipped if Docker is not available (using `testing.Short()` or build tags).

### 2. Implement Service Integration Tests
*   [x] Port key `NoteService` tests to run against the real DB.
*   [x] Port `GroupService` (logic in handlers currently, might need refactoring or direct handler integration testing) tests. (Covered by Contract Tests in `groups_test.go`)
*   [x] Verify complex queries (e.g., full-text search, date filtering) which are hard to mock accurately with `pgxmock`.

### 3. Implement "Frontend Contract" API Tests
*   [x] Create a new test suite `api/tests/contract` (or within `integration`).
*   [x] These tests will use `httptest` to spin up the `chi` router.
*   [x] **Crucial**: The test payloads (JSON bodies) must be copied *verbatim* or derived directly from `web/src/services/api.ts` logic.
*   [x] **Crucial**: The assertions must verify that the response JSON structure matches exactly what `web/src/services/api.ts` expects (e.g., field names, date formats, nullability).
*   [x] Cover the "Happy Path" for all major entities: Notes (Done), Groups (Done), Connections (Done), Reading Plans (Done).

### 4. Review and Refine Unit Tests
*   [ ] Audit existing `api/handlers/*_test.go`.
*   [ ] Ensure `pgxmock` expectations match the actual SQL used in production (drift is a common cause of failure).
*   [ ] Add tests for edge cases: invalid JSON, missing required fields, permission denied.

## Implementation Details

### Integration Test Helper (Draft)
```go
package integration

import (
    "context"
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/jackc/pgx/v5/pgxpool"
)

func SetupDB(t *testing.T) *pgxpool.Pool {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    // Spin up container
    // Run migrations
    // Return pool
}
```

### Contract Test Example (Draft)
```go
func TestCreateNote_Contract(t *testing.T) {
    // 1. Setup Real DB
    db := SetupDB(t)
    // 2. Setup Router
    r := api.SetupRouter(db)

    // 3. Prepare Payload matching api.ts: createNote(title, content)
    // api.ts sends: { "title": "...", "content": {...} }
    payload := `{"title": "Test", "content": {"text": "Hello"}}`

    // 4. Execute Request
    w := httptest.NewRecorder()
    req := httptest.NewRequest("POST", "/api/notes", strings.NewReader(payload))
    // ... add auth headers ...
    r.ServeHTTP(w, req)

    // 5. Verify Response matches api.ts expectation (res.json())
    // api.ts returns: { "id": "..." } (based on handler)
    var resp map[string]string
    json.Unmarshal(w.Body.Bytes(), &resp)
    if resp["id"] == "" {
        t.Fatal("Expected ID in response")
    }
}
```

## Outcome
By running these tests in CI (where Docker is available), we will catch:
*   SQL syntax errors specific to Postgres (which mocks might miss).
*   Data type mismatches.
*   JSON marshalling/unmarshalling issues.
*   Logic errors in complex queries.

## Next Steps
1.  Audit existing `api/handlers/*_test.go` and refine unit tests.
2.  Add tests for edge cases (invalid JSON, missing fields, permission denied).
