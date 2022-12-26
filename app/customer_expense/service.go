package customer_expense

import (
	_ "database/sql"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func AddExpense(ce CustomerExpenses) (CustomerExpenses, error) {
	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4::text[])  RETURNING id", ce.Title, ce.Amount, ce.Note, pq.Array(ce.Tags))
	err := row.Scan(&ce.ID)
	if err != nil {
		return ce, err
	}

	return ce, nil
}
