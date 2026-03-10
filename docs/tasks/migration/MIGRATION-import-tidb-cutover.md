---
id: MIGRATION-import-tidb-cutover
status: todo
title: Staging Data Import and Final Cutover
priority: high
created: 2026-03-03
category: migration
dependencies: MIGRATION-extract-postgres-data
type: task
---

# Staging Data Import and Final Cutover

## Objective
Once we have verifiable MySQL-compatible CSVs extracted from Postgres, we need to inject them into the TiDB cluster, functionally test the API against it, and perform the final production application cutover.

## Requirements
*   **Staging Import:** Execute an import of the `.csv` files into the Staging TiDB cluster (using `TiDB Lightning` or standard scripts).
*   **Verification:** Deploy the newly refactored Go backend (using the `mysql` driver) to the Staging Cloud Run instance and confirm that API endpoints (like fetching/saving journal entries and memory verses) work perfectly.
*   **Cutover Runbook:**
    *   1. Place the current API into maintenance mode to freeze PostgreSQL writes.
    *   2. Run the final Postgres-to-CSV extraction script.
    *   3. Run the TiDB Lightning import script to push the final data to the Production TiDB cluster.
    *   4. Trigger the GitHub Actions workflow to deploy the refactored `main` branch. This will deploy the Go API using the TiDB configurations located in Google Secret Manager (`PROD_DJ_DB_HOST`, etc.).
    *   5. Verify Production functionality and take off maintenance mode.

## Next Steps
*   [ ] Populate Staging TiDB database with the converted CSV extracts.
*   [ ] Test Staging frontend/backend stack fully against the TiDB schemas.
*   [ ] Coordinate the downtime window and execute the Cutover Runbook.
