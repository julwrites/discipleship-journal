import os
import csv
import sys
import argparse

try:
    import pymysql
except ImportError:
    print("Error: PyMySQL is required. Install it via 'pip install pymysql'")
    sys.exit(1)

def import_table(conn, table_name: str, csv_file_path: str):
    print(f"Importing table: {table_name} from {csv_file_path}...")
    
    if not os.path.exists(csv_file_path):
        print(f"Skipping {table_name} - CSV file not found.")
        return False
        
    with open(csv_file_path, 'r', encoding='utf-8') as f:
        reader = csv.reader(f)
        try:
            columns = next(reader)
        except StopIteration:
            print(f"Skipping {table_name} - empty CSV.")
            return False
            
        # Optional: You can map column names if schema changed between Postgres and MySQL,
        # but if they remain identical, we just use the original names
        col_string = ", ".join([f"`{c}`" for c in columns])
        placeholders = ", ".join(["%s" for _ in columns])
        
        insert_query = f"INSERT INTO `{table_name}` ({col_string}) VALUES ({placeholders})"
        
        # We might want to handle duplicates gracefully (ON DUPLICATE KEY UPDATE or IGNORE)
        # depending on if we already have partial data. Using IGNORE for safety.
        insert_query = f"INSERT IGNORE INTO `{table_name}` ({col_string}) VALUES ({placeholders})"
        
        with conn.cursor() as cursor:
            # First, check if TiDB table exists 
            cursor.execute("SHOW TABLES LIKE %s", (table_name,))
            if not cursor.fetchone():
               print(f"Error: Table `{table_name}` does not exist in TiDB schema! Skipping...")
               return False

            # Disable foreign key checks for the session so we don't fail if we insert out of order
            cursor.execute("SET FOREIGN_KEY_CHECKS=0;")
            
            # Read rows and insert in batches
            batch = []
            batch_size = 500
            total_inserted = 0
            
            for row in reader:
                # Replace r"\N" (Postgres NULL export representation) with actual None
                sanitized_row = [None if val == r'\N' or val == '' else val for val in row]
                batch.append(sanitized_row)
                
                if len(batch) >= batch_size:
                    cursor.executemany(insert_query, batch)
                    total_inserted += len(batch)
                    batch = []
                    
            if batch:
                cursor.executemany(insert_query, batch)
                total_inserted += len(batch)
                
            cursor.execute("SET FOREIGN_KEY_CHECKS=1;")
            conn.commit()
            
            print(f"Successfully imported {total_inserted} rows into `{table_name}`")
            return True

def main():
    parser = argparse.ArgumentParser(description="Import Postgres-extracted CSVs into TiDB Serverless")
    parser.add_argument("--host", required=True, help="TiDB Host endpoint")
    parser.add_argument("--port", default="4000", help="TiDB Port (Default: 4000)")
    parser.add_argument("--dbname", required=True, help="TiDB Database Name")
    parser.add_argument("--user", required=True, help="TiDB Username (Make sure to wrap in quotes if it has prefix)")
    parser.add_argument("--password", required=True, help="TiDB Password")
    parser.add_argument("--csv-dir", required=True, help="Relative path to unzipped folder containing CSV files")

    args = parser.parse_args()

    # Define the precise order of importing to respect dependencies (even though we disabled FK checks, it's good practice)
    ORDERED_TABLES = [
        "users",
        "groups",
        "connections",
        "group_members",
        "reading_plans",
        "reading_plan_days",
        "user_reading_plans",
        "user_reading_plan_progress",
        "notes",
        "group_shares",
        "memory_verses",
        "verse_packs",
        "study_templates",
        "tags",
        "note_tags",
        "user_devices",
        "user_verse_preferences"
    ]

    print(f"Attempting to connect to TiDB Serverless cluster...")
    try:
        # Connect initially without specifying a database to ensure we can create it if it's missing
        conn = pymysql.connect(
            host=args.host,
            user=args.user,
            password=args.password,
            port=int(args.port),
            ssl={'ssl': {'ca': ''}} # PyMySQL will generally default trust the system certs for TLS
        )
        print("Connected successfully!")
        
        # Provision the database explicitly so we never hit an "Unknown Database" fatal error
        with conn.cursor() as cursor:
             cursor.execute(f"CREATE DATABASE IF NOT EXISTS `{args.dbname}`")
             cursor.execute(f"USE `{args.dbname}`")
             
    except Exception as e:
        print(f"Failed to connect to TiDB: {e}")
        sys.exit(1)

    successful_imports = 0
    # Walk through the tables requested
    for table_name in ORDERED_TABLES:
        csv_path = os.path.join(args.csv_dir, f"{table_name}.csv")
        # Ensure we only try tables that were actually extracted (we gracefully skip missing ones above)
        if os.path.exists(csv_path):
            if import_table(conn, table_name, csv_path):
                successful_imports += 1

    conn.close()
    print(f"\nMigration Complete! Imported {successful_imports} tables to TiDB.")

if __name__ == "__main__":
    main()
