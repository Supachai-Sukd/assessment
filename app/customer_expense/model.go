package customer_expense

type CustomerExpenses struct {
	ID     int      `json:"id"`
	Title  string   `db:"title" json:"title"`
	Amount float64  `db:"amount" json:"amount"`
	Note   string   `db:"note" json:"note"`
	Tags   []string `db:"tags" json:"tags"`
}

func NewCusExpenses() *CustomerExpenses {
	return &CustomerExpenses{}
}

//type CreateCustomerExpense struct {
//	Title  string   `json:"title"`
//	Amount float64  `json:"amount"`
//	Note   string   `json:"note"`
//	Tags   []string `json:"tags"`
//}
//
//type UpdateCustomerExpense struct {
//	Title  string   `json:"title"`
//	Amount float64  `json:"amount"`
//	Note   string   `json:"note"`
//	Tags   []string `json:"tags"`
//}
