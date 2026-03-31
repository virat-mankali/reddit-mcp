package api

import (
	"context"
	"net/url"
)

func (c *Client) ListMessages(ctx context.Context, folder string, limit int) ([]Message, error) {
	var resp listingResponse[Message]
	if err := c.getJSON(ctx, "/message/"+folder, listingParams(limit, ""), &resp); err != nil {
		return nil, err
	}
	return extractChildren(resp.Data.Children), nil
}

func (c *Client) SendMessage(ctx context.Context, to, subject, body string) error {
	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("to", to)
	form.Set("subject", subject)
	form.Set("text", body)
	return c.postForm(ctx, "/api/compose", form, nil)
}

func (c *Client) MarkMessageRead(ctx context.Context, messageID string) error {
	form := url.Values{}
	form.Set("id", normalizeThingID(messageID, "t4"))
	return c.postForm(ctx, "/api/read_message", form, nil)
}
