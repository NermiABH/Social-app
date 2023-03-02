package apiserver

import (
	"Social-app/internal/jwt"
	"Social-app/internal/store"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var (
	ctxAuth *ctxAuthStruct
)

type ctxAuthStruct struct {
	UserID int
	Err    error
}

type Server struct {
	router *mux.Router
	store  *store.Store
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store *store.Store) *Server {
	s := &Server{
		router: mux.NewRouter(),
		store:  store,
	}
	s.configureRouter()
	return s
}

func (s *Server) configureRouter() {
	s.router.Use(s.isAuth)
	s.router.HandleFunc("/users", s.HandleUsersGet()).Methods("GET")
	s.router.HandleFunc("/users", s.HandleUserCreate()).Methods("POST")
	s.router.HandleFunc("/users/login", s.HandleUserLogin()).Methods("POST")
	s.router.HandleFunc("/users/refresh", s.HandleUserRecreateTokens()).Methods("POST")
	s.router.HandleFunc("/post", s.HandlePostCreate()).Methods("POST")

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

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.response(w, r, code, map[string]string{"errors": err.Error()})
}
func (s *Server) errorOfAny(w http.ResponseWriter, r *http.Request, code int, err any) {
	s.response(w, r, code, map[string]any{"errors": err})
}
