package customer_expense

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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

func TestListNews(t *testing.T) {
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
