# pgx-crud

## Testing

Create a test database and populate it with the test schema.

    createdb pgxdata
    psql pgxdata -f test/structure.sql

Copy `test/data/config.toml.example` to `test/data/config.toml` and enter database connection information.

Set PG* envvars when running tests to pass connection information to tests.

    PGHOST=/var/run/postgresql rake

Regenerate the test app in test/data as needed (probably need to automate this).
