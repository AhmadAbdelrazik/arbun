package handlers

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/pkg/assert"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ProductToPostProductInput(p domain.Product) postProductInput {
	return postProductInput{
		Name:        p.Name,
		Description: p.Description,
		Vendor:      p.Vendor,
		Amount:      p.AvailableAmount,
		Properties:  p.Properties,
		Price:       p.Price,
	}
}

func AddCustomer(t *testing.T, ts *TestClient, c domain.Customer, password string) *http.Cookie {
	t.Helper()

	body := struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"type"`
	}{
		FullName: c.Name,
		Email:    c.Email,
		Password: password,
		UserType: domain.TypeCustomer,
	}
	res, err := ts.Post("/signup", body)
	assert.Nil(t, err)
	cookie := ts.GetCookie(res, AuthCookie)

	assert.Nil(t, cookie.Valid())
	assert.True(t, len(cookie.Value) == 26)

	return cookie
}

func AddAdmin(t *testing.T, ts *TestClient, a domain.Admin, password string) *http.Cookie {
	t.Helper()

	body := struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"type"`
	}{
		FullName: a.Name,
		Email:    a.Email,
		Password: password,
		UserType: domain.TypeAdmin,
	}
	res, err := ts.Post("/signup", body)
	assert.Nil(t, err)
	cookie := ts.GetCookie(res, AuthCookie)

	assert.Nil(t, cookie.Valid())
	assert.True(t, len(cookie.Value) == 26)

	return cookie
}

func AddProduct(t *testing.T, ts *TestClient, p domain.Product, adminCookie *http.Cookie) {
	t.Helper()
	res, err := ts.PostWithCookies("/products", ProductToPostProductInput(p), adminCookie)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)
}

type TestClient struct {
	server *httptest.Server
}

func NewTestClient() *TestClient {
	app := NewApplication()
	return &TestClient{
		server: httptest.NewServer(app.routes()),
	}
}

func (c *TestClient) Close() {
	c.server.Close()
}

func (c *TestClient) Get(endpoint string) (*http.Response, error) {
	return http.Get(c.server.URL + endpoint)
}

func (c *TestClient) Post(endpoint string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return http.Post(c.server.URL+endpoint, "application/json", bytes.NewBuffer(jsonBody))
}

func (c *TestClient) Put(endpoint string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, c.server.URL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func (c *TestClient) Patch(endpoint string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, c.server.URL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func (c *TestClient) Delete(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, c.server.URL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func (c *TestClient) ReadResponseBody(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func (c *TestClient) GetCookie(res *http.Response, cookieName string) *http.Cookie {
	for _, c := range res.Cookies() {
		if c.Name == cookieName {
			return c
		}
	}

	return nil
}

func (c *TestClient) Do(method, endpoint string, body any, cookies ...*http.Cookie) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.server.URL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	for _, c := range cookies {
		req.AddCookie(c)
	}

	return http.DefaultClient.Do(req)
}

func (c *TestClient) GetWithCookies(endpoint string, cookies ...*http.Cookie) (*http.Response, error) {
	return c.Do(http.MethodGet, endpoint, nil, cookies...)
}
func (c *TestClient) PostWithCookies(endpoint string, body any, cookies ...*http.Cookie) (*http.Response, error) {
	return c.Do(http.MethodPost, endpoint, body, cookies...)
}

func (c *TestClient) PatchWithCookies(endpoint string, body any, cookies ...*http.Cookie) (*http.Response, error) {
	return c.Do(http.MethodPatch, endpoint, body, cookies...)
}

func (c *TestClient) DeleteWithCookies(endpoint string, body any, cookies ...*http.Cookie) (*http.Response, error) {
	return c.Do(http.MethodDelete, endpoint, body, cookies...)
}
