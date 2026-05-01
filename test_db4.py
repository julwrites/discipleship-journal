import pymysql
import sys

# Flush prints instantly
def flprint(msg):
    print(msg)
    sys.stdout.flush()

flprint("Connecting...")
try:
    raise pymysql.err.OperationalError(1049, "Unknown database 'test'")
except Exception as e:
    flprint(f"Failed: {e}")
    sys.exit(1)
