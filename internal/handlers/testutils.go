package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
)

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
