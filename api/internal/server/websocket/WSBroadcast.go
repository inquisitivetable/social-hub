package websocket

import (
	"SocialNetworkRestApi/api/pkg/models"
	"encoding/json"
)

func (w *WebsocketServer) BroadcastFollowRequest(c *Client, followRequestId int64, otherId int64) error {
	userData, err := w.userService.GetUserByID(int64(c.clientID))
	if err != nil {
		return err
	}

	recipientClient := w.getClientByUserID(otherId)

	if recipientClient == nil {
		w.Logger.Printf("Recipient client not found (recipient offline)")
		return nil
	}

	w.Logger.Printf("Recipient client found (recipient online)")

	dataToSend, err := json.Marshal(
		&NotificationPayload{
			NotificationType: "follow_request",
			NotificationID:   int(followRequestId),
			SenderID:         int(c.clientID),
			SenderName:       userData.FirstName + " " + userData.LastName,
		},
	)

	if err != nil {
		return err
	}

	recipientClient.gate <- Payload{
		Type: "notification",
		Data: dataToSend,
	}

	w.Logger.Printf("Sent notification to recipient")

	return nil
}

func (w *WebsocketServer) BroadcastSingleMessage(c *Client, message *models.Message) error {

	recipientClient := w.getClientByUserID(message.RecipientId)

	if recipientClient == nil {
		w.Logger.Printf("Recipient client not found (recipient offline)")
	} else {
		w.Logger.Printf("Recipient client found (recipient online)")

		userData, err := w.userService.GetUserByID(c.clientID)
		if err != nil {
			return err
		}

		if userData.Nickname == "" {
			userData.Nickname = userData.FirstName + " " + userData.LastName
		}

		recipientData, err := w.userService.GetUserData(c.clientID, message.RecipientId)
		if err != nil {
			return err
		}

		if recipientData.Nickname == "" {
			recipientData.Nickname = recipientData.FirstName + " " + recipientData.LastName
		}

		dataToSend, err := json.Marshal(
			&MessagePayload{
				MessageID:     int(message.Id),
				SenderID:      int(c.clientID),
				SenderName:    userData.Nickname,
				SenderImage:   userData.ImagePath,
				RecipientID:   recipientData.UserID,
				RecipientName: recipientData.Nickname,
				//GroupID:       data.GroupID,
				//GroupName:     data.GroupName,
				Content:   message.Content,
				Timestamp: message.SentAt,
			},
		)

		if err != nil {
			return err
		}

		recipientClient.gate <- Payload{
			Type: "message",
			Data: dataToSend,
		}

		w.Logger.Printf("Sent message to recipient")

	}

	return nil
}

func (w *WebsocketServer) BroadcastGroupMessage(c *Client, message *models.Message) error {

	recipientUsers, err := w.groupMemberService.GetGroupMembers(int64(message.GroupId))
	if err != nil {
		return err
	}

	recipientClients := []*Client{}

	for _, member := range recipientUsers {
		getGlient := w.getClientByUserID(int64(member.Id))
		if getGlient != nil {
			recipientClients = append(recipientClients, getGlient)
		}
	}

	groupName, err := w.groupService.GetGroupById(message.GroupId)
	if err != nil {
		return err
	}

	if len(recipientClients) == 0 {
		w.Logger.Printf("Recipient clients not found (all recipients offline)")
	} else {
		w.Logger.Printf("Recipient clients found (%d recipients online)", len(recipientClients))

		userData, err := w.userService.GetUserByID(c.clientID)
		if err != nil {
			return err
		}

		if userData.Nickname == "" {
			userData.Nickname = userData.FirstName + " " + userData.LastName
		}

		for _, recipientClient := range recipientClients {

			if recipientClient.clientID == c.clientID {
				// skip the sender
				continue
			}

			// pick the recipient user from the list of group members
			var recipientUser *models.User
			for _, user := range recipientUsers {
				if int64(user.Id) == recipientClient.clientID {
					recipientUser = &models.User{
						Id:       int64(user.Id),
						Nickname: user.Nickname,
					}
					break
				}
			}

			if recipientUser.Nickname == "" {
				recipientUser.Nickname = recipientUser.FirstName + " " + recipientUser.LastName
			}

			dataToSend, err := json.Marshal(
				&MessagePayload{
					MessageID:     int(message.Id),
					SenderID:      int(c.clientID),
					SenderName:    userData.Nickname,
					SenderImage:   userData.ImagePath,
					RecipientID:   int(recipientUser.Id),
					RecipientName: recipientUser.Nickname,
					GroupID:       int(message.GroupId),
					GroupName:     groupName.Title,
					Content:       message.Content,
					Timestamp:     message.SentAt,
				},
			)

			if err != nil {
				return err
			}

			recipientClient.gate <- Payload{
				Type: "message",
				Data: dataToSend,
			}
		}

		w.Logger.Printf("Sent message to recipient")

	}

	return nil
}

func (w *WebsocketServer) BroadcastGroupNotifications(notifications []*models.NotificationJSON) error {

	for _, notification := range notifications {

		recipientClient := w.getClientByUserID(notification.ReceiverId)

		if recipientClient == nil {
			w.Logger.Printf("Recipient client not found (recipient offline)")
		} else {
			w.Logger.Printf("Recipient client found (recipient online)")

			dataToSend, err := json.Marshal(
				&NotificationPayload{
					NotificationType: notification.NotificationType,
					NotificationID:   int(notification.NotificationId),
					SenderID:         int(notification.SenderId),
					SenderName:       notification.SenderName,
					GroupID:          int(notification.GroupId),
					GroupName:        notification.GroupName,
					EventID:          int(notification.EventId),
					EventName:        notification.EventName,
					EventDate:        notification.EventDate,
				},
			)

			if err != nil {
				return err
			}

			recipientClient.gate <- Payload{
				Type: "notification",
				Data: dataToSend,
			}

			w.Logger.Printf("Sent event notification to recipient")

		}
	}

	return nil
}

func (w *WebsocketServer) BroadcastGroupJoinRequest(c *Client, groupReqId int64, groupId int64) error {
	userData, err := w.userService.GetUserByID(c.clientID)
	if err != nil {
		return err
	}

	if userData.Nickname == "" {
		userData.Nickname = userData.FirstName + " " + userData.LastName
	}

	creatorUser, err := w.groupService.GetGroupCreator(groupId)
	if err != nil {
		return err
	}

	creatorClient := w.getClientByUserID(creatorUser.Id)

	if creatorClient == nil {
		w.Logger.Printf("Group creator client not found (creator offline)")
	} else {
		w.Logger.Printf("Group creator client found (creator online)")

		groupData, err := w.groupService.GetGroupById(groupId)
		if err != nil {
			return err
		}

		dataToSend, err := json.Marshal(
			&NotificationPayload{
				NotificationType: "group_request",
				NotificationID:   int(groupReqId),
				SenderID:         int(c.clientID),
				SenderName:       userData.Nickname,
				GroupID:          int(groupId),
				GroupName:        groupData.Title,
			},
		)

		if err != nil {
			return err
		}

		creatorClient.gate <- Payload{
			Type: "notification",
			Data: dataToSend,
		}

		w.Logger.Printf("Sent group join request to creator")

	}

	return nil
}
