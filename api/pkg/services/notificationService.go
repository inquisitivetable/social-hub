package services

import (
	"SocialNetworkRestApi/api/pkg/models"
	"database/sql"
	"errors"
	"log"
	"time"
)

type INotificationService interface {
	GetById(notificationId int64) (*models.Notification, error)
	GetDetailsById(notificationId int64) (*models.NotificationDetails, error)
	GetByEventAndUserId(eventId int64, userId int64) (*models.Notification, error)
	GetUserNotifications(userId int64) ([]*models.NotificationJSON, error)
	CreateFollowRequest(followerId int64, followingId int64) (int64, error)
	HandleFollowRequest(notificationId int64, accepted bool) error
	CreateGroupRequest(senderId int64, groupId int64) (int64, error)
	HandleGroupRequest(creatorID int64, notificationID int64, accepted bool) error
	HandleEventInvite(notificationID int64, accepted bool) error
	CreateGroupInvite(creatorId int64, groupId int64, membersToAdd []int64) ([]*models.NotificationJSON, error)
	HandleGroupInvite(notificationID int64, accepted bool) error
}

type NotificationService struct {
	Logger                 *log.Logger
	UserRepo               models.IUserRepository
	FollowerRepo           models.IFollowerRepository
	NotificationRepository models.INotificationRepository
	GroupRepo              models.IGroupRepository
	GroupMemberRepo        models.IGroupMemberRepository
	EventRepo              models.IEventRepository
	EventAttendanceRepo    models.IEventAttendanceRepository
}

func InitNotificationService(
	logger *log.Logger,
	userRepo *models.UserRepository,
	followerRepo *models.FollowerRepository,
	notificationRepo *models.NotificationRepository,
	groupRepo *models.GroupRepository,
	groupMemberRepo *models.GroupMemberRepository,
	eventRepo *models.EventRepository,
	eventAttendanceRepo *models.EventAttendanceRepository,
) *NotificationService {
	return &NotificationService{
		Logger:                 logger,
		UserRepo:               userRepo,
		FollowerRepo:           followerRepo,
		NotificationRepository: notificationRepo,
		GroupRepo:              groupRepo,
		GroupMemberRepo:        groupMemberRepo,
		EventRepo:              eventRepo,
		EventAttendanceRepo:    eventAttendanceRepo,
	}
}

func (s *NotificationService) GetById(notificationId int64) (*models.Notification, error) {

	notification, err := s.NotificationRepository.GetById(notificationId)
	if err != nil {
		s.Logger.Printf("Cannot get notification: %s", err)
		return nil, err
	}

	s.Logger.Printf("Notification returned: %d", notification.Id)

	return notification, nil
}

func (s *NotificationService) GetByEventAndUserId(eventId int64, userId int64) (*models.Notification, error) {

	notification, err := s.NotificationRepository.GetByEventAndUserId(eventId, userId)
	if err != nil {
		s.Logger.Printf("Cannot get notification: %s", err)
		return nil, err
	}

	return notification, nil

}

func (s *NotificationService) GetDetailsById(notificationId int64) (*models.NotificationDetails, error) {

	notificationDetails, err := s.NotificationRepository.GetDetailsById(notificationId)
	if err != nil {
		s.Logger.Printf("Cannot get notification details: %s", err)
		return nil, err
	}

	s.Logger.Printf("Notification details returned: %d", notificationDetails.Id)

	return notificationDetails, nil
}

func (s *NotificationService) GetUserNotifications(userId int64) ([]*models.NotificationJSON, error) {

	notifications, err := s.NotificationRepository.GetByReceiverId(userId)
	if err != nil {
		s.Logger.Printf("Cannot get user notifications: %s", err)
		return nil, err
	}

	NotificationJSON := []*models.NotificationJSON{}

	for _, notification := range notifications {

		if notification.Reaction.Valid {
			s.Logger.Printf("Notification already processed: %d", notification.Id)
			continue
		}

		notificationDetails, err := s.NotificationRepository.GetDetailsById(notification.NotificationDetailsId)
		if err != nil {
			s.Logger.Printf("Cannot get notification details: %s", err)
			return nil, err
		}

		singleNotification := &models.NotificationJSON{
			ReceiverId:       userId,
			NotificationType: notificationDetails.NotificationType,
			NotificationId:   notification.Id,
			SenderId:         notificationDetails.SenderId,
		}
		sender, err := s.UserRepo.GetById(notificationDetails.SenderId)
		if err != nil {
			s.Logger.Printf("Cannot get sender: %s", err)
			return nil, err
		}
		if sender.Nickname == "" {
			singleNotification.SenderName = sender.FirstName + " " + sender.LastName
		} else {
			singleNotification.SenderName = sender.Nickname
		}

		if notificationDetails.NotificationType == "group_invite" {
			group, err := s.GroupRepo.GetById(notificationDetails.EntityId)
			if err != nil {
				s.Logger.Printf("Cannot get group: %s", err)
				return nil, err
			}
			singleNotification.GroupId = group.Id
			singleNotification.GroupName = group.Title
		}

		if notificationDetails.NotificationType == "group_request" {
			member, err := s.GroupMemberRepo.GetById(notificationDetails.EntityId)

			if err != nil {
				s.Logger.Printf("Cannot get group member: %s", err)
				return nil, err
			}

			group, err := s.GroupRepo.GetById(member.GroupId)
			if err != nil {
				s.Logger.Printf("Cannot get group: %s", err)
				return nil, err
			}
			singleNotification.GroupId = group.Id
			singleNotification.GroupName = group.Title
		}

		if notificationDetails.NotificationType == "event_invite" {
			//s.Logger.Printf("Getting event: %d", notificationDetails.EntityId)
			event, err := s.EventRepo.GetById(notificationDetails.EntityId)
			if err != nil {
				s.Logger.Printf("Cannot get event: %s", err)
				return nil, err
			}
			singleNotification.GroupId = event.GroupId
			group, err := s.GroupRepo.GetById(event.GroupId)
			if err != nil {
				s.Logger.Printf("Cannot get group: %s", err)
				return nil, err
			}
			singleNotification.GroupName = group.Title
			singleNotification.EventId = event.Id
			singleNotification.EventName = event.Title
			singleNotification.EventDate = event.EventTime
		}

		NotificationJSON = append(NotificationJSON, singleNotification)
	}

	s.Logger.Printf("User notifications returned: %d", len(NotificationJSON))

	return NotificationJSON, nil
}

func (s *NotificationService) CreateFollowRequest(followerId int64, followingId int64) (int64, error) {

	// check if follower and following exist
	_, err := s.UserRepo.GetById(followerId)
	if err != nil {
		s.Logger.Printf("Follower not found: %s", err)
		return -1, err
	}
	following, err := s.UserRepo.GetById(followingId)
	if err != nil {
		s.Logger.Printf("Following not found: %s", err)
		return -1, err
	}

	// check if follow request already exists
	_, err = s.FollowerRepo.GetByFollowerAndFollowing(followerId, followingId)
	if err == nil {
		return -1, errors.New("follow request already exists")
	}

	accepted := sql.NullBool{Bool: false, Valid: false}
	if following.IsPublic {
		accepted = sql.NullBool{Bool: true, Valid: true}
	}
	follower := &models.Follower{
		FollowerId:  followerId,
		FollowingId: followingId,
		Accepted:    accepted,
	}

	// create follow request
	lastID, err := s.FollowerRepo.Insert(follower)
	if err != nil {
		s.Logger.Printf("Cannot insert follow request: %s", err)
		return -1, err
	}

	s.Logger.Printf("Follow request created: %d", lastID)

	// check if following is private
	if following.IsPublic {
		return -1, nil
	}

	// create notification
	notificationDetails := models.NotificationDetails{
		SenderId:         followerId,
		NotificationType: "follow_request",
		EntityId:         lastID,
		CreatedAt:        time.Now(),
	}

	notificationDetailsId, err := s.NotificationRepository.InsertDetails(&notificationDetails)
	if err != nil {
		return -1, err
	}

	notification := models.Notification{
		ReceiverId:            followingId,
		NotificationDetailsId: notificationDetailsId,
		Reaction:              sql.NullBool{Bool: false, Valid: false},
	}

	notificationId, err := s.NotificationRepository.InsertNotification(&notification)
	if err != nil {
		return -1, err
	}

	return notificationId, nil
}

func (s *NotificationService) HandleFollowRequest(notificationId int64, accepted bool) error {

	notification, err := s.NotificationRepository.GetById(notificationId)
	if err != nil {
		s.Logger.Printf("Cannot get notification: %s", err)
		return err
	}

	notificationDetails, err := s.NotificationRepository.GetDetailsById(notification.NotificationDetailsId)
	if err != nil {
		s.Logger.Printf("Cannot get notification details: %s", err)
		return err
	}

	// check if follow request already handled
	if notification.Reaction.Valid {
		return errors.New("follow request already handled")
	}

	// check if follow request exists
	follower, err := s.FollowerRepo.GetById(notificationDetails.EntityId)
	if err != nil {
		s.Logger.Printf("Cannot get follow request: %s", err)
		return err
	}

	// check if follow request is accepted
	if follower.Accepted.Valid {
		return errors.New("follow request already accepted")
	}

	// update follow request
	if accepted {
		follower.Accepted = sql.NullBool{Bool: true, Valid: true}
		err = s.FollowerRepo.Update(follower)
		if err != nil {
			s.Logger.Printf("Cannot update follow request: %s", err)
			return err
		}
	} else {
		err = s.FollowerRepo.Delete(follower)
		if err != nil {
			s.Logger.Printf("Cannot delete follow request: %s", err)
			return err
		}
	}

	s.Logger.Printf("Follow request updated: %d", follower.Id)

	// update notification
	notification.Reaction = sql.NullBool{Bool: true, Valid: true}
	err = s.NotificationRepository.Update(notification)
	if err != nil {
		s.Logger.Printf("Cannot update notification: %s", err)
		return err
	}

	s.Logger.Printf("Notification updated: %d", notification.Id)

	return nil
}

func (s *NotificationService) CreateGroupRequest(senderId int64, groupId int64) (int64, error) {

	// check if sender and group exist
	_, err := s.UserRepo.GetById(senderId)
	if err != nil {
		s.Logger.Printf("Sender not found: %s", err)
		return -1, err
	}
	groupData, err := s.GroupRepo.GetById(groupId)
	if err != nil {
		s.Logger.Printf("Group not found: %s", err)
		return -1, err
	}

	// check if user is already member of group
	member, err := s.GroupMemberRepo.GetMemberByGroupId(groupId, senderId)
	if err != sql.ErrNoRows {
		if err != nil {
			s.Logger.Printf("Cannot validate user: %s", err)
			return -1, errors.New("error in checking if user is already member of group")
		} else if member.Accepted {
			s.Logger.Printf("User %d is already a member of this group", senderId)
			return -1, errors.New("already a member of this group")
		} else if !member.Accepted {
			s.Logger.Printf("User %d already has a pending request for this group", senderId)
			return -1, errors.New("already has a pending request for this group")
		}
	}

	// add member to group with joined at Zero
	groupMember := &models.GroupMember{
		UserId:  senderId,
		GroupId: groupId,
	}

	lastID, err := s.GroupMemberRepo.Insert(groupMember)
	if err != nil {
		s.Logger.Printf("Cannot insert group request: %s", err)
		return -1, err
	}

	s.Logger.Printf("Member added: %d", lastID)

	// create notification

	notifcationDetails := models.NotificationDetails{
		SenderId:         senderId,
		NotificationType: "group_request",
		EntityId:         lastID,
		CreatedAt:        time.Now(),
	}

	notifcationDetailsId, err := s.NotificationRepository.InsertDetails(&notifcationDetails)
	if err != nil {
		return -1, err
	}

	notification := models.Notification{
		ReceiverId:            groupData.CreatorId,
		NotificationDetailsId: notifcationDetailsId,
		Reaction:              sql.NullBool{Bool: false, Valid: false},
	}

	notificationId, err := s.NotificationRepository.InsertNotification(&notification)
	if err != nil {
		return -1, err
	}

	return notificationId, nil
}

func (s *NotificationService) HandleGroupRequest(creatorID int64, notificationID int64, accepted bool) error {

	notification, err := s.NotificationRepository.GetById(notificationID)
	if err != nil {
		s.Logger.Printf("Cannot get notification: %s", err)
		return err
	}

	// check if group request already handled
	if notification.Reaction.Valid {
		return errors.New("group request already handled")
	}

	notificationDetails, err := s.NotificationRepository.GetDetailsById(notification.NotificationDetailsId)
	if err != nil {
		s.Logger.Printf("Cannot get notification details: %s", err)
		return err
	}

	// check if group request exists
	groupMember, err := s.GroupMemberRepo.GetById(notificationDetails.EntityId)
	if err == sql.ErrNoRows {
		s.Logger.Printf("Group request not found: %s", err)
		return err
	}

	if err != nil {
		s.Logger.Printf("Cannot validate request: %s", err)
		return err
	}

	if groupMember.Accepted {
		s.Logger.Printf("Group request already accepted: %d", notificationDetails.EntityId)
		return errors.New("group request already accepted")
	}

	// check if user is creator of group
	group, err := s.GroupRepo.GetById(groupMember.GroupId)
	if err != nil {
		s.Logger.Printf("Cannot get group: %s", err)
		return err
	}

	if group.CreatorId != creatorID {
		return errors.New("user is not creator of group")
	}

	// update group request
	if accepted {
		groupMember.JoinedAt = time.Now()
		groupMember.Accepted = true
		err = s.GroupMemberRepo.Update(groupMember)
		if err != nil {
			s.Logger.Printf("Cannot update group request: %s", err)
			return err
		}
	} else {
		err = s.GroupMemberRepo.Delete(groupMember)
		if err != nil {
			s.Logger.Printf("Cannot delete group request: %s", err)
			return err
		}
	}

	s.Logger.Printf("Group request updated: %d", notificationDetails.EntityId)

	// update notification
	notification.Reaction = sql.NullBool{Bool: true, Valid: true}
	err = s.NotificationRepository.Update(notification)
	if err != nil {
		s.Logger.Printf("Cannot update notification: %s", err)
		return err
	}

	s.Logger.Printf("Notification updated: %d", notification.Id)

	return nil

}

func (s *NotificationService) HandleEventInvite(notificationID int64, accepted bool) error {

	notification, err := s.NotificationRepository.GetById(notificationID)
	if err != nil {
		s.Logger.Printf("Cannot get notification: %s", err)
		return err
	}

	// check if event invite already handled
	if notification.Reaction.Valid {
		return errors.New("event invite already handled")
	}

	notificationDetails, err := s.NotificationRepository.GetDetailsById(notification.NotificationDetailsId)
	if err != nil {
		s.Logger.Printf("Cannot get notification details: %s", err)
		return err
	}

	// check if event exists
	event, err := s.EventRepo.GetById(notificationDetails.EntityId)
	if err != nil {
		s.Logger.Printf("Cannot get event invite: %s", err)
		return err
	}

	// check if event invite has been processed
	attendees, err := s.EventAttendanceRepo.GetAttendeesByEventId(event.Id)
	if err != nil {
		s.Logger.Printf("Cannot get event attendees: %s", err)
		return err
	}

	for _, attendee := range attendees {
		if attendee.UserId == notification.ReceiverId {
			return errors.New("event invite already processed")
		}
	}

	// update event invite
	eventAttendance := &models.EventAttendance{
		UserId:      notification.ReceiverId,
		EventId:     event.Id,
		IsAttending: accepted,
	}
	_, err = s.EventAttendanceRepo.Insert(eventAttendance)
	if err != nil {
		s.Logger.Printf("Cannot insert event attendance: %s", err)
		return err
	}

	// update notification
	notification.Reaction = sql.NullBool{Bool: true, Valid: true}
	err = s.NotificationRepository.Update(notification)
	if err != nil {
		s.Logger.Printf("Cannot update notification: %s", err)
		return err
	}

	s.Logger.Printf("Event attendance and notification updated: %d", notificationDetails.EntityId)

	return nil

}

func (s *NotificationService) CreateGroupInvite(creatorId int64, groupId int64, membersToAdd []int64) ([]*models.NotificationJSON, error) {

	// check if user is creator of group
	group, err := s.GroupRepo.GetById(groupId)
	if err != nil {
		s.Logger.Printf("Cannot get group: %s", err)
		return nil, err
	}

	if group.CreatorId != creatorId {
		return nil, errors.New("user is not creator of group")
	}

	// create notification

	notificationDetails := &models.NotificationDetails{
		SenderId:         creatorId,
		NotificationType: "group_invite",
		EntityId:         groupId,
		CreatedAt:        time.Now(),
	}

	detailsId, err := s.NotificationRepository.InsertDetails(notificationDetails)
	if err != nil {
		s.Logger.Printf("Cannot insert notification details: %s", err)
		return nil, err
	}

	creatorData, err := s.UserRepo.GetById(creatorId)
	if err != nil {
		s.Logger.Printf("Cannot get creator data: %s", err)
		return nil, err
	}

	if creatorData.Nickname != "" {
		creatorData.Nickname = creatorData.FirstName + " " + creatorData.LastName
	}

	notificationsToBroadcast := []*models.NotificationJSON{}

	for _, memberToAdd := range membersToAdd {
		notification := &models.Notification{
			ReceiverId:            memberToAdd,
			NotificationDetailsId: detailsId,
		}

		notificationId, err := s.NotificationRepository.InsertNotification(notification)

		if err != nil {
			s.Logger.Printf("Cannot insert notification: %s", err)
			return nil, err
		}

		// broadcast notification to users

		notificationJSON := &models.NotificationJSON{
			ReceiverId:       memberToAdd,
			NotificationType: notificationDetails.NotificationType,
			NotificationId:   notificationId,
			SenderId:         creatorId,
			SenderName:       creatorData.Nickname,
			GroupId:          groupId,
			GroupName:        group.Title,
		}

		notificationsToBroadcast = append(notificationsToBroadcast, notificationJSON)

	}

	return notificationsToBroadcast, nil
}

func (s *NotificationService) HandleGroupInvite(notificationID int64, accepted bool) error {
	notification, err := s.NotificationRepository.GetById(notificationID)
	if err != nil {
		s.Logger.Printf("Cannot get notification: %s", err)
		return err
	}

	// check if group invite already handled
	if notification.Reaction.Valid {
		return errors.New("group invite already handled")
	}

	notificationDetails, err := s.NotificationRepository.GetDetailsById(notification.NotificationDetailsId)
	if err != nil {
		s.Logger.Printf("Cannot get notification details: %s", err)
		return err
	}

	// check if group invite exists
	groupMember, err := s.GroupMemberRepo.GetMemberByGroupId(notificationDetails.EntityId, notification.ReceiverId)
	if err != sql.ErrNoRows {
		if err != nil {
			s.Logger.Printf("Cannot validate user: %s", err)
			return err
		} else if groupMember.Accepted {
			s.Logger.Printf("Group invite already accepted: %d", notificationDetails.EntityId)
			return errors.New("group invite already accepted")
		}
	}

	// update group invite
	if accepted {
		groupMember.JoinedAt = time.Now()
		groupMember.Accepted = true
		err = s.GroupMemberRepo.Update(groupMember)
		if err != nil {
			s.Logger.Printf("Cannot update group invite: %s", err)
			return err
		}
	} else {
		err = s.GroupMemberRepo.Delete(groupMember)
		if err != nil {
			s.Logger.Printf("Cannot delete group invite: %s", err)
			return err
		}
	}

	s.Logger.Printf("Group invite updated: %d", notificationDetails.EntityId)

	// update notification
	notification.Reaction = sql.NullBool{Bool: true, Valid: true}
	err = s.NotificationRepository.Update(notification)
	if err != nil {
		s.Logger.Printf("Cannot update notification: %s", err)
		return err
	}

	s.Logger.Printf("Notification updated: %d", notification.Id)

	return nil
}
