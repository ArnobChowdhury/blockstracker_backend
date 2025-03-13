package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestOption func(*http.Request)

func WithAccessToken(token string) RequestOption {
	return func(req *http.Request) { req.Header.Set("Authorization", "Bearer "+token) }
}
func CreateRequest(method, path string, body interface{}, options ...RequestOption) (*http.Request, error) {
	var jsonBody []byte
	var err error

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for _, option := range options {
		option(req)
	}
	return req, nil
}
