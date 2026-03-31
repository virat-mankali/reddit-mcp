package mcpserver

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/virat-mankali/reddit-mcp/internal/api"
)

func registerTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("get_me", mcp.WithDescription("Get the authenticated Reddit profile")), handleGetMe)
	s.AddTool(mcp.NewTool(
		"get_feed",
		mcp.WithDescription("Get the authenticated user's home feed"),
		mcp.WithString("sort", mcp.Description("best, hot, new, or top")),
		mcp.WithString("time", mcp.Description("hour, day, week, month, year, all")),
		mcp.WithNumber("limit", mcp.Description("Result count")),
	), handleGetFeed)

	s.AddTool(mcp.NewTool(
		"list_posts",
		mcp.WithDescription("List posts from a subreddit"),
		mcp.WithString("subreddit", mcp.Required(), mcp.Description("Subreddit name, with or without r/ prefix")),
		mcp.WithString("sort", mcp.Description("hot, new, or top")),
		mcp.WithString("time", mcp.Description("hour, day, week, month, year, all")),
		mcp.WithNumber("limit", mcp.Description("Result count")),
	), handleListPosts)
	s.AddTool(mcp.NewTool(
		"get_post",
		mcp.WithDescription("Get a post and its comments"),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("Reddit post ID")),
		mcp.WithNumber("comments", mcp.Description("Number of comments to return")),
	), handleGetPost)
	s.AddTool(mcp.NewTool(
		"submit_post",
		mcp.WithDescription("Submit a post to a subreddit"),
		mcp.WithString("subreddit", mcp.Required(), mcp.Description("Subreddit name")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Post title")),
		mcp.WithString("type", mcp.Description("text or link"), mcp.Enum("text", "link")),
		mcp.WithString("text", mcp.Description("Text body for a self post")),
		mcp.WithString("url", mcp.Description("URL for a link post")),
	), handleSubmitPost)
	s.AddTool(mcp.NewTool(
		"delete_post",
		mcp.WithDescription("Delete a post"),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("Reddit post ID")),
	), handleDeletePost)
	s.AddTool(mcp.NewTool(
		"vote_post",
		mcp.WithDescription("Vote on a post"),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("Reddit post ID")),
		mcp.WithString("direction", mcp.Required(), mcp.Description("up, down, or clear"), mcp.Enum("up", "down", "clear")),
	), handleVotePost)
	s.AddTool(mcp.NewTool(
		"save_post",
		mcp.WithDescription("Save or unsave a post"),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("Reddit post ID")),
		mcp.WithBoolean("save", mcp.Description("true to save, false to unsave")),
	), handleSavePost)

	s.AddTool(mcp.NewTool(
		"list_comments",
		mcp.WithDescription("List comments for a post"),
		mcp.WithString("post_id", mcp.Required(), mcp.Description("Reddit post ID")),
		mcp.WithNumber("limit", mcp.Description("Result count")),
	), handleListComments)
	s.AddTool(mcp.NewTool(
		"reply_comment",
		mcp.WithDescription("Reply to a post or comment"),
		mcp.WithString("parent_id", mcp.Required(), mcp.Description("Post or comment fullname/id")),
		mcp.WithString("text", mcp.Required(), mcp.Description("Reply body")),
	), handleReplyComment)
	s.AddTool(mcp.NewTool(
		"delete_comment",
		mcp.WithDescription("Delete a comment"),
		mcp.WithString("comment_id", mcp.Required(), mcp.Description("Reddit comment ID")),
	), handleDeleteComment)
	s.AddTool(mcp.NewTool(
		"vote_comment",
		mcp.WithDescription("Vote on a comment"),
		mcp.WithString("comment_id", mcp.Required(), mcp.Description("Reddit comment ID")),
		mcp.WithString("direction", mcp.Required(), mcp.Description("up, down, or clear"), mcp.Enum("up", "down", "clear")),
	), handleVoteComment)

	s.AddTool(mcp.NewTool(
		"list_subreddits",
		mcp.WithDescription("List subscribed subreddits"),
		mcp.WithNumber("limit", mcp.Description("Result count")),
	), handleListSubreddits)
	s.AddTool(mcp.NewTool(
		"get_subreddit",
		mcp.WithDescription("Get details for a subreddit"),
		mcp.WithString("subreddit", mcp.Required(), mcp.Description("Subreddit name")),
	), handleGetSubreddit)
	s.AddTool(mcp.NewTool(
		"search_subreddits",
		mcp.WithDescription("Search subreddits by name or topic"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithNumber("limit", mcp.Description("Result count")),
	), handleSearchSubreddits)
	s.AddTool(mcp.NewTool(
		"subscribe_subreddit",
		mcp.WithDescription("Subscribe or unsubscribe from a subreddit"),
		mcp.WithString("subreddit", mcp.Required(), mcp.Description("Subreddit name")),
		mcp.WithBoolean("subscribe", mcp.Description("true to subscribe, false to unsubscribe")),
	), handleSubscribeSubreddit)

	s.AddTool(mcp.NewTool(
		"search_posts",
		mcp.WithDescription("Search Reddit posts"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithString("subreddit", mcp.Description("Optional subreddit filter")),
		mcp.WithString("sort", mcp.Description("relevance, new, top, comments")),
		mcp.WithString("time", mcp.Description("hour, day, week, month, year, all")),
		mcp.WithNumber("limit", mcp.Description("Result count")),
	), handleSearchPosts)
	s.AddTool(mcp.NewTool(
		"search_comments",
		mcp.WithDescription("Search Reddit comments"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithString("subreddit", mcp.Description("Optional subreddit filter")),
		mcp.WithString("sort", mcp.Description("relevance, new, top, comments")),
		mcp.WithString("time", mcp.Description("hour, day, week, month, year, all")),
		mcp.WithNumber("limit", mcp.Description("Result count")),
	), handleSearchComments)
	s.AddTool(mcp.NewTool(
		"search_users",
		mcp.WithDescription("Search Reddit users"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithNumber("limit", mcp.Description("Result count")),
	), handleSearchUsers)

	s.AddTool(mcp.NewTool(
		"list_messages",
		mcp.WithDescription("List inbox, unread, or sent messages"),
		mcp.WithString("folder", mcp.Description("inbox, unread, or sent")),
		mcp.WithNumber("limit", mcp.Description("Result count")),
	), handleListMessages)
	s.AddTool(mcp.NewTool(
		"send_message",
		mcp.WithDescription("Send a Reddit direct message"),
		mcp.WithString("to", mcp.Required(), mcp.Description("Recipient username")),
		mcp.WithString("subject", mcp.Required(), mcp.Description("Message subject")),
		mcp.WithString("body", mcp.Required(), mcp.Description("Message body")),
	), handleSendMessage)
	s.AddTool(mcp.NewTool(
		"read_message",
		mcp.WithDescription("Mark a message as read"),
		mcp.WithString("message_id", mcp.Required(), mcp.Description("Reddit message ID")),
	), handleReadMessage)
}

func handleGetMe(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		me, err := client.Me(ctx)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(me)
	})
}

func handleGetFeed(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		posts, err := client.Feed(ctx, req.GetString("sort", "best"), req.GetString("time", "day"), listLimit(req, 10))
		if err != nil {
			return toolError(err)
		}
		return jsonResult(posts)
	})
}

func handleListPosts(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subreddit, bad, err := requiredString(req, "subreddit")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		posts, err := client.ListPosts(ctx, subreddit, req.GetString("sort", "hot"), req.GetString("time", "day"), listLimit(req, 10))
		if err != nil {
			return toolError(err)
		}
		return jsonResult(posts)
	})
}

func handleGetPost(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	postID, bad, err := requiredString(req, "post_id")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		post, err := client.GetPost(ctx, postID, req.GetInt("comments", 10))
		if err != nil {
			return toolError(err)
		}
		return jsonResult(post)
	})
}

func handleSubmitPost(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subreddit, bad, err := requiredString(req, "subreddit")
	if bad != nil || err != nil {
		return bad, err
	}
	title, bad, err := requiredString(req, "title")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		postType := req.GetString("type", "text")
		result, err := client.SubmitPost(ctx, subreddit, title, req.GetString("text", ""), req.GetString("url", ""), postType)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(result)
	})
}

func handleDeletePost(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	postID, bad, err := requiredString(req, "post_id")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		if err := client.DeleteThing(ctx, "t3_"+api.BareThingIDForMCP(postID)); err != nil {
			return toolError(err)
		}
		return jsonResult(map[string]any{"ok": true, "post_id": postID})
	})
}

func handleVotePost(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	postID, bad, err := requiredString(req, "post_id")
	if bad != nil || err != nil {
		return bad, err
	}
	direction, bad, err := requiredString(req, "direction")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		if err := client.Vote(ctx, "t3_"+api.BareThingIDForMCP(postID), voteDirection(direction)); err != nil {
			return toolError(err)
		}
		return jsonResult(map[string]any{"ok": true, "post_id": postID, "direction": direction})
	})
}

func handleSavePost(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	postID, bad, err := requiredString(req, "post_id")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		save := req.GetBool("save", true)
		if err := client.Save(ctx, "t3_"+api.BareThingIDForMCP(postID), save); err != nil {
			return toolError(err)
		}
		return jsonResult(map[string]any{"ok": true, "post_id": postID, "save": save})
	})
}

func handleListComments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	postID, bad, err := requiredString(req, "post_id")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		comments, err := client.ListComments(ctx, postID, listLimit(req, 20))
		if err != nil {
			return toolError(err)
		}
		return jsonResult(comments)
	})
}

func handleReplyComment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	parentID, bad, err := requiredString(req, "parent_id")
	if bad != nil || err != nil {
		return bad, err
	}
	text, bad, err := requiredString(req, "text")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		parent := parentID
		if api.BareThingIDForMCP(parentID) == parentID {
			parent = "t3_" + parentID
		}
		if err := client.Reply(ctx, parent, text); err != nil {
			return toolError(err)
		}
		return jsonResult(map[string]any{"ok": true, "parent_id": parentID})
	})
}

func handleDeleteComment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	commentID, bad, err := requiredString(req, "comment_id")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		if err := client.DeleteThing(ctx, "t1_"+api.BareThingIDForMCP(commentID)); err != nil {
			return toolError(err)
		}
		return jsonResult(map[string]any{"ok": true, "comment_id": commentID})
	})
}

func handleVoteComment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	commentID, bad, err := requiredString(req, "comment_id")
	if bad != nil || err != nil {
		return bad, err
	}
	direction, bad, err := requiredString(req, "direction")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		if err := client.Vote(ctx, "t1_"+api.BareThingIDForMCP(commentID), voteDirection(direction)); err != nil {
			return toolError(err)
		}
		return jsonResult(map[string]any{"ok": true, "comment_id": commentID, "direction": direction})
	})
}

func handleListSubreddits(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		items, err := client.ListMySubreddits(ctx, listLimit(req, 25))
		if err != nil {
			return toolError(err)
		}
		return jsonResult(items)
	})
}

func handleGetSubreddit(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subreddit, bad, err := requiredString(req, "subreddit")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		item, err := client.GetSubreddit(ctx, subreddit)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(item)
	})
}

func handleSearchSubreddits(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, bad, err := requiredString(req, "query")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		items, err := client.SearchSubredditsByName(ctx, query, listLimit(req, 25))
		if err != nil {
			return toolError(err)
		}
		return jsonResult(items)
	})
}

func handleSubscribeSubreddit(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	subreddit, bad, err := requiredString(req, "subreddit")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		subscribe := req.GetBool("subscribe", true)
		if err := client.Subscribe(ctx, subreddit, subscribe); err != nil {
			return toolError(err)
		}
		return jsonResult(map[string]any{"ok": true, "subreddit": subreddit, "subscribe": subscribe})
	})
}

func handleSearchPosts(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	opts, bad, err := searchOptionsFromRequest(req)
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		items, err := client.SearchPosts(ctx, opts)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(items)
	})
}

func handleSearchComments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	opts, bad, err := searchOptionsFromRequest(req)
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		items, err := client.SearchComments(ctx, opts)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(items)
	})
}

func handleSearchUsers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, bad, err := requiredString(req, "query")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		items, err := client.SearchUsers(ctx, api.SearchOptions{Query: query, Limit: listLimit(req, 10)})
		if err != nil {
			return toolError(err)
		}
		return jsonResult(items)
	})
}

func handleListMessages(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		folder := req.GetString("folder", "inbox")
		items, err := client.ListMessages(ctx, folder, listLimit(req, 10))
		if err != nil {
			return toolError(err)
		}
		return jsonResult(items)
	})
}

func handleSendMessage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	to, bad, err := requiredString(req, "to")
	if bad != nil || err != nil {
		return bad, err
	}
	subject, bad, err := requiredString(req, "subject")
	if bad != nil || err != nil {
		return bad, err
	}
	body, bad, err := requiredString(req, "body")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		if err := client.SendMessage(ctx, to, subject, body); err != nil {
			return toolError(err)
		}
		return jsonResult(map[string]any{"ok": true, "to": to, "subject": subject})
	})
}

func handleReadMessage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	messageID, bad, err := requiredString(req, "message_id")
	if bad != nil || err != nil {
		return bad, err
	}
	return withClient(ctx, func(client *api.Client) (*mcp.CallToolResult, error) {
		if err := client.MarkMessageRead(ctx, messageID); err != nil {
			return toolError(err)
		}
		return jsonResult(map[string]any{"ok": true, "message_id": messageID})
	})
}

func searchOptionsFromRequest(req mcp.CallToolRequest) (api.SearchOptions, *mcp.CallToolResult, error) {
	query, bad, err := requiredString(req, "query")
	if bad != nil || err != nil {
		return api.SearchOptions{}, bad, err
	}
	return api.SearchOptions{
		Query:     query,
		Subreddit: req.GetString("subreddit", ""),
		Sort:      req.GetString("sort", "relevance"),
		Time:      req.GetString("time", "all"),
		Limit:     listLimit(req, 10),
	}, nil, nil
}

func voteDirection(direction string) int {
	switch direction {
	case "up":
		return 1
	case "down":
		return -1
	default:
		return 0
	}
}
