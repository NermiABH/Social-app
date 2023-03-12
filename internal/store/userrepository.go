package store

import (
	"Social-app/internal/model"
	"database/sql"
	"fmt"
	"reflect"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	q := `INSERT INTO users (username, email, encrypted_password) 
		VALUES($1, $2, $3)
		RETURNING id`
	err := u.BeforeCreate()
	if err != nil {
		return err
	}
	return r.store.db.QueryRow(q,
		u.Username, u.Email, u.EncryptedPassword).Scan(&u.ID)
}

func (r *UserRepository) GetSeveralByUsername(part string, offset, limit int) ([]model.User, error) {
	q := `SELECT id, username, userpic
		FROM users 
		WHERE username LIKE $1
		OFFSET $2 LIMIT $3`
	part = "%" + part + "%"
	rows, err := r.store.db.Query(q, part, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	uSlice := make([]model.User, 0)
	for rows.Next() {
		var u model.User
		var userpic sql.NullString
		if err = rows.Scan(&u.ID, &u.Username, &userpic); err != nil {
			return nil, err
		}
		u.Userpic = userpic.String
		uSlice = append(uSlice, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return uSlice, nil
}

func (r *UserRepository) GetPasswordByUsername(u *model.User) error {
	q := `SELECT id, email, encrypted_password
		FROM users
		WHERE username = $1`
	return r.store.db.QueryRow(q,
		u.Username).Scan(&u.ID, &u.Password, &u.EncryptedPassword)
}
func (r *UserRepository) GetPasswordByEmail(u *model.User) error {
	q := `SELECT id, username, encrypted_password
		FROM users
		WHERE email = $1`
	return r.store.db.QueryRow(q,
		u.Email).Scan(&u.ID, &u.Password, &u.EncryptedPassword)
}

func (r *UserRepository) GetByID(u *model.User) error {
	q := `SELECT u.username, u.email, u.userpic, u.name, u.surname, u.date_of_birth, u.date_of_create,
    	count(ss1) as subscriptions_count,
    	count(ss2) as subscribers_count,
    	count(p) as posts_count,
    	count(ufp) as favorite_posts
		FROM users u
		LEFT JOIN subscription_subscriber ss1 on u.id = ss1.subscription_id
		LEFT JOIN subscription_subscriber ss2 on u.id = ss2.subscriber_id
		LEFT JOIN post p on u.id = p.author_id
		LEFT JOIN user_favorite_post ufp on p.id = ufp.user_id
		WHERE u.id = $1
		GROUP BY u.username, u.email, u.userpic, u.name, u.surname, u.date_of_birth, u.date_of_create`
	var (
		userpic   sql.NullString
		name      sql.NullString
		surname   sql.NullString
		dateBirth sql.NullString
	)
	if err := r.store.db.QueryRow(q,
		u.ID).Scan(&u.Username, &u.Email, &userpic, &name, &surname, &dateBirth, &u.DateCreation,
		&u.SubscriptionsCount, &u.SubscribersCount, &u.PostsCount, &u.FavoritesPosts); err != nil {
		return err
	}
	u.Userpic = userpic.String
	u.Name = name.String
	u.Surname = name.String
	u.DateBirth = dateBirth.String
	return nil
}

func (r *UserRepository) Update(id int, fields []string, values []any) error {
	var s string
	for i, field := range fields {
		value := values[i]
		if reflect.TypeOf(value).String() == "string" {
			value = "'" + value.(string) + "'"
		}
		s += fmt.Sprintf("%s=%s,", field, value)
	}
	s = s[:len(s)-1]
	q := fmt.Sprintf("UPDATE users SET %s WHERE id = $1", s)
	_, err := r.store.db.Exec(q, id)
	return err
}

func (r *UserRepository) Delete(id int) error {
	q := `DELETE FROM users WHERE id = $1`
	return r.store.db.QueryRow(q, id).Err()
}

func (r *UserRepository) SubscribeUser(subscriptionId, subscriberId int) error {
	q := `INSERT INTO subscription_subscriber
		VALUES ($1, $2)`
	return r.store.db.QueryRow(q, subscriptionId, subscriberId).Err()
}

func (r *UserRepository) UnSubscribeUser(subscriptionId, subscriberId int) error {
	q := `DELETE FROM subscription_subscriber
		WHERE subscription_id=$1 AND subscriber_id = $2`
	return r.store.db.QueryRow(q, subscriptionId, subscriberId).Err()
}

func (r *UserRepository) IsExist(id int) (bool, error) {
	q := `SELECT EXISTS (SELECT * FROM users WHERE id = $1)`
	var exist bool
	err := r.store.db.QueryRow(q, id).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (r *UserRepository) Relation(u *model.User, anID int) error {
	q := `SELECT EXISTS (SELECT * FROM subscription_subscriber WHERE subscription_id = $1 AND subscriber_id = $2)`
	if u.ID == 0 {
		return nil
	} else if u.ID == anID {
		u.IsOwn = true
		return nil
	}
	var fBool, sBool bool
	if err := r.store.db.QueryRow(q, anID, u.ID).Scan(&fBool); err != nil {
		return err
	} else if fBool {
		u.IsSubscription = true
	}
	if err := r.store.db.QueryRow(q, u.ID, anID).Scan(&sBool); err != nil {
		return err
	} else if sBool {
		u.IsSubscriber = true
	}
	return nil
}

func (r *UserRepository) IsExistByUsername(username string) (bool, error) {
	q := `SELECT EXISTS (SELECT * FROM users WHERE username = $1)`
	var exist bool
	err := r.store.db.QueryRow(q, username).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (r *UserRepository) IsExistByEmail(email string) (bool, error) {
	q := `SELECT EXISTS (SELECT * FROM users WHERE email = $1)`
	var exist bool
	err := r.store.db.QueryRow(q, email).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}
