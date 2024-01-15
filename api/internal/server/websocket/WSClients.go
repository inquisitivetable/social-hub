package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	clientID   int64
	manager    *WebsocketServer
	gate       chan Payload
}

var (
	// turn ping messages on to monitor connections
	pingOn bool = false
	// pongWait is how long we will await a pong response from client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval         = (pongWait * 9) / 10
	maxMessageSize int64 = 512
)

func NewClient(conn *websocket.Conn, userID int64, manager *WebsocketServer) *Client {
	return &Client{
		connection: conn,
		clientID:   userID,
		manager:    manager,
		gate:       make(chan Payload),
	}
}

func (c *Client) monitor() {
	defer func() {
		c.manager.Logger.Println("Closing connection")
		c.manager.removeClient(c)
	}()

	c.connection.SetReadLimit(maxMessageSize)

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.manager.Logger.Printf("Error setting read deadline: %v", err)
		return
	}

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.manager.Logger.Printf("Error reading message: %v", err)
			}
			break
		}

		var request Payload
		if err := json.Unmarshal(payload, &request); err != nil {
			c.manager.Logger.Printf("Error unmarshalling payload: %v", err)
			continue
		}

		if err := c.manager.routePayloads(request, c); err != nil {
			c.manager.Logger.Printf("Error routing payload: %v", err)
			continue
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	if pingOn {
		c.manager.Logger.Printf("Received pong from client %v", c.clientID)
	}
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.manager.Logger.Printf("Error setting read deadline: %v", err)
		return err
	}
	return nil
}

func (c *Client) write() {
	ticker := time.NewTicker(pingInterval)

	defer func() {
		c.manager.Logger.Printf("Closing connection for client %v", c.clientID)
		ticker.Stop()
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.gate:
			// OK is false when channel is closed
			if !ok {
				// Server closed the channel
				if err := c.connection.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.manager.Logger.Printf("Error writing close message: %v", err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				c.manager.Logger.Printf("Error marshalling message: %v", err)
				return
			}
			c.manager.Logger.Printf("Writing message '%v' to client %v", message.Type, c.clientID)
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				c.manager.Logger.Printf("Error writing message: %v", err)
				return
			}
		case <-ticker.C:
			if pingOn {
				c.manager.Logger.Printf("Sending ping to client %v", c.clientID)
			}
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.manager.Logger.Printf("Error writing ping message: %v", err)
				return
			}
		}
	}
}
