package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func (c *Client) Feed(ctx context.Context, sort, timeFilter string, limit int) ([]Post, error) {
	path := "/" + normalizeListingSort(sort)
	query := listingParams(limit, timeFilter)
	return c.listPosts(ctx, path, query)
}

func (c *Client) ListPosts(ctx context.Context, subreddit, sort, timeFilter string, limit int) ([]Post, error) {
	subreddit = normalizeSubreddit(subreddit)
	path := fmt.Sprintf("/r/%s/%s", subreddit, normalizeListingSort(sort))
	query := listingParams(limit, timeFilter)
	return c.listPosts(ctx, path, query)
}

func (c *Client) GetPost(ctx context.Context, postID string, commentsLimit int) (*PostDetails, error) {
	query := url.Values{}
	query.Set("limit", fmt.Sprintf("%d", max(1, commentsLimit)))

	var raw []json.RawMessage
	if err := c.getJSON(ctx, "/comments/"+bareThingID(postID), query, &raw); err != nil {
		return nil, err
	}
	if len(raw) < 2 {
		return nil, fmt.Errorf("post not found")
	}

	var postListing listingResponse[Post]
	if err := json.Unmarshal(raw[0], &postListing); err != nil {
		return nil, err
	}
	if len(postListing.Data.Children) == 0 {
		return nil, fmt.Errorf("post not found")
	}

	var commentListing listingResponse[Comment]
	if err := json.Unmarshal(raw[1], &commentListing); err != nil {
		return nil, err
	}

	post := postListing.Data.Children[0].Data
	comments := flattenComments(commentListing.Data.Children, 0)
	return &PostDetails{Post: post, Comments: comments}, nil
}

func (c *Client) SubmitPost(ctx context.Context, subreddit, title, text, linkURL, postType string) (*SubmitResult, error) {
	form := url.Values{}
	form.Set("api_type", "json")
	form.Set("sr", normalizeSubreddit(subreddit))
	form.Set("title", title)
	form.Set("resubmit", "true")

	switch postType {
	case "text", "self", "":
		form.Set("kind", "self")
		form.Set("text", text)
	case "link":
		form.Set("kind", "link")
		form.Set("url", linkURL)
	default:
		return nil, fmt.Errorf("unsupported post type: %s", postType)
	}

	var resp submitAPIResponse
	if err := c.postForm(ctx, "/api/submit", form, &resp); err != nil {
		return nil, err
	}
	if len(resp.JSON.Errors) > 0 {
		return nil, fmt.Errorf("submit failed: %v", resp.JSON.Errors)
	}

	return &SubmitResult{
		URL:       resp.JSON.Data.URL,
		Name:      resp.JSON.Data.Name,
		ID:        strings.TrimPrefix(resp.JSON.Data.Name, "t3_"),
		Subreddit: normalizeSubreddit(subreddit),
	}, nil
}

func (c *Client) listPosts(ctx context.Context, path string, query url.Values) ([]Post, error) {
	var resp listingResponse[Post]
	if err := c.getJSON(ctx, path, query, &resp); err != nil {
		return nil, err
	}
	return extractChildren(resp.Data.Children), nil
}
