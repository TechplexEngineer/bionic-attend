package server

import (
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	// shared state between routes
	router     *mux.Router
	templateFs fs.FS
}

// ServeHTTP implements the http.Handler interface which allows the server to be passed to http.ListenAndServe
func (s Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) SetupTemplateFs(fs fs.FS) {
	s.templateFs = fs
}
