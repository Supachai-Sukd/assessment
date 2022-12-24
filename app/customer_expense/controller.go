package customer_expense

import (
	"database/sql"
	_ "database/sql"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
)

var db *sql.DB

type Err struct {
	Message string `json:"message"`
}

func AddExpenses(c echo.Context) error {
	ce := CustomerExpenses{}

	err := c.Bind(&ce)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	//var sliceString = strings.Join(ce.Tags, ",")
	//var sliceString = "'{" + strings.Join(ce.Tags, ",") + "}'"

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4::text[])  RETURNING id", ce.Title, ce.Amount, ce.Note, pq.Array(ce.Tags))
	err = row.Scan(&ce.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ce)
}
