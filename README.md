# pgx-crud

## Testing

Create a test database and populate it with the test schema.

    createdb pgxdata
    psql pgxdata -f test/structure.sql

Set PG* envvars when running tests to pass connection information to tests.

    PGHOST=/var/run/postgresql PGDATABASE=pgxdata rake

Regenerate the test app in test/data as needed (probably need to automate this).
