package apiserver

import (
	"Social-app/internal/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

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
	s.router.HandleFunc("/users", s.HandleUserCreate()).Methods("POST")
	s.router.HandleFunc("/users/login", s.HandleUserLogin()).Methods("POST")
	s.router.HandleFunc("/users/refresh", s.HandleUserRecreateTokens()).Methods("POST")
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
