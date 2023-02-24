package apiserver

import (
	"Social-app/internal/store"
	"github.com/gorilla/mux"
	"net/http"
)

type server struct {
	router *mux.Router
	store  *store.Store
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store *store.Store) *server {
	return &server{
		router: mux.NewRouter(),
		store:  store,
	}
}
