package customer_expense

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestAddExpense(t *testing.T) {

	body := bytes.NewBufferString(`{
		"title": "bank",
		"amount": 19,
		"note": "banana",
		"tags": ["food", "beverage"]
	}`)
	var c CustomerExpenses

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	assert.NotEqual(t, 0, c.ID)
	assert.Equal(t, "bank", c.Title)
	assert.Equal(t, "banana", c.Note)
}

func TestGetExpensesById(t *testing.T) {
	ce := seedExpensesInformation(t)

	var latest CustomerExpenses
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(ce.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, ce.ID, latest.ID)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)
}

func TestUpdateExpensesById(t *testing.T) {
	ce := seedExpensesInformation(t)

	body := bytes.NewBufferString(`{
		"title": "bank_kub",
		"amount": 99,
		"note": "banana grape",
		"tags": ["food", "fruit"]
	}`)

	var latest CustomerExpenses
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(ce.ID)), body)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, ce.ID, latest.ID)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)
}

func TestGetAllExpenses(t *testing.T) {
	_ = seedExpensesInformation(t)

	var latest []CustomerExpenses
	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	assert.NotEmpty(t, latest)
	assert.True(t, len(latest) > 0)              // check that the array has more than zero element
	assert.NotEmpty(t, latest[len(latest)-1].ID) // check that the last element has a non-empty
	assert.NotEmpty(t, latest[len(latest)-1].Title)
	assert.NotEmpty(t, latest[len(latest)-1].Amount)
	assert.NotEmpty(t, latest[len(latest)-1].Note)
	assert.NotEmpty(t, latest[len(latest)-1].Tags)

}

func seedExpensesInformation(t *testing.T) CustomerExpenses {
	var ce CustomerExpenses
	body := bytes.NewBufferString(`{
		"title": "bank",
		"amount": 19,
		"note": "banana",
		"tags": ["food", "beverage"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&ce)
	if err != nil {
		t.Fatal("can't add expenses:", err)
	}
	return ce
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
