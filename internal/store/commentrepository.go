package store

import (
	"Social-app/internal/model"
	"fmt"
)

type CommentRepository struct {
	store *Store
}

func (r *CommentRepository) Create(c *model.Comment) error {
	q := `INSERT INTO comment (author_id, post_id, parent_id, text)
		VALUES ($1, $2, $3, $4) RETURNING id;`
	return r.store.db.QueryRow(q,
		c.AuthorID, c.PostID, c.ParentID, c.Text,
	).Scan(&c.ID)
}

func (r *CommentRepository) GetSeveral(pID, offset, limit int) ([]*model.Comment, error) {
	q := `SELECT sc.id, sc.author_id, sc.username, sc.userpic, sc.post_id,
       sc.parent_id, sc.text, sc.date_of_creation, sc.changed
FROM (SELECT * FROM (SELECT id FROM comment
WHERE parent_id isnull AND post_id=$1) as c OFFSET $2 LIMIT $3) as ci
INNER JOIN select_comment(ci.id) sc ON true`
	row, err := r.store.db.Query(q, pID, offset, limit)
	if err != nil {
		return nil, err
	}
	cSlice := make([]*model.Comment, 0)
	for row.Next() {
		c := &model.Comment{}
		if err = row.Scan(&c.ID, &c.AuthorID, &c.AuthorUsername, &c.AuthorUserpic,
			&c.PostID, &c.ParentID, &c.Text, &c.DateCreation, &c.Change); err != nil {
			return nil, err
		}
		cSlice = append(cSlice, c)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return cSlice, nil
}
func (r *CommentRepository) GetByID(c *model.Comment) error {
	q := `SELECT c.author_id, u.username, u.userpic, c.post_id, c.parent_id,
    	c.text, c.date_of_creation, c.changed,
    	COUNT(ulc) as likes, COUNT(udc) as dislikes
		FROM comment c
        LEFT JOIN users u on u.id = c.author_id
        LEFT JOIN user_like_comment ulc on c.id = ulc.comment_id
        LEFT JOIN user_dislike_comment udc on c.id = udc.comment_id
		WHERE c.id = $1
		GROUP BY c.author_id, c.post_id, u.username, u.userpic, 
				 c.parent_id, c.text, c.date_of_creation, c.changed`
	return r.store.db.QueryRow(q, c.ID).Scan(
		&c.AuthorID, &c.AuthorUsername, &c.AuthorUserpic, &c.PostID, &c.ParentID,
		&c.Text, &c.DateCreation, &c.Change, &c.Likes, &c.Dislikes)
}

func (r *CommentRepository) Update(id int, text string) error {
	q := `UPDATE comment SET text = $1, changed = true WHERE id=$2`
	return r.store.db.QueryRow(q, text, id).Err()
}
func (r *CommentRepository) Delete(id int) error {
	q := `DELETE FROM comment WHERE id = $1`
	return r.store.db.QueryRow(q, id).Err()
}

func (r *CommentRepository) Like(userID, commentID int) error {
	q := `INSERT INTO user_like_comment VALUES ($1, $2)`
	if err := r.store.db.QueryRow(q, userID, commentID).Err(); err != nil {
		return err
	}
	return r.UnDislike(userID, commentID)
}

func (r *CommentRepository) UnLike(userID, commentID int) error {
	q := `DELETE FROM user_like_comment WHERE user_id = $1 AND comment_id = $2`
	return r.store.db.QueryRow(q, userID, commentID).Err()
}

func (r *CommentRepository) Dislike(userID, commentID int) error {
	q := `INSERT INTO user_dislike_comment VALUES ($1, $2)`
	if err := r.store.db.QueryRow(q, userID, commentID).Err(); err != nil {
		return err
	}
	return r.UnLike(userID, commentID)
}

func (r *CommentRepository) UnDislike(userID, commentID int) error {
	q := `DELETE FROM user_dislike_comment WHERE user_id = $1 AND comment_id = $2`
	return r.store.db.QueryRow(q, userID, commentID).Err()
}

func (r *CommentRepository) IsExist(id int) (bool, error) {
	q := `SELECT EXISTS (SELECT * FROM comment WHERE id = $1)`
	var exist bool
	err := r.store.db.QueryRow(q, id).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (r *CommentRepository) IsOwner(uID, cID int) (bool, error) {
	q := `SELECT EXISTS (SELECT * FROM comment WHERE author_id = $1 AND id = $2)`
	var owner bool
	if err := r.store.db.QueryRow(q, uID, cID).Scan(&owner); err != nil {
		return false, err
	}
	return owner, nil
}

func (r *CommentRepository) LikedOrDisliked(c *model.Comment, uID int) error {
	q := `SELECT EXISTS (SELECT * FROM user_like_comment WHERE user_id = $1 AND comment_id = $2)`
	var boolean bool
	if err := r.store.db.QueryRow(q, uID, c.ID).Scan(&boolean); err != nil {
		return err
	} else if boolean {
		c.IsLiked = true
	}
	fmt.Println(boolean)
	q = `SELECT EXISTS (SELECT * FROM user_dislike_comment WHERE user_id = $1 AND comment_id = $2)`
	if err := r.store.db.QueryRow(q, uID, c.ID).Scan(&boolean); err != nil {
		return err
	} else if boolean {
		c.IsDisliked = true
	}
	return nil
}
