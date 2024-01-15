package handlers

import (
	"SocialNetworkRestApi/api/pkg/enums"
	"SocialNetworkRestApi/api/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type createPostJSON struct {
	UserId      int      `json:"userId"`
	Content     string   `json:"content"`
	ImagePath   string   `json:"imagePath"`
	PrivacyType int      `json:"privacyType"`
	Receivers   []string `json:"selectedReceivers"`
}

func (app *Application) Post(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		// Limit the size of the request body to 5MB
		r.Body = http.MaxBytesReader(rw, r.Body, 20<<18+512)

		err := r.ParseMultipartForm(20 << 18)

		if err != nil {
			app.Logger.Printf("Failed parsing form: %v", err)
			http.Error(rw, "Parsing form error", http.StatusRequestEntityTooLarge)
		}

		// decoder := json.NewDecoder(r.Body)
		// decoder.DisallowUnknownFields()

		content := r.FormValue("content")
		privacyTypeStr := r.FormValue("privacyType")

		var privacyType enums.PrivacyType
		switch privacyTypeStr {
		case "1":
			privacyType = 1
		case "2":
			privacyType = 2
		case "3":
			privacyType = 3
		default:
			app.Logger.Printf("Invalid privacyType value: %s", privacyTypeStr)
			http.Error(rw, "Invalid privacyType value", http.StatusBadRequest)
			return
		}

		receiversStr := r.FormValue("selectedReceivers")
		receivers := strings.Split(receiversStr, ",")

		if err != nil {
			app.Logger.Printf("JSON error: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
		}

		file, header, err := r.FormFile("image")
		var imagePath string

		if err == nil {
			defer file.Close()

			imagePath, err = app.PostService.SavePostImage(file, header)
			if err != nil {
				app.Logger.Printf("Failed saving image: %v", err)
				http.Error(rw, "Save image error", http.StatusBadRequest)
				return
			}
		}

		userId, err := app.UserService.GetUserID(r)

		if err != nil {
			app.Logger.Printf("Failed fetching user: %v", err)
			http.Error(rw, "Get user error", http.StatusBadRequest)
			return
		}

		post := &models.Post{
			UserId:      userId,
			ImagePath:   imagePath,
			Content:     content,
			PrivacyType: privacyType,
			Receivers:   receivers,
		}

		err = app.PostService.CreatePost(post)

		if err != nil {
			app.Logger.Printf("Cannot create post: %s", err)
			http.Error(rw, "err", http.StatusBadRequest)
			return
		}

		rw.Write([]byte("ok"))

	default:
		http.Error(rw, "err", http.StatusBadRequest)
		return
	}

}

func (app *Application) GroupPost(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		vars := mux.Vars(r)

		groupIdStr := vars["groupId"]
		groupId, err := strconv.ParseInt(groupIdStr, 10, 64)

		if groupId < 0 || err != nil {
			app.Logger.Printf("DATA PARSE error: %v", err)
			http.Error(rw, "DATA PARSE error", http.StatusBadRequest)
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		JSONdata := &createPostJSON{}
		err = decoder.Decode(&JSONdata)

		if err != nil {
			app.Logger.Printf("JSON error: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
		}

		userId, err := app.UserService.GetUserID(r)

		if err != nil {
			app.Logger.Printf("Failed fetching user: %v", err)
			http.Error(rw, "Get user error", http.StatusBadRequest)
		}

		post := &models.Post{
			UserId:      userId,
			PrivacyType: enums.PrivacyType(enums.None),
			Content:     JSONdata.Content,
			ImagePath:   JSONdata.ImagePath,
			GroupId:     groupId,
		}

		err = app.PostService.CreateGroupPost(post)

		if err != nil {
			app.Logger.Printf("Cannot create post: %s", err)
			http.Error(rw, "err", http.StatusBadRequest)
			return
		}

		rw.Write([]byte("ok"))

	default:
		http.Error(rw, "err", http.StatusBadRequest)
		return
	}

}
