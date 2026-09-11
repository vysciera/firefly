package flowerpress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Project struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Status      string `json:"status"`

	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (c *Client) ListProjects(ctx context.Context, session string) ([]Project, error) {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet,
		c.baseURL+"/api/projects",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create projects request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: session,
	})

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("flowerpress projects request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthenticated
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flowerpress projectsr returned %s", err)
	}

	var projects []Project

	if err := json.NewDecoder(res.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("decode projects response: %w", err)
	}

	return projects, nil
}
