package apiserver

import (
	"Social-app/internal/jwt"
	"Social-app/internal/store"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
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

type Server struct {
	router *mux.Router
	store  *store.Store
	logger *logrus.Logger
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store *store.Store) *Server {
	s := &Server{
		router: mux.NewRouter(),
		store:  store,
		logger: logrus.New(),
	}
	s.configureRouter()
	return s
}
func (s *Server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedMethods([]string{"*"})))
	role := s.router.PathPrefix("/api/{role:admin|public}").Subrouter()
	s.router.Use(s.isAuth)
	role.HandleFunc("/users", s.HandleUsersSeveralGet).Methods("GET")
	role.HandleFunc("/users", s.HandleUserCreate).Methods("POST")
	role.HandleFunc("/users/login", s.HandleUserLogin).Methods("POST")
	role.HandleFunc("/users/refresh", s.HandleUserRecreateTokens).Methods("POST")

	user := role.PathPrefix("/user/{id:[0-9]+}").Subrouter()
	user.Use(s.Exist("id", "user"))
	user.Use(s.IsOwner("user"))
	user.HandleFunc("", s.HandleUserGet).Methods("GET")
	user.HandleFunc("", s.HandleUserUpdate).Methods("PATCH")
	user.HandleFunc("", s.HandleUserDelete).Methods("DELETE")
	user.HandleFunc("/subscribe", s.HandleUserSubscribe).Methods("POST")
	user.HandleFunc("/unsubscribe", s.HandleUserUnSubscribe).Methods("POST")

	role.HandleFunc("/posts", s.HandlePostsSeveralGet).Methods("GET")
	role.HandleFunc("/posts", s.HandlePostCreate).Methods("POST")

	post := role.PathPrefix("/post/{id:[0-9]+}").Subrouter()
	post.Use(s.Exist("id", "post"), s.IsOwner("post"))
	post.HandleFunc("", s.HandlePostGet).Methods("GET")
	post.HandleFunc("", s.HandlePostUpdate).Methods("PATCH")
	post.HandleFunc("", s.HandlePostDelete).Methods("DELETE")

	post.HandleFunc("/like", s.HandlePostLike).Methods("POST")
	post.HandleFunc("/unlike", s.HandlePostUnLike).Methods("DELETE")
	post.HandleFunc("/dislike", s.HandlePostDislike).Methods("POST")
	post.HandleFunc("/undislike", s.HandlePostUnDislike).Methods("DELETE")
	post.HandleFunc("/favorite", s.HandlePostFavorite).Methods("POST")
	post.HandleFunc("/unfavorite", s.HandlePostUnFavorite).Methods("DELETE")

	post.HandleFunc("/comments", s.HandleCommentsSeveralGet).Methods("GET")
	post.HandleFunc("/comments", s.HandleCommentCreate).Methods("POST")

	comment := post.PathPrefix("/comment/{c_id:[0-9]+}").Subrouter()
	comment.Use(s.Exist("c_id", "comment"), s.IsOwner("comment"))
	comment.HandleFunc("", s.HandleCommentGet).Methods("GET")
	comment.HandleFunc("", s.HandleCommentUpdate).Methods("PATCH")
	comment.HandleFunc("", s.HandleCommentDelete).Methods("DELETE")

	comment.HandleFunc("/like", s.HandleCommentLike).Methods("POST")
	comment.HandleFunc("/unlike", s.HandleCommentUnLike).Methods("DELETE")
	comment.HandleFunc("/dislike", s.HandleCommentDislike).Methods("POST")
	comment.HandleFunc("/undislike", s.HandleCommentUnDislike).Methods("DELETE")
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
				exist, err = s.store.User().IsExist(id)
			case "post":
				exist, err = s.store.Post().IsExist(id)
			case "comment":
				exist, err = s.store.Comment().IsExist(id)
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
				owner, err = s.store.Post().IsOwner(auth.UserID, id)
			case "comment":
				owner, err = s.store.Comment().IsOwner(auth.UserID, id)
			default:
				s.error(w, r, http.StatusInternalServerError, nil)
				return
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
func (s *Server) response(w http.ResponseWriter, _ *http.Request, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err any) {
	st := reflect.TypeOf(err)
	if _, ok := st.MethodByName("Error"); ok {
		s.response(w, r, code, map[string]string{"error": err.(error).Error()})
		return
	}
	s.response(w, r, code, map[string]any{"errors": err})
}

func (s *Server) correctRequest(r *http.Request, request any) (int, any) {
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return http.StatusBadRequest, err
	}
	if err := Validate(request); err != nil {
		return http.StatusUnprocessableEntity, err
	}
	return 0, nil
}
