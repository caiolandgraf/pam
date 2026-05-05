# Database Support
To create and save a database connection to be used with other `pam` commands, use the `init` command
```bash
init <name> <type> <conn-string> [schema]
```
If the `<type>` is omitted, PAM will try to infer the database type from the connection string.

> **Environment variable expansion:** Connection strings support `${VAR}` substitution. PAM expands environment variables at runtime, so you can safely store credentials outside your config file:
> ```bash
> pam init prod postgres "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:5432/mydb"
> ```

## Init Examples

### PostgreSQL

```bash
pam init pg-prod postgres postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable

# or connect to a specific schema:
pam init pg-prod postgres postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable schema-name
```

### MySQL / MariaDB

```bash
pam init mysql-dev mysql 'myuser:mypassword@tcp(127.0.0.1:3306)/mydb'

pam init mariadb-docker mariadb "root:MyStrongPass123@tcp(localhost:3306)/forestgrove"
```

### SQL Server


```bash
pam init sqlserver-docker sqlserver "sqlserver://sa:MyStrongPass123@localhost:1433/master"
```

### SQLite

```bash
pam init sqlite-local sqlite file:///home/eduardo/dbeesly/sqlite/mydb.sqlite
```

### DuckDB

> Requires CGO (included by default; excluded by building with `CGO_ENABLED=0`).

```bash
pam init duckdb-local duckdb /path/to/mydb.db
```

DuckDB can also query CSV and JSON files directly:

```bash
pam init duckdb-local duckdb employees.json
pam init duckdb-local duckdb employees.csv
```

For json and csv, the data will be available to query as a view, with the same name as the file
(eg. `pam run "select * from employees" from the employees.json file`)

### Oracle

```bash
pam init oracle-stg oracle "oracle://myuser:mypassword@localhost:1521/XEPDB1"

# or connect to a specific schema:
pam init oracle-stg oracle "oracle://myuser:mypassword@localhost:1521/XEPDB1" schema-name
```

### ClickHouse

```bash
pam init clickhouse-docker clickhouse "clickhouse://myuser:mypassword@localhost:9000/forestgrove"
```

### FireBird

```bash
pam init firebird-docker firebird user:masterkey@localhost:3050//var/lib/firebird/data/the_office
```

### Snowflake

> Supports keypair authentication for secure, password-free connections.

```bash
# Username/password
pam init sf-prod snowflake "user:password@account/dbname/schema"

# Keypair authentication (recommended)
pam init sf-prod snowflake "user@account/dbname/schema?authenticator=snowflake_jwt&privateKeyPath=/path/to/rsa_key.p8"

# With environment variable expansion
pam init sf-prod snowflake "${SF_USER}:${SF_PASS}@${SF_ACCOUNT}/mydb"
```

---

## 🐝 Dbeesly

To run containerized test database servers for all supported databases, use the sister project [dbeesly](https://github.com/eduardofuncao/dbeesly)

<img width="879" height="571" alt="image" src="https://github.com/user-attachments/assets/c0a131eb-ea95-4523-86ac-cd00a561a5e0" />
