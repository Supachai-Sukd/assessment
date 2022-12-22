package customer_expense

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddExpenses(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, echo.Map{
		"msg": "abc",
	})
}
