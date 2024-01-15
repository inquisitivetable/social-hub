package handlers

import (
	"SocialNetworkRestApi/api/pkg/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *Application) GroupEvents(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)

		groupIdStr := vars["groupId"]
		groupId, err := strconv.ParseInt(groupIdStr, 10, 64)

		if groupId < 0 || err != nil {
			app.Logger.Printf("DATA PARSE error: %v", err)
			http.Error(rw, "DATA PARSE error", http.StatusBadRequest)
		}

		groupEvents, err := app.GroupEventService.GetGroupEvents(groupId)

		if err != nil {
			app.Logger.Printf("Failed fetching groups: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
		}

		json.NewEncoder(rw).Encode(&groupEvents)

	default:
		http.Error(rw, "method is not supported", http.StatusNotFound)
		return
	}

}

func (app *Application) CreateGroupEvent(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		//Create a post method here
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		JSONdata := &models.CreateGroupEventFormData{}
		err := decoder.Decode(&JSONdata)

		if err != nil {
			app.Logger.Printf("JSON error: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
			return
		}

		userId, err := app.UserService.GetUserID(r)

		if err != nil {
			app.Logger.Printf("Failed fetching user: %v", err)
			http.Error(rw, "Get user error", http.StatusBadRequest)
			return
		}

		notifications, err := app.GroupEventService.CreateGroupEvent(JSONdata, userId)

		if err != nil {
			http.Error(rw, "err", http.StatusBadRequest)
			return
		}

		err = app.WS.BroadcastGroupNotifications(notifications)

		if err != nil {
			http.Error(rw, "err", http.StatusBadRequest)
			return
		}

		rw.Write([]byte("ok"))

	default:
		http.Error(rw, "err", http.StatusBadRequest)
		return
	}

}

func (app *Application) Event(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)

		eventIdStr := vars["eventId"]
		eventId, err := strconv.ParseInt(eventIdStr, 10, 64)

		if eventId < 0 || err != nil {
			app.Logger.Printf("DATA PARSE error: %v", err)
			http.Error(rw, "DATA PARSE error", http.StatusBadRequest)
		}

		userId, err := app.UserService.GetUserID(r)

		if err != nil {
			app.Logger.Printf("Failed fetching user: %v", err)
			http.Error(rw, "Get user error", http.StatusBadRequest)
			return
		}

		event, err := app.GroupEventService.GetEventById(eventId)

		if err != nil {
			app.Logger.Printf("Failed fetching event: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
		}

		member, err := app.GroupMemberService.GetMemberById(event.GroupId, userId)

		if err != nil && err != sql.ErrNoRows {
			app.Logger.Printf("Failed checking group membership: %v", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err == sql.ErrNoRows {
			app.Logger.Printf("User %d is not a member of this group", userId)
			http.Error(rw, "Not a member of this group", http.StatusForbidden)
			return
		}

		if !member.Accepted {
			app.Logger.Printf("User %d is not a member of this group", userId)
			http.Error(rw, "Not a member of this group", http.StatusForbidden)
			return
		}

		json.NewEncoder(rw).Encode(&event)

	default:
		http.Error(rw, "method is not supported", http.StatusNotFound)
		return
	}

}

func (app *Application) EventReaction(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		JSONdata := &models.EventAttendance{}
		err := decoder.Decode(&JSONdata)

		if err != nil {
			app.Logger.Printf("JSON error: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
			return
		}

		userId, err := app.UserService.GetUserID(r)

		if err != nil {
			app.Logger.Printf("Failed fetching user: %v", err)
			http.Error(rw, "Get user error", http.StatusBadRequest)
			return
		}

		JSONdata.UserId = userId

		notification, err := app.NotificationService.GetByEventAndUserId(JSONdata.EventId, userId)

		if err != nil && err != sql.ErrNoRows {
			app.Logger.Printf("Failed fetching notification: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
		}

		if err == sql.ErrNoRows {
			app.Logger.Printf("No notification found")
		}

		if notification != nil {
			err = app.NotificationService.HandleEventInvite(notification.Id, JSONdata.IsAttending)

			if err != nil && err.Error() != "event invite already handled" {
				app.Logger.Printf("Failed updating notification: %v", err)
				http.Error(rw, "JSON error", http.StatusBadRequest)
			}
		}

		err = app.GroupEventService.UpdateEventAttendance(JSONdata)

		if err != nil {
			app.Logger.Printf("Failed updating event attendance: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
		}

		rw.Write([]byte("ok"))

	default:
		http.Error(rw, "method is not supported", http.StatusNotFound)
		return

	}
}
