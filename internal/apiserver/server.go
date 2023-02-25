package apiserver

import (
	"Social-app/internal/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type server struct {
	router *mux.Router
	store  *store.Store
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store *store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
	}
	s.configureRouter()
	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/user", s.HandleCreateUser()).Methods("POST")
}

func (s *server) response(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.response(w, r, code, map[string]any{"errors": err.Error()})
}

func (s *server) errorAny(w http.ResponseWriter, r *http.Request, code int, err any) {
	s.response(w, r, code, map[string]any{"errors": err})
}
