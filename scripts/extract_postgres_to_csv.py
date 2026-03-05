import os
import csv
import json
import uuid
import sys
from typing import Any, List, Dict
import argparse

try:
    import psycopg2
    from psycopg2.extras import RealDictCursor
except ImportError:
    print("Error: psycopg2-binary is required. Install it via 'pip install psycopg2-binary'")
    sys.exit(1)

# List of all tables we want to extract
USER_CENTRIC_TABLES = [
    "users",
    "journal_entries",
    "connections",
    "groups",
    "group_shares",
    "reading_plans",
    "reading_plan_bookmarks",  # common pattern, but we'll try to extract what exists
    "memory_verses",
    "user_devices",
    "tags",
    "note_tags",
    "study_templates",
    "user_verse_preferences"
]

def serialize_value(val: Any) -> Any:
    """
    Convert Postgres specific python types into generic strings for MySQL.
    - uuid.UUID -> str
    - dict / list (from JSONB) -> json string
    """
    if isinstance(val, uuid.UUID):
        return str(val)
    elif isinstance(val, (dict, list)):
        return json.dumps(val)
    elif val is None:
        return r"\N"  # TiDB / MySQL standard NULL representation in some CSV loads, or just use empty / let csv handle
    return val

def extract_table_to_csv(conn, table_name: str, output_dir: str):
    print(f"Extracting table: {table_name}...")
    
    with conn.cursor() as cursor:
        try:
            # We skip information_schema because Cloud SQL 'postgres' user might not see tables 
            # if they are owned by the 'cloudrun-service' IAM identity without explicit GRANTS.
            # We rely on catching the native failure.
            # Count for progress
            cursor.execute(f'SELECT count(*) FROM "{table_name}"')
            total_rows = cursor.fetchone()[0]
            print(f"Found {total_rows} rows in {table_name}")
            
            if total_rows == 0:
                print(f"Skipping empty table: {table_name}")
                return True

            cursor.execute(f'SELECT * FROM "{table_name}"')
            
            # Get column names
            col_names = [desc[0] for desc in cursor.description]
            
            output_file = os.path.join(output_dir, f"{table_name}.csv")
            with open(output_file, 'w', newline='', encoding='utf-8') as csvfile:
                # We use strict quoting for accurate JSON handling
                writer = csv.writer(csvfile, quoting=csv.QUOTE_MINIMAL)
                writer.writerow(col_names)
                
                rows_written = 0
                for row in cursor:
                    serialized_row = [serialize_value(val) for val in row]
                    writer.writerow(serialized_row)
                    rows_written += 1
                    
                    if rows_written % 1000 == 0:
                        print(f"  ... {rows_written}/{total_rows} rows extracted")
                        
            print(f"Finished extracting {table_name} to {output_file}")
            return True
            
        except Exception as e:
            print(f"Error extracting table {table_name}: {e}")
            conn.rollback()
            raise

def main():
    parser = argparse.ArgumentParser(description="Extract Postgres tables to CSV for TiDB Import")
    parser.add_argument("--host", default=os.getenv("PG_HOST", "localhost"), help="Postgres Host")
    parser.add_argument("--port", default=os.getenv("PG_PORT", "5432"), help="Postgres Port")
    parser.add_argument("--dbname", default=os.getenv("PG_DBNAME", "discipleship_journal"), help="Postgres DB Name")
    parser.add_argument("--user", default=os.getenv("PG_USER", "postgres"), help="Postgres Username")
    parser.add_argument("--password", default=os.getenv("PG_PASSWORD", "postgres"), help="Postgres Password")
    parser.add_argument("--out-dir", default="./csv_extracts", help="Output directory for CSVs")
    parser.add_argument("--all-tables", action="store_true", help="Extract all tables in the public schema instead of hardcoded user tables")

    args = parser.parse_args()

    if not os.path.exists(args.out_dir):
        os.makedirs(args.out_dir)

    print(f"Connecting to Postgres: {args.user}@{args.host}:{args.port}/{args.dbname}")
    
    try:
        conn_kwargs = {
            "host": args.host,
            "port": args.port,
            "dbname": args.dbname,
            "user": args.user,
        }
        
        # If using IAM authentication, sometimes the password should be passed or omitted
        # The proxy auto-iam-authn handles generating the token, but we still have to give it to psycopg.
        # Ensure password is included in kwargs if provided
        if args.password:
             conn_kwargs["password"] = args.password
             
        conn = psycopg2.connect(**conn_kwargs)
        
    except Exception as e:
        print(f"Failed to connect to database: {e}")
        sys.exit(1)

    # Diagnostic Table Logging
    print("--- Diagnostic: What tables actually exist in this database? ---")
    with conn.cursor() as cur:
        try:
            cur.execute("SELECT schemaname, tablename, tableowner FROM pg_catalog.pg_tables WHERE schemaname NOT IN ('pg_catalog', 'information_schema')")
            all_tables = cur.fetchall()
            for t in all_tables:
                print(f"Schema: {t[0]}, Table: {t[1]}, Owner: {t[2]}")
        except Exception as e:
            print(f"Could not load pg_tables catalog: {e}")
            conn.rollback()
    print("----------------------------------------------------------------")

    # Track metrics
    extracted_tables = 0
    tables_to_extract = USER_CENTRIC_TABLES
    
    if args.all_tables:
        print("Using --all-tables rule")
        with conn.cursor() as cursor:
            cursor.execute("""
                SELECT table_name 
                FROM information_schema.tables 
                WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
            """)
            tables_to_extract = [row[0] for row in cursor.fetchall()]

    for table in tables_to_extract:
        result = extract_table_to_csv(conn, table, args.out_dir)
        if result:
            extracted_tables += 1
            
    conn.close()
    
    if extracted_tables == 0:
        print("MIGRATION FAILED: No tables were extracted. Please check database connectivity and target Schema.")
        sys.exit(1)
        
    print(f"Migration extraction complete. Total tables: {extracted_tables}")

if __name__ == "__main__":
    main()
