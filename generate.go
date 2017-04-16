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
	Name       string
	ValueType  string
	FormatCode string
}

type intBoxType struct {
	Name    string
	BitSize int
}

type initData struct {
	PkgName string
	Version string
}

var pgToBoxTypeMap = map[string]string{
	"bigint":            "pgtype.Int8",
	"integer":           "pgtype.Int4",
	"smallint":          "pgtype.Int2",
	"character varying": "pgtype.Varchar",
	"text":              "pgtype.Text",
	"date":              "pgtype.Date",
	"timestamp with time zone": "pgtype.Timestamptz",
	"inet":  "pgtype.Inet",
	"cidr":  "pgtype.Cidr",
	"bytea": "pgtype.Bytea",
}

var pgToGoTypeMap = map[string]string{
	"bigint":            "int64",
	"integer":           "int32",
	"smallint":          "int16",
	"character varying": "string",
	"text":              "string",
	"date":              "time.Time",
	"timestamp with time zone": "time.Time",
	"bytea":                    "[]byte",
}

var acronyms = map[string]bool{
	"id":  true,
	"ip":  true,
	"url": true,
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
	GoBoxType string

	VarName string
	GoType  string
}

type ColumnConfig struct {
	ColumnName string `toml:"column_name"`
	FieldName  string `toml:"field_name"`
}

type UniqueKey struct {
	Name    string
	Columns []*Column
}

type Table struct {
	TableName             string         `toml:"table_name"`
	StructName            string         `toml:"struct_name"`
	PrimaryKeyColumnNames []string       `toml:"primary_key"`
	ColumnConfigs         []ColumnConfig `toml:"columns"`
	Columns               []Column
	PrimaryKeyColumns     []*Column
	CandidateKeys         []UniqueKey
}

func generateCmd(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Fprintln(os.Stderr, "generate does not take any arguments")
		os.Exit(1)
	}

	var c Config
	var err error
	c.Database, err = pgx.ParseEnvLibpq()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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

	templates := loadTemplates()

	supportData := initData{
		PkgName: c.Package,
		Version: VERSION,
	}

	supportFiles := []struct {
		path string
		tmpl *template.Template
	}{
		{"pgxdata_db.go", templates.Lookup("db")},
	}
	for _, f := range supportFiles {
		err := writeSupportFile(f.path, f.tmpl, supportData)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	for _, t := range c.Tables {
		file, err := os.Create("pgxdata_" + goCaseToFileCase(t.StructName) + ".go")
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
		PkgName           string
		TableName         string
		StructName        string
		Columns           []Column
		PrimaryKeyColumns []*Column
	}{
		PkgName:           pkgName,
		TableName:         table.TableName,
		StructName:        table.StructName,
		Columns:           table.Columns,
		PrimaryKeyColumns: table.PrimaryKeyColumns,
	})
}

func inspectDatabase(db Queryer, tables []Table) error {
	for i := range tables {
		schemaName, tableName,err  := splitSchemaFromTableName(tables[i].TableName)
		if err != nil {
			return err
		}

		rows, err := db.Query(`
			SELECT  column_name, data_type, ordinal_position
			FROM information_schema.columns
			WHERE table_schema=$1 AND table_name=$2`, schemaName, tableName)
		if err != nil {
			return err
		}

		var columns []Column
		for rows.Next() {
			var c Column
			rows.Scan(&c.ColumnName, &c.DataType, &c.OrdinalPosition)
			c.FieldName = pgCaseToGoPublicCase(c.ColumnName)
			c.GoBoxType = pgTypeToGoBoxType(c.DataType)
			c.VarName = pgCaseToGoPrivateCase(c.ColumnName)
			c.GoType = pgTypeToGoType(c.DataType)

			columns = append(columns, c)
		}

		if rows.Err() != nil {
			return err
		}

		tables[i].Columns = columns

		if len(tables[i].PrimaryKeyColumnNames) == 0 {
			tables[i].PrimaryKeyColumnNames = []string{"id"}
		}

		for _, columnName := range tables[i].PrimaryKeyColumnNames {
			var found bool
			for j := range tables[i].Columns {
				if tables[i].Columns[j].ColumnName == columnName {
					tables[i].PrimaryKeyColumns = append(tables[i].PrimaryKeyColumns, &tables[i].Columns[j])
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("table %s primary_key column %s not found", tables[i].TableName, columnName)
			}
		}

		for _, cc := range tables[i].ColumnConfigs {
			var found bool
			for j := range tables[i].Columns {
				if tables[i].Columns[j].ColumnName == cc.ColumnName {
					if cc.FieldName != "" {
						tables[i].Columns[j].FieldName = cc.FieldName
					}
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("table %s column %s not found", tables[i].TableName, cc.ColumnName)
			}
		}
	}

	return nil
}

const selectorInferenceQuery = /* language=PostgreSQL */ `
SELECT tc.constraint_type, tc.constraint_name, c.column_name
FROM
  information_schema.table_constraints tc
  INNER JOIN information_schema.constraint_column_usage ccu
    ON
      ccu.table_schema = tc.table_schema AND
      ccu.table_name = tc.table_name AND
      ccu.constraint_name = tc.constraint_name
  INNER JOIN information_schema.columns c
    ON
      c.table_schema = tc.table_schema AND
      c.table_name = tc.table_name AND
      c.column_name = ccu.column_name
WHERE "tc"."table_schema" = $1 AND
      tc.table_name = $2 AND
      tc.constraint_type IN ('PRIMARY KEY', 'UNIQUE')
ORDER BY tc.constraint_type, tc.constraint_name, c.ordinal_position ASC;
`

type candidateKeyType string

const (
	unique     candidateKeyType = "UNIQUE"
	primaryKey candidateKeyType = "PRIMARY KEY"
)

type candidateKeyConstraint struct {
	Name    string
	Columns []string
}

type keyConstraints struct {
	Primary *candidateKeyConstraint
	Unique  []*candidateKeyConstraint
}

func inferKeyConstraints(db Queryer, schema string, table string) (*keyConstraints, error) {
	rows, err := db.Query(selectorInferenceQuery, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keySet := &keyConstraints{}

	var ctype, cname, columnName string
	var currentUnique *candidateKeyConstraint

	for rows.Next() {
		rows.Scan(&ctype, &cname, &columnName)
		switch ctype {
		case primaryKey:
			if keySet.Primary == nil {
				keySet.Primary = &candidateKeyConstraint{cname, []string{columnName}}
			} else {
				keySet.Primary.Columns = append(keySet.Primary.Columns, columnName)
			}
		case unique:
			if currentUnique == nil || currentUnique.Name != cname {
				currentUnique = &candidateKeyConstraint{cname, []string{columnName}}
				keySet.Unique = append(keySet.Unique, currentUnique)
			} else {
				currentUnique.Columns = append(currentUnique.Columns, columnName)
			}
		}
	}
	return keySet, nil
}

// splitSchemaFromTableName splits off the name from the table name. if the name is not of the form "<schema>.<table>"
// it will return public. More than 2 periods in a table identifier are invalid.
func splitSchemaFromTableName(tableName string) (schema string, table string, err error) {
	schema = "public"
	table = tableName
	if parts := strings.Split(tableName, "."); len(parts) > 2 {
		err = fmt.Errorf("table name %s contain more than one period", tableName)
	} else if len(parts) == 2 {
		schema = parts[0]
		table = parts[1]
	}
	return
}

func pgCaseToGoPublicCase(pg string) string {
	parts := strings.Split(pg, "_")
	buf := &bytes.Buffer{}
	for _, s := range parts {
		if acronyms[s] {
			s = strings.ToUpper(s)
		} else {
			s = strings.Title(s)
		}
		buf.WriteString(s)
	}

	return buf.String()
}

func pgCaseToGoPrivateCase(pg string) string {
	parts := strings.Split(pg, "_")
	buf := &bytes.Buffer{}
	for i, s := range parts {
		if i == 0 {
			s = strings.ToLower(s)
		} else if acronyms[s] {
			s = strings.ToUpper(s)
		} else {
			s = strings.Title(s)
		}
		buf.WriteString(s)
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

func pgTypeToGoBoxType(pg string) string {
	if t, ok := pgToBoxTypeMap[pg]; ok {
		return t
	} else {
		return "String"
	}
}

func pgTypeToGoType(pg string) string {
	if t, ok := pgToGoTypeMap[pg]; ok {
		return t
	} else {
		return "string"
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
