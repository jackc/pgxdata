package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx"
)

func main() {
	config, err := extractConfig()
	if err != nil {
		log.Fatalln(err)
	}

	pool, err := pgx.NewConnPool(config)
	if err != nil {
		log.Fatalln(err)
	}

	newCr := &CustomerRow{
		FirstName: String{"John", Present},
		LastName:  String{"Smith", Present},
	}
	err = InsertCustomer(pool, newCr)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(newCr)

	var n int64
	n, err = CountCustomer(pool)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("customer count", n)

	crs, err := SelectAllCustomer(pool)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("all customers", crs)

	cr, err := SelectCustomerByID(pool, 3)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("customer 3:", cr)

	err = UpdateCustomer(pool, cr.ID.Value, &CustomerRow{FirstName: String{"Phil", Present}})
	if err != nil {
		log.Fatalln(err)
	}

	err = DeleteCustomer(pool, newCr.ID.Value)
	if err != nil {
		log.Fatalln(err)
	}
}

func extractConfig() (config pgx.ConnPoolConfig, err error) {
	config.ConnConfig, err = pgx.ParseEnvLibpq()
	if err != nil {
		return config, err
	}

	if config.Host == "" {
		config.Host = "localhost"
	}

	if config.User == "" {
		config.User = os.Getenv("USER")
	}

	if config.Database == "" {
		config.Database = "crud"
	}

	config.TLSConfig = nil
	config.UseFallbackTLS = false

	config.MaxConnections = 10

	return config, nil
}
