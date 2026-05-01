import yaml
with open('.github/workflows/migration_postgres_extract.yml', 'r') as file:
    data = yaml.safe_load(file)
print("YAML parsed successfully")
