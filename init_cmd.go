package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var defaultConfig = `package = "%s"

[database]
host = "127.0.0.1"
port = 5432
database = "myapp_development"
user = "myuser"
password = "secret"

[[tables]]
table_name = "customer"
# struct_name = "CustomerRow"
`

func initCmd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "init requires exactly one argument")
		os.Exit(1)
	}
	pkgName := args[0]

	err := os.Mkdir(pkgName, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	file, err := os.Create(filepath.Join(pkgName, "config.toml"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, defaultConfig, pkgName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
