package api

import (
	"encoding/json"
	"time"
)

type Me struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	LinkKarma     int     `json:"link_karma"`
	CommentKarma  int     `json:"comment_karma"`
	TotalKarma    int     `json:"total_karma"`
	CreatedUTC    float64 `json:"created_utc"`
	HasVerifiedEM bool    `json:"has_verified_email"`
	IsGold        bool    `json:"is_gold"`
}

func (m Me) CreatedAt() time.Time {
	return time.Unix(int64(m.CreatedUTC), 0)
}

type Post struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	Subreddit     string  `json:"subreddit"`
	Permalink     string  `json:"permalink"`
	URL           string  `json:"url"`
	SelfText      string  `json:"selftext"`
	Score         int     `json:"score"`
	NumComments   int     `json:"num_comments"`
	Saved         bool    `json:"saved"`
	Likes         *bool   `json:"likes"`
	IsSelf        bool    `json:"is_self"`
	CreatedUTC    float64 `json:"created_utc"`
	Over18        bool    `json:"over_18"`
	LinkFlair     string  `json:"link_flair_text"`
	SubredditID   string  `json:"subreddit_id"`
	AuthorFull    string  `json:"author_fullname"`
	Distinguished string  `json:"distinguished"`
}

func (p Post) CreatedAt() time.Time {
	return time.Unix(int64(p.CreatedUTC), 0)
}

type Comment struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Author     string      `json:"author"`
	Body       string      `json:"body"`
	Score      int         `json:"score"`
	Permalink  string      `json:"permalink"`
	ParentID   string      `json:"parent_id"`
	LinkID     string      `json:"link_id"`
	Subreddit  string      `json:"subreddit"`
	Saved      bool        `json:"saved"`
	Likes      *bool       `json:"likes"`
	CreatedUTC float64     `json:"created_utc"`
	Replies    commentTree `json:"replies"`
	Depth      int         `json:"-"`
}

func (c Comment) CreatedAt() time.Time {
	return time.Unix(int64(c.CreatedUTC), 0)
}

type Subreddit struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	DisplayName       string  `json:"display_name"`
	Title             string  `json:"title"`
	PublicDescription string  `json:"public_description"`
	URL               string  `json:"url"`
	Subscribers       int     `json:"subscribers"`
	CreatedUTC        float64 `json:"created_utc"`
	UserIsSubscriber  bool    `json:"user_is_subscriber"`
	Over18            bool    `json:"over18"`
}

func (s Subreddit) CreatedAt() time.Time {
	return time.Unix(int64(s.CreatedUTC), 0)
}

type Message struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Author     string  `json:"author"`
	Dest       string  `json:"dest"`
	Subject    string  `json:"subject"`
	Body       string  `json:"body"`
	Subreddit  string  `json:"subreddit"`
	WasComment bool    `json:"was_comment"`
	New        bool    `json:"new"`
	CreatedUTC float64 `json:"created_utc"`
}

func (m Message) CreatedAt() time.Time {
	return time.Unix(int64(m.CreatedUTC), 0)
}

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	IconImg      string `json:"icon_img"`
	CommentKarma int    `json:"comment_karma"`
	LinkKarma    int    `json:"link_karma"`
}

type PostDetails struct {
	Post     Post      `json:"post"`
	Comments []Comment `json:"comments"`
}

type SubmitResult struct {
	URL       string `json:"url"`
	Name      string `json:"name"`
	ID        string `json:"id"`
	Subreddit string `json:"subreddit"`
}

type listingResponse[T any] struct {
	Data struct {
		After    string            `json:"after"`
		Before   string            `json:"before"`
		Children []listingChild[T] `json:"children"`
	} `json:"data"`
}

type listingChild[T any] struct {
	Kind string `json:"kind"`
	Data T      `json:"data"`
}

type commentTree struct {
	Data struct {
		Children []listingChild[Comment] `json:"children"`
	} `json:"data"`
}

func (c *commentTree) UnmarshalJSON(data []byte) error {
	if string(data) == `""` || string(data) == "null" {
		return nil
	}
	type alias commentTree
	var out alias
	if err := json.Unmarshal(data, &out); err != nil {
		return err
	}
	*c = commentTree(out)
	return nil
}

type submitAPIResponse struct {
	JSON struct {
		Data struct {
			Name string `json:"name"`
			URL  string `json:"url"`
			ID   string `json:"id"`
		} `json:"data"`
		Errors [][]any `json:"errors"`
	} `json:"json"`
}
