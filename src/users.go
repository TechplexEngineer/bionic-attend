package src

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/techplexengineer/go-frc-attend/src/data"
	"log"
	"net/http"
)

func (s Server) handleNewUser() http.HandlerFunc {
	// one time handler setup work can go here

	return func(w http.ResponseWriter, req *http.Request) {
		// generate data

		d := struct {
			UserID string
		}{
			UserID: mux.Vars(req)[RouteCreateVarUser],
		}

		// handle the request
		s.handleTemplate("userCreate.html", d)(w, req) //@todo
	}
}

func (s Server) handleNewUserPOST() http.HandlerFunc {
	// one time handler setup work can go here

	return func(w http.ResponseWriter, req *http.Request) {
		// generate data
		err := req.ParseForm()
		if err != nil {
			err = fmt.Errorf("error ParseForm - %w", err)
			s.handleInternalError(err)(w, req)
			return
		}

		// check if user exists
		count, err := s.db.GetUserByName(context.Background(), data.GetUserByNameParams{
			FirstName: req.FormValue("firstname"),
			LastName:  req.FormValue("lastname"),
		})
		if err != nil {
			err = fmt.Errorf("error GetUserByName - %s", err)
			s.handleInternalError(err)(w, req)
			return
		}
		if count > 0 {
			SetFlash(w, fmt.Sprintf("User '%s %s' already exists", req.FormValue("firstname"), req.FormValue("lastname")))
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}

		// check if userid is unique
		count, err = s.db.UserIDExists(context.Background(), req.FormValue("userid"))
		if err != nil {
			err = fmt.Errorf("error UserIDExists - %s", err)
			s.handleInternalError(err)(w, req)
			return
		}
		if count > 0 {
			SetFlash(w, fmt.Sprintf("User ID '%s' is already in use", req.FormValue("userid")))
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}

		params := data.CreateUserParams{
			Userid:    req.FormValue("userid"),
			FirstName: req.FormValue("firstname"),
			LastName:  req.FormValue("lastname"),
			Data:      "{}", // for future expansion
		}
		err = s.db.CreateUser(context.Background(), params)
		if err != nil {
			err = fmt.Errorf("error CreateUser(%v) - %s", params, err)
			s.handleInternalError(err)(w, req)
			return
		}
		log.Printf("CreateUser(%#v)", params)

		SetFlash(w, "User created")
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
