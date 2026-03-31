package api

import "context"

func (c *Client) Me(ctx context.Context) (*Me, error) {
	var me Me
	if err := c.getJSON(ctx, "/api/v1/me", nil, &me); err != nil {
		return nil, err
	}
	return &me, nil
}
