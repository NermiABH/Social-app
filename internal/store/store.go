package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	db                *sql.DB
	userRepository    *UserRepository
	postRepository    *PostRepository
	commentRepository *CommentRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() *UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}
	return s.userRepository
}

func (s *Store) Post() *PostRepository {
	if s.postRepository == nil {
		s.postRepository = &PostRepository{
			store: s,
		}
	}
	return s.postRepository
}

func (s *Store) Comment() *CommentRepository {
	if s.commentRepository == nil {
		s.commentRepository = &CommentRepository{
			store: s,
		}
	}
	return s.commentRepository
}
