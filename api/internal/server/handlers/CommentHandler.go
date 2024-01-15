package handlers

import (
	"SocialNetworkRestApi/api/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *Application) Comments(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		vars := mux.Vars(r)

		//Get postId from endpoint and parse
		postIdStr := vars["postId"]
		postId, err := strconv.ParseInt(postIdStr, 10, 64)

		// fmt.Println("postIdStr", postIdStr)

		if postId < 0 || err != nil {
			app.Logger.Printf("DATA PARSE error: %v", err)
			http.Error(rw, "DATA PARSE error", http.StatusBadRequest)
		}

		//Get offset from endpoint and parse
		offsetStr := vars["offset"]
		offset, err := strconv.ParseInt(offsetStr, 10, 64)

		if offset < 0 || err != nil {
			app.Logger.Printf("DATA PARSE error: %v", err)
			http.Error(rw, "DATA PARSE error", http.StatusBadRequest)
		}

		comments, err := app.CommentService.GetPostComments(postId, offset)

		if err != nil {
			app.Logger.Printf("JSON error: %v", err)
			http.Error(rw, "JSON error", http.StatusBadRequest)
		}

		json.NewEncoder(rw).Encode(&comments)

	default:
		http.Error(rw, "method is not supported", http.StatusNotFound)
		return
	}

}

func (app *Application) Comment(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":

		// Limit the size of the request body to 5MB
		r.Body = http.MaxBytesReader(rw, r.Body, 20<<18+512)

		err := r.ParseMultipartForm(20 << 18)

		if err != nil {
			app.Logger.Printf("Failed parsing form: %v", err)
			http.Error(rw, "Parsing form error", http.StatusRequestEntityTooLarge)
		}

		postIdStr := r.FormValue("postId")
		postId, err := strconv.ParseInt(postIdStr, 10, 64)
		if err != nil {
			app.Logger.Printf("DATA PARSE error: %v", err)
			http.Error(rw, "DATA PARSE error", http.StatusBadRequest)
		}

		content := r.FormValue("content")

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
		}

		comment := &models.Comment{
			PostId:    postId,
			UserId:    userId,
			Content:   content,
			ImagePath: imagePath,
		}

		err = app.CommentService.CreateComment(comment)

		if err != nil {
			app.Logger.Printf("Creating comment failed: %v", err)
			http.Error(rw, "Error", http.StatusBadRequest)
		}

		rw.Write([]byte("ok"))

	default:
		http.Error(rw, "method is not supported", http.StatusNotFound)
		return
	}

}
