package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/techplexengineer/go-frc-attend/data"
	"log"
	"net/http"
)

const (
	TemplateEdit = "edit.html"
)

func (s Server) handleEdit() http.HandlerFunc {
	// one time handler setup work can go here

	return func(w http.ResponseWriter, req *http.Request) {
		// generate data

		userID := mux.Vars(req)[RouteVarUser]

		user, err := s.queries.GetUser(context.Background(), userID)
		if err != nil {
			err = fmt.Errorf("error GetUser - %w", err)
			s.handleInternalError(err)(w, req)
			return
		}

		d := struct {
			User data.User
		}{
			User: user,
		}

		// handle the request
		s.handleTemplate(TemplateEdit, d)(w, req)
	}
}

func (s Server) handleEditPOST() http.HandlerFunc {
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
		newUserID := req.FormValue("newuserid")
		existingUserID := req.FormValue("existinguserid")

		handleUserMsg := func(message string, redirectPath Route) error {
			SetFlash(w, message)
			if string(redirectPath) == "" {
				redirectPath = RouteHome
			}
			http.Redirect(w, req, string(redirectPath), http.StatusSeeOther)
			return nil
		}
		tx, err := s.db.Begin() // start a transaction
		if err != nil {
			err = fmt.Errorf("error begintx - %w", err)
			s.handleInternalError(err)(w, req)
			return
		}
		defer func() {
			err = tx.Rollback()
			if err != nil {
				log.Printf("editUser tx rollback error %s", err)
			}
		}()

		log.Printf("fn:%s ln:%s ex:%s nid:%s", firstName, lastName, existingUserID, newUserID)

		if err := UpdateUser(s.queries.WithTx(tx), handleUserMsg, firstName, lastName, existingUserID, newUserID); err != nil {
			err = fmt.Errorf("error UpdateUser - %w", err)
			s.handleInternalError(err)(w, req)
			return
		}
		if err := tx.Commit(); err != nil {
			err = fmt.Errorf("error txcommit - %w", err)
			s.handleInternalError(err)(w, req)
			return
		}
	}
}

func (s Server) handleHidePOST() http.HandlerFunc {
	// one time handler setup work can go here

	return func(w http.ResponseWriter, req *http.Request) {
		// generate data
		err := req.ParseForm()
		if err != nil {
			err = fmt.Errorf("error ParseForm - %w", err)
			s.handleInternalError(err)(w, req)
			return
		}

		userid := req.FormValue("userid")
		if len(userid) < 2 {
			err = fmt.Errorf("error missing userid")
			s.handleInternalError(err)(w, req)
			return
		}

		//log.Printf("If it were implemented, the user would be hidden!")
		err = s.queries.SoftDeleteUser(context.Background(), userid)
		if err != nil {
			err = fmt.Errorf("unable to soft delete user: %w", err)
			s.handleInternalError(err)(w, req)
			return
		}

		SetFlash(w, "User successfully hidden")
		http.Redirect(w, req, string(RouteReport), http.StatusSeeOther)
	}
}

func UpdateUser(db *data.Queries, handleUserMsg func(msg string, path Route) error, firstName, lastName, existingUserID, newUserID string) error {
	if len(newUserID) <= 2 {
		return handleUserMsg("UserID must be 3 characters or longer.", NewRoute(RouteEdit, existingUserID))
	}
	if len(firstName) == 0 {
		return handleUserMsg("First Name must be 1 character or longer.", NewRoute(RouteEdit, existingUserID)) // not ideal as data entered is lost
	}
	if len(lastName) == 0 {
		return handleUserMsg("Last Name must be 1 character or longer.", NewRoute(RouteEdit, existingUserID)) // not ideal as data entered is lost
	}
	ctx := context.Background()

	params := data.UpdateUserParams{
		Userid:    existingUserID,
		FirstName: firstName,
		LastName:  lastName,
		Data:      "{}",
	}

	if err := db.UpdateUser(ctx, params); err != nil {
		return fmt.Errorf("error UpdateUser(%v) - %s", params, err)
	}

	if existingUserID == newUserID {
		return handleUserMsg("User Updated", RouteHome)
	}

	// check if userid is unique
	count, err := db.UserIDExists(ctx, newUserID)
	if err != nil {
		return fmt.Errorf("error UserIDExists %s - %s", newUserID, err)
	}
	if count > 0 {
		return handleUserMsg(fmt.Sprintf("User ID '%s' is already in use", newUserID), NewRoute(RouteEdit, existingUserID))
	}

	{
		params := data.UpdateUserIDinUsersParams{
			Userid:   newUserID,
			Userid_2: existingUserID,
		}
		err := db.UpdateUserIDinUsers(ctx, params)
		if err != nil {
			return fmt.Errorf("error UpdateUserIDinUsersParams(%v) - %s", params, err)
		}
	}

	{
		params := data.UpdateUserIDinAttendanceParams{
			Userid:   newUserID,
			Userid_2: existingUserID,
		}
		err := db.UpdateUserIDinAttendance(ctx, params)
		if err != nil {
			return fmt.Errorf("error UpdateUserIDinAttendance(%v) - %s", params, err)
		}
	}

	//params := data.CreateUserParams{
	//	Userid:    userID,
	//	FirstName: firstName,
	//	LastName:  lastName,
	//	Data:      "{}", // for future expansion
	//}
	//
	//if err := queries.CreateUser(ctx, params); err != nil {
	//	return fmt.Errorf("error CreateUser(%v) - %s", params, err)
	//}
	//log.Printf("CreateUser(%#v)", params)

	return handleUserMsg("User Updated", RouteHome)
}
