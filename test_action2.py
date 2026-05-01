import urllib.request
import json
url = "https://api.github.com/repos/julwrites/discipleship-journal/actions/workflows/migration_postgres_extract.yml"
req = urllib.request.Request(url)
try:
    with urllib.request.urlopen(req) as response:
        print(response.read().decode('utf-8'))
except Exception as e:
    print(e)
