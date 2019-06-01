package main

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

func TestInspectDatabase(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

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
				{
					TableName:  "widget",
					StructName: "WidgetRow",
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
							GoBoxType:       "pgtype.Int4",
						},
						{
							ColumnName:      "first_name",
							DataType:        "character varying",
							OrdinalPosition: 2,
							FieldName:       "FirstName",
							GoBoxType:       "pgtype.Varchar",
						},
						{
							ColumnName:      "last_name",
							DataType:        "character varying",
							OrdinalPosition: 3,
							FieldName:       "LastName",
							GoBoxType:       "pgtype.Varchar",
						},
						{
							ColumnName:      "birth_date",
							DataType:        "date",
							OrdinalPosition: 4,
							FieldName:       "BirthDate",
							GoBoxType:       "pgtype.Date",
						},
						{
							ColumnName:      "creation_time",
							DataType:        "timestamp with time zone",
							OrdinalPosition: 5,
							FieldName:       "CreationTime",
							GoBoxType:       "pgtype.Timestamptz",
						},
					},
				},
				{
					TableName:  "widget",
					StructName: "WidgetRow",
					Columns: []Column{
						{
							ColumnName:      "id",
							DataType:        "bigint",
							OrdinalPosition: 1,
							FieldName:       "ID",
							GoBoxType:       "pgtype.Int8",
						},
						{
							ColumnName:      "name",
							DataType:        "character varying",
							OrdinalPosition: 2,
							FieldName:       "Name",
							GoBoxType:       "pgtype.Varchar",
						},
						{
							ColumnName:      "weight",
							DataType:        "smallint",
							OrdinalPosition: 3,
							FieldName:       "Weight",
							GoBoxType:       "pgtype.Int2",
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
				t.Errorf("%d:%d. expected %d columns, got %d", testIdx, tableIdx, len(expectedTable.Columns), len(inputTable.Columns))
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
				if expectedColumn.GoBoxType != inputColumn.GoBoxType {
					t.Errorf("%d:%d:%d. expected GoBoxType %s, got %s", testIdx, tableIdx, columnIdx, expectedColumn.GoBoxType, inputColumn.GoBoxType)
				}
			}
		}
	}
}

func TestPgCaseToGoPublicCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"first_name", "FirstName"},
		{"id", "ID"},
		{"person_id", "PersonID"},
		{"person_ideal", "PersonIdeal"},
		{"ip", "IP"},
		{"ip_source", "IPSource"},
		{"url", "URL"},
		{"url_base", "URLBase"},
		{"curl", "Curl"},
	}

	for i, tt := range tests {
		actual := pgCaseToGoPublicCase(tt.input)
		if actual != tt.expected {
			t.Errorf(`%d. Given "%s", expected "%s", but got "%s"`, i, tt.input, tt.expected, actual)
		}
	}
}

func TestPgCaseToGoPrivateCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"first_name", "firstName"},
		{"id", "id"},
		{"person_id", "personID"},
		{"person_ideal", "personIdeal"},
		{"url", "url"},
		{"url_base", "urlBase"},
		{"curl", "curl"},
	}

	for i, tt := range tests {
		actual := pgCaseToGoPrivateCase(tt.input)
		if actual != tt.expected {
			t.Errorf(`%d. Given "%s", expected "%s", but got "%s"`, i, tt.input, tt.expected, actual)
		}
	}
}

func TestGoCaseToFileCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"Widget", "widget"},
		{"WidgetRow", "widget_row"},
	}

	for i, tt := range tests {
		actual := goCaseToFileCase(tt.input)
		if actual != tt.expected {
			t.Errorf(`%d. Given "%s", expected "%s", but got "%s"`, i, tt.input, tt.expected, actual)
		}
	}
}
