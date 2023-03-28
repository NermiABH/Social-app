package store

import (
	"Social-app/internal/model"
	"fmt"
)

type PostRepository struct {
	store *Store
}

func (r *PostRepository) Create(p *model.Post) error {
	q := `INSERT INTO post (author_id, text, media) VALUES ($1, $2, $3) RETURNING id`
	if err := r.store.db.QueryRow(q, p.AuthorID, p.Text, p.Media).Scan(&p.ID); err != nil {
		return err
	}
	return r.GetByID(p)
}

var LimitPost = 15

func (r *PostRepository) GetSeveralByAuthor(page int, authorID *int) (model.Posts, error) {
	q := `SELECT p.id, p.text, p.media, p.created,
		COUNT(c) as comments_count, COUNT(ulp) as likes, COUNT(udp) as dislikes
		FROM post p
		LEFT JOIN comment c ON p.id = c.post_id
		LEFT JOIN user_like_post ulp ON p.id = ulp.post_id
		LEFT JOIN user_dislike_post udp on p.id = udp.post_id 
		%s
		GROUP BY p.id, p.text, p.media, p.created
		ORDER BY p.created DESC
		OFFSET $1 LIMIT $2;`
	if authorID == nil {
		q = fmt.Sprintf(q, "")
	} else {
		q = fmt.Sprintf(q, fmt.Sprintf("WHERE p.authorID=%v", *authorID))
	}
	rows, err := r.store.db.Query(q, page*LimitPost, LimitPost)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pSlice := make(model.Posts, 0)
	for rows.Next() {
		var p model.Post
		if err = rows.Scan(&p.ID, &p.Text, &p.Media, &p.Created, &p.CommentCount,
			&p.Liked, &p.Disliked); err != nil {
			return nil, err
		}
		pSlice = append(pSlice, &p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pSlice, nil
}

func (r *PostRepository) GetByID(p *model.Post) error {
	q := `SELECT p.author_id, u.username, u.userpic, p.text, p.media, p.created,
		COUNT(c) as comments_count, COUNT(ulp) as likes, COUNT(udp) as dislikes
		FROM post p
		LEFT JOIN users u on p.author_id = u.id
		LEFT JOIN comment c on p.id = c.post_id
		LEFT JOIN user_like_post ulp on p.id = ulp.post_id
		LEFT JOIN user_dislike_post udp on p.id = udp.post_id
		WHERE p.id = $1
		GROUP BY p.author_id, u.username, u.userpic, p.text, p.media, p.created`
	return r.store.db.QueryRow(q, p.ID).Scan(&p.AuthorID, &p.AuthorUsername, &p.AuthorUserpic,
		&p.Text, &p.Media, &p.Created, &p.CommentCount, &p.LikeCount, &p.DislikeCount,
	)
}

//func (r *PostRepository) Update(id int, fields []string, values []any) error {
//	var s string
//	for i, field := range fields {
//		value := values[i]
//		if reflect.TypeOf(value).String() == "string" {
//			value = "'" + value.(string) + "'"
//		}
//		s += fmt.Sprintf("%s=%s,", field, value)
//	}
//	s = s[:len(s)-1]
//	q := fmt.Sprintf("UPDATE post SET %s WHERE id = $1", s)
//	_, err := r.store.db.Exec(q, id)
//	return err
//}

func (r *PostRepository) Delete(id int) error {
	q := `DELETE FROM post WHERE id=$1`
	return r.store.db.QueryRow(q, id).Err()
}

func (r *PostRepository) Like(userID, postID int) error {
	q := `INSERT INTO user_like_post VALUES ($1, $2)`
	if err := r.store.db.QueryRow(q, userID, postID).Err(); err != nil {
		return err
	}
	return r.UnDislike(userID, postID)
}

func (r *PostRepository) UnLike(userID, postID int) error {
	q := `DELETE FROM user_like_post WHERE user_id = $1 AND post_id = $2`
	return r.store.db.QueryRow(q, userID, postID).Err()
}

func (r *PostRepository) Dislike(userID, postID int) error {
	q := `INSERT INTO user_dislike_post VALUES ($1, $2)`
	if err := r.store.db.QueryRow(q, userID, postID).Err(); err != nil {
		return err
	}
	return r.UnLike(userID, postID)
}

func (r *PostRepository) UnDislike(userID, postID int) error {
	q := `DELETE FROM user_dislike_post WHERE user_id = $1 AND post_id = $2`
	return r.store.db.QueryRow(q, userID, postID).Err()
}

func (r *PostRepository) Favorite(userID, postID int) error {
	q := `INSERT INTO user_favorite_post VALUES ($1, $2)`
	return r.store.db.QueryRow(q, userID, postID).Err()
}

func (r *PostRepository) UnFavorite(userID, postID int) error {
	q := `DELETE FROM user_favorite_post WHERE user_id = $1 AND post_id = $2`
	return r.store.db.QueryRow(q, userID, postID).Err()
}

func (r *PostRepository) ExistByID(id int) (bool, error) {
	q := `SELECT EXISTS (SELECT * FROM post WHERE id = $1)`
	var exist bool
	err := r.store.db.QueryRow(q, id).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (r *PostRepository) IsOwnerByID(userID, postID int) (bool, error) {
	q := `SELECT EXISTS (SELECT * FROM post WHERE id = $2 AND author_id = $1)`
	var owner bool
	if err := r.store.db.QueryRow(q, userID, postID).Scan(&owner); err != nil {
		return false, err
	}
	return owner, nil
}

func (r *PostRepository) LikedOrDislikedOrFavorited(p *model.Post, uID int) error {
	if uID == 0 {
		return nil
	}
	var boolean bool
	q := "SELECT EXISTS (SELECT * FROM user_favorite_post WHERE user_id=$1 AND post_id=$2)"
	if err := r.store.db.QueryRow(q, p.ID, uID).Scan(&boolean); err != nil {
		return nil
	} else if boolean {
		p.Favorited = true
	}
	q = "SELECT EXISTS (SELECT * FROM user_like_post WHERE user_id=$1 AND post_id=$2)"
	if err := r.store.db.QueryRow(q, p.ID, uID).Scan(&boolean); err != nil {
		return err
	} else if boolean {
		p.Liked = true
		return nil
	}
	q = "SELECT EXISTS (SELECT * FROM user_dislike_post WHERE user_id=$1 AND post_id=$2)"
	if err := r.store.db.QueryRow(q, p.ID, uID).Scan(&boolean); err != nil {
		return err
	} else if boolean {
		p.Disliked = true
	}
	return nil
}

func (r *PostRepository) Count() (int, error) {
	q := "SELECT count(*) FROM post"
	var count int
	if err := r.store.db.QueryRow(q).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
