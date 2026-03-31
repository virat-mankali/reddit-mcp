package api

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) ListMySubreddits(ctx context.Context, limit int) ([]Subreddit, error) {
	var resp listingResponse[Subreddit]
	if err := c.getJSON(ctx, "/subreddits/mine/subscriber", listingParams(limit, ""), &resp); err != nil {
		return nil, err
	}
	return extractChildren(resp.Data.Children), nil
}

func (c *Client) GetSubreddit(ctx context.Context, subreddit string) (*Subreddit, error) {
	var wrapped struct {
		Kind string    `json:"kind"`
		Data Subreddit `json:"data"`
	}
	if err := c.getJSON(ctx, fmt.Sprintf("/r/%s/about", normalizeSubreddit(subreddit)), nil, &wrapped); err != nil {
		return nil, err
	}
	return &wrapped.Data, nil
}

func (c *Client) SearchSubredditsByName(ctx context.Context, queryString string, limit int) ([]Subreddit, error) {
	query := url.Values{}
	query.Set("q", queryString)
	query.Set("limit", fmt.Sprintf("%d", max(1, limit)))

	var resp listingResponse[Subreddit]
	if err := c.getJSON(ctx, "/subreddits/search", query, &resp); err != nil {
		return nil, err
	}
	return extractChildren(resp.Data.Children), nil
}

func (c *Client) Subscribe(ctx context.Context, subreddit string, subscribe bool) error {
	form := url.Values{}
	if subscribe {
		form.Set("action", "sub")
	} else {
		form.Set("action", "unsub")
	}
	form.Set("sr_name", normalizeSubreddit(subreddit))
	return c.postForm(ctx, "/api/subscribe", form, nil)
}
