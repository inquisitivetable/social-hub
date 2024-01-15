package router

import (
	"SocialNetworkRestApi/api/internal/server/handlers"
	"SocialNetworkRestApi/api/internal/server/utils"

	"github.com/gorilla/mux"
)

func New(app *handlers.Application) *mux.Router {
	r := mux.NewRouter()

	r.Use(utils.CorsMiddleware)

	r.HandleFunc("/", app.UserService.Authenticate(app.Home)).Methods("GET")
	r.HandleFunc("/ws", app.UserService.Authenticate(app.WS.WShandler))
	//Session
	r.HandleFunc("/auth", app.UserService.Authenticate(nil)).Methods("GET")
	r.HandleFunc("/login", app.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", app.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/logout", app.UserService.Authenticate(app.Logout)).Methods("GET")
	//Profile
	r.HandleFunc("/profile", app.UserService.Authenticate(app.Profile)).Methods("GET")
	r.HandleFunc("/profile/{id:[0-9]+?}", app.UserService.Authenticate(app.Profile)).Methods("GET")
	r.HandleFunc("/profile/update", app.UserService.Authenticate(app.UpdateProfile)).Methods("POST", "OPTIONS")
	r.HandleFunc("/profile/update/avatar", app.UserService.Authenticate(app.UpdateUserImage)).Methods("POST")
	r.HandleFunc("/following", app.UserService.Authenticate(app.Following)).Methods("GET")
	r.HandleFunc("/following/{id:[0-9]+?}", app.UserService.Authenticate(app.OtherFollowing)).Methods("GET")
	r.HandleFunc("/followers", app.UserService.Authenticate(app.Followers)).Methods("GET")
	r.HandleFunc("/followers/{id:[0-9]+?}", app.UserService.Authenticate(app.OtherFollowers)).Methods("GET")
	//Posts
	r.HandleFunc("/feedposts/{offset:[0-9]+?}", app.UserService.Authenticate(app.FeedPosts)).Methods("GET")
	r.HandleFunc("/comments/{postId:[0-9]+?}/{offset:[0-9]+?}", app.UserService.Authenticate(app.Comments)).Methods("GET")
	r.HandleFunc("/insertcomment", app.UserService.Authenticate(app.Comment)).Methods("POST")
	r.HandleFunc("/post", app.UserService.Authenticate(app.Post)).Methods("POST")
	r.HandleFunc("/profileposts/{offset:[0-9]+?}", app.UserService.Authenticate(app.ProfilePosts)).Methods("GET")
	r.HandleFunc("/userposts/{userId:[0-9]+?}/{offset:[0-9]+?}", app.UserService.Authenticate(app.UserPosts)).Methods("GET")
	r.HandleFunc("/groups/{groupId:[0-9]+?}/post", app.UserService.Authenticate(app.GroupPost)).Methods("POST")
	//Groups
	r.HandleFunc("/creategroup", app.UserService.Authenticate(app.CreateGroup)).Methods("POST")
	r.HandleFunc("/usergroups", app.UserService.Authenticate(app.UserGroups)).Methods("GET")
	r.HandleFunc("/mygroups", app.UserService.Authenticate(app.MyGroups)).Methods("GET")
	r.HandleFunc("/groups/{groupId:[0-9]+?}", app.UserService.Authenticate(app.Group)).Methods("GET")
	r.HandleFunc("/groups/{groupId:[0-9]+?}/avatar", app.UserService.Authenticate(app.UpdateGroupImage)).Methods("POST")
	r.HandleFunc("/groupmembers/{groupId:[0-9]+?}", app.UserService.Authenticate(app.GroupMembers)).Methods("GET")
	r.HandleFunc("/addmembers", app.UserService.Authenticate(app.AddMembers)).Methods("POST")
	r.HandleFunc("/addmembers/{groupId:[0-9]+?}", app.UserService.Authenticate(app.GetMembersToAdd)).Methods("GET")
	r.HandleFunc("/groupfeed/{groupId:[0-9]+?}/{offset:[0-9]+?}", app.UserService.Authenticate(app.GroupPosts)).Methods("GET")
	//Events
	r.HandleFunc("/userevents", app.UserService.Authenticate(app.UserEvents)).Methods("GET")
	r.HandleFunc("/creategroupevent", app.UserService.Authenticate(app.CreateGroupEvent)).Methods("POST")
	r.HandleFunc("/groupevents/{groupId:[0-9]+?}", app.UserService.Authenticate(app.GroupEvents)).Methods("GET")
	r.HandleFunc("/event/{eventId:[0-9]+?}", app.UserService.Authenticate(app.Event)).Methods("GET")
	r.HandleFunc("/eventreaction", app.UserService.Authenticate(app.EventReaction)).Methods("POST")
	//Search
	r.HandleFunc("/search/{searchcriteria}", app.UserService.Authenticate(app.Search)).Methods("GET")
	r.HandleFunc("/notifications", app.UserService.Authenticate(app.Notifications)).Methods("GET")
	return r
}
