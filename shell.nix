let
  pkgs = import <nixpkgs> { config.allowUnfree = true; };
  oracleLib = pkgs.oracle-instantclient.lib;
in
pkgs.mkShell {
  buildInputs = with pkgs; [
    oracle-instantclient
    postgresql_15
    go
    sqlite

    vhs
  ];

  hardeningDisable = [ "fortify" ];

  shellHook = ''
    export PGDATA=$PWD/postgres_data
    export PGHOST=$PWD/postgres
    export LOG_PATH=$PWD/postgres/LOG
    export PGDATABASE=postgres
    export DATABASE_URL="postgresql:///sampledb?host=$PGHOST"

    export SQLITE_DB_PATH="$PWD/sample_sqlite.db"
    export SQLITE_CONNECTION_STRING="file:$SQLITE_DB_PATH"

    export LD_LIBRARY_PATH=${oracleLib}/lib:$LD_LIBRARY_PATH
    export ORACLE_HOME=${oracleLib}

    # PostgreSQL setup
    if [ ! -d $PGHOST ]; then
      mkdir -p $PGHOST
    fi

    if [ ! -d $PGDATA ]; then
      echo "Initializing PostgreSQL database..."
      initdb $PGDATA --auth=trust >/dev/null
    fi

    pg_ctl start -l $LOG_PATH -o "-c unix_socket_directories=$PGHOST -c listen_addresses= -c port=5432"
    echo "PostgreSQL started successfully!"

    if ! psql -lqt | cut -d \| -f 1 | grep -qw sampledb; then
      echo ""
      echo "Setting up sample database with initial data..."
      psql -f init.sql
    else
      echo "Sample database already exists."
    fi

    echo ""
    echo "========================================="
    echo "PostgreSQL is ready!"
    echo "========================================="
    echo "Database URL: $DATABASE_URL"
    echo ""

    # SQLite setup
    if [ ! -f "$SQLITE_DB_PATH" ]; then
      echo "Creating SQLite sample database..."
      if [ -f init.sql ]; then
        sqlite3 "$SQLITE_DB_PATH" < init.sql
      else
        sqlite3 "$SQLITE_DB_PATH" "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT);"
        sqlite3 "$SQLITE_DB_PATH" "INSERT INTO users (username) VALUES ('alice'), ('bob');"
      fi
    fi

    echo "========================================="
    echo "SQLite sample database ready!"
    echo "========================================="
    echo "SQLite connection string:"
    echo "  $SQLITE_CONNECTION_STRING"
    echo ""
    echo "Useful SQLite commands:"
    echo "  sqlite3 $SQLITE_DB_PATH"
    echo "  SELECT * FROM users;"
    echo ""
  '';

  exitHook = ''
    pg_ctl stop
    echo "PostgreSQL stopped"
  '';
}
