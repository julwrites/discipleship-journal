import pymysql

try:
    conn = pymysql.connect(
        host='dummy',
        user='dummy',
        password='dummy',
        port=4000,
        ssl={'ssl': {'ca': ''}}
    )
except pymysql.err.OperationalError as e:
    code, msg = e.args
    print(f"OperationalError: code={code}, msg={msg}")
except Exception as e:
    print(e)
