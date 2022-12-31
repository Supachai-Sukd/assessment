package customer_expense

import (
	"bytes"
	"encoding/json"
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

func TestListExpenses(t *testing.T) {
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
	//h := handler{db}
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

func TestAddExpenses(t *testing.T) {
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

func TestUpdateExpenses(t *testing.T) {
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

func TestGetExpensesById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ce := CustomerExpenses{
		ID:     1,
		Title:  "Test Expense",
		Amount: 100,
		Note:   "Test note",
		Tags:   []string{"tag1", "tag2"},
	}

	tags := "{" + strings.Join(ce.Tags, ",") + "}"

	// Set up expected rows to be returned from the mock database
	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(ce.ID, ce.Title, ce.Amount, ce.Note, tags)

	// Set up the mock to expect a query and return the expected rows
	mock.ExpectQuery("SELECT id, title, amount, note, tags FROM expenses WHERE id = \\$1").
		WithArgs("1").
		WillReturnRows(rows)

	e := echo.New()
	c := e.NewContext(nil, nil)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the GetExpensesById function with the mock database and context
	if err := GetExpensesById(c); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Check the response from the function to make sure it matches the expected result
	if c.Response().Status != http.StatusOK {
		t.Errorf("unexpected status code: got %d, want %d", c.Response().Status, http.StatusOK)
	}
	var resp CustomerExpenses
	if err := c.JSON(http.StatusOK, &resp); err != nil {
		t.Errorf("unexpected error unmarshalling response: %s", err)
	}
	if resp.ID != ce.ID {
		t.Errorf("unexpected ID field in response: got %d, want %d", resp.ID, ce.ID)
	}

	if resp.Title != ce.Title {
		t.Errorf("unexpected Title field in response: got %s, want %s", resp.Title, ce.Title)
	}
}
