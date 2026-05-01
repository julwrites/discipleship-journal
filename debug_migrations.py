import os
import pymysql

# We use test_tidb logic since we validated it
from test_tidb import get_secret

host = get_secret("STG_DJ_DB_HOST")
user = get_secret("STG_DJ_DB_USERNAME")
pwd = get_secret("STG_DJ_DB_PASSWORD")
dbname = get_secret("STG_DJ_DB_NAME")

conn = pymysql.connect(
    host=host, user=user, password=pwd, database=dbname,
    port=4000, ssl={'ssl': {'ca': ''}},
    client_flag=pymysql.constants.CLIENT.MULTI_STATEMENTS
)

migrations_dir = "api/migrations"
up_files = sorted([f for f in os.listdir(migrations_dir) if f.endswith('.up.sql')])

with conn.cursor() as cursor:
    cursor.execute("SET FOREIGN_KEY_CHECKS=0;")
    for f in up_files:
        if "reseed" in f:
            continue
        sql_path = os.path.join(migrations_dir, f)
        with open(sql_path, 'r', encoding='utf-8') as sql_file:
            sql_content = sql_file.read().strip()
            if sql_content:
                try:
                    cursor.execute(sql_content)
                    print(f"SUCCESS: {f}")
                except Exception as e:
                    print(f"ERROR: {f} -> {e}")
    cursor.execute("SET FOREIGN_KEY_CHECKS=1;")
    
    # We purposefully don't commit so we don't pollute the db completely while debugging
