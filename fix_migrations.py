import os
import glob
import re

mig_dir = "api/migrations"
for f_path in glob.glob(f"{mig_dir}/*.up.sql"):
    with open(f_path, 'r') as f:
        content = f.read()
    
    # 1. Fix ADD COLUMN ... UNIQUE
    def repl_unique(m):
        table = m.group(1)
        col = m.group(2)
        rest = m.group(3)
        return f"ALTER TABLE {table} ADD COLUMN {col} {rest};\nALTER TABLE {table} ADD UNIQUE INDEX idx_{table}_{col} ({col});"

    # Match: ALTER TABLE users ADD COLUMN username VARCHAR(50) UNIQUE;
    pattern = re.compile(r"ALTER TABLE\s+([A-Za-z0-9_]+)\s+ADD COLUMN\s+([A-Za-z0-9_]+)\s+([^;]*?)\s+UNIQUE;", re.IGNORECASE)
    content = pattern.sub(repl_unique, content)
    
    # 2. Fix ALTER TABLE memory_verses MODIFY verse_pack_id VARCHAR(36) NOT NULL;
    # (TiDB/MySQL requires MODIFY COLUMN or MODIFY)
    # Actually wait, MySQL uses `MODIFY verse_pack_id VARCHAR(36) NOT NULL` which is completely fine!
    
    with open(f_path, 'w') as f:
        f.write(content)

print("Migrations patched")
