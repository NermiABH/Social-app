package model

type Comment struct {
	ID             int
	Text           string
	Created        string
	Changed        bool
	LikeCount      int
	DislikeCount   int
	AuthorID       int
	AuthorUsername string
	AuthorUserpic  string
	PostID         int
	ParentID       *int
	Liked          bool
	Disliked       bool
	Own            bool
}

func (c *Comment) ConvertMap() *map[string]any {
	return &map[string]any{
		"type": "comments",
		"id":   c.ID,
		"attributes": map[string]any{
			"text":     c.Text,
			"created":  c.Created,
			"changed":  c.Changed,
			"likes":    c.Liked,
			"dislikes": c.Disliked,
		},
		"relationships": map[string]any{
			"author": map[any]any{
				"id":       c.AuthorID,
				"username": c.AuthorUsername,
				"userpic":  c.AuthorUserpic,
			},
			"post:id":   c.PostID,
			"parent:id": c.ParentID,
			"self": map[string]any{
				"liked":    c.Liked,
				"disliked": c.Disliked,
				"own":      c.Own,
			},
		},
	}
}
