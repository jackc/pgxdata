package data_test

import (
	"testing"

	"github.com/jackc/pgx-crud/test/data"
)

func TestCount(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	customerCount, err := data.CountCustomerRow(tx)
	if err != nil {
		t.Fatalf("CountCustomerRow unexpectedly failed: %v", err)
	}
	if customerCount != 0 {
		t.Fatalf("Expected CountCustomerRow to return %v, but is was %v", 0, customerCount)
	}

	err = data.InsertCustomerRow(tx, &data.CustomerRow{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
		BirthDate: data.Time{Status: data.Null},
	})
	if err != nil {
		t.Fatalf("InsertCustomerRow unexpectedly failed: %v", err)
	}

	customerCount, err = data.CountCustomerRow(tx)
	if err != nil {
		t.Fatalf("CountCustomerRow unexpectedly failed: %v", err)
	}
	if customerCount != 1 {
		t.Fatalf("Expected CountCustomerRow to return %v, but is was %v", 1, customerCount)
	}
}

func TestSelectAll(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	customers, err := data.SelectAllCustomerRow(tx)
	if err != nil {
		t.Fatalf("SelectAllCustomerRow unexpectedly failed: %v", err)
	}
	if len(customers) != 0 {
		t.Fatalf("Expected SelectAllCustomerRow to return %d rows, but is was %d", 0, len(customers))
	}

	insertedRow := data.CustomerRow{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
		BirthDate: data.Time{Status: data.Null},
	}

	err = data.InsertCustomerRow(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomerRow unexpectedly failed: %v", err)
	}

	customers, err = data.SelectAllCustomerRow(tx)
	if err != nil {
		t.Fatalf("SelectAllCustomerRow unexpectedly failed: %v", err)
	}
	if len(customers) != 1 {
		t.Fatalf("Expected SelectAllCustomerRow to return %d rows, but is was %d", 1, len(customers))
	}

	if customers[0].FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", insertedRow.FirstName, customers[0].FirstName)
	}
	if customers[0].LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", insertedRow.LastName, customers[0].LastName)
	}
}

func TestSelectByID(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	customer, err := data.SelectCustomerRowByID(tx, -1)
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectCustomerRowByID to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.CustomerRow{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
		BirthDate: data.Time{Status: data.Null},
	}

	err = data.InsertCustomerRow(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomerRow unexpectedly failed: %v", err)
	}

	customer, err = data.SelectCustomerRowByID(tx, insertedRow.ID.Value)
	if err != nil {
		t.Fatalf("SelectCustomerRowByID unexpectedly failed: %v", err)
	}

	if customer.FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", data.String{Value: "John", Status: data.Present}, customer.FirstName)
	}
	if customer.LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", data.String{Value: "Smith", Status: data.Present}, customer.FirstName)
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	insertedRow := data.CustomerRow{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
	}

	err := data.InsertCustomerRow(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomerRow unexpectedly failed: %v", err)
	}

	customer, err := data.SelectCustomerRowByID(tx, insertedRow.ID.Value)
	if err != nil {
		t.Fatalf("SelectCustomerRowByID unexpectedly failed: %v", err)
	}

	if customer.FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", data.String{Value: "John", Status: data.Present}, customer.FirstName)
	}
	if customer.LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", data.String{Value: "Smith", Status: data.Present}, customer.FirstName)
	}
}

func TestInsertOverridingPK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	insertedRow := data.CustomerRow{
		ID:        data.Int32{Value: -2, Status: data.Present},
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
	}

	err := data.InsertCustomerRow(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomerRow unexpectedly failed: %v", err)
	}

	customer, err := data.SelectCustomerRowByID(tx, -2)
	if err != nil {
		t.Fatalf("SelectCustomerRowByID unexpectedly failed: %v", err)
	}

	if customer.FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", data.String{Value: "John", Status: data.Present}, customer.FirstName)
	}
	if customer.LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", data.String{Value: "Smith", Status: data.Present}, customer.FirstName)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	insertedRow := data.CustomerRow{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
	}

	err := data.InsertCustomerRow(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomerRow unexpectedly failed: %v", err)
	}

	customer, err := data.SelectCustomerRowByID(tx, insertedRow.ID.Value)
	if err != nil {
		t.Fatalf("SelectCustomerRowByID unexpectedly failed: %v", err)
	}

	if customer.FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", data.String{Value: "John", Status: data.Present}, customer.FirstName)
	}
	if customer.LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", data.String{Value: "Smith", Status: data.Present}, customer.FirstName)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	insertedRow := data.CustomerRow{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
	}

	err := data.InsertCustomerRow(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomerRow unexpectedly failed: %v", err)
	}

	_, err = data.SelectCustomerRowByID(tx, insertedRow.ID.Value)
	if err != nil {
		t.Fatalf("SelectCustomerRowByID unexpectedly failed: %v", err)
	}

	err = data.DeleteCustomerRow(tx, insertedRow.ID.Value)
	if err != nil {
		t.Fatalf("DeleteCustomerRow unexpectedly failed: %v", err)
	}

	_, err = data.SelectCustomerRowByID(tx, insertedRow.ID.Value)
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectCustomerRowByID to return err data.ErrNotFound but it was: %v", err)
	}
}
