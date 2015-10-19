package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx"
	"github.com/spf13/cobra"
)

var dataTypeMap = map[string]string{
	"bigint":            "Int64",
	"integer":           "Int32",
	"smallint":          "Int16",
	"character varying": "String",
	"date":              "Time",
	"timestamp with time zone": "Time",
}

type Config struct {
	Package  string
	Database pgx.ConnConfig
	Tables   []Table
}

type Column struct {
	ColumnName      string
	DataType        string
	OrdinalPosition int32

	FieldName string
	GoType    string
}

type Table struct {
	TableName  string `toml:"table_name"`
	StructName string `toml:"struct_name"`
	Columns    []Column
}

func generateCmd(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Fprintln(os.Stderr, "generate does not take any arguments")
		os.Exit(1)
	}

	templates, err := template.ParseGlob(filepath.Join(appArgs.templatesPath, "*"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var c Config

	_, err = toml.DecodeFile("config.toml", &c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	conn, err := pgx.Connect(c.Database)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = inspectDatabase(conn, c.Tables)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, t := range c.Tables {
		file, err := os.Create(t.StructName + ".go")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = writeTableCrud(file, templates, c.Package, t)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		file.Close()
	}
}

type Queryer interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (pgx.CommandTag, error)
}

func writeTableCrud(w io.Writer, templates *template.Template, pkgName string, table Table) error {
	return templates.ExecuteTemplate(w, "row", struct {
		PkgName    string
		TableName  string
		StructName string
		Columns    []Column
	}{
		PkgName:    pkgName,
		TableName:  table.TableName,
		StructName: table.StructName,
		Columns:    table.Columns,
	})
}

func inspectDatabase(db Queryer, tables []Table) error {
	for i := range tables {
		rows, err := db.Query(`select column_name, data_type, ordinal_position from information_schema.columns where table_name=$1`, tables[i].TableName)
		if err != nil {
			return err
		}

		var columns []Column
		for rows.Next() {
			var c Column
			rows.Scan(&c.ColumnName, &c.DataType, &c.OrdinalPosition)
			c.FieldName = pgCaseToGoCase(c.ColumnName)
			c.GoType = pgTypeToGoType(c.DataType)
			columns = append(columns, c)
		}

		if rows.Err() != nil {
			return err
		}

		tables[i].Columns = columns
	}

	return nil
}

func pgCaseToGoCase(pg string) string {
	parts := strings.Split(pg, "_")
	buf := &bytes.Buffer{}
	for _, s := range parts {
		if s == "id" {
			buf.WriteString("ID")
		} else {
			buf.WriteString(strings.Title(s))
		}
	}

	return buf.String()
}

func pgTypeToGoType(pg string) string {
	if t, ok := dataTypeMap[pg]; ok {
		return t
	} else {
		return "String"
	}
}
