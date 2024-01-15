package handlers

import (
	"fmt"
	"net/http"
)

func (app *Application) Home(rw http.ResponseWriter, r *http.Request) {

	_, err := fmt.Fprintf(rw, "Homepage hit")
	if err != nil {
		app.Logger.Println("Cannot access homepage")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
