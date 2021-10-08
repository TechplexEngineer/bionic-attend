package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/techplexengineer/go-frc-attend/data"
	"log"
	"net/http"
)

func (s Server) handleNewUser() http.HandlerFunc {
	// one time handler setup work can go here

	return func(w http.ResponseWriter, req *http.Request) {
		// generate data

		data := struct {
			UserId string
		}{
			UserId: mux.Vars(req)[RouteCreateVarUser],
		}

		// handle the request
		s.handleTemplate("create.html", data)(w, req) //@todo
	}
}

func (s Server) handleNewUserPOST() http.HandlerFunc {
	// one time handler setup work can go here

	return func(w http.ResponseWriter, req *http.Request) {
		// generate data
		err := req.ParseForm()
		if err != nil {
			log.Printf("error ParseForm - %s", err)
			s.handleInternalError(err)(w, req)
			return
		}

		params := data.CreateUserParams{
			Userid:    req.FormValue("userid"),
			FirstName: req.FormValue("firstname"),
			LastName:  req.FormValue("lastname"),
			Data:      "{}",
		}
		err = s.db.CreateUser(context.Background(), params)
		if err != nil {
			log.Printf("error CreateUser(%v) - %s", params, err)
			s.handleInternalError(err)(w, req)
			return
		}
		log.Printf("CreateUser(%#v)", params)

		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
