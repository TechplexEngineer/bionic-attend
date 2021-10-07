package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) SetupRoutes() {
	if s.router == nil {
		s.router = mux.NewRouter()
	}

	s.router.PathPrefix("/static/").Handler(http.FileServer(http.FS(s.templateFs))).Methods(http.MethodGet)

	s.router.HandleFunc("/", s.handleIndex())

	s.router.HandleFunc("/checkin", s.handleCheckIn())

	// s.router.HandleFunc("/admin", s.adminOnly(s.handleAdminIndex()))
}

//func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
//    return func(w http.ResponseWriter, r *http.Request) {
//        if !currentUser(r).IsAdmin {
//            http.NotFound(w, r)
//            return
//        }
//        h(w, r)
//    }
//}
