package customer_expense

import (
	"database/sql"
	//"database/sql"
	//_ "database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type Err struct {
	Message string `json:"message"`
}

func AddExpenses(c echo.Context) error {

	ce := CustomerExpenses{}

	err := c.Bind(&ce)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	resp, errs := AddExpenseService(ce)
	if errs != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func GetExpensesById(c echo.Context) error {
	id := c.Param("id")
	ce, err := GetExpensesByIdService(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, Err{Message: "expenses information not found"})
		}
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses information:" + err.Error()})
	}
	return c.JSON(http.StatusOK, ce)
}

func UpdateExpenses(c echo.Context) error {
	id := c.Param("id")

	ce := CustomerExpenses{}
	if err := c.Bind(&ce); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request body"})
	}

	updatedCE, err := UpdateExpensesService(id, ce)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, updatedCE)
}

func GetAllExpenses(c echo.Context) error {
	stmt, err := db.Prepare("SELECT * FROM expenses ")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query all expenses statment:" + err.Error()})
	}

	rows, errs := stmt.Query()
	if errs != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query all expenses:" + err.Error()})
	}
	ce := []CustomerExpenses{}

	for rows.Next() {
		cust := CustomerExpenses{}
		var tags sql.NullString
		err := rows.Scan(&cust.ID, &cust.Title, &cust.Amount, &cust.Note, &tags)
		if tags.Valid {
			cust.Tags = strings.Split(tags.String, ",")
		}

		for i, tag := range cust.Tags {
			cust.Tags[i] = strings.Trim(tag, "{}")
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses:" + err.Error()})
		}
		ce = append(ce, cust)
	}

	return c.JSON(http.StatusOK, ce)
}
