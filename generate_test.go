package main

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
		config.Database = "pgx_crud"
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

func TestInspectDatabase(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	tests := []struct {
		input    []Table
		expected []Table
	}{
		{
			input: []Table{
				{
					TableName:  "customer",
					StructName: "CustomerRow",
				},
			},
			expected: []Table{
				{
					TableName:  "customer",
					StructName: "CustomerRow",
					Columns: []Column{
						{
							ColumnName:      "id",
							DataType:        "integer",
							OrdinalPosition: 1,
							FieldName:       "ID",
							GoType:          "Int32",
						},

						{
							ColumnName:      "first_name",
							DataType:        "character varying",
							OrdinalPosition: 2,
							FieldName:       "FirstName",
							GoType:          "String",
						},

						{
							ColumnName:      "last_name",
							DataType:        "character varying",
							OrdinalPosition: 3,
							FieldName:       "LastName",
							GoType:          "String",
						},

						{
							ColumnName:      "birth_date",
							DataType:        "date",
							OrdinalPosition: 4,
							FieldName:       "BirthDate",
							GoType:          "Time",
						},

						{
							ColumnName:      "creation_time",
							DataType:        "timestamp with time zone",
							OrdinalPosition: 5,
							FieldName:       "CreationTime",
							GoType:          "Time",
						},
					},
				},
			},
		},
	}

	for testIdx, tt := range tests {
		err := inspectDatabase(tx, tt.input)
		if err != nil {
			t.Errorf("%d. inspectDatabase failed: %v", testIdx, err)
			continue
		}

		if len(tt.expected) != len(tt.input) {
			t.Errorf("%d. expected %d tables, got %d", testIdx, len(tt.expected), len(tt.input))
			continue
		}

		for tableIdx := range tt.input {
			expectedTable := tt.expected[tableIdx]
			inputTable := tt.input[tableIdx]

			if expectedTable.TableName != inputTable.TableName {
				t.Errorf("%d:%d. expected TableName to be %s, got %s", testIdx, tableIdx, expectedTable.TableName, inputTable.TableName)
			}
			if expectedTable.StructName != inputTable.StructName {
				t.Errorf("%d:%d. expected StructName to be %s, got %s", testIdx, tableIdx, expectedTable.StructName, inputTable.StructName)
			}

			if len(expectedTable.Columns) != len(inputTable.Columns) {
				t.Errorf("%d:%d. expected %d tables, got %d", testIdx, tableIdx, len(expectedTable.Columns), len(inputTable.Columns))
				continue
			}

			for columnIdx := range inputTable.Columns {
				expectedColumn := expectedTable.Columns[columnIdx]
				inputColumn := inputTable.Columns[columnIdx]

				if expectedColumn.ColumnName != inputColumn.ColumnName {
					t.Errorf("%d:%d:%d. expected ColumnName %s, got %s", testIdx, tableIdx, columnIdx, expectedColumn.ColumnName, inputColumn.ColumnName)
				}
				if expectedColumn.DataType != inputColumn.DataType {
					t.Errorf("%d:%d:%d. expected DataType %s, got %s", testIdx, tableIdx, columnIdx, expectedColumn.DataType, inputColumn.DataType)
				}
				if expectedColumn.OrdinalPosition != inputColumn.OrdinalPosition {
					t.Errorf("%d:%d:%d. expected OrdinalPosition %d, got %d", testIdx, tableIdx, columnIdx, expectedColumn.OrdinalPosition, inputColumn.OrdinalPosition)
				}
				if expectedColumn.FieldName != inputColumn.FieldName {
					t.Errorf("%d:%d:%d. expected FieldName %s, got %s", testIdx, tableIdx, columnIdx, expectedColumn.FieldName, inputColumn.FieldName)
				}
				if expectedColumn.GoType != inputColumn.GoType {
					t.Errorf("%d:%d:%d. expected GoType %s, got %s", testIdx, tableIdx, columnIdx, expectedColumn.GoType, inputColumn.GoType)
				}
			}
		}
	}
}

func TestPgCaseToGoCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"first_name", "FirstName"},
		{"id", "ID"},
		{"person_id", "PersonID"},
		{"person_ideal", "PersonIdeal"},
	}

	for i, tt := range tests {
		actual := pgCaseToGoCase(tt.input)
		if actual != tt.expected {
			t.Errorf(`%d. Given "%s", expected "%s", but got "%s"`, i, tt.input, tt.expected, actual)
		}
	}
}
