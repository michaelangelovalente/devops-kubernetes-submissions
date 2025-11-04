package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TodoClient struct {
	HTTPClient
}

func NewTodoClient(baseURL string, timeout time.Duration) *TodoClient {
	return &TodoClient{
		HTTPClient: HTTPClient{
			client: &http.Client{
				Timeout: timeout,
			},
			baseURL: baseURL,
		},
	}
}

func (tc *TodoClient) GetClientBaseURL() string {
	return tc.baseURL
}

type Todo struct {
	ID        int
	Task      string
	Completed bool
}

func (tc *TodoClient) GetTodos() ([]Todo, error) {
	reqURL := fmt.Sprintf("%s/todos", tc.GetClientBaseURL())
	resp, err := tc.client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("could not perform get request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var todoResp struct {
		Data []Todo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&todoResp); err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}

	return todoResp.Data, nil
}

func (tc *TodoClient) AddTodo(todoEntry string) (Todo, error) {
	reqURL := fmt.Sprintf("%s/%s", tc.GetClientBaseURL(), "todo")
	todoReq := struct {
		Data string `json:"data"`
	}{
		Data: todoEntry,
	}

	fmt.Println("Sending request to: ", reqURL)
	fmt.Println("BODY: ", todoReq)

	jsonReq, err := json.Marshal(todoReq)
	if err != nil {
		return Todo{}, fmt.Errorf("could not marshall request: %v", err)
	}

	resp, err := tc.client.Post(reqURL, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return Todo{}, fmt.Errorf("could not create post request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Todo{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var todoResp struct {
		Data []Todo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&todoResp); err != nil {
		return Todo{}, fmt.Errorf("could not decode response: %v", err)
	}

	return todoResp.Data[len(todoResp.Data)-1], nil
}
