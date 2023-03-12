package model

type Post struct {
	ID            int    `json:"id"`
	AuthorID      int    `json:"author_id"`
	Text          string `json:"text"`
	Object        string `json:"object"`
	CommentsCount int    `json:"commentsCount"`
	DateCreation  string `json:"date_of_creation"`
	Likes         int    `json:"likes"`
	Dislikes      int    `json:"dislikes"`
	IsOwn         bool   `json:"is_own"`
	IsLiked       bool   `json:"is_like"`
	IsDisliked    bool   `json:"is_disliked"`
	IsFavorited   bool   `json:"is_favorited"`
}
