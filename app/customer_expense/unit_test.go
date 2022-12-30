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
	// Set up the mock database and a mock row with the expected ID
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := "123"
	mockRow := sqlmock.NewRows([]string{"id"}).AddRow(id)

	// Set up the mock preparer to return the mock row
	prep := mock.ExpectPrepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5::text[] WHERE id=$1 RETURNING id")
	prep.ExpectQuery().WithArgs(id, "title", 100, "note", pq.Array([]string{"tag1", "tag2"})).WillReturnRows(mockRow)

	// Set up the HTTP request and response recorder
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses/"+id, strings.NewReader(`{"title":"title","amount":100,"note":"note","tags":["tag1","tag2"]}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues(id)

	// Set the DB variable to the mock DB
	config.DB = db

	// Call the UpdateExpenses function
	if assert.NoError(t, UpdateExpenses(c)) {
		// Assert that the mock preparer and query were called as expected
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		// Assert that the response has the expected status code and body
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"id":"123","title":"title","amount":100,"note":"note","tags":["tag1","tag2"]}`, rec.Body.String())
	}
}
