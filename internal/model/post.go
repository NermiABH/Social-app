package model

type Post struct {
	ID             int
	Text           *string
	Media          *[]string
	ViewCount      int
	Created        *string
	Changed        bool
	CommentCount   int
	LikeCount      int
	DislikeCount   int
	AuthorID       int
	AuthorUsername string
	AuthorUserpic  string
	Liked          bool
	Disliked       bool
	Favorited      bool
	Own            bool
}

func (p *Post) ConvertMap() map[string]any {
	return map[string]any{
		"type": "post",
		"id":   p.ID,
		"attributes": map[string]any{
			"text":         p.Text,
			"media":        p.Media,
			"created":      p.Created,
			"changed":      p.Changed,
			"viewCount":    p.ViewCount,
			"likeCount":    p.LikeCount,
			"dislikeCount": p.DislikeCount,
			"commentCount": p.CommentCount,
		},
		"relationships": map[string]any{
			"author": map[string]any{
				"id":       p.AuthorID,
				"username": p.AuthorUsername,
				"userpic":  p.AuthorUserpic,
			},
			"self": map[string]any{
				"liked":     p.Liked,
				"disliked":  p.Disliked,
				"favorited": p.Favorited,
				"own":       p.Own,
			},
		},
	}
}

type Posts []*Post

func (p Posts) ConvertMap() []map[string]any {
	posts := make([]map[string]any, len(p))
	for i, post := range p {
		posts[i] = post.ConvertMap()
	}
	return posts
}
