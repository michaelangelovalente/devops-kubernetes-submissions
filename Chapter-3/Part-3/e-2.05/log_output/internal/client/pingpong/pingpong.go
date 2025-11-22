package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	// GetCount(ctx context.Context) (int, error)
	GetCount() (int, error)
}

type httpClient struct {
	client  *http.Client
	baseUrl string
}

func NewClient(baseUrl string, timeout time.Duration) Client {
	return &httpClient{
		client: &http.Client{
			Timeout: timeout,
		},
		baseUrl: baseUrl,
	}
}

func (c *httpClient) GetCount() (int, error) {
	url := fmt.Sprintf("%s/pingpong", c.baseUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return -1, fmt.Errorf("failed to create request to pingpong count: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return -1, fmt.Errorf("failed to send request to pingpong count: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("received status code non-OK: %s", resp.Status)
	}

	var payload struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return -1, fmt.Errorf("could not decode response body: %w", err)
	}
	return payload.Count, nil
}
