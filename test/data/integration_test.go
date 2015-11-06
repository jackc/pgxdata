package data_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx"
)

var pool *pgx.ConnPool

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

func createConnPool() (*pgx.ConnPool, error) {
	var config pgx.ConnPoolConfig
	var err error
	config.ConnConfig, err = pgx.ParseEnvLibpq()
	if err != nil {
		return nil, err
	}

	if config.Host == "" {
		config.Host = "localhost"
	}

	if config.User == "" {
		config.User = os.Getenv("USER")
	}

	if config.Database == "" {
		config.Database = "pgxdata"
	}

	config.TLSConfig = nil
	config.UseFallbackTLS = false
	config.MaxConnections = 10

	return pgx.NewConnPool(config)
}

func begin(t *testing.T) *pgx.Tx {
	tx, err := pool.Begin()
	if err != nil {
		t.Fatal(t)
	}

	return tx
}
