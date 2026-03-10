# TiDB Migration Docker Environment

This folder contains the configuration to spin up a Docker container packed with all the tools required for completing the PostgreSQL to TiDB migration, including:

*   **Google Cloud SDK (`gcloud`)**
*   **Google Cloud SQL Auth Proxy (`cloud-sql-proxy`)**
*   **Python 3 & `psycopg2-binary`**
*   **TiDB Lightning (`tiup tidb-lightning`)**

## Usage Instructions

### 1. Build and Start the Container

From your host machine, navigate to this directory and run:

```bash
cd scripts/migration_env
docker-compose up -d --build
```

### 2. Enter the Environment

Exec into the running container with:

```bash
docker exec -it dj_migrator bash
```

### 3. Authenticate and Connect to Google Cloud SQL

Once you are inside the container bash shell:

1.  **Authenticate (if it didn't map up properly from your host machine):**
    ```bash
    gcloud auth login
    gcloud config set project YOUR_DEVELOPMENT_PROJECT_ID
    ```

2.  **Start the Cloud SQL Proxy:**
    Run the proxy in the background to tunnel the connection on port `5432`:
    ```bash
    cloud-sql-proxy YOUR_PROJECT:YOUR_REGION:YOUR_INSTANCE_NAME &
    ```
    *Ensure you see "Ready for new connections" in the logs before proceeding.*

### 4. Run the Data Extraction

You are now ready to dump the data using the Python script!

```bash
# Set your target database details (if not mapped via docker-compose environment vars)
export PG_USER="postgres"
export PG_PASSWORD="YOUR_PASSWORD"
export PG_DBNAME="discipleship_journal"

# Run the python script to generate CSV files in /app/csv_extracts
python3 /app/scripts/extract_postgres_to_csv.py
```

### 5. TiDB Lightning Import

After reviewing your CSV extracts, you can use `tiup tidb-lightning` directly inside this container to load the CSV files into your Staging or Production TiDB Serverless clusters!
