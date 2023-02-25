package apiserver

import (
	"Social-app/internal/model"
	"encoding/json"
	"net/http"
)

func (s *server) HandleCreateUser() http.HandlerFunc {
	type request struct {
		Email    string
		Password string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := &model.User{
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
