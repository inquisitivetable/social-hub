package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"SocialNetworkRestApi/api/pkg/models"
)

type signinJSON struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

func (app *Application) Login(rw http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(rw, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
	}
	r.Body = http.MaxBytesReader(rw, r.Body, 1024)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	JSONdata := &signinJSON{}
	err := decoder.Decode(JSONdata)

	if err != nil {
		app.Logger.Printf("JSON error: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if JSONdata.Email == "" || JSONdata.Password == "" {
		app.Logger.Printf("Incomplete login credentials: %s", JSONdata)
		http.Error(rw, "Incomplete credentials", http.StatusBadRequest)
		return
	}

	userData := &models.User{
		Email:    JSONdata.Email,
		Password: JSONdata.Password,
	}

	sessionToken, err := app.UserService.UserLogin(userData)
	if err != nil {
		app.Logger.Printf("Cannot login user: %s", err)
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	app.UserService.SetCookie(rw, sessionToken)

	_, err = fmt.Fprintf(rw, "Successful login, cookie set")
	if err != nil {
		app.Logger.Printf("Cannot access login page: %s", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}
