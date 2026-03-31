package api

import (
	"fmt"
	"net/url"
	"strings"
)

func listingParams(limit int, timeFilter string) url.Values {
	query := url.Values{}
	query.Set("limit", fmt.Sprintf("%d", max(1, limit)))
	if timeFilter != "" {
		query.Set("t", timeFilter)
	}
	return query
}

func normalizeListingSort(sort string) string {
	switch sort {
	case "", "best":
		return "best"
	case "hot", "new", "top":
		return sort
	default:
		return sort
	}
}

func normalizeSubreddit(subreddit string) string {
	subreddit = strings.TrimSpace(subreddit)
	subreddit = strings.TrimPrefix(subreddit, "r/")
	return subreddit
}

func normalizeThingID(id, prefix string) string {
	id = strings.TrimSpace(id)
	if strings.HasPrefix(id, "t1_") || strings.HasPrefix(id, "t3_") || strings.HasPrefix(id, "t4_") || strings.HasPrefix(id, "t5_") {
		return id
	}
	return prefix + "_" + id
}

func bareThingID(id string) string {
	id = strings.TrimSpace(id)
	for _, prefix := range []string{"t1_", "t3_", "t4_", "t5_"} {
		id = strings.TrimPrefix(id, prefix)
	}
	return id
}

func BareThingIDForMCP(id string) string {
	return bareThingID(id)
}

func flattenComments(children []listingChild[Comment], depth int) []Comment {
	var out []Comment
	for _, child := range children {
		if child.Kind != "t1" {
			continue
		}
		comment := child.Data
		comment.Depth = depth
		out = append(out, comment)
		if len(comment.Replies.Data.Children) > 0 {
			out = append(out, flattenComments(comment.Replies.Data.Children, depth+1)...)
		}
	}
	return out
}

func extractChildren[T any](children []listingChild[T]) []T {
	items := make([]T, 0, len(children))
	for _, child := range children {
		items = append(items, child.Data)
	}
	return items
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
