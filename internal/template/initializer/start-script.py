import os, sys, time, MySQLdb

print("Initialization task starts.")

host = os.environ.get('FE_SVC')
port = os.environ.get('FE_QUERY_PORT')
acc_user = os.environ.get('ACC_USER')
acc_password = os.environ.get('ACC_PWD')

# connect to fe
retry_count = 0
for i in range(0, 10):
    try:
        conn = MySQLdb.connect(host=host,
                               port=int(port),
                               user=acc_user,
                               passwd=acc_password,
                               connect_timeout=5,
                               charset='utf8')
    except MySQLdb.OperationalError as e:
        print(e)
        retry_count += 1
        time.sleep(1)
        continue
    break
if retry_count == 10:
    print("Failed to connect to FE.")
    sys.exit(1)

print("Create FE query connection successfully.")

# modify the password of default user
password_dir = '/etc/doris/password'

for file in os.listdir(password_dir):
    if file.startswith('.'):
        continue
    user = file
    with open(os.path.join(password_dir, file), 'r') as f:
        lines = f.read().splitlines()
        password = lines[0] if len(lines) > 0 else ""
    if user == 'admin':
        conn.cursor().execute("set password for %s = password(%s);", (user, password))
        print("Reset password for user: admin")
    elif user == 'root':
        try:
            _conn = MySQLdb.connect(host=host, port=int(port), user='root', connect_timeout=5)
            _conn.cursor().execute("set password for %s = password(%s);", (user, password))
            _conn.commit()
            _conn.close()
            print("Reset password for user: root")
        except MySQLdb.OperationalError as e:
            print('Skip setting password for user: root')
    else:
        print("Skip setting password for user: %s" % user)

# execute init sql scripts
if os.path.isfile('/etc/doris/init.sql'):
    with open('/etc/doris/init.sql', 'r') as f:
        sql = f.read()
        # split sqls
        stmts = sql.split(';')
        sql_stmts = [stmt.strip() for stmt in stmts if stmt.strip()]

        # execute sql
        for stmt in sql_stmts:
            print("Execute sql: " + stmt)
            try:
                conn.cursor().execute(stmt)
                conn.commit()
            except MySQLdb.OperationalError as e:
                print("Fail to execute sql:" + stmt)
                print(e)
                conn.close()
                sys.exit(1)
else:
    print("init.sql is empty, skip executing init sql scripts.")

conn.commit()
conn.close()
print("Initialization task finished.")
