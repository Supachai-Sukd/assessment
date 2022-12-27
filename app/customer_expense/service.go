package customer_expense

import (
	"database/sql"
	_ "database/sql"
	"fmt"
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

func UpdateExpensesService(id string, ce CustomerExpenses) (CustomerExpenses, error) {
	stmt, err := db.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5::text[] WHERE id=$1 RETURNING id")
	if err != nil {
		return CustomerExpenses{}, fmt.Errorf("can't prepare query expenses information statement: %v", err)
	}
	row := stmt.QueryRow(id, ce.Title, ce.Amount, ce.Note, pq.Array(ce.Tags))

	err = row.Scan(&ce.ID)
	switch err {
	case sql.ErrNoRows:
		return CustomerExpenses{}, fmt.Errorf("expenses information not found")
	case nil:
		return ce, nil
	default:
		return CustomerExpenses{}, fmt.Errorf("can't scan expenses information: %v", err)
	}
}

func GetAllExpensesService(db *sql.DB) ([]CustomerExpenses, error) {
	stmt, err := db.Prepare("SELECT * FROM expenses")
	if err != nil {
		return nil, err
	}

	rows, errs := stmt.Query()
	if errs != nil {
		return nil, errs
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
			return nil, err
		}
		ce = append(ce, cust)
	}

	return ce, nil
}
