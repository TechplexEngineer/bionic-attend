package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route string

const (
	RouteVarUser = "userid"

	RouteHome    Route = "/"
	RouteCheckin Route = "/checkin"
	RouteCreate  Route = "/create"
	RouteReport  Route = "/report"
	RouteEdit    Route = "/edit"
)

func (s *Server) SetupRoutes() {
	if s.router == nil {
		s.router = mux.NewRouter()
	}

	s.router.PathPrefix("/static/").Handler(http.FileServer(http.FS(s.templateFs))).Methods(http.MethodGet)

	s.router.HandleFunc(string(RouteHome), s.handleHome())

	s.router.HandleFunc(string(RouteCheckin), s.handleCheckIn())
	s.router.HandleFunc(string(RouteCreate), s.handleNewUserPOST()).Methods(http.MethodPost)
	s.router.HandleFunc(string(RouteCreate), s.handleNewUser()).Methods(http.MethodGet)
	s.router.HandleFunc(string(RouteCreate+"{"+RouteVarUser+"}"), s.handleNewUser())
	s.router.HandleFunc(string(RouteReport), s.handleReport())
	s.router.HandleFunc(string(RouteEdit+"{"+RouteVarUser+"}"), s.handleReport())

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
