package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx"
	"github.com/spf13/cobra"
)

var dataTypeMap = map[string]string{
	"integer":                  "Int32",
	"character varying":        "String",
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

	err = templates.ExecuteTemplate(os.Stdout, "row", struct {
		PkgName    string
		StructName string
		Columns    []Column
	}{
		PkgName:    c.Package,
		StructName: c.Tables[0].StructName,
		Columns:    c.Tables[0].Columns,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = templates.ExecuteTemplate(os.Stdout, "count_func", struct {
		TableName  string
		StructName string
	}{
		TableName:  c.Tables[0].TableName,
		StructName: c.Tables[0].StructName,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = templates.ExecuteTemplate(os.Stdout, "select_all_func", struct {
		TableName  string
		StructName string
		Columns    []Column
	}{
		TableName:  c.Tables[0].TableName,
		StructName: c.Tables[0].StructName,
		Columns:    c.Tables[0].Columns,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = templates.ExecuteTemplate(os.Stdout, "select_by_id", struct {
		TableName  string
		StructName string
		Columns    []Column
	}{
		TableName:  c.Tables[0].TableName,
		StructName: c.Tables[0].StructName,
		Columns:    c.Tables[0].Columns,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = templates.ExecuteTemplate(os.Stdout, "insert", struct {
		TableName  string
		StructName string
		Columns    []Column
	}{
		TableName:  c.Tables[0].TableName,
		StructName: c.Tables[0].StructName,
		Columns:    c.Tables[0].Columns,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = templates.ExecuteTemplate(os.Stdout, "update", struct {
		TableName  string
		StructName string
		Columns    []Column
	}{
		TableName:  c.Tables[0].TableName,
		StructName: c.Tables[0].StructName,
		Columns:    c.Tables[0].Columns,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = templates.ExecuteTemplate(os.Stdout, "delete", struct {
		TableName  string
		StructName string
		Columns    []Column
	}{
		TableName:  c.Tables[0].TableName,
		StructName: c.Tables[0].StructName,
		Columns:    c.Tables[0].Columns,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Queryer interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (pgx.CommandTag, error)
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
