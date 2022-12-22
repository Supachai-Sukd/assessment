package customer_expense

type CustomerExpenses struct {
	Title  string   `db:"picture" json:"picture"`
	Amount float64  `db:"amount" json:"amount"`
	Note   string   `db:"note" json:"note"`
	Tags   []string `db:"tags" json:"tags"`
}

func NewCusExpenses() *CustomerExpenses {
	return &CustomerExpenses{}
}

type CreateCustomerExpense struct {
	Title  string   `json:"picture"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type UpdateCustomerExpense struct {
	Title  string   `json:"picture"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}
