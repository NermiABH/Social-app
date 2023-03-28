package apiserver

import (
	"Social-app/internal/jwt"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ctxKeyRequestID ctxKey = iota
	ctxObject       ctxKey = iota
)

var (
	ctxAuth                       *CtxAuth
	ctxOwner                      *CtxOwner
	ErrorIncorrectLoginOrPassword = errors.New("incorrect password or login")
	ErrorPermissionDenied         = errors.New("permission denied")
)

type ctxKey int8

type CtxAuth struct {
	UserID int
	Err    error
}
type CtxOwner struct {
	UserID   int
	ObjectID int
	Code     int
	Err      error
}

func (s *Server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_add": r.RemoteAddr,
			"request_id": r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)
		start := time.Now()
		next.ServeHTTP(w, r)
		logrus.Infof("completed %s", time.Now().Sub(start))
	})
}

func (s *Server) isAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access := r.Header.Get("Authorization")
		splitToken := strings.Split(access, "Bearer")
		if len(splitToken) != 2 {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(),
				ctxAuth, &CtxAuth{Err: errors.New("access is empty")})))
			return
		}
		access = strings.TrimSpace(splitToken[1])
		payload, err := jwt.Verify(access)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxAuth, &CtxAuth{Err: err})))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxAuth, &CtxAuth{UserID: payload.UserID})))
	})
}

func (s *Server) Exist(Var string, table string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, _ := strconv.Atoi(mux.Vars(r)[Var])
			var (
				exist bool
				err   error
			)
			switch table {
			case "user":
				exist, err = s.store.User().ExistByID(id)
			case "post":
				exist, err = s.store.Post().ExistByID(id)
			case "comment":
				exist, err = s.store.Comment().ExistByID(id)
			default:
				s.error(w, r, http.StatusInternalServerError, nil)
				return
			}
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			} else if !exist {
				s.error(w, r, http.StatusNotFound, errors.New(fmt.Sprintf("%s is not exist", table)))
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxObject, id)))
		})
	}
}
func (s *Server) IsOwner(table string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(ctxObject).(int)
			auth := r.Context().Value(ctxAuth).(*CtxAuth)
			if auth.Err != nil {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(),
					ctxOwner, &CtxOwner{ObjectID: id, Code: http.StatusUnauthorized, Err: auth.Err})))
				return
			}
			var (
				owner bool
				err   error
			)
			switch table {
			case "user":
				owner = auth.UserID == id
			case "post":
				owner, err = s.store.Post().IsOwnerByID(auth.UserID, id)
			case "comment":
				owner, err = s.store.Comment().IsOwnerByID(auth.UserID, id)
			}
			if err != nil {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(),
					ctxOwner, &CtxOwner{UserID: auth.UserID, ObjectID: id, Code: http.StatusInternalServerError, Err: err})))
				return
			} else if !owner {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(),
					ctxOwner, &CtxOwner{UserID: auth.UserID, ObjectID: id, Code: http.StatusForbidden, Err: ErrorPermissionDenied})))
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(),
				ctxOwner, &CtxOwner{UserID: auth.UserID, ObjectID: id})))
		})
	}
}
