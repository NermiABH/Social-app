package store

import (
	"Social-app/internal/model"
	"database/sql"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	q := `INSERT INTO users (username, email, password) VALUES($1, $2, $3) RETURNING id, created`
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	return r.store.db.QueryRow(q,
		u.Username, u.Email, u.EncryptedPassword).Scan(&u.ID, &u.Created)
}

var LimitUser = 15

func (r *UserRepository) GetSeveralByUsername(part string, page int) (model.Users, error) {
	q := `SELECT id, username, userpic FROM users WHERE username LIKE $1 ORDER BY created DESC OFFSET $2 LIMIT $3 `
	rows, err := r.store.db.Query(q, "%"+part+"%", page*LimitUser, LimitUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	uSlice := make(model.Users, 0)
	for rows.Next() {
		var u model.User
		var userpic sql.NullString
		if err = rows.Scan(&u.ID, &u.Username, &userpic); err != nil {
			return nil, err
		}
		u.Userpic = userpic.String
		uSlice = append(uSlice, &u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return uSlice, nil
}

func (r *UserRepository) GetPasswordByUsername(u *model.User) error {
	q := `SELECT id, email, password FROM users WHERE username = $1`
	return r.store.db.QueryRow(q,
		u.Username).Scan(&u.ID, &u.Password, &u.EncryptedPassword)
}

func (r *UserRepository) GetPasswordByEmail(u *model.User) error {
	q := `SELECT id, username, password FROM users WHERE email = $1`
	return r.store.db.QueryRow(q,
		u.Email).Scan(&u.ID, &u.Password, &u.EncryptedPassword)
}

func (r *UserRepository) GetByID(u *model.User) error {
	q := `SELECT u.username, u.email, u.userpic, u.name, u.surname, u.born, u.created,
   		count(ss1) as subscription_count, count(ss2) as subscriber_count,
   		count(p) as post_count, count(ufp) as favorite_posts
		FROM users u
		LEFT JOIN subscription_subscriber ss1 on u.id = ss1.subscription_id
		LEFT JOIN subscription_subscriber ss2 on u.id = ss2.subscriber_id
		LEFT JOIN post p on u.id = p.author_id
		LEFT JOIN user_favorite_post ufp on p.id = ufp.user_id
		WHERE u.id = $1
		GROUP BY u.username, u.email, u.userpic, u.name, u.surname, u.born, u.created`
	if err := r.store.db.QueryRow(q,
		u.ID).Scan(&u.Username, &u.Email, &u.Userpic, &u.Name, &u.Surname, &u.Born, &u.Created,
		&u.SubscriptionCount, &u.SubscriberCount, &u.PostCount, &u.FavoritePostCount); err != nil {
		return err
	}
	return nil
}

//func (r *UserRepository) Update(id int, fields []string, values []any) error {
//	var s string
//	for i, field := range fields {
//		value := values[i]
//		if reflect.TypeOf(value).String() == "string" {
//			value = "'" + value.(string) + "'"
//		}
//		s += fmt.Sprintf("%s=%s,", field, value)
//	}
//	s = s[:len(s)-1]
//	q := fmt.Sprintf("UPDATE users SET %s WHERE id = $1", s)
//	_, err := r.store.db.Exec(q, id)
//	return err
//}

func (r *UserRepository) Delete(id int) error {
	q := `DELETE FROM users WHERE id = $1`
	return r.store.db.QueryRow(q, id).Err()
}

func (r *UserRepository) SubscribeUser(subscriptionID, subscriberID int) error {
	q := `INSERT INTO subscription_subscriber VALUES ($1, $2)`
	return r.store.db.QueryRow(q, subscriptionID, subscriberID).Err()
}

func (r *UserRepository) UnSubscribeUser(subscriptionID, subscriberID int) error {
	q := `DELETE FROM subscription_subscriber WHERE subscription_id=$1 AND subscriber_id = $2`
	return r.store.db.QueryRow(q, subscriptionID, subscriberID).Err()
}

func (r *UserRepository) ExistByID(id int) (bool, error) {
	q := `SELECT EXISTS (SELECT * FROM users WHERE id = $1)`
	var exist bool
	if err := r.store.db.QueryRow(q, id).Scan(&exist); err != nil {
		return false, err
	}
	return exist, nil
}

func (r *UserRepository) ExistByUsername(username string) (bool, error) {
	q := `SELECT EXISTS (SELECT * FROM users WHERE username = $1)`
	var exist bool
	if err := r.store.db.QueryRow(q, username).Scan(&exist); err != nil {
		return false, err
	}
	return exist, nil
}

func (r *UserRepository) ExistByEmail(email string) (bool, error) {
	q := `SELECT EXISTS (SELECT * FROM users WHERE email = $1)`
	var exist bool
	if err := r.store.db.QueryRow(q, email).Scan(&exist); err != nil {
		return false, err
	}
	return exist, nil
}

func (r *UserRepository) SubscribedOrSubscriber(u *model.User, anotherID int) error {
	q := `SELECT EXISTS (SELECT * FROM subscription_subscriber WHERE subscription_id = $1 AND subscriber_id = $2)`
	var boolean bool
	if err := r.store.db.QueryRow(q, u.ID, anotherID).Scan(&boolean); err != nil {
		return err
	} else if boolean {
		u.Subscribed = true
	}
	if err := r.store.db.QueryRow(q, anotherID, u.ID).Scan(&boolean); err != nil {
		return err
	} else if boolean {
		u.Subscriber = true
	}
	return nil
}

func (r *UserRepository) Count() (int, error) {
	q := "SELECT count(*) FROM users"
	var count int
	if err := r.store.db.QueryRow(q).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
