package store

import "Social-app/internal/model"

type PostRepository struct {
	store *Store
}

func (r *PostRepository) CreatePost(p *model.Post) error {
	return r.store.db.QueryRow(
		"INSERT INTO post (author_id, text, object, views) VALUES ($1, $2, $3, $4) RETURNING id",
		p.AuthorID, p.Text, p.Object, p.Views).Scan(&p.ID)
}
