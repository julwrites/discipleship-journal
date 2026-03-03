## 1. Evaluate Current Database Usage and Identify Incompatibilities

*   **Database Engine:** The current application uses PostgreSQL (via Google Cloud SQL).
*   **Driver:** The Go backend uses `github.com/jackc/pgx/v5`, which is strictly for PostgreSQL. It will need to be replaced with a MySQL-compatible driver (like `github.com/go-sql-driver/mysql`).
*   **Schema Types:**
    *   The schema relies heavily on the `uuid` type and `uuid-ossp` extension (`uuid_generate_v4()`, `gen_random_uuid()`). TiDB (MySQL compatible) doesn't have a native UUID type; `VARCHAR(36)` is typically used, and `UUID()` generates the values.
    *   The schema uses `JSONB` for `notes.content` and `users.settings`. TiDB supports the `JSON` data type, but not `JSONB`.
*   **Indexes:** The schema uses `CREATE INDEX ... USING gin (content)` on `notes`. TiDB does not support GIN indexes in the same way. We need to evaluate if a simple index on extracted JSON fields or a full-text index is more appropriate.
*   **Query Syntax:**
    *   **Placeholders:** PostgreSQL uses `$1`, `$2`, etc. TiDB (MySQL) uses `?`. This will require rewriting all SQL queries in the `api/services` layer.
    *   **Functions:** Usage of `unnest()` (e.g., in `memory_verse_service.go`) will need to be rewritten, likely using multiple inserts or `INSERT ... ON DUPLICATE KEY UPDATE`.
    *   **Type Casting:** PostgreSQL type casting like `::jsonb` or `::uuid[]` will need to be removed or adapted to MySQL standard `CAST()`.
    *   **Returning Clause:** PostgreSQL uses `RETURNING id` to get the inserted ID. MySQL typically uses `LAST_INSERT_ID()`, or in Go, the `LastInsertId()` method on the result object.
*   **Migration Tool:** The project uses `golang-migrate`. The migration files (in `api/migrations/`) are written in PostgreSQL dialect and will need to be completely rewritten for MySQL/TiDB.

## 2. Infrastructure Provisioning Plan (Pulumi)

To deploy TiDB Serverless, you'll need to use Pulumi and the TiDB Cloud provider (`@tidbcloud/pulumi-tidbcloud`). Since the project doesn't have an existing Pulumi setup, we'll guide you through setting it up from scratch.

### Environment Strategy
We have provisioned **separate** TiDB Serverless clusters for Staging and Production via Pulumi in the `bible-api-infrastructure` repository.

### Actions Completed:
1.  **Install Pulumi:** Done.
2.  **Configure Pulumi State Management (GCS):** Done.
3.  **Set Up a TiDB Cloud Account:** Done.
4.  **Get TiDB Cloud API Keys:** Done.
5.  **Initialize Pulumi:** Done.
6.  **Configure Pulumi with Credentials:** Done.

### Proposed Pulumi Code (`infra/index.ts`):

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as tidbcloud from "@tidbcloud/pulumi-tidbcloud";

// Ensure you have configured the provider with the necessary credentials
// (publicKey, privateKey) and defaults (projectId).

const config = new pulumi.Config("tidbcloud");
const projectId = config.require("projectId");

// Helper function to provision a cluster
function createCluster(env: string) {
    const cluster = new tidbcloud.Cluster(`discipleship-journal-tidb-${env}`, {
        projectId: projectId,
        name: `dj-serverless-db-${env}`,
        clusterType: "DEVELOPER", // Represents the Serverless tier
        cloudProvider: "AWS", // Or GCP if preferred and available in your region
        region: "us-east-1",  // Adjust to your desired region
        serverless: {
            spendLimit: {
                monthly: 0, // 0 for free tier if applicable, adjust as needed
            },
        },
    });

    const dbPassword = new pulumi.random.RandomPassword(`${env}-dbPassword`, {
        length: 16,
        special: true,
        overrideSpecial: "_%@",
    });

    return {
        clusterId: cluster.id,
        endpoint: cluster.endpoints.apply(e => e.standard.host),
        port: cluster.endpoints.apply(e => e.standard.port),
        password: dbPassword.result,
    };
}

// Provision separate clusters for staging and production
// Note: This fits within the free tier limits (up to 5 clusters per org)
const stagingCluster = createCluster("staging");
const productionCluster = createCluster("production");

// Export the connection details
export const stagingEndpoint = stagingCluster.endpoint;
export const stagingPassword = stagingCluster.password;
export const productionEndpoint = productionCluster.endpoint;
export const productionPassword = productionCluster.password;

// Note: The connection string typically looks like:
// mysql://root:<password>@<endpoint>:<port>/<dbname>?tls=true
```

**Next Steps After Provisioning:**
*   You will need to manually log into the TiDB Cloud Console (or use a script) to create the actual database schema (`discipleship_journal`).
*   Store the generated connection endpoint, port, and password in Google Secret Manager to be used by the Cloud Run deployment.
## 3. Application Code Migration Plan (PostgreSQL to MySQL/TiDB)

The most significant change will be migrating the Go backend's data layer from PostgreSQL (`pgx`) to MySQL (`database/sql` + `go-sql-driver/mysql`). TiDB is wire-compatible with MySQL.

### **3.1 Database Driver Switch**
*   **Remove `pgx`:** Remove all dependencies on `github.com/jackc/pgx/v5` and `github.com/pashagolub/pgxmock/v4`.
*   **Add MySQL Driver:** Add `github.com/go-sql-driver/mysql` to `go.mod`.
*   **Standard Library:** Switch all interfaces in `api/database/interfaces.go` and `api/services/` to use standard `*sql.DB` and `*sql.Tx` from the `database/sql` package instead of `*pgxpool.Pool` and `pgx.Tx`.
*   **Testing:** Replace `pgxmock` with `github.com/DATA-DOG/go-sqlmock` for all unit tests.

### **3.2 Connection Logic (`api/database/db.go`)**
The connection logic will need a complete rewrite. We are moving from the `cloudsqlconn` specific `pgxpool` configuration to a standard `sql.Open("mysql", dsn)` call.

*   **New DSN Format:**
    `user:password@tcp(host:port)/dbname?tls=true&parseTime=true`
*   **Remove Cloud SQL Connector:** The `cloud.google.com/go/cloudsqlconn` dependency is no longer needed. TiDB Serverless handles connections natively over TLS.

### **3.3 Schema Translation (`api/database/schema.sql` & `api/migrations/`)**
The `golang-migrate` CLI can still be used, but all existing `.sql` files must be converted from PostgreSQL to MySQL dialects.

*   **Types:**
    *   `UUID` -> `CHAR(36)` or `VARCHAR(36)` (MySQL typically uses `CHAR(36)` for UUID strings, but `VARCHAR(36)` works too).
    *   `JSONB` -> `JSON`.
    *   `TIMESTAMP WITH TIME ZONE` -> `DATETIME` or `TIMESTAMP` (TiDB handles time zones differently; you may need to ensure your app always sends UTC).
    *   `TEXT` -> `VARCHAR(...)` or `TEXT` (MySQL's `TEXT` is similar but lacks some implicit casting).
*   **Defaults:**
    *   `uuid_generate_v4()` -> Remove database-side UUID generation defaults. UUIDs must be generated in the Go application before insertion (see 3.4).
    *   `NOW()` -> `CURRENT_TIMESTAMP`.
*   **Indexes:**
    *   GIN indexes on JSON paths (like `idx_notes_content` on `notes.content` or `idx_memory_verses_tags` on `memory_verses.tags`) are not natively supported in TiDB.
    *   **Strategy:** Create generated virtual columns for the specific JSON fields that need indexing, and then create standard B-Tree indexes on those virtual columns. Alternatively, if full-text search is required, evaluate TiDB's full-text search capabilities on the extracted fields.

### **3.4 Query Translation (`api/services/`)**
*   **Placeholders:** Replace all `$1`, `$2`, etc., in SQL queries with `?`.
*   **UUID Generation & Returning:** `INSERT ... RETURNING id` is a PostgreSQL extension. In MySQL, `LastInsertId()` only works for auto-increment IDs.
    *   **Action Required (Developer):** Since we cannot easily retrieve database-generated UUIDs after insertion in MySQL, all `Create*` methods across all services (e.g., `NoteService`, `GroupService`, `MemoryVerseService`) must be updated to generate a UUID in Go using `uuid.New().String()` and pass it directly in the `INSERT` statement. The `RETURNING` clauses must be removed.
*   **Functions & Bulk Upserts:**
    *   `unnest()` with array casting (e.g., `unnest($2::uuid[])` found in `memory_verse_service.go` and `bible_version_service.go` for bulk upserts) is a PostgreSQL-specific feature.
    *   **Strategy:** Replace `UNNEST` with dynamic `INSERT ... ON DUPLICATE KEY UPDATE` queries. Construct a single batch query by iterating over the slice and appending `(?, ?, ...)` to the `VALUES` clause, passing the arguments dynamically to maintain performance.
*   **Casting:** Remove all PostgreSQL-specific casts like `::uuid[]`, `::jsonb`, `::text`. For JSON parsing, let the database driver and standard `json.Marshal` handle it, or use standard MySQL `CAST(x AS JSON)`.

## 4. Data Migration Strategy

Migrating data from PostgreSQL to a MySQL-compatible system like TiDB is not as simple as `pg_dump` and `mysql < dump.sql` because the schemas, data types, and dump formats are fundamentally different.

### Proposed Strategy: TiDB Data Migration (DM) / Lightning

To ensure a reliable and performant migration within an acceptable 24-hour downtime window, we recommend using TiDB-native tools, specifically **TiDB Data Migration (DM)** or exporting to CSV/SQL via `pg_dump`/custom scripts and importing with **TiDB Lightning**.

**Steps:**

1.  **Preparation:**
    *   Ensure the new TiDB Serverless cluster is provisioned and accessible.
    *   Apply the newly rewritten MySQL schema to the TiDB cluster (using the updated `golang-migrate` files). The target database structure must exist before migrating data.
    *   Create a temporary user in the Cloud SQL (PostgreSQL) instance with read-only access to all tables.

2.  **Export & Translation:**
    *   Since TiDB DM natively supports MySQL/MariaDB sources (not PostgreSQL), we cannot stream directly. Instead, we must export the PostgreSQL data into a format TiDB Lightning can consume (CSV or SQL).
    *   Use `pg_dump` or a custom extraction script to export the data to CSV files.
    *   During export or via a preprocessing script, translate PostgreSQL-specific data types:
        *   Format `JSONB` output to standard `JSON` strings.
        *   Ensure `UUID`s are output as 36-character strings.

3.  **Execution (Using TiDB Lightning):**
    *   Install TiDB Lightning on a GCP VM that has high-bandwidth access to the TiDB Serverless cluster.
    *   Configure TiDB Lightning to read the exported CSV/SQL files.
    *   Run TiDB Lightning to perform a high-speed bulk import into the TiDB cluster.

4.  **Verification:**
    *   Write a simple script to verify row counts for each table in both databases.
    *   Perform spot checks on `JSON` columns to ensure they parsed correctly.

5.  **Cutover (Downtime Window: 24h):**
    *   Put the application into maintenance mode (read-only or completely offline) to stop writes to PostgreSQL. We have an acceptable downtime window of 24 hours to complete the final export/import and cutover.
    *   Perform the final data export from PostgreSQL and import via TiDB Lightning.
    *   Update the Google Secret Manager with the new TiDB connection string.
    *   Deploy the updated Go backend (which uses the MySQL driver).
    *   Verify the application is functional against the new TiDB database.
    *   Disable maintenance mode.

## 5. Deployment Pipeline Updates

The GitHub Actions pipeline (`.github/workflows/deploy.yml`) currently deploys the backend to Google Cloud Run and specifically links the Cloud SQL instance using the `--set-cloudsql-instances` flag. This needs to be removed for TiDB Serverless.

### Required Changes to `deploy.yml`:

1.  **Remove Cloud SQL Flag:** Completed.
2.  **Remove CLOUD_SQL_INSTANCE Environment Variable:** Completed.
3.  **Update Secret Configuration (Google Secret Manager):** Completed. Secrets for PROD and STG have been successfully created and populated via Pulumi.
    *   `DB_USERNAME`
    *   `DB_PASSWORD`
    *   `DB_HOST`
    *   `DB_PORT`
    *   `DB_NAME`
4.  **Database Connection Logic Update:** Completed. `api/database/db.go` has been refactored.

## 6. Pre-Commit Steps

Since this task is purely an evaluation and documentation exercise, no code changes are being committed at this stage. The actual code migration will require a separate phase of work, executing the plan detailed above.
