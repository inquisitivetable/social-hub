package websocket

import (
	"SocialNetworkRestApi/api/pkg/services"
	"encoding/json"
	"time"
)

type Payload struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type RequestPayload struct {
	ID          int  `json:"id"`
	Reaction    bool `json:"reaction"`
	GroupID     int  `json:"group_id"`
	LastMessage int  `json:"last_message"`
}

type NotificationPayload struct {
	NotificationType string    `json:"notification_type"`
	NotificationID   int       `json:"notification_id"`
	SenderID         int       `json:"sender_id"`
	SenderName       string    `json:"sender_name"`
	GroupID          int       `json:"group_id"`
	GroupName        string    `json:"group_name"`
	EventID          int       `json:"event_id"`
	EventName        string    `json:"event_name"`
	EventDate        time.Time `json:"event_datetime"`
}

type MessagePayload struct {
	MessageID     int       `json:"id"`
	SenderID      int       `json:"sender_id"`
	SenderName    string    `json:"sender_name"`
	SenderImage   string    `json:"avatar_image"`
	RecipientID   int       `json:"recipient_id"`
	RecipientName string    `json:"recipient_name"`
	GroupID       int       `json:"group_id"`
	GroupName     string    `json:"group_name"`
	Content       string    `json:"body"`
	Timestamp     time.Time `json:"timestamp"`
}

type ChatListPayload struct {
	UserID        int                      `json:"user_id"`
	UserChatlist  []services.UserChatList  `json:"user_chatlist"`
	GroupChatlist []services.GroupChatList `json:"group_chatlist"`
}
