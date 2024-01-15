package websocket

import (
	"SocialNetworkRestApi/api/pkg/models"
	"encoding/json"
	"errors"
	"time"
)

type PayloadHandler func(payload Payload, client *Client) error

var (
	ErrPayloadTypeNotSupported = errors.New("this payload type is not supported")
	ErrorInvalidPayload        = errors.New("invalid payload")
	ErrorInvalidNotification   = errors.New("invalid notification")
)

const (
	FollowRequest   = "follow_request"
	Unfollow        = "unfollow"
	RequestChatlist = "request_chatlist"
	MessageHistory  = "request_message_history"
	Message         = "message"
	GroupRequest    = "group_request"
	Response        = "response"
	MessagesRead    = "messages_read"
)

func (w *WebsocketServer) setupHandlers() {
	w.handlers[FollowRequest] = w.FollowRequestHandler
	w.handlers[Unfollow] = w.UnfollowHandler
	w.handlers[RequestChatlist] = w.RequestChatlistHandler
	w.handlers[MessageHistory] = w.MessageHistoryHandler
	w.handlers[Message] = w.NewMessageHandler
	w.handlers[GroupRequest] = w.GroupRequestHandler
	w.handlers[Response] = w.ResponseHandler
	w.handlers[MessagesRead] = w.MessagesReadHandler
}

func (w *WebsocketServer) routePayloads(payload Payload, client *Client) error {
	handler, ok := w.handlers[payload.Type]
	if !ok {
		w.Logger.Printf("No handler for event %s", payload.Type)
		return ErrPayloadTypeNotSupported
	}
	if err := handler(payload, client); err != nil {
		return err
	}
	return nil
}

func (w *WebsocketServer) ResponseHandler(p Payload, c *Client) error {
	data := &RequestPayload{}
	err := json.Unmarshal(p.Data, &data)
	if err != nil {
		return err
	}
	w.Logger.Printf("User %v responded to notification %v with %v", c.clientID, data.ID, data.Reaction)

	notification, err := w.notificationService.GetById(int64(data.ID))
	if err != nil {
		return err
	}

	if notification == nil {
		w.Logger.Printf("Notification not found")
		return ErrorInvalidNotification
	}

	if notification.ReceiverId != int64(c.clientID) {
		w.Logger.Printf("Notification does not belong to user")
		return ErrorInvalidNotification
	}

	NotificationDetails, err := w.notificationService.GetDetailsById(notification.NotificationDetailsId)
	if err != nil {
		w.Logger.Printf("Error getting notification details: %s", err.Error())
		return err
	}

	// perhaps case switch here?
	if NotificationDetails.NotificationType == "follow_request" {
		w.Logger.Printf("User %v reacted to follow request %v", c.clientID, data.ID)
		err = w.notificationService.HandleFollowRequest(int64(data.ID), data.Reaction)
		if err != nil {
			return err
		}
		return nil
	}

	if NotificationDetails.NotificationType == "group_invite" {
		w.Logger.Printf("User %v reacted to group invite %v", c.clientID, data.ID)
		err = w.notificationService.HandleGroupInvite(int64(data.ID), data.Reaction)
		if err != nil {
			return err
		}
		return nil
	}

	if NotificationDetails.NotificationType == "group_request" {
		w.Logger.Printf("User %v reacted to group request %v", c.clientID, data.ID)
		err = w.notificationService.HandleGroupRequest(c.clientID, int64(data.ID), data.Reaction)
		if err != nil {
			return err
		}
		return nil
	}

	if NotificationDetails.NotificationType == "event_invite" {
		w.Logger.Printf("User %v reacted to event invite %v", c.clientID, data.ID)
		err = w.notificationService.HandleEventInvite(int64(data.ID), data.Reaction)
		if err != nil && err.Error() != "event invite already handled" {
			return err
		}
		if err != nil && err.Error() == "event invite already handled" {
			attendance := &models.EventAttendance{
				EventId:     NotificationDetails.EntityId,
				UserId:      int64(c.clientID),
				IsAttending: data.Reaction,
			}
			w.Logger.Printf("attendance: %+v", attendance)
			w.groupEventService.UpdateEventAttendance(attendance)
		}
		return nil
	}

	w.Logger.Printf("Notification type %v not handled", NotificationDetails.NotificationType)

	return errors.New("unknown notification type: " + NotificationDetails.NotificationType)
}

func (w *WebsocketServer) FollowRequestHandler(p Payload, c *Client) error {
	data := &RequestPayload{}
	err := json.Unmarshal(p.Data, &data)
	if err != nil {
		return err
	}
	w.Logger.Printf("User %v wants to start following user %v", c.clientID, data.ID)

	followRequestId, err := w.notificationService.CreateFollowRequest(int64(c.clientID), int64(data.ID))
	if err != nil {
		return err
	}

	if followRequestId == -1 {
		w.Logger.Printf("User %v now follows public user %v", c.clientID, data.ID)

		// sendNewChatlist
		userChatList, groupChatList, err := w.chatService.GetChatlist(int64(c.clientID))
		if err != nil {
			return err
		}

		dataToSend, err := json.Marshal(
			&ChatListPayload{
				UserID:        int(c.clientID),
				UserChatlist:  userChatList,
				GroupChatlist: groupChatList,
			},
		)

		if err != nil {
			return err
		}

		c.gate <- Payload{
			Type: "chatlist",
			Data: dataToSend,
		}

		w.Logger.Printf("Sent new chatlist to sender %v", c.clientID)

		return nil
	}

	w.Logger.Printf("Created follow request with id %v", followRequestId)

	// broadcast to recipient

	err = w.BroadcastFollowRequest(c, followRequestId, int64(data.ID))
	if err != nil {
		return err
	}

	return nil
}

func (w *WebsocketServer) UnfollowHandler(p Payload, c *Client) error {
	data := &RequestPayload{}
	err := json.Unmarshal(p.Data, &data)
	if err != nil {
		return err
	}
	w.Logger.Printf("User %v wants to unfollow user %v", c.clientID, data.ID)

	err = w.userService.Unfollow(int64(c.clientID), int64(data.ID))
	if err != nil {
		return err
	}
	w.Logger.Printf("User successfully %v unfollowed user %v", c.clientID, data.ID)

	return nil
}

func (w *WebsocketServer) RequestChatlistHandler(p Payload, c *Client) error {
	w.Logger.Printf("User %v has requested chatlist", c.clientID)

	userChatList, groupChatList, err := w.chatService.GetChatlist(int64(c.clientID))
	if err != nil {
		return err
	}

	w.Logger.Printf("Chatlist successfully retrieved (%v user chats, %v group chats)", len(userChatList), len(groupChatList))

	dataToSend, err := json.Marshal(
		&ChatListPayload{
			UserID:        int(c.clientID),
			UserChatlist:  userChatList,
			GroupChatlist: groupChatList,
		},
	)

	if err != nil {
		return err
	}

	c.gate <- Payload{
		Type: "chatlist",
		Data: dataToSend,
	}

	w.Logger.Printf("Sent chatlist to user %v", c.clientID)

	return nil
}

func (w *WebsocketServer) MessageHistoryHandler(p Payload, c *Client) error {
	data := &RequestPayload{}
	err := json.Unmarshal(p.Data, &data)
	if err != nil {
		return err
	}

	if data.ID == 0 && data.GroupID > 0 {
		w.Logger.Printf("User %v requests message history with group %v starting from %v", c.clientID, data.GroupID, data.LastMessage)
	} else if data.GroupID == 0 && data.ID > 0 {
		w.Logger.Printf("User %v requests message history with user %v starting from %v", c.clientID, data.ID, data.LastMessage)
	} else {
		w.Logger.Printf("Invalid request payload")
		return ErrorInvalidPayload
	}

	messages, err := w.chatService.GetMessageHistory(int64(c.clientID), int64(data.ID), int64(data.GroupID), int64(data.LastMessage))
	if err != nil {
		return err
	}

	w.Logger.Printf("Message history successfully retrieved (%v messages)", len(messages))

	dataToSend, err := json.Marshal(messages)

	if err != nil {
		return err
	}

	c.gate <- Payload{
		Type: "message_history",
		Data: dataToSend,
	}

	return nil
}

func (w *WebsocketServer) NewMessageHandler(p Payload, c *Client) error {
	data := &MessagePayload{}
	err := json.Unmarshal(p.Data, &data)
	if err != nil {
		return err
	}

	messageData := &models.Message{
		SenderId:    c.clientID,
		RecipientId: int64(data.RecipientID),
		GroupId:     int64(data.GroupID),
		Content:     data.Content,
		SentAt:      time.Now(),
	}

	messageID, err := w.chatService.CreateMessage(messageData)
	if err != nil {
		return err
	}

	messageData.Id = messageID

	if data.GroupID == 0 && data.RecipientID > 0 {

		w.Logger.Printf("User %v sent message %v to user %v", c.clientID, data.MessageID, data.RecipientID)
		defer func() {
			err = w.BroadcastSingleMessage(c, messageData)
			if err != nil {
				w.Logger.Printf("Error broadcasting message: %v", err)
			}
		}()

	} else if data.RecipientID == 0 && data.GroupID > 0 {

		w.Logger.Printf("User %v sent message %v to group %v", c.clientID, data.MessageID, data.GroupID)
		defer func() {
			err = w.BroadcastGroupMessage(c, messageData)
			if err != nil {
				w.Logger.Printf("Error broadcasting message: %v", err)
			}
		}()

	} else {

		w.Logger.Printf("Invalid request payload")
		return ErrorInvalidPayload

	}

	w.Logger.Printf("Message successfully created with id %v", messageID)

	return nil
}

func (w *WebsocketServer) GroupRequestHandler(p Payload, c *Client) error {

	//w.Logger.Printf("Payload: %s", p)

	data := &RequestPayload{}
	err := json.Unmarshal(p.Data, &data)
	if err != nil {
		return err
	}
	w.Logger.Printf("User %v wants to join group %v", c.clientID, data.GroupID)

	groupRequestId, err := w.notificationService.CreateGroupRequest(int64(c.clientID), int64(data.GroupID))

	if err != nil {
		return err
	}

	w.Logger.Printf("Created group request with id %v", groupRequestId)

	// broadcast to group owner

	err = w.BroadcastGroupJoinRequest(c, groupRequestId, int64(data.GroupID))
	if err != nil {
		return err
	}

	return nil
}

func (w *WebsocketServer) MessagesReadHandler(p Payload, c *Client) error {
	data := &RequestPayload{}
	err := json.Unmarshal(p.Data, &data)
	if err != nil {
		return err
	}

	w.Logger.Printf("Payload: %+v", *data)

	err = w.chatService.HandleMessagesRead(c.clientID, int64(data.LastMessage))
	if err != nil {
		if err.Error() == "not recipient" {
			// do not mark messages as read if user is not recipient
			return nil
		}
		return err
	}

	w.Logger.Printf("User %v has read message %v from user %v", c.clientID, data.LastMessage, data.ID)

	return nil
}
