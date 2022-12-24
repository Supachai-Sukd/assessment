package customer_expense

import (
	"database/sql"
	_ "database/sql"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
	"strings"
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

func GetExpensesById(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses information statment:" + err.Error()})
	}

	ce := CustomerExpenses{}
	row := stmt.QueryRow(id)

	var tags sql.NullString
	err = row.Scan(&ce.ID, &ce.Title, &ce.Amount, &ce.Note, &tags)
	if tags.Valid {
		ce.Tags = strings.Split(tags.String, ",")
	}

	for i, tag := range ce.Tags {
		ce.Tags[i] = strings.Trim(tag, "{}")
	}

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expenses information not found"})
	case nil:
		return c.JSON(http.StatusOK, ce)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses information:" + err.Error()})
	}

}

func UpdateExpenses(c echo.Context) error {
	id := c.Param("id")

	ce := CustomerExpenses{}
	if err := c.Bind(&ce); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request body"})
	}

	stmt, err := db.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5::text[] WHERE id=$1 RETURNING id")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses information statement: " + err.Error()})
	}
	row := stmt.QueryRow(id, ce.Title, ce.Amount, ce.Note, pq.Array(ce.Tags))

	err = row.Scan(&ce.ID)
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expenses information not found"})
	case nil:
		return c.JSON(http.StatusOK, ce)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses information: " + err.Error()})
	}
}
