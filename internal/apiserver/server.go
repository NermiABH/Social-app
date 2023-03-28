package apiserver

import (
	"Social-app/internal/store"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

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
	//user.HandleFunc("", s.HandleUserUpdate).Methods("PATCH")
	user.HandleFunc("", s.HandleUserDelete).Methods("DELETE")
	user.HandleFunc("/subscribe", s.HandleUserSubscribe).Methods("POST")
	user.HandleFunc("/subscribe", s.HandleUserUnSubscribe).Methods("DELETE")

	role.HandleFunc("/posts", s.HandlePostsSeveralGet).Methods("GET")
	role.HandleFunc("/posts", s.HandlePostCreate).Methods("POST")

	post := role.PathPrefix("/post/{id:[0-9]+}").Subrouter()
	post.Use(s.Exist("id", "post"), s.IsOwner("post"))
	post.HandleFunc("", s.HandlePostGet).Methods("GET")
	//post.HandleFunc("", s.HandlePostUpdate).Methods("PATCH")
	post.HandleFunc("", s.HandlePostDelete).Methods("DELETE")

	post.HandleFunc("/like", s.HandlePostLike).Methods("POST")
	post.HandleFunc("/like", s.HandlePostUnLike).Methods("DELETE")
	post.HandleFunc("/dislike", s.HandlePostDislike).Methods("POST")
	post.HandleFunc("/dislike", s.HandlePostUnDislike).Methods("DELETE")
	post.HandleFunc("/favorite", s.HandlePostFavorite).Methods("POST")
	post.HandleFunc("/favorite", s.HandlePostUnFavorite).Methods("DELETE")

	//post.HandleFunc("/comments", s.HandleCommentsSeveralGet).Methods("GET")
	post.HandleFunc("/comments", s.HandleCommentCreate).Methods("POST")

	comment := post.PathPrefix("/comment/{c_id:[0-9]+}").Subrouter()
	comment.Use(s.Exist("c_id", "comment"), s.IsOwner("comment"))
	comment.HandleFunc("", s.HandleCommentGet).Methods("GET")
	comment.HandleFunc("", s.HandleCommentUpdate).Methods("PATCH")
	comment.HandleFunc("", s.HandleCommentDelete).Methods("DELETE")

	comment.HandleFunc("/like", s.HandleCommentLike).Methods("POST")
	comment.HandleFunc("/like", s.HandleCommentUnLike).Methods("DELETE")
	comment.HandleFunc("/dislike", s.HandleCommentDislike).Methods("POST")
	comment.HandleFunc("/dislike", s.HandleCommentUnDislike).Methods("DELETE")
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

func returnPointer[T any](t T) *T {
	return &t
}
