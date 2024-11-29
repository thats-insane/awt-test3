package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/thats-insane/awt-test3/internal/data"
	"github.com/thats-insane/awt-test3/internal/validator"
)

func (a *appDependencies) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := a.readJSON(w, r, incomingData)
	if err != nil {
		a.badRequest(w, r, err)
		return
	}

	user := &data.User{
		Username:  incomingData.Username,
		Email:     incomingData.Email,
		Activated: false,
	}

	err = user.Password.Set(incomingData.Password)

	if err != nil {
		a.serverErr(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateUser(v, user)

	if !v.IsEmpty() {
		a.failedValidation(w, r, v.Errors)
		return
	}

	err = a.userModel.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			a.failedValidation(w, r, v.Errors)
		default:
			a.serverErr(w, r, err)
		}
		return
	}

	token, err := a.tokenModel.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		a.serverErr(w, r, err)
		return
	}

	data := envelope{
		"user": user,
	}

	a.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		err = a.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			a.logger.Error(err.Error())
		}
	})

	err = a.writeJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		a.serverErr(w, r, err)
		return
	}
}

func (a *appDependencies) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Plaintext string `json:"token"`
	}
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequest(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateTokenPlaintext(v, incomingData.Plaintext)
	if !v.IsEmpty() {
		a.failedValidation(w, r, v.Errors)
		return
	}

	user, err := a.userModel.GetForToken(data.ScopeActivation, incomingData.Plaintext)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid/expired activation token")
			a.failedValidation(w, r, v.Errors)
		default:
			a.serverErr(w, r, err)
		}
		return
	}

	user.Activated = true
	err = a.userModel.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			a.editConflict(w, r)
		default:
			a.serverErr(w, r, err)
		}
		return
	}

	data := envelope{
		"user": user,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErr(w, r, err)
	}
}

func (a *appDependencies) displayUserListsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFound(w, r)
		return
	}

	booklist, err := a.listModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFound(w, r)
		default:
			a.serverErr(w, r, err)
		}
		return
	}

	books, err := a.listModel.GetBooks(booklist.BookListID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFound(w, r)
		default:
			a.serverErr(w, r, err)
		}
		return
	}

	data := envelope{
		"lists": books,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErr(w, r, err)
		return
	}
}
