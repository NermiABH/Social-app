package apiserver

import (
	"Social-app/internal/jwt"
	"Social-app/internal/store"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const ctxKeyRequestID ctxKey = iota

var (
	ctxAuth *ctxAuthStruct
)

type ctxKey int8
type ctxAuthStruct struct {
	UserID int
	Err    error
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
	s.router.Use(s.isAuth)
	s.router.HandleFunc("/users", s.HandleUsersGet).Methods("GET")
	s.router.HandleFunc("/users", s.HandleUserCreate).Methods("POST")
	s.router.HandleFunc("/users/login", s.HandleUserLogin).Methods("POST")
	s.router.HandleFunc("/users/refresh", s.HandleUserRecreateTokens).Methods("POST")
	s.router.HandleFunc("/post", s.HandlePostCreate).Methods("POST")

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
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxAuth, &ctxAuthStruct{Err: errors.New(access)})))
			return
		}
		access = strings.TrimSpace(splitToken[1])
		payload, err := jwt.Verify(access)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxAuth, &ctxAuthStruct{Err: err})))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxAuth, &ctxAuthStruct{UserID: payload.UserID})))
	})
}

func (s *Server) response(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err any) {
	if reflect.TypeOf(err).String() == "*errors.errorString" {
		s.response(w, r, code, map[string]string{"errors": err.(error).Error()})
		return
	}
	s.response(w, r, code, map[string]any{"errors": err})
}
