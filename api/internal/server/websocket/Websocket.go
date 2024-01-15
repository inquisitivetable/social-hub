package websocket

import (
	"SocialNetworkRestApi/api/pkg/services"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Mostly done using https://programmingpercy.tech/blog/mastering-websockets-with-go/

type IWebsocketServer interface {
	WShandler(rw http.ResponseWriter, r *http.Request)
}

type WebsocketServer struct {
	Logger              *log.Logger
	upgrader            websocket.Upgrader
	clients             ClientList
	handlers            map[string]PayloadHandler
	userService         services.IUserService
	notificationService services.INotificationService
	chatService         services.IChatService
	groupService        services.IGroupService
	groupMemberService  services.IGroupMemberService
	groupEventService   services.IGroupEventService
	sync.RWMutex
}

func (w *WebsocketServer) WShandler(rw http.ResponseWriter, r *http.Request) {
	conn, err := w.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		w.Logger.Println("Cannot upgrade:", err)
		rw.Write([]byte(err.Error()))
		return
	}

	w.Logger.Println("Successfully upgraded connection")

	userID, err := w.userService.GetUserID(r)
	if err != nil {
		w.Logger.Printf("Cannot get user ID: %s", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	client := NewClient(conn, userID, w)

	w.addClient(client)

	go client.monitor()
	go client.write()

}

func InitWebsocket(
	logger *log.Logger,
	userService *services.UserService,
	notificationService *services.NotificationService,
	chatService *services.ChatService,
	groupService *services.GroupService,
	groupMemberService *services.GroupMemberService,
	groupEventService *services.GroupEventService,
) *WebsocketServer {
	w := &WebsocketServer{
		Logger: logger,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     checkOrigin,
		},
		clients:             make(ClientList),
		handlers:            make(map[string]PayloadHandler),
		userService:         userService,
		notificationService: notificationService,
		chatService:         chatService,
		groupService:        groupService,
		groupMemberService:  groupMemberService,
		groupEventService:   groupEventService,
	}
	w.setupHandlers()
	return w
}

func checkOrigin(r *http.Request) bool {
	// required for CORS
	// should return true if origin is allowed
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://localhost:3000":
		return true
	default:
		return false
	}
}

func (w *WebsocketServer) addClient(client *Client) {
	w.Lock()
	defer w.Unlock()
	w.Logger.Printf("Adding client %v", client.clientID)
	w.clients[client] = true
}

func (w *WebsocketServer) removeClient(client *Client) {
	w.Lock()
	defer w.Unlock()
	// Check if Client exists, then delete it
	if _, ok := w.clients[client]; ok {
		w.Logger.Printf("Removing client %v", client.clientID)
		client.connection.Close()
		delete(w.clients, client)
	}
}

func (w *WebsocketServer) getClientByUserID(TargetID int64) *Client {
	w.RLock()
	defer w.RUnlock()
	for client := range w.clients {
		if client.clientID == TargetID {
			return client
		}
	}
	return nil
}
