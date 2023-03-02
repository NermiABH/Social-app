package apiserver

import (
	"Social-app/internal/jwt"
	"Social-app/internal/model"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var (
	ErrorIncorectLoginOrPassword = errors.New("incorrect password or login")
)

func (s *Server) HandleUsersGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.User().GetUsers()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.response(w, r, http.StatusOK, map[string][]model.User{"users": users})
	}
}

func (s *Server) HandleUserCreate() http.HandlerFunc {
	type request struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := Validate(req); err != nil {
			s.errorOfAny(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u := &model.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.User().CreateUser(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitize()
		s.response(w, r, http.StatusCreated, u)
	}
}

func (s *Server) HandleUserLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username" validate:"required_if=Email ''"`
		Email    string `json:"email" validate:"required_if=Username ''"`
		Password string `json:"password" validate:"required"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := Validate(req); err != nil {
			s.errorOfAny(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		var (
			err error
			u   *model.User
		)
		if req.Username != "" {
			u, err = s.store.User().FindByUsername(req.Email)
		} else {
			u, err = s.store.User().FindByEmail(req.Email)
		}
		if err != nil {
			s.error(w, r, http.StatusNotFound, ErrorIncorectLoginOrPassword)
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(req.Password)); err != nil {
			s.error(w, r, http.StatusNotFound, ErrorIncorectLoginOrPassword)
			return
		}
		tokens := jwt.New()
		tokens.CreateTokens(u.ID)
		s.response(w, r, http.StatusOK, tokens)
	}
}

func (s *Server) HandleUserRecreateTokens() http.HandlerFunc {
	type request struct {
		Refresh string `json:"refresh" validate:"required"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := Validate(req); err != nil {
			s.errorOfAny(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		tokens := jwt.Tokens{}
		tokens.Refresh = req.Refresh
		err := tokens.RecreateTokens()
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		s.response(w, r, http.StatusCreated, tokens)
	}
}

func (s *Server) HandlePostCreate() http.HandlerFunc {
	type request struct {
		Text   string `json:"text" validate:"required_if=Object ''"`
		Object string `json:"object" validate:"required_if=Text ''"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxAuth).(*ctxAuthStruct)
		if user.Err != nil {
			s.error(w, r, http.StatusUnauthorized, user.Err)
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := Validate(req); err != nil {
			s.errorOfAny(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		p := &model.Post{
			AuthorID: user.UserID,
			Text:     req.Text,
			Object:   req.Text,
		}
		if err := s.store.Post().CreatePost(p); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.response(w, r, http.StatusCreated, p)
	}
}
