package flowerpress

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const SessionCookieName = "flowerpress_session"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthenticated    = errors.New("authentication required")
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Client) Login(ctx context.Context, username, password string) (*User, []*http.Cookie, error) {
	body, err := json.Marshal(loginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("encode login request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/auth/login",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("flowerpress login request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return nil, res.Cookies(), ErrInvalidCredentials
	}

	if res.StatusCode != http.StatusOK {
		return nil, res.Cookies(), fmt.Errorf(
			"flowerpress login returned %s",
			res.Status,
		)
	}

	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, res.Cookies(), fmt.Errorf(
			"decode login response: %w",
			err,
		)
	}

	return &user, res.Cookies(), nil
}

func (c *Client) Me(ctx context.Context, session string) (*User, []*http.Cookie, error) {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet,
		c.baseURL+"/api/auth/me",
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create me request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: session,
	})

	res, err := c.http.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"flowerpress me request: %w",
			err,
		)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return nil, res.Cookies(), ErrUnauthenticated
	}

	if res.StatusCode != http.StatusOK {
		return nil, res.Cookies(), fmt.Errorf(
			"flowerpress me returned %s",
			res.Status,
		)
	}

	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, res.Cookies(), fmt.Errorf(
			"decode me response: %w",
			err,
		)
	}

	return &user, res.Cookies(), nil
}

func (c *Client) Logout(ctx context.Context, session string) ([]*http.Cookie, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/auth/logout", nil)
	if err != nil {
		return nil, fmt.Errorf("create logout request: %w", err)
	}

	if session != "" {
		req.AddCookie(&http.Cookie{
			Name:  SessionCookieName,
			Value: session,
		})
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"flowerpress logout request: %w",
			err,
		)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return res.Cookies(), fmt.Errorf(
			"flowerpress logout returned %s",
			res.Status,
		)
	}

	return res.Cookies(), nil
}
