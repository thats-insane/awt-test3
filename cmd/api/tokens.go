package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/thats-insane/awt-test3/internal/data"
	"github.com/thats-insane/awt-test3/internal/validator"
)

func (a *appDependencies) createAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequest(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, incomingData.Email)
	data.ValidatePasswordPlaintext(v, incomingData.Password)
	if !v.IsEmpty() {
		a.failedValidation(w, r, v.Errors)
		return
	}

	user, err := a.userModel.GetByEmail(incomingData.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.invalidCredentials(w, r)
		default:
			a.serverErr(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(incomingData.Password)
	if err != nil {
		a.serverErr(w, r, err)
		return
	}

	if !match {
		a.invalidCredentials(w, r)
		return
	}

	token, err := a.tokenModel.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		a.serverErr(w, r, err)
		return
	}

	data := envelope{
		"authenticationToken": token,
	}

	err = a.writeJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		a.serverErr(w, r, err)
	}
}
