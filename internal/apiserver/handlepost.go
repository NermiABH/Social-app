package apiserver

import (
	"Social-app/internal/model"
	"Social-app/internal/store"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (s *Server) HandlePostsSeveralGet(w http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(ctxAuth).(*CtxAuth)
	var authorID *int
	if queryID := r.URL.Query().Get("author_id"); queryID != "" {
		var err error
		if *authorID, err = strconv.Atoi(queryID); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		if exist, err := s.store.User().ExistByID(*authorID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		} else if !exist {
			s.error(w, r, http.StatusNotFound, errors.New("author does not exist"))
			return
		}
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 0 {
		page = 0
	}
	posts, err := s.store.Post().GetSeveralByAuthor(page, authorID)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	if auth.Err == nil {
		for _, post := range posts {
			if owner, err := s.store.Post().IsOwnerByID(auth.UserID, post.ID); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			} else if owner {
				post.Own = true
			}
			if err := s.store.Post().LikedOrDislikedOrFavorited(post, post.ID); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		}
	}
	var linkNext, linkPrev *string
	if page != 0 {
		linkPrev = returnPointer(fmt.Sprintf("%s?page=%v", r.URL.Path, page-1))
	}
	if postCount, err := s.store.Post().Count(); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	} else if postCount > (page+1)*store.LimitPost {
		linkNext = returnPointer(fmt.Sprintf("%s?page=%v", r.URL.Path, page+1))
	}
	s.response(w, r, http.StatusOK, map[string]any{
		"links": map[string]any{
			"self": r.RequestURI,
			"next": linkNext,
			"prev": linkPrev,
		},
		"data": posts.ConvertMap(),
	})
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
		Media:    req.Media,
	}
	if err := s.store.Post().Create(p); err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	p.Own = true
	s.response(w, r, http.StatusCreated, map[string]any{
		"links": map[string]string{
			"self": r.RequestURI,
		},
		"data": p.ConvertMap(),
	})
}

func (s *Server) HandlePostGet(w http.ResponseWriter, r *http.Request) {
	owner := r.Context().Value(ctxOwner).(*CtxOwner)
	p := &model.Post{ID: owner.ObjectID}
	if err := s.store.Post().GetByID(p); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	if owner.Err == nil {
		p.Own = true
	} else if err := s.store.Post().LikedOrDislikedOrFavorited(p, owner.UserID); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.response(w, r, http.StatusCreated, map[string]any{
		"links": map[string]string{
			"self": r.RequestURI,
		},
		"data": p.ConvertMap(),
	})
}

//func (s *Server) HandlePostUpdate(w http.ResponseWriter, r *http.Request) {
//	owner := r.Context().Value(ctxOwner).(*CtxOwner)
//	if owner.Err != nil {
//		s.error(w, r, owner.Code, owner.Err)
//		return
//	}
//	pID, req := owner.ObjectID, &PostCreateUpdate{}
//	if code, err := s.correctRequest(r, req); err != nil {
//		s.error(w, r, code, err)
//		return
//	} else if reflect.DeepEqual(*req, PostCreateUpdate{}) {
//		s.error(w, r, http.StatusUnprocessableEntity, errors.New("empty struct"))
//		return
//	}
//	var fields []string
//	valuesArray := make([]any, 0)
//	values := reflect.ValueOf(*req)
//	types := values.Type()
//	for i := 0; i < values.NumField(); i++ {
//		if values.Field(i).IsZero() {
//			fmt.Println(values.Field(i).IsZero())
//			continue
//		}
//		fields = append(fields, types.Field(i).Tag.Get("json"))
//		valuesArray = append(valuesArray, values.Field(i).Interface())
//	}
//	if err := s.store.Post().Update(pID, fields, valuesArray); err != nil {
//		s.error(w, r, http.StatusInternalServerError, err)
//		return
//	}
//	s.HandlePostGet(w, r)
//}

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
