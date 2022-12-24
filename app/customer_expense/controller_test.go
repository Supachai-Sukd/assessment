package customer_expense

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestExpense(t *testing.T) {

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
