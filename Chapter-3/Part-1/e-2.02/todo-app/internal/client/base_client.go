package client

import (
	"net/http"
)

type Client interface {
	GetClientBaseURL() string
}

type HTTPClient struct {
	client  *http.Client
	baseURL string
}

// func NewClient(baseURL string, timeout time.Duration) *Client {
// 	return &httpClient{
// 		client: &http.Client{
// 			Timeout: timeout,
// 		},
// 		baseURL: baseURL,
// 	}
// }
