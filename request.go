package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	rjson "github.com/ramin0/json"
)

// Client type
type Client struct {
	BaseURL        string
	DefaultParams  url.Values
	DefaultHeaders url.Values
}

type kvs interface {
	Add(string, string)
}

// NewClient func
func NewClient() *Client {
	return &Client{
		BaseURL:        "",
		DefaultParams:  url.Values{},
		DefaultHeaders: url.Values{},
	}
}

func (c *Client) httpClient() *http.Client {
	return http.DefaultClient
}

// Get func
func (c *Client) Get(path string, params, headers url.Values) (*rjson.JSON, error) {
	return c.request(http.MethodGet, path, params, nil, headers)
}

// Post func
func (c *Client) Post(path string, params url.Values, bodyJSON *rjson.JSON, headers url.Values) (*rjson.JSON, error) {
	return c.request(http.MethodPost, path, params, bodyJSON, headers)
}

func merge(base kvs, extra url.Values) {
	if extra != nil {
		for k, vs := range extra {
			for _, v := range vs {
				base.Add(k, v)
			}
		}
	}
}

func (c *Client) request(method, path string, params url.Values, bodyJSON *rjson.JSON, headers url.Values) (*rjson.JSON, error) {
	uri, _ := url.Parse(fmt.Sprintf("%s/%s", c.BaseURL, path))

	query := url.Values{}
	merge(query, c.DefaultParams)
	merge(query, params)
	uri.RawQuery = query.Encode()

	requestBody := []byte(nil)
	if bodyJSON != nil {
		body, err := json.Marshal(bodyJSON.Raw)
		if err != nil {
			return nil, err
		}
		requestBody = body
	}

	request, _ := http.NewRequest(method, uri.String(), bytes.NewReader(requestBody))

	merge(request.Header, c.DefaultHeaders)
	merge(request.Header, headers)

	response, err := c.httpClient().Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode > 399 {
		return nil, fmt.Errorf(response.Status)
	}
	var responseJSON rjson.JSON
	if err := json.NewDecoder(response.Body).Decode(&responseJSON.Raw); err != nil {
		return nil, err
	}

	return &responseJSON, nil
}
