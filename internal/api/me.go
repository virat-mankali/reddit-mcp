package api

import (
	"context"
	"encoding/json"
	"net/http"
)

type Me struct {
	Name         string `json:"name"`
	LinkKarma    int    `json:"link_karma"`
	CommentKarma int    `json:"comment_karma"`
}

func (c *Client) Me(ctx context.Context) (*Me, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/me", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var me Me
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
		return nil, err
	}
	return &me, nil
}
