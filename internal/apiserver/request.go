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

type PostCreate struct {
	Text   string `json:"text" validate:"required_if=Object ''"`
	Object string `json:"object" validate:"required_if=Text ''"`
}
