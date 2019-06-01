package data_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	pgxpool "github.com/jackc/pgx/v4/pool"
)

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	flag.Parse()

	var err error
	pool, err = createConnPool()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create connection pool:", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func createConnPool() (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), "")
}

func begin(t *testing.T) *pgxpool.Tx {
	tx, err := pool.Begin(context.Background(), nil)
	if err != nil {
		t.Fatal(t)
	}

	return tx
}
