package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type CustomerExpenses struct {
	Title  string   `db:"picture" json:"picture"`
	Amount int      `db:"amount" json:"amount"`
	Note   string   `db:"note" json:"note"`
	Tags   []string `db:"tags" json:"tags"`
}

func (c CustomerExpenses) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *CustomerExpenses) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &c)
}

func NewCusExpenses() *CustomerExpenses {
	return &CustomerExpenses{}
}
