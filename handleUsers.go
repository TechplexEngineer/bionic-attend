package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/techplexengineer/go-frc-attend/data"
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
			UserID: mux.Vars(req)[RouteVarUser],
		}

		// handle the request
		s.handleTemplate("create.html", d)(w, req) //@todo
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

		firstName := req.FormValue("firstname")
		lastName := req.FormValue("lastname")
		userID := req.FormValue("userid")

		handleUserMsg := func(message string, redirectPath Route) error {
			SetFlash(w, message)
			if string(redirectPath) != "" {
				redirectPath = RouteCreate
			}
			http.Redirect(w, req, string(redirectPath), http.StatusSeeOther)
			return nil
		}

		if err := CreateNewUser(s.db, handleUserMsg, firstName, lastName, userID); err != nil {
			err = fmt.Errorf("error CreateNewUser - %w", err)
			s.handleInternalError(err)(w, req)
			return
		}

	}
}

func CreateNewUser(db *data.Queries, handleUserMsg func(msg string, path Route) error, firstName, lastName, userID string) error {

	if len(userID) <= 3 {
		return handleUserMsg("UserID must be longer than 3 characters.", RouteCreate)
	}
	if len(firstName) <= 3 {
		return handleUserMsg("First Name must be longer than 3 characters.", RouteCreate) // not ideal as data entered is lost
	}
	if len(lastName) <= 3 {
		return handleUserMsg("Last Name must be longer than 3 characters.", RouteCreate) // not ideal as data entered is lost
	}

	// check if user exists
	count, err := db.GetUserByName(context.Background(), data.GetUserByNameParams{
		FirstName: firstName,
		LastName:  lastName,
	})
	if err != nil {
		return fmt.Errorf("error GetUserByName - %s", err)
	}
	if count > 0 {
		return handleUserMsg(fmt.Sprintf("User '%s %s' already exists", firstName, lastName), "/")
	}

	// check if userid is unique
	count, err = db.UserIDExists(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("error UserIDExists - %s", err)
	}
	if count > 0 {
		return handleUserMsg(fmt.Sprintf("User ID '%s' is already in use", userID), "/")
	}

	params := data.CreateUserParams{
		Userid:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Data:      "{}", // for future expansion
	}

	if err := db.CreateUser(context.Background(), params); err != nil {
		return fmt.Errorf("error CreateUser(%v) - %s", params, err)
	}
	log.Printf("CreateUser(%#v)", params)

	return handleUserMsg("User created", RouteHome)
}
