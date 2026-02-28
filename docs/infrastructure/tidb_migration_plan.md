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

### Actions Required From You:
1.  **Install Pulumi:** Download and install the Pulumi CLI from [pulumi.com](https://www.pulumi.com/docs/install/).
2.  **Create a Pulumi Account:** Sign up or log in via the CLI (`pulumi login`).
3.  **Set Up a TiDB Cloud Account:** Create an account at [tidbcloud.com](https://tidbcloud.com) if you don't have one.
4.  **Get TiDB Cloud API Keys:** Generate a Public Key and a Private Key in the TiDB Cloud Console (Project Settings > API Access).
5.  **Initialize Pulumi:** Create a new directory for infrastructure management (e.g., `infra/` at the root of the project).
    *   Navigate to the directory: `mkdir infra && cd infra`
    *   Run `pulumi new typescript`
    *   Install the TiDB Cloud provider: `npm install @tidbcloud/pulumi-tidbcloud`
6.  **Configure Pulumi with Credentials:** Set the configuration values (in your terminal within the `infra/` directory):
    *   `pulumi config set tidbcloud:publicKey <your-public-key> --secret`
    *   `pulumi config set tidbcloud:privateKey <your-private-key> --secret`
    *   `pulumi config set tidbcloud:projectId <your-tidb-project-id>`

### Proposed Pulumi Code (`infra/index.ts`):

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as tidbcloud from "@tidbcloud/pulumi-tidbcloud";

// Ensure you have configured the provider with the necessary credentials
// (publicKey, privateKey) and defaults (projectId).

const config = new pulumi.Config("tidbcloud");
const projectId = config.require("projectId");

// Provision a TiDB Cloud Serverless Cluster
const serverlessCluster = new tidbcloud.Cluster("discipleship-journal-tidb", {
    projectId: projectId,
    name: "dj-serverless-db",
    clusterType: "DEVELOPER", // Represents the Serverless tier
    cloudProvider: "AWS", // Or GCP if preferred and available in your region
    region: "us-east-1",  // Adjust to your desired region
    serverless: {
        spendLimit: {
            monthly: 0, // 0 for free tier if applicable, adjust as needed
        },
    },
});

// Create a database password (store this securely, Pulumi will manage it as a secret)
const dbPassword = new pulumi.random.RandomPassword("dbPassword", {
    length: 16,
    special: true,
    overrideSpecial: "_%@",
});

// Since the DB user creation and permission granting is not natively supported
// by the standard Pulumi TiDB provider at this moment, you will typically use
// the default 'root' user provided by the serverless cluster or manually create one.
// For production, you should create a dedicated user via a custom resource or a script.

// Export the connection details
export const clusterId = serverlessCluster.id;
export const databaseEndpoint = serverlessCluster.endpoints.apply(e => e.standard.host);
export const databasePort = serverlessCluster.endpoints.apply(e => e.standard.port);
export const rootPassword = dbPassword.result;

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
    *   `uuid_generate_v4()` -> `UUID()`.
    *   `NOW()` -> `CURRENT_TIMESTAMP`.
*   **Indexes:**
    *   GIN indexes on JSON paths (like `notes.content`) are not natively supported in TiDB the same way. You may need to create generated columns on specific JSON fields and index those instead.

### **3.4 Query Translation (`api/services/`)**
*   **Placeholders:** Replace all `$1`, `$2`, etc., in SQL queries with `?`.
*   **Returning:** `INSERT ... RETURNING id` is a PostgreSQL extension. In MySQL, use `LastInsertId()` on the `sql.Result` object. *Wait, this only works for auto-increment IDs. For UUIDs, you must generate the UUID in Go before inserting it, as MySQL's `UUID()` function cannot be retrieved via `LastInsertId()`.*
    *   **Action Required (Developer):** All `Create*` methods in the services will need to generate a `uuid.New().String()` in Go and insert it directly, instead of relying on the database default.
*   **Functions:**
    *   `UNNEST($1::TEXT[], $2::TEXT[])` (found in `BibleVersionService.SyncVersions` according to memories) must be completely rewritten. Since TiDB doesn't support array types or unnesting them, you will likely need to use multiple individual `INSERT` statements or a single `INSERT ... VALUES (?, ?), (?, ?)` built dynamically, or use `INSERT ... ON DUPLICATE KEY UPDATE` instead of `ON CONFLICT DO UPDATE`.
*   **Casting:** Remove all `::uuid`, `::jsonb`, `::text[]` casts.

## 4. Data Migration Strategy

Migrating data from PostgreSQL to a MySQL-compatible system like TiDB is not as simple as `pg_dump` and `mysql < dump.sql` because the schemas, data types, and dump formats are fundamentally different.

### Proposed Strategy: Logical Replication via pgloader

The most reliable approach is to use a tool that translates data on the fly. We recommend `pgloader`, an open-source data migration tool designed exactly for this purpose.

**Steps:**

1.  **Preparation:**
    *   Ensure the new TiDB Serverless cluster is provisioned and accessible.
    *   Apply the newly rewritten MySQL schema to the TiDB cluster (using the updated `golang-migrate` files). The target database structure must exist before migrating data.
    *   Create a temporary user in the Cloud SQL (PostgreSQL) instance with read-only access to all tables.

2.  **Execution (Using pgloader):**
    *   You can run `pgloader` locally (if you have network access to both databases) or from a VM within GCP that has access to both the Cloud SQL instance and the public internet (for TiDB Serverless).
    *   **Action Required (You):** Run a command similar to the following:
        ```bash
        pgloader postgresql://<pg_user>:<pg_pass>@<pg_host>:5432/discipleship_journal \
                 mysql://<tidb_user>:<tidb_pass>@<tidb_endpoint>:4000/discipleship_journal
        ```
    *   *Note:* You may need to write a custom `.load` file for `pgloader` to handle specific type conversions (like PostgreSQL `UUID` to MySQL `CHAR(36)` or `JSONB` to `JSON`).

3.  **Verification:**
    *   Write a simple script to verify row counts for each table in both databases.
    *   Perform spot checks on `JSON` columns to ensure they parsed correctly.

4.  **Cutover:**
    *   Put the application into a maintenance mode (read-only or completely offline) to stop writes to PostgreSQL.
    *   Run `pgloader` one final time to sync any last-minute changes (if you didn't set up continuous replication).
    *   Update the Google Secret Manager with the new TiDB connection string.
    *   Deploy the updated Go backend (which uses the MySQL driver).
    *   Disable maintenance mode.

## 5. Deployment Pipeline Updates

The GitHub Actions pipeline (`.github/workflows/deploy.yml`) currently deploys the backend to Google Cloud Run and specifically links the Cloud SQL instance using the `--set-cloudsql-instances` flag. This needs to be removed for TiDB Serverless.

### Required Changes to `deploy.yml`:

1.  **Remove Cloud SQL Flag:** In the `deploy-backend` job, locate the `google-github-actions/deploy-cloudrun@v2` step.
    *   **Change From:**
        ```yaml
        flags: '--service-account=${{ secrets.GCP_SERVICE_ACCOUNT }} --allow-unauthenticated --set-cloudsql-instances=${{ steps.vars.outputs.CLOUD_SQL_INSTANCE }}'
        ```
    *   **Change To:**
        ```yaml
        flags: '--service-account=${{ secrets.GCP_SERVICE_ACCOUNT }} --allow-unauthenticated'
        ```
2.  **Remove CLOUD_SQL_INSTANCE Environment Variable:** Remove the `CLOUD_SQL_INSTANCE` variable from the `env_vars` section in the same step.
3.  **Update Secret Configuration (Google Secret Manager):** The deployment script itself doesn't need to change how secrets are loaded, but the *contents* of the secrets in GCP must change.
    *   **Action Required (You):** Update the following secrets in Google Secret Manager for both `main` and `staging` (prefixed with `STG_`):
        *   `DB_USERNAME` -> TiDB Username (e.g., `root`)
        *   `DB_PASSWORD` -> TiDB Password
        *   `DB_HOST` -> TiDB Endpoint (e.g., `gateway01.us-east-1.prod.aws.tidbcloud.com`)
        *   `DB_PORT` -> `4000` (TiDB default)
        *   `DB_NAME` -> `discipleship_journal`
4.  **Database Connection Logic Update:** As mentioned in Section 3, the Go application's `api/database/db.go` must be updated to construct a standard MySQL DSN using these variables, rather than relying on the Cloud SQL Connector.

## 6. Pre-Commit Steps

Since this task is purely an evaluation and documentation exercise, no code changes are being committed at this stage. The actual code migration will require a separate phase of work, executing the plan detailed above.
