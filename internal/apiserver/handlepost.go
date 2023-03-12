package apiserver

import (
	"Social-app/internal/model"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

func (s *Server) HandlePostsSeveralGet(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("author_id")
	uID, err := strconv.Atoi(tag)
	if err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, errors.New("author_id_tag must be not empty and natural number"))
		return
	}
	if exist, err := s.store.User().IsExist(uID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	} else if !exist {
		fmt.Println(3)
		s.error(w, r, http.StatusNotFound, errors.New("user is not exist"))
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		limit = 10
	}
	posts, err := s.store.Post().GetSeveralByAuthor(uID, offset, limit)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.response(w, r, http.StatusOK, map[string][]model.Post{"posts": posts})
}

func (s *Server) HandlePostCreate(w http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(ctxAuth).(*CtxAuth)
	if auth.Err != nil {
		s.error(w, r, http.StatusUnauthorized, auth.Err)
		return
	}
	pID, req := auth.UserID, &PostCreateUpdate{}
	if code, err := s.correctRequest(r, req); err != nil {
		s.error(w, r, code, err)
		return
	}
	p := &model.Post{
		AuthorID: pID,
		Text:     req.Text,
		Object:   req.Object,
		IsOwn:    true,
	}
	if err := s.store.Post().Create(p); err != nil {
		fmt.Println(2)
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	s.response(w, r, http.StatusCreated, map[string]*model.Post{"post": p})
}

func (s *Server) HandlePostGet(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	p := &model.Post{ID: owner.ObjectID}
	if err := s.store.Post().GetByID(p); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	if owner.Err == nil {
		p.IsOwn = true
	} else if err := s.store.Post().LikedOrDislikedOrFavorited(p, owner.UserID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.response(w, r, http.StatusOK, map[string]*model.Post{"post": p})
}

func (s *Server) HandlePostUpdate(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Err != nil {
		s.error(w, r, owner.Code, owner.Err)
		return
	}
	pID, req := owner.ObjectID, &PostCreateUpdate{}
	if code, err := s.correctRequest(r, req); err != nil {
		s.error(w, r, code, err)
		return
	} else if reflect.DeepEqual(*req, PostCreateUpdate{}) {
		s.error(w, r, http.StatusUnprocessableEntity, errors.New("empty struct"))
		return
	}
	var fields []string
	valuesArray := make([]any, 0)
	values := reflect.ValueOf(*req)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).IsZero() {
			fmt.Println(values.Field(i).IsZero())
			continue
		}
		fields = append(fields, types.Field(i).Tag.Get("json"))
		valuesArray = append(valuesArray, values.Field(i).Interface())
	}
	if err := s.store.Post().Update(pID, fields, valuesArray); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandlePostGet(w, r)
}

func (s *Server) HandlePostDelete(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Err != nil {
		s.error(w, r, owner.Code, owner.Err)
		return
	}
	pID := owner.ObjectID
	if err := s.store.Post().Delete(pID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.response(w, r, http.StatusOK, nil)
}

func (s *Server) HandlePostLike(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	uID, pID := owner.UserID, owner.ObjectID
	if err := s.store.Post().Like(uID, pID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandlePostGet(w, r)
}

func (s *Server) HandlePostUnLike(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Err != nil {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	uID, pID := owner.UserID, owner.ObjectID
	if err := s.store.Post().UnLike(uID, pID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandleUserGet(w, r)
}

func (s *Server) HandlePostDislike(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	uID, pID := owner.UserID, owner.ObjectID
	if err := s.store.Post().Dislike(uID, pID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandlePostGet(w, r)
}

func (s *Server) HandlePostUnDislike(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	uID, pID := owner.UserID, owner.ObjectID
	if err := s.store.Post().UnDislike(uID, pID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandlePostGet(w, r)
}

func (s *Server) HandlePostFavorite(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	uID, pID := owner.UserID, owner.ObjectID
	if err := s.store.Post().Favorite(uID, pID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandlePostGet(w, r)
}

func (s *Server) HandlePostUnFavorite(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	uID, pID := owner.UserID, owner.ObjectID
	if err := s.store.Post().UnFavorite(uID, pID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandlePostGet(w, r)
}
