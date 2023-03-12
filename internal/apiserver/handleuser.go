package apiserver

import (
	"Social-app/internal/jwt"
	"Social-app/internal/model"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"reflect"
	"strconv"
)

func (s *Server) HandleUsersSeveralGet(w http.ResponseWriter, r *http.Request) {
	partUsername := r.URL.Query().Get("username")
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		limit = 10
	}
	users, err := s.store.User().GetSeveralByUsername(partUsername, offset, limit)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.response(w, r, http.StatusOK, map[string][]model.User{"users": users})
}

func (s *Server) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	req := &UserCreate{}
	if code, err := s.correctRequest(r, req); err != nil {
		s.error(w, r, code, err)
		return
	}
	if vucu, err := s.ValidateUserCreateUpdate(req.Username, req.Email); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	} else if len(vucu) != 0 {
		s.error(w, r, http.StatusUnprocessableEntity, vucu)
		return
	}
	u := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsOwn:    true,
	}
	if err := s.store.User().Create(u); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	u.Sanitize()
	s.response(w, r, http.StatusCreated, map[string]*model.User{"user": u})
}

func (s *Server) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	req := &UserLogin{}
	if code, err := s.correctRequest(r, req); err != nil {
		s.error(w, r, code, err)
		return
	}
	var err error
	u := &model.User{}
	if req.Username != "" {
		u.Username = req.Username
		err = s.store.User().GetPasswordByUsername(u)
	} else {
		u.Email = req.Email
		err = s.store.User().GetPasswordByEmail(u)
	}
	if err != nil {
		s.error(w, r, http.StatusNotFound, ErrorIncorrectLoginOrPassword)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(req.Password)); err != nil {
		s.error(w, r, http.StatusNotFound, ErrorIncorrectLoginOrPassword)
		return
	}
	tokens := jwt.New()
	tokens.CreateTokens(u.ID)
	s.response(w, r, http.StatusOK, tokens)
}

func (s *Server) HandleUserRecreateTokens(w http.ResponseWriter, r *http.Request) {
	req := &UserRecreateTokens{}
	if code, err := s.correctRequest(r, req); err != nil {
		s.error(w, r, code, err)
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

func (s *Server) HandleUserGet(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	u := &model.User{ID: owner.ObjectID}
	if err := s.store.User().GetByID(u); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	err := s.store.User().Relation(u, owner.UserID)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
	}
	s.response(w, r, http.StatusOK, map[string]*model.User{"user": u})
}

func (s *Server) HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Err != nil {
		s.error(w, r, owner.Code, owner.Err)
		return
	}
	uID, req := owner.ObjectID, &UserUpdate{}
	if code, err := s.correctRequest(r, req); err != nil {
		s.error(w, r, code, err)
		return
	}
	if vucu, err := s.ValidateUserCreateUpdate(req.Username, req.Email); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	} else if len(vucu) != 0 {
		s.error(w, r, http.StatusUnprocessableEntity, vucu)
		return
	}
	var fields []string
	valuesArray := make([]any, 0)
	values := reflect.ValueOf(*req)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).IsZero() {
			continue
		}
		fields = append(fields, types.Field(i).Tag.Get("json"))
		valuesArray = append(valuesArray, values.Field(i).Interface())
	}
	if err := s.store.User().Update(uID, fields, valuesArray); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandleUserGet(w, r)
}

func (s *Server) HandleUserDelete(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Err != nil {
		s.error(w, r, owner.Code, owner.Err)
		return
	}
	if err := s.store.User().Delete(owner.ObjectID); err != nil {
		s.error(w, r, http.StatusNotFound, err)
		return
	}
	s.response(w, r, http.StatusOK, nil)
}

func (s *Server) HandleUserSubscribe(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, owner.Code, owner.Err)
		return
	}
	if owner.Err == nil {
		s.error(w, r, http.StatusForbidden, errors.New("you can't subscribe to yourself"))
		return
	}
	subscriptionID, subscriberID := owner.UserID, owner.ObjectID
	if err := s.store.User().SubscribeUser(subscriptionID, subscriberID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandleUserGet(w, r)
}

func (s *Server) HandleUserUnSubscribe(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, owner.Code, owner.Err)
		return
	}
	if owner.Err == nil {
		s.error(w, r, owner.Code, errors.New("you can't subscribe to yourself"))
		return
	}
	subscriptionID, subscriberID := owner.UserID, owner.ObjectID
	if err := s.store.User().UnSubscribeUser(subscriptionID, subscriberID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandleUserGet(w, r)
}
