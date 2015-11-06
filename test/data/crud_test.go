package data_test

import (
	"testing"

	"github.com/jackc/pgxdata/test/data"
)

func TestCount(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	customerCount, err := data.CountCustomer(tx)
	if err != nil {
		t.Fatalf("CountCustomer unexpectedly failed: %v", err)
	}
	if customerCount != 0 {
		t.Fatalf("Expected CountCustomer to return %v, but is was %v", 0, customerCount)
	}

	err = data.InsertCustomer(tx, &data.Customer{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
		BirthDate: data.Time{Status: data.Null},
	})
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customerCount, err = data.CountCustomer(tx)
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
	defer tx.Rollback()

	customers, err := data.SelectAllCustomer(tx)
	if err != nil {
		t.Fatalf("SelectAllCustomer unexpectedly failed: %v", err)
	}
	if len(customers) != 0 {
		t.Fatalf("Expected SelectAllCustomer to return %d rows, but is was %d", 0, len(customers))
	}

	insertedRow := data.Customer{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
		BirthDate: data.Time{Status: data.Null},
	}

	err = data.InsertCustomer(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customers, err = data.SelectAllCustomer(tx)
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
	defer tx.Rollback()

	customer, err := data.SelectCustomerByPK(tx, -1)
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectCustomerByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Customer{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
		BirthDate: data.Time{Status: data.Null},
	}

	err = data.InsertCustomer(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customer, err = data.SelectCustomerByPK(tx, insertedRow.ID.Value)
	if err != nil {
		t.Fatalf("SelectCustomerByPK unexpectedly failed: %v", err)
	}

	if customer.FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", data.String{Value: "John", Status: data.Present}, customer.FirstName)
	}
	if customer.LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", data.String{Value: "Smith", Status: data.Present}, customer.FirstName)
	}
}

func TestSelectByPKWithInt64PK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	widget, err := data.SelectWidgetByPK(tx, -1)
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectWidgetByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Widget{
		Name:   data.String{Value: "Foozle", Status: data.Present},
		Weight: data.Int16{Value: 20, Status: data.Present},
	}

	err = data.InsertWidget(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertWidget unexpectedly failed: %v", err)
	}

	widget, err = data.SelectWidgetByPK(tx, insertedRow.ID.Value)
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
	defer tx.Rollback()

	part, err := data.SelectPartByPK(tx, "E100")
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectPartByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Part{
		Code:        data.String{Value: "E100", Status: data.Present},
		Description: data.String{Value: "Engine 100", Status: data.Present},
	}

	err = data.InsertPart(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertPart unexpectedly failed: %v", err)
	}

	part, err = data.SelectPartByPK(tx, insertedRow.Code.Value)
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
	defer tx.Rollback()

	semester, err := data.SelectSemesterByPK(tx, 1999, "Fall")
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectSemesterByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Semester{
		Year:        data.Int16{Value: 1999, Status: data.Present},
		Season:      data.String{Value: "Fall", Status: data.Present},
		Description: data.String{Value: "Last of the century", Status: data.Present},
	}

	err = data.InsertSemester(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertSemeseter unexpectedly failed: %v", err)
	}

	semester, err = data.SelectSemesterByPK(tx, insertedRow.Year.Value, insertedRow.Season.Value)
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
	defer tx.Rollback()

	insertedRow := data.Customer{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
	}

	err := data.InsertCustomer(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customer, err := data.SelectCustomerByPK(tx, insertedRow.ID.Value)
	if err != nil {
		t.Fatalf("SelectCustomerByPK unexpectedly failed: %v", err)
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

	insertedRow := data.Customer{
		ID:        data.Int32{Value: -2, Status: data.Present},
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
	}

	err := data.InsertCustomer(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customer, err := data.SelectCustomerByPK(tx, -2)
	if err != nil {
		t.Fatalf("SelectCustomerByPK unexpectedly failed: %v", err)
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

	insertedRow := data.Customer{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
	}

	err := data.InsertCustomer(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	customer, err := data.SelectCustomerByPK(tx, insertedRow.ID.Value)
	if err != nil {
		t.Fatalf("SelectCustomerByPK unexpectedly failed: %v", err)
	}

	if customer.FirstName != insertedRow.FirstName {
		t.Errorf("Expected FirstName to be %v, but it was %v", data.String{Value: "John", Status: data.Present}, customer.FirstName)
	}
	if customer.LastName != insertedRow.LastName {
		t.Errorf("Expected LastName to be %v, but it was %v", data.String{Value: "Smith", Status: data.Present}, customer.FirstName)
	}
}

func TestUpdateWithCompositePK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	semester, err := data.SelectSemesterByPK(tx, 1999, "Fall")
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectSemesterByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Semester{
		Year:        data.Int16{Value: 1999, Status: data.Present},
		Season:      data.String{Value: "Fall", Status: data.Present},
		Description: data.String{Value: "Last of the century", Status: data.Present},
	}

	err = data.InsertSemester(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertSemeseter unexpectedly failed: %v", err)
	}

	updateAttrs := &data.Semester{
		Description: data.String{Value: "New value", Status: data.Present},
	}

	data.UpdateSemester(tx,
		insertedRow.Year.Value,
		insertedRow.Season.Value,
		updateAttrs,
	)

	semester, err = data.SelectSemesterByPK(tx, insertedRow.Year.Value, insertedRow.Season.Value)
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
	defer tx.Rollback()

	insertedRow := data.Customer{
		FirstName: data.String{Value: "John", Status: data.Present},
		LastName:  data.String{Value: "Smith", Status: data.Present},
	}

	err := data.InsertCustomer(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertCustomer unexpectedly failed: %v", err)
	}

	_, err = data.SelectCustomerByPK(tx, insertedRow.ID.Value)
	if err != nil {
		t.Fatalf("SelectCustomerByPK unexpectedly failed: %v", err)
	}

	err = data.DeleteCustomer(tx, insertedRow.ID.Value)
	if err != nil {
		t.Fatalf("DeleteCustomer unexpectedly failed: %v", err)
	}

	_, err = data.SelectCustomerByPK(tx, insertedRow.ID.Value)
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectCustomerByPK to return err data.ErrNotFound but it was: %v", err)
	}
}

func TestDeleteWithCompositePK(t *testing.T) {
	t.Parallel()

	tx := begin(t)
	defer tx.Rollback()

	_, err := data.SelectSemesterByPK(tx, 1999, "Fall")
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectSemesterByPK to return err data.ErrNotFound but it was: %v", err)
	}

	insertedRow := data.Semester{
		Year:        data.Int16{Value: 1999, Status: data.Present},
		Season:      data.String{Value: "Fall", Status: data.Present},
		Description: data.String{Value: "Last of the century", Status: data.Present},
	}

	err = data.InsertSemester(tx, &insertedRow)
	if err != nil {
		t.Fatalf("InsertSemeseter unexpectedly failed: %v", err)
	}

	_, err = data.SelectSemesterByPK(tx, insertedRow.Year.Value, insertedRow.Season.Value)
	if err != nil {
		t.Fatalf("SelectSemesterByPK unexpectedly failed: %v", err)
	}

	data.DeleteSemester(tx,
		insertedRow.Year.Value,
		insertedRow.Season.Value,
	)

	_, err = data.SelectSemesterByPK(tx, insertedRow.Year.Value, insertedRow.Season.Value)
	if err != data.ErrNotFound {
		t.Fatalf("Expected SelectSemesterByPK to return err data.ErrNotFound but it was: %v", err)
	}
}
