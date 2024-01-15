package handlers

import "net/http"

func (app *Application) Logout(rw http.ResponseWriter, r *http.Request) {

	err := app.UserService.UserLogout(r)
	if err != nil {
		app.Logger.Printf("Cannot logout user: %s", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	app.UserService.ClearCookie(rw)
	_, err = rw.Write([]byte("Successful logout, cookie cleared"))
	if err != nil {
		app.Logger.Printf("Cannot access logout page: %s", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
