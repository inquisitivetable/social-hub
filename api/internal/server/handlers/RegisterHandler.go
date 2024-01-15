package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"SocialNetworkRestApi/api/pkg/models"
)

type signupJSON struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Birthday        string `json:"dateOfBirth"`
	Nickname        string `json:"nickname"`
	About           string `json:"about"`
}

func (app *Application) Register(rw http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(rw, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	r.Body = http.MaxBytesReader(rw, r.Body, 1024)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	JSONdata := &signupJSON{}
	err := decoder.Decode(JSONdata)

	if err != nil {
		app.Logger.Printf("JSON error: %v", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	birthday, err := time.Parse("2006-01-02", JSONdata.Birthday)

	if err != nil {
		app.Logger.Printf("Cannot parse birthday: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	userData := &models.User{
		FirstName: JSONdata.FirstName,
		LastName:  JSONdata.LastName,
		Email:     JSONdata.Email,
		Password:  JSONdata.Password,
		Birthday:  birthday,
		Nickname:  JSONdata.Nickname,
		About:     JSONdata.About,
	}

	sessionToken, err := app.UserService.UserRegister(userData)
	if err != nil {
		app.Logger.Printf("Cannot register user: %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	app.UserService.SetCookie(rw, sessionToken)

	_, err = fmt.Fprintf(rw, "Successful registration")
	if err != nil {
		app.Logger.Printf("Cannot access register page: %s", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
