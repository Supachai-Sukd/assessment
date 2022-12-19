package repository

import (
	"github.com/supachai-sukd/assessment/app/model"
)

type CustomerExpenseRepository interface {
	Create(c *model.CreateCustomerExpense) error
	Update(ID int, cust *model.UpdateCustomerExpense) error
	GetAll() ([]*model.CustomerExpenses, error)
	GetByID(ID int) (*model.CreateCustomerExpense, error)
}
