package customer_expense

type CustomerExpenses struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
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
