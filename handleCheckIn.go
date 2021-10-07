package server

import (
	"log"
	"net/http"
)

func (s Server) handleCheckIn() http.HandlerFunc {
	// one time handler setup work can go here

	return func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			s.handleInternalError(err)(w, req)
			return
		}

		log.Printf("userid: %s", req.FormValue("userid"))

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return

		// handle the request
		//s.handleTemplate("index.html", nil)(w, req)
	}
}
