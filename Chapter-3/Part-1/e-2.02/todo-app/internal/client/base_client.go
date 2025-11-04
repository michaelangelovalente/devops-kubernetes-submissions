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
