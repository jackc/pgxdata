package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx"
	"github.com/spf13/cobra"
)

type boxType struct {
	Name      string
	ValueType string
}

type intBoxType struct {
	Name    string
	BitSize int
}

type initData struct {
	PkgName     string
	Version     string
	BoxTypes    []boxType
	IntBoxTypes []intBoxType
}

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

	var c Config
	_, err := toml.DecodeFile("config.toml", &c)
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

	templates := loadTemplates()

	supportData := initData{
		PkgName: c.Package,
		Version: VERSION,
		BoxTypes: []boxType{
			{Name: "Bool", ValueType: "bool"},
			{Name: "Int16", ValueType: "int16"},
			{Name: "Int32", ValueType: "int32"},
			{Name: "Int64", ValueType: "int64"},
			{Name: "String", ValueType: "string"},
			{Name: "Time", ValueType: "time.Time"},
		},
		IntBoxTypes: []intBoxType{
			{Name: "Int16", BitSize: 16},
			{Name: "Int32", BitSize: 32},
			{Name: "Int64", BitSize: 64},
		},
	}

	supportFiles := []struct {
		path string
		tmpl *template.Template
	}{
		{"attribute.go", templates.Lookup("attribute")},
		{"db.go", templates.Lookup("db")},
	}
	for _, f := range supportFiles {
		err := writeSupportFile(f.path, f.tmpl, supportData)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	for _, t := range c.Tables {
		file, err := os.Create(goCaseToFileCase(t.StructName) + ".go")
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

func goCaseToFileCase(g string) string {
	buf := &bytes.Buffer{}

	for i, r := range g {
		if unicode.IsUpper(r) && i != 0 {
			buf.WriteRune('_')
		}
		buf.WriteRune(unicode.ToLower(r))
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

func writeSupportFile(path string, tmpl *template.Template, data initData) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
