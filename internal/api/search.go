package api

import (
	"context"
	"fmt"
	"net/url"
)

type SearchOptions struct {
	Query     string
	Type      string
	Subreddit string
	Sort      string
	Time      string
	Limit     int
}

func (c *Client) SearchPosts(ctx context.Context, opts SearchOptions) ([]Post, error) {
	query := searchQuery(opts)
	query.Set("type", "link")
	return c.listPosts(ctx, searchPath(opts.Subreddit), query)
}

func (c *Client) SearchComments(ctx context.Context, opts SearchOptions) ([]Comment, error) {
	query := searchQuery(opts)
	query.Set("type", "comment")

	var resp listingResponse[Comment]
	if err := c.getJSON(ctx, searchPath(opts.Subreddit), query, &resp); err != nil {
		return nil, err
	}
	return extractChildren(resp.Data.Children), nil
}

func (c *Client) SearchSubreddits(ctx context.Context, opts SearchOptions) ([]Subreddit, error) {
	query := searchQuery(opts)
	query.Set("type", "sr")

	var resp listingResponse[Subreddit]
	if err := c.getJSON(ctx, "/search", query, &resp); err != nil {
		return nil, err
	}
	return extractChildren(resp.Data.Children), nil
}

func (c *Client) SearchUsers(ctx context.Context, opts SearchOptions) ([]User, error) {
	query := url.Values{}
	query.Set("q", opts.Query)
	query.Set("limit", fmt.Sprintf("%d", max(1, opts.Limit)))

	var resp listingResponse[User]
	if err := c.getJSON(ctx, "/users/search", query, &resp); err != nil {
		return nil, err
	}
	return extractChildren(resp.Data.Children), nil
}

func searchQuery(opts SearchOptions) url.Values {
	query := listingParams(opts.Limit, opts.Time)
	query.Set("q", opts.Query)
	if opts.Sort != "" {
		query.Set("sort", opts.Sort)
	}
	if opts.Subreddit != "" {
		query.Set("restrict_sr", "1")
	}
	return query
}

func searchPath(subreddit string) string {
	if subreddit == "" {
		return "/search"
	}
	return fmt.Sprintf("/r/%s/search", normalizeSubreddit(subreddit))
}
