package flowerpress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type Health struct {
	Status string `json:"status"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Health(ctx context.Context) (*Health, error) {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet,
		c.baseURL+"/health",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create health request: %w", err)
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("flowerpress health request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flowerpress health returned %s", res.Status)
	}

	var health Health
	if err := json.NewDecoder(res.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("decode flowerpress health: %w", err)
	}

	return &health, nil
}
