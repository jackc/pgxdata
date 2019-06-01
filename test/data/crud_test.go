package data_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgxdata/test/data"
)

func TestCount(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	customerCount, err := data.CountCustomer(context.Background(), tx)
	if err != nil {
		t.Fatalf("CountCustomer unexpectedly failed: %v", err)
	}
	if customerCount != 0 {
		t.Fatalf("Expected CountCustomer to return %v, but is was %v", 0, customerCount)
	}

	err = data.InsertCustomer(context.Background(), tx, &data.Customer{
		FirstName: pgtype.Varchar{String: "John", Status: pgtype.Present},
		LastName:  pgtype.Varchar{String: "Smith", Status: pgtype.Present},
		BirthDate: pgtype.Date{Status: pgtype.Null},
	})
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customerCount, err = data.CountCustomer(context.Background(), tx)
	if err != nil {
		t.Fatalf("CountCustomer unexpectedly failed: %v", err)
	}
	if customerCount != 1 {
		t.Fatalf("Expected CountCustomer to return %v, but is was %v", 1, customerCount)
	}
}

func TestSelectAll(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	customers, err := data.SelectAllCustomer(context.Background(), tx)
	if err != nil {
		t.Fatalf("SelectAllCustomer unexpectedly failed: %v", err)
	}
	if len(customers) != 0 {
		t.Fatalf("Expected SelectAllCustomer to return %d rows, but is was %d", 0, len(customers))
	}

	insertedRow := data.Customer{
		FirstName: pgtype.Varchar{String: "John", Status: pgtype.Present},
		LastName:  pgtype.Varchar{String: "Smith", Status: pgtype.Present},
		BirthDate: pgtype.Date{Status: pgtype.Null},
	}

	err = data.InsertCustomer(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customers, err = data.SelectAllCustomer(context.Background(), tx)
	if err != nil {
		t.Fatalf("SelectAllCustomer unexpectedly failed: %v", err)
	}
	if len(customers) != 1 {
		t.Fatalf("Expected SelectAllCustomer to return %d rows, but is was %d", 1, len(customers))
	}

	if customers[0].FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", insertedRow.FirstName, customers[0].FirstName)
	}
	if customers[0].LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", insertedRow.LastName, customers[0].LastName)
	}
}

func TestSelectByPK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	customer, err := data.SelectCustomerByPK(context.Background(), tx, -1)
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectCustomerByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Customer{
		FirstName: pgtype.Varchar{String: "John", Status: pgtype.Present},
		LastName:  pgtype.Varchar{String: "Smith", Status: pgtype.Present},
		BirthDate: pgtype.Date{Status: pgtype.Null},
	}

	err = data.InsertCustomer(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customer, err = data.SelectCustomerByPK(context.Background(), tx, insertedRow.ID.Int)
	if err != nil {
		t.Fatalf("SelectCustomerByPK unexpectedly failed: %v", err)
	}

	if customer.FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", pgtype.Varchar{String: "John", Status: pgtype.Present}, customer.FirstName)
	}
	if customer.LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", pgtype.Varchar{String: "Smith", Status: pgtype.Present}, customer.LastName)
	}
}

func TestSelectByPKWithInt64PK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	widget, err := data.SelectWidgetByPK(context.Background(), tx, -1)
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectWidgetByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Widget{
		Name:   pgtype.Varchar{String: "Foozle", Status: pgtype.Present},
		Weight: pgtype.Int2{Int: 20, Status: pgtype.Present},
	}

	err = data.InsertWidget(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertWidget unexpectedly failed: %v", err)
	}

	widget, err = data.SelectWidgetByPK(context.Background(), tx, insertedRow.ID.Int)
	if err != nil {
		t.Fatalf("SelectWidgetByPK unexpectedly failed: %v", err)
	}

	if widget.Name != insertedRow.Name {
		t.Errorf("Expected Name to be %v, but it was %v", insertedRow.Name, widget.Name)
	}
	if widget.Weight != insertedRow.Weight {
		t.Errorf("Expected Weight to be %v, but it was %v", insertedRow.Weight, widget.Weight)
	}
}

func TestSelectByPKWithVarcharNotNamedIDAsPK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	part, err := data.SelectPartByPK(context.Background(), tx, "E100")
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectPartByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Part{
		Code:        pgtype.Varchar{String: "E100", Status: pgtype.Present},
		Description: pgtype.Text{String: "Engine 100", Status: pgtype.Present},
	}

	err = data.InsertPart(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertPart unexpectedly failed: %v", err)
	}

	part, err = data.SelectPartByPK(context.Background(), tx, insertedRow.Code.String)
	if err != nil {
		t.Fatalf("SelectPartByPK unexpectedly failed: %v", err)
	}

	if part.Code != insertedRow.Code {
		t.Errorf("Expected Code to be %v, but it was %v", insertedRow.Code, part.Code)
	}
	if part.Description != insertedRow.Description {
		t.Errorf("Expected Description to be %v, but it was %v", insertedRow.Description, part.Description)
	}
}

func TestSelectByPKWithCompositePK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	semester, err := data.SelectSemesterByPK(context.Background(), tx, 1999, "Fall")
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectSemesterByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Semester{
		Year:        pgtype.Int2{Int: 1999, Status: pgtype.Present},
		Season:      pgtype.Varchar{String: "Fall", Status: pgtype.Present},
		Description: pgtype.Text{String: "Last of the century", Status: pgtype.Present},
	}

	err = data.InsertSemester(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertSemeseter unexpectedly failed: %v", err)
	}

	semester, err = data.SelectSemesterByPK(context.Background(), tx, insertedRow.Year.Int, insertedRow.Season.String)
	if err != nil {
		t.Fatalf("SelectSemesterByPK unexpectedly failed: %v", err)
	}

	if semester.Year != insertedRow.Year {
		t.Errorf("Expected Year to be %v, but it was %v", insertedRow.Year, semester.Year)
	}
	if semester.Season != insertedRow.Season {
		t.Errorf("Expected Season to be %v, but it was %v", insertedRow.Season, semester.Season)
	}
	if semester.Description != insertedRow.Description {
		t.Errorf("Expected Description to be %v, but it was %v", insertedRow.Description, semester.Description)
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	insertedRow := data.Customer{
		FirstName: pgtype.Varchar{String: "John", Status: pgtype.Present},
		LastName:  pgtype.Varchar{String: "Smith", Status: pgtype.Present},
	}

	err := data.InsertCustomer(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customer, err := data.SelectCustomerByPK(context.Background(), tx, insertedRow.ID.Int)
	if err != nil {
		t.Fatalf("SelectCustomerByPK unexpectedly failed: %v", err)
	}

	if customer.FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", pgtype.Varchar{String: "John", Status: pgtype.Present}, customer.FirstName)
	}
	if customer.LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", pgtype.Varchar{String: "Smith", Status: pgtype.Present}, customer.LastName)
	}
}

func TestInsertOverridingPK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	insertedRow := data.Customer{
		ID:        pgtype.Int4{Int: -2, Status: pgtype.Present},
		FirstName: pgtype.Varchar{String: "John", Status: pgtype.Present},
		LastName:  pgtype.Varchar{String: "Smith", Status: pgtype.Present},
	}

	err := data.InsertCustomer(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customer, err := data.SelectCustomerByPK(context.Background(), tx, -2)
	if err != nil {
		t.Fatalf("SelectCustomerByPK unexpectedly failed: %v", err)
	}

	if customer.FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", pgtype.Varchar{String: "John", Status: pgtype.Present}, customer.FirstName)
	}
	if customer.LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", pgtype.Varchar{String: "Smith", Status: pgtype.Present}, customer.LastName)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	insertedRow := data.Customer{
		FirstName: pgtype.Varchar{String: "John", Status: pgtype.Present},
		LastName:  pgtype.Varchar{String: "Smith", Status: pgtype.Present},
	}

	err := data.InsertCustomer(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customer, err := data.SelectCustomerByPK(context.Background(), tx, insertedRow.ID.Int)
	if err != nil {
		t.Fatalf("SelectCustomerByPK unexpectedly failed: %v", err)
	}

	if customer.FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", pgtype.Varchar{String: "John", Status: pgtype.Present}, customer.FirstName)
	}
	if customer.LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", pgtype.Varchar{String: "Smith", Status: pgtype.Present}, customer.FirstName)
	}
}

func TestUpdateWithCompositePK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	semester, err := data.SelectSemesterByPK(context.Background(), tx, 1999, "Fall")
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectSemesterByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Semester{
		Year:        pgtype.Int2{Int: 1999, Status: pgtype.Present},
		Season:      pgtype.Varchar{String: "Fall", Status: pgtype.Present},
		Description: pgtype.Text{String: "Last of the century", Status: pgtype.Present},
	}

	err = data.InsertSemester(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertSemeseter unexpectedly failed: %v", err)
	}

	updateAttrs := &data.Semester{
		Description: pgtype.Text{String: "New value", Status: pgtype.Present},
	}

	data.UpdateSemester(context.Background(), tx,
		insertedRow.Year.Int,
		insertedRow.Season.String,
		updateAttrs,
	)

	semester, err = data.SelectSemesterByPK(context.Background(), tx, insertedRow.Year.Int, insertedRow.Season.String)
	if err != nil {
		t.Fatalf("SelectSemesterByPK unexpectedly failed: %v", err)
	}

	if semester.Description != updateAttrs.Description {
		t.Errorf("Expected Description to be %v, but it was %v", updateAttrs.Description, semester.Description)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	insertedRow := data.Customer{
		FirstName: pgtype.Varchar{String: "John", Status: pgtype.Present},
		LastName:  pgtype.Varchar{String: "Smith", Status: pgtype.Present},
	}

	err := data.InsertCustomer(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	_, err = data.SelectCustomerByPK(context.Background(), tx, insertedRow.ID.Int)
	if err != nil {
		t.Fatalf("SelectCustomerByPK unexpectedly failed: %v", err)
	}

	err = data.DeleteCustomer(context.Background(), tx, insertedRow.ID.Int)
	if err != nil {
		t.Fatalf("DeleteCustomer unexpectedly failed: %v", err)
	}

	_, err = data.SelectCustomerByPK(context.Background(), tx, insertedRow.ID.Int)
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectCustomerByPK to return err data.ErrNotFound but it was: %v", err)
	}
}

func TestDeleteWithCompositePK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	_, err := data.SelectSemesterByPK(context.Background(), tx, 1999, "Fall")
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectSemesterByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Semester{
		Year:        pgtype.Int2{Int: 1999, Status: pgtype.Present},
		Season:      pgtype.Varchar{String: "Fall", Status: pgtype.Present},
		Description: pgtype.Text{String: "Last of the century", Status: pgtype.Present},
	}

	err = data.InsertSemester(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertSemeseter unexpectedly failed: %v", err)
	}

	_, err = data.SelectSemesterByPK(context.Background(), tx, insertedRow.Year.Int, insertedRow.Season.String)
	if err != nil {
		t.Fatalf("SelectSemesterByPK unexpectedly failed: %v", err)
	}

	data.DeleteSemester(context.Background(),
		tx,
		insertedRow.Year.Int,
		insertedRow.Season.String,
	)

	_, err = data.SelectSemesterByPK(context.Background(), tx, insertedRow.Year.Int, insertedRow.Season.String)
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectSemesterByPK to return err data.ErrNotFound but it was: %v", err)
	}
}

func TestMappingOfRenamedField(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	insertedRow := data.RenamedFieldCustomer{
		FName:    pgtype.Varchar{String: "John", Status: pgtype.Present},
		LastName: pgtype.Varchar{String: "Smith", Status: pgtype.Present},
	}

	err := data.InsertRenamedFieldCustomer(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertRenamedFieldCustomer unexpectedly failed: %v", err)
	}

	customer, err := data.SelectRenamedFieldCustomerByPK(context.Background(), tx, insertedRow.ID.Int)
	if err != nil {
		t.Fatalf("SelectRenamedFieldCustomerByPK unexpectedly failed: %v", err)
	}

	if customer.FName != insertedRow.FName {
		t.Errorf("Expected FName to be %v, but it was %v", pgtype.Varchar{String: "John", Status: pgtype.Present}, customer.FName)
	}
}

func TestByteaByteSliceMapping(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback(context.Background())

	insertedRow := data.Blob{
		Payload: pgtype.Bytea{Bytes: []byte("Hello"), Status: pgtype.Present},
	}

	err := data.InsertBlob(context.Background(), tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertBlob unexpectedly failed: %v", err)
	}

	blob, err := data.SelectBlobByPK(context.Background(), tx, insertedRow.ID.Int)
	if err != nil {
		t.Fatalf("SelectBlobByPK unexpectedly failed: %v", err)
	}

	if bytes.Compare(blob.Payload.Bytes, insertedRow.Payload.Bytes) != 0 {
		t.Errorf("Expected Payload to be %v, but it was %v", insertedRow.Payload, blob.Payload)
	}
}
