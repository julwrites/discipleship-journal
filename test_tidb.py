import os
import subprocess
import pymysql

def get_secret(name):
    try:
        res = subprocess.run(
            ["gcloud", "secrets", "versions", "access", "latest", "--secret", name],
            capture_output=True, text=True, check=True
        )
        return res.stdout.strip()
    except subprocess.CalledProcessError as e:
        print(f"Failed to get secret {name}: {e.stderr}")
        return None

def main():
    host = get_secret("STG_DJ_DB_HOST")
    user = get_secret("STG_DJ_DB_USERNAME")
    pwd = get_secret("STG_DJ_DB_PASSWORD")
    dbname = get_secret("STG_DJ_DB_NAME")
    
    if not (host and user and pwd):
        print("Missing credentials")
        return
        
    print(f"Connecting to {host} as {user} for {dbname}")
    
    try:
        conn = pymysql.connect(
            host=host,
            user=user,
            password=pwd,
            port=4000,
            ssl={'ssl': {'ca': ''}}
        )
        with conn.cursor() as cursor:
            cursor.execute("SHOW DATABASES;")
            dbs = cursor.fetchall()
            print("Databases:", dbs)
            
            cursor.execute(f"CREATE DATABASE IF NOT EXISTS `{dbname}`;")
            print("Create DB returned.")
            
            cursor.execute(f"USE `{dbname}`;")
            print("USE DB returned.")
            
    except Exception as e:
        print("Error:", repr(e))

if __name__ == "__main__":
    main()
