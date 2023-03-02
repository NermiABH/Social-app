package model

type Post struct {
	ID            int    `json:"id"`
	AuthorID      int    `json:"author_id,omitempty"`
	Text          string `json:"text,omitempty"`
	Object        string `json:"object,omitempty"`
	Views         int    `json:"views,omitempty"`
	CommentsCount int    `json:"commentsCount,omitempty"`
	DateCreation  string `json:"date_of_creation,omitempty"`
	Likes         int    `json:"likes"`
	Dislikes      int    `json:"dislikes"`
}
