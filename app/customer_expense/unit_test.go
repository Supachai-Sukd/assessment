//go:build unit

package customer_expense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/supachai-sukd/assessment/pkg/config"
	"net/http"
	_ "net/http"
	"net/http/httptest"
	_ "net/http/httptest"
	"strings"
	_ "strings"
	"testing"
	_ "testing"

	_ "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/labstack/echo/v4"
	_ "github.com/stretchr/testify/assert"
)

func TestListExpensesUnit(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow("1", "strawberry smoothie", 79, "night market promotion discount 10 bath", "{food, beverage}")

	_, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses").WillReturnRows(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	c := e.NewContext(req, rec)
	expected := "[{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}]"

	// Act
	err = GetAllExpenses(c)
	bodyArray := strings.Split(rec.Body.String(), ",")
	idZeroIdx := bodyArray[0]
	titleZeroIdx := bodyArray[1]
	amountZeroIdx := bodyArray[2]
	noteZeroIdx := bodyArray[3]
	tagsZeroIdx := bodyArray[4]
	tagsZeroIdx2 := bodyArray[5]

	var combineRecBody = idZeroIdx + "," + titleZeroIdx + "," + amountZeroIdx + "," + noteZeroIdx + "," + tagsZeroIdx + "," + tagsZeroIdx2 + "]"

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, combineRecBody)
	}
}

func TestAddExpensesUnit(t *testing.T) {
	// Set up the mock DB and defer the close to after the test
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expenseInsert := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO expenses").WillReturnRows(expenseInsert)

	config.DB = db

	requestBody, _ := json.Marshal(CustomerExpenses{
		Title:  "Test Expense",
		Amount: 100,
		Note:   "Test note",
		Tags:   []string{"test", "expense"},
	})
	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Set up the HTTP recorder to capture the response
	recorder := httptest.NewRecorder()
	c := echo.New().NewContext(req, recorder)

	// Call the function and check the response
	if assert.NoError(t, AddExpenses(c)) {
		assert.Equal(t, http.StatusCreated, recorder.Code)
		var resp CustomerExpenses
		json.Unmarshal(recorder.Body.Bytes(), &resp)
		assert.Equal(t, 1, resp.ID)
		assert.Equal(t, "Test Expense", resp.Title)
		assert.Equal(t, float64(100), resp.Amount)
		assert.Equal(t, "Test note", resp.Note)
		assert.Equal(t, []string{"test", "expense"}, resp.Tags)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateExpensesUnit(t *testing.T) {
	// Set up mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	config.DB = db

	// Define mock query and response
	mock.ExpectPrepare("UPDATE expenses").
		ExpectQuery().
		WithArgs("1", "test title", float64(100), "test note", pq.Array([]string{"tag1", "tag2"})).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Set up HTTP request and context
	req := httptest.NewRequest(http.MethodPut, "/expenses/1", strings.NewReader(`{"title": "test title", "amount": 100, "note": "test note", "tags": ["tag1", "tag2"]}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call UpdateExpenses and check response
	err = UpdateExpenses(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"id": 1, "title": "test title", "amount": 100, "note": "test note", "tags": ["tag1", "tag2"]}`, rec.Body.String())

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetExpensesByIdUnit(t *testing.T) {
	// Set up mock data
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	expenseID := "abc123"
	mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses WHERE id = \\$1").
		WithArgs(expenseID).
		WillReturnError(fmt.Errorf("sql: database is closed"))

	// Set up mock context
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses/"+expenseID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues(expenseID)

	// Call GetExpensesById
	if err := GetExpensesById(c); err != nil {
		t.Errorf("GetExpensesById returned unexpected error: %v", err)
	}

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	expectedBody := `{"message":"can't scan expenses information:sql: database is closed"}`
	assert.JSONEq(t, expectedBody, rec.Body.String())
}
