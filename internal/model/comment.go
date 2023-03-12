package model

type Comment struct {
	ID             int     `json:"id"`
	AuthorID       int     `json:"author_id"`
	AuthorUsername string  `json:"author_username"`
	AuthorUserpic  *string `json:"author_userpic"`
	PostID         int     `json:"post_id"`
	ParentID       *int    `json:"parent_id"`
	Text           string  `json:"text"`
	DateCreation   string  `json:"date_of_creation,omitempty"`
	Change         bool    `json:"change,omitempty"`
	Likes          int     `json:"likes,omitempty"`
	Dislikes       int     `json:"dislikes,omitempty"`
	IsOwn          bool    `json:"is_own,omitempty"`
	IsLiked        bool    `json:"is_like,omitempty"`
	IsDisliked     bool    `json:"is_dislike,omitempty"`
}
