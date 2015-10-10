# pgx-crud

## Testing

Create a test database and populate it with the test schema.

    createdb pgx_crud
    psql pgx_crud -f testdata/structure.sql

Set PG* envvars when running tests to pass connection information to tests.

    PGHOST=/var/run/postgresql go test
