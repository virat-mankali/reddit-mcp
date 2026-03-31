package api

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) DeleteThing(ctx context.Context, thingID string) error {
	form := url.Values{}
	form.Set("id", thingID)
	return c.postForm(ctx, "/api/del", form, nil)
}

func (c *Client) Vote(ctx context.Context, thingID string, dir int) error {
	form := url.Values{}
	form.Set("id", thingID)
	form.Set("dir", fmt.Sprintf("%d", dir))
	return c.postForm(ctx, "/api/vote", form, nil)
}

func (c *Client) Save(ctx context.Context, thingID string, save bool) error {
	form := url.Values{}
	form.Set("id", thingID)
	path := "/api/save"
	if !save {
		path = "/api/unsave"
	}
	return c.postForm(ctx, path, form, nil)
}
