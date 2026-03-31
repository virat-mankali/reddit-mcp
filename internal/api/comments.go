package api

import (
	"context"
	"net/url"
)

func (c *Client) ListComments(ctx context.Context, postID string, limit int) ([]Comment, error) {
	details, err := c.GetPost(ctx, postID, limit)
	if err != nil {
		return nil, err
	}
	return details.Comments, nil
}

func (c *Client) Reply(ctx context.Context, parentID, text string) error {
	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("thing_id", parentID)
	form.Set("text", text)
	return c.postForm(ctx, "/api/comment", form, nil)
}
