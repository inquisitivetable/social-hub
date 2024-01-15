package handlers

import (
	"encoding/json"
	"net/http"
)

func (app *Application) Notifications(rw http.ResponseWriter, r *http.Request) {

	userID, err := app.UserService.GetUserID(r)
	if err != nil {
		app.Logger.Printf("Cannot get user ID: %s", err)
		http.Error(rw, "Cannot get user ID", http.StatusUnauthorized)
		return
	}

	notifications, err := app.NotificationService.GetUserNotifications(int64(userID))

	if err != nil {
		app.Logger.Printf("Cannot get user notifications: %s", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(notifications)

	if err != nil {
		app.Logger.Printf("Cannot encode user notifications: %s", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
