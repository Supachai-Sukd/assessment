package customer_expense

type CustomerExpenseRepository interface {
	Create(c *CreateCustomerExpense) error
	Update(ID int, cust *UpdateCustomerExpense) error
	GetAll() ([]*CustomerExpenses, error)
	GetByID(ID int) (*CreateCustomerExpense, error)
}
