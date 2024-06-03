package organizze

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BaseURL = "https://api.organizze.com.br/rest/v2"

type Client struct {
	name       string
	email      string
	apiKey     string
	httpClient *http.Client
}

func NewClient(name string, email string, apiKey string) *Client {
	return &Client{
		name:       name,
		email:      email,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c Client) Categories() ([]Category, error) {
	res, err := c.request(http.MethodGet, "categories", nil, nil)
	if err != nil {
		return nil, err
	}

	return parseResponse[[]Category](res)
}

func (c Client) Invoices(creditCardId int64) ([]InvoiceHeader, error) {
	endpoint := fmt.Sprintf("credit_cards/%v/invoices", creditCardId)
	res, err := c.request(http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	return parseResponse[[]InvoiceHeader](res)
}

func (c Client) Invoice(creditCardId int64, invoiceId int64) (*Invoice, error) {
	endpoint := fmt.Sprintf("credit_cards/%v/invoices/%v", creditCardId, invoiceId)
	res, err := c.request(http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	i, err := parseResponse[Invoice](res)
	return &i, err
}

func (c Client) CreateTransaction(t Transaction) (*Transaction, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	res, err := c.request(http.MethodPost, "transactions", nil, b)
	if err != nil {
		return nil, err
	}

	t, err = parseResponse[Transaction](res)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c Client) UpdateTransaction(t Transaction) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("transactions/%v", t.Id)
	res, err := c.request(http.MethodPut, endpoint, nil, b)
	fmt.Println(res)

	if err != nil {
		return err
	}

	return nil
}

func (c Client) DeleteTransaction(t Transaction) error {
	endpoint := fmt.Sprintf("transactions/%v", t.Id)
	res, err := c.request(http.MethodDelete, endpoint, nil, nil)
	fmt.Println(res)

	if err != nil {
		return err
	}

	return nil
}

func (c Client) request(method, endpoint string, params map[string]string, body []byte) (*http.Response, error) {
	buf := bytes.NewBuffer(body)

	url := fmt.Sprintf("%s/%s", BaseURL, endpoint)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}
	req.URL.RawQuery = q.Encode()

	authorization := "Basic " + base64.StdEncoding.EncodeToString([]byte(c.email+":"+c.apiKey))
	req.Header.Add("Authorization", authorization)

	userAgent := fmt.Sprintf("%s (%s)", c.name, c.email)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func parseResponse[T any](r *http.Response) (t T, err error) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return t, err
	}

	if err := json.Unmarshal([]byte(body), &t); err != nil {
		return t, err
	}

	return t, nil
}
