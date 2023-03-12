package apiserver

type UserCreate struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required_if=Email ''"`
	Email    string `json:"email" validate:"required_if=Username ''"`
	Password string `json:"password" validate:"required"`
}

type UserRecreateTokens struct {
	Refresh string `json:"refresh" validate:"required"`
}

type UserUpdate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Userpic  string `json:"userpic"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type PostCreateUpdate struct {
	Text   string `json:"text"`
	Object string `json:"object"`
}

type CommentCreateUpdate struct {
	ParentID *int   `json:"parent_id"`
	Text     string `json:"text" validate:"required"`
}
