package customer_expense

import (
	"database/sql"
	_ "database/sql"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strings"
)

// รอสร้าง NewExpenses ที่ pointer ไปหา Struct ที่ใส้ในมี db *sql.db

func AddExpenseService(ce CustomerExpenses) (CustomerExpenses, error) {
	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4::text[])  RETURNING id", ce.Title, ce.Amount, ce.Note, pq.Array(ce.Tags))
	err := row.Scan(&ce.ID)
	if err != nil {
		return ce, err
	}

	return ce, nil
}

func GetExpensesByIdService(db *sql.DB, id string) (CustomerExpenses, error) {
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return CustomerExpenses{}, err
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

	return ce, err
}
