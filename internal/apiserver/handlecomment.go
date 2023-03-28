package apiserver

import (
	"Social-app/internal/model"
	"fmt"
	"net/http"
)

//func (s *Server) HandleCommentsSeveralGet(w http.ResponseWriter, r *http.Request) {
//	pOwner := r.Context().Value(ctxOwner).(*CtxOwner)
//	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
//	if err != nil || offset < 0 {
//		offset = 0
//	}
//	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
//	if err != nil || limit < 0 {
//		limit = 10
//	}
//	cSlice, err := s.store.Comment().GetSeveral(pOwner.ObjectID, offset, limit)
//	if err != nil {
//		s.error(w, r, http.StatusInternalServerError, err)
//		return
//	}
//	uID := pOwner.UserID
//	for _, c := range cSlice {
//		if c.AuthorID == uID {
//			c.IsOwn = true
//		}
//		if err = s.store.Comment().LikedOrDisliked(c, uID); err != nil {
//			s.error(w, r, http.StatusInternalServerError, err)
//			return
//		}
//	}
//	s.response(w, r, http.StatusOK, map[string][]*model.Comment{"comments": cSlice})
//}

func (s *Server) HandleCommentGet(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	fmt.Println(3)
	c := &model.Comment{ID: owner.ObjectID}
	if err := s.store.Comment().GetByID(c); err != nil {
		fmt.Println(err, 3)
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	if owner.Err == nil {
		c.Own = true
	}
	if err := s.store.Comment().LikedOrDisliked(c, owner.UserID); err != nil {
		fmt.Println(err)
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.response(w, r, http.StatusOK, map[string]*model.Comment{"comment": c})
}

func (s *Server) HandleCommentCreate(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	req := &CommentCreateUpdate{}
	if code, err := s.correctRequest(r, req); err != nil {
		s.error(w, r, code, err)
		return
	}
	c := &model.Comment{
		PostID:   owner.ObjectID,
		AuthorID: owner.UserID,
		Text:     req.Text,
		Own:      true,
	}
	if err := s.store.Comment().Create(c); err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	s.response(w, r, http.StatusCreated, map[string]*model.Comment{"comment": c})
}

func (s *Server) HandleCommentUpdate(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	fmt.Println(1)
	if owner.Err != nil {
		s.error(w, r, owner.Code, owner.Err)
		return
	}
	fmt.Println(2)
	cID, req := owner.ObjectID, &CommentCreateUpdate{}
	if code, err := s.correctRequest(r, req); err != nil {
		fmt.Println(err, 1)
		s.error(w, r, code, err)
		return
	}
	if err := s.store.Comment().Update(cID, req.Text); err != nil {
		fmt.Println(err, 2)
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandleCommentGet(w, r)
}

func (s *Server) HandleCommentDelete(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Err != nil {
		s.error(w, r, owner.Code, owner.Err)
		return
	}
	id := owner.ObjectID
	if err := s.store.Comment().Delete(id); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.response(w, r, http.StatusOK, nil)
}

func (s *Server) HandleCommentLike(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	uID, cID := owner.UserID, owner.ObjectID
	if err := s.store.Comment().Like(uID, cID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandleCommentGet(w, r)
}

func (s *Server) HandleCommentUnLike(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	uID, cID := owner.UserID, owner.ObjectID
	if err := s.store.Comment().UnLike(uID, cID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandleCommentGet(w, r)
}

func (s *Server) HandleCommentDislike(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	uID, cID := owner.UserID, owner.ObjectID
	if err := s.store.Comment().Dislike(uID, cID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandleCommentGet(w, r)
}

func (s *Server) HandleCommentUnDislike(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	if owner.Code == http.StatusUnauthorized {
		s.error(w, r, http.StatusUnauthorized, owner.Err)
		return
	}
	uID, cID := owner.UserID, owner.ObjectID
	if err := s.store.Comment().UnDislike(uID, cID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.HandleCommentGet(w, r)
}
