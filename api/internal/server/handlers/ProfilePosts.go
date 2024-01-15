package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *Application) ProfilePosts(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		offset := vars["offset"]
		offsetInt, err := strconv.ParseInt(offset, 10, 64)

		if offsetInt < 0 || err != nil {
			app.Logger.Printf("DATA PARSE error: %v", err)
			http.Error(rw, "DATA PARSE error", http.StatusBadRequest)
		}

		userID, err := app.UserService.GetUserID(r)

		if err != nil {
			app.Logger.Printf("Failed fetching user: %v", err)
			http.Error(rw, "Get user error", http.StatusBadRequest)
		}

		feed, err := app.PostService.GetProfilePosts(userID, offsetInt)

		if err != nil {
			app.Logger.Printf("JSON error: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
		}

		json.NewEncoder(rw).Encode(&feed)
	// case "OPTIONS":
	// 	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	default:
		http.Error(rw, "method is not supported", http.StatusNotFound)
		return
	}

}

func (app *Application) UserPosts(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		vars := mux.Vars(r)

		//Get userId from endpoint and parse
		userIdStr := vars["userId"]
		userId, err := strconv.Atoi(userIdStr)

		if userId < 0 || err != nil {
			app.Logger.Printf("DATA PARSE error: %v", err)
			http.Error(rw, "DATA PARSE error", http.StatusBadRequest)
		}

		//Get offset from endpoint and parse
		offsetStr := vars["offset"]
		offset, err := strconv.Atoi(offsetStr)

		if offset < 0 || err != nil {
			app.Logger.Printf("DATA PARSE error: %v", err)
			http.Error(rw, "DATA PARSE error", http.StatusBadRequest)
		}

		requestingUserId, err := app.UserService.GetUserID(r)

		if err != nil {
			app.Logger.Printf("Failed fetching user: %v", err)
			http.Error(rw, "Get user error", http.StatusBadRequest)
		}

		posts, err := app.PostService.GetUserPosts(int64(userId), int64(offset), requestingUserId)

		if err != nil {
			app.Logger.Printf("JSON error: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
		}

		json.NewEncoder(rw).Encode(&posts)

	default:
		http.Error(rw, "method is not supported", http.StatusNotFound)
		return
	}

}
