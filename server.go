package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/supachai-sukd/assessment/app/customer_expense"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	customer_expense.InitDB()
	//db := database.GetInstance()
	//
	//defer db.Close()
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	e.POST("/expenses", customer_expense.AddExpenses)
	e.GET("/expenses/:id", customer_expense.GetExpensesById)

	// os.Getenv("PORT") Use after refactor.
	go func() {
		if err := e.Start(":2565"); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("shutting down the server")
		}
	}()
	fmt.Printf("start at port: %v\n", 2565) // Port 2565

	// Create channel shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	<-shutdown

	// Shutdown in 10 second
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
