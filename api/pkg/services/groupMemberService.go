package services

import (
	"SocialNetworkRestApi/api/pkg/models"
	"database/sql"
	"errors"
	"log"
	"time"
)

type IGroupMemberService interface {
	GetGroupMembers(groupId int64) ([]*models.SimpleUserJSON, error)
	GetMemberById(groupId int64, userId int64) (*models.GroupMember, error)
	AddMembers(userId int64, members models.GroupMemberJSON) ([]*models.NotificationJSON, error)
	GetMembersToAdd(groupId int64, userId int64) ([]*models.SimpleUserJSON, error)
}

type GroupMemberService struct {
	Logger                 *log.Logger
	UserRepository         models.IUserRepository
	NotificationRepository models.INotificationRepository
	GroupRepository        models.IGroupRepository
	GroupMemberRepository  models.IGroupMemberRepository
}

func InitGroupMemberService(
	logger *log.Logger,
	userRepo *models.UserRepository,
	notificationsRepo *models.NotificationRepository,
	groupRepository *models.GroupRepository,
	groupMemberRepo *models.GroupMemberRepository) *GroupMemberService {
	return &GroupMemberService{
		Logger:                 logger,
		UserRepository:         userRepo,
		NotificationRepository: notificationsRepo,
		GroupRepository:        groupRepository,
		GroupMemberRepository:  groupMemberRepo,
	}
}

func (s *GroupMemberService) GetGroupMembers(groupId int64) ([]*models.SimpleUserJSON, error) {

	members, err := s.GroupMemberRepository.GetGroupMembersByGroupId(groupId)

	if err != nil {
		s.Logger.Printf("Failed fetching group members: %s", err)
		return nil, err
	}

	simpleMembers := []*models.SimpleUserJSON{}

	for _, member := range members {

		if !member.Accepted {
			continue
		}

		userData, err := s.UserRepository.GetById(member.UserId)

		if err != nil {
			s.Logger.Printf("Failed fetching user data: %s", err)
			return nil, err
		}

		simpleMember := &models.SimpleUserJSON{
			Id:        int(member.UserId),
			Nickname:  userData.Nickname,
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
			ImagePath: userData.ImagePath,
		}

		simpleMembers = append(simpleMembers, simpleMember)
	}

	return simpleMembers, nil
}

func (s *GroupMemberService) GetMemberById(groupId int64, userId int64) (*models.GroupMember, error) {
	member, err := s.GroupMemberRepository.GetMemberByGroupId(groupId, userId)

	if err != nil && err != sql.ErrNoRows {
		s.Logger.Printf("Cannot validate user: %s", err)
		return nil, err
	}

	return member, err
}

func (s *GroupMemberService) AddMembers(userId int64, members models.GroupMemberJSON) ([]*models.NotificationJSON, error) {

	member, err := s.GroupMemberRepository.GetMemberByGroupId(int64(members.GroupId), userId)

	if err == sql.ErrNoRows {
		s.Logger.Printf("Cannot find user in group: %s", err)
		return nil, errors.New("not a member of this group")
	}

	if err != nil {
		s.Logger.Printf("Cannot validate user: %s", err)
		return nil, errors.New("error in checking if user is already member of group")
	}

	if !member.Accepted {
		s.Logger.Printf("User %d is not a member of this group", userId)
		return nil, errors.New("not a member of this group")
	}

	notificationDetails := &models.NotificationDetails{
		SenderId:         userId,
		NotificationType: "group_invite",
		EntityId:         int64(members.GroupId),
		CreatedAt:        time.Now(),
	}

	detailsId, err := s.NotificationRepository.InsertDetails(notificationDetails)

	if err != nil {
		s.Logger.Printf("Cannot insert notification details: %s", err)
		return nil, err
	}

	group, err := s.GroupRepository.GetById(int64(members.GroupId))
	if err != nil {
		s.Logger.Printf("Cannot get group: %s", err)
		return nil, err
	}

	userData, err := s.UserRepository.GetById(userId)
	if err != nil {
		s.Logger.Printf("Cannot get user: %s", err)
		return nil, err
	}

	if userData.Nickname == "" {
		userData.Nickname = userData.FirstName + " " + userData.LastName
	}

	notificationsToBroadcast := []*models.NotificationJSON{}

	for _, userIdToAdd := range members.UserIds {

		member, err := s.GroupMemberRepository.GetMemberByGroupId(int64(members.GroupId), int64(userIdToAdd))
		if err != sql.ErrNoRows {
			if err != nil {
				s.Logger.Printf("Cannot validate user: %s", err)
				return nil, err
			} else if member.Accepted {
				s.Logger.Printf("User %d is already a member of this group", userIdToAdd)
				continue
			}
		}

		groupMember := &models.GroupMember{
			UserId:   int64(userIdToAdd),
			GroupId:  int64(members.GroupId),
			JoinedAt: time.Now(),
			Accepted: false,
		}

		_, err = s.GroupMemberRepository.Insert(groupMember)

		if err != nil {
			s.Logger.Printf("Cannot add user %d to group %d: %s", userIdToAdd, members.GroupId, err)
			return nil, err
		}

		s.Logger.Printf("User %d added to group %d", userIdToAdd, members.GroupId)

		// send notification to user

		notification := &models.Notification{
			ReceiverId:            int64(userIdToAdd),
			NotificationDetailsId: detailsId,
		}

		notificationId, err := s.NotificationRepository.InsertNotification(notification)

		if err != nil {
			s.Logger.Printf("Cannot insert notification: %s", err)
			return nil, err
		}

		// broadcast notification to users

		notificationJSON := &models.NotificationJSON{
			ReceiverId:       int64(userIdToAdd),
			NotificationType: notificationDetails.NotificationType,
			NotificationId:   notificationId,
			SenderId:         userId,
			SenderName:       userData.Nickname,
			GroupId:          int64(members.GroupId),
			GroupName:        group.Title,
		}

		notificationsToBroadcast = append(notificationsToBroadcast, notificationJSON)

	}

	return notificationsToBroadcast, nil
}

func (s *GroupMemberService) GetMembersToAdd(groupId int64, userId int64) ([]*models.SimpleUserJSON, error) {

	publicUsers, err := s.UserRepository.GetAllUsers(userId)

	if err != nil {
		s.Logger.Printf("Failed fetching public users: %s", err)
		return nil, err
	}

	simpleMembers := map[int64]*models.SimpleUserJSON{}

	for _, user := range publicUsers {

		if !user.IsPublic {
			continue
		}

		member, err := s.GroupMemberRepository.GetMemberByGroupId(groupId, user.Id)

		if err != sql.ErrNoRows {
			if err != nil {
				s.Logger.Printf("Cannot validate user: %s", err)
				return nil, err
			} else if member.Accepted {
				continue
			}
		}

		simpleMember := &models.SimpleUserJSON{
			Id:        int(user.Id),
			Nickname:  user.Nickname,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			ImagePath: user.ImagePath,
		}

		simpleMembers[user.Id] = simpleMember

	}

	followers, err := s.UserRepository.GetAllUserFollowers(userId)

	if err != nil {
		s.Logger.Printf("Failed fetching user followers: %s", err)
		return nil, err
	}

	for _, follower := range followers {

		if simpleMembers[follower.Id] != nil {
			continue
		}

		member, err := s.GroupMemberRepository.GetMemberByGroupId(groupId, follower.Id)

		if err != sql.ErrNoRows {
			if err != nil {
				s.Logger.Printf("Cannot validate user: %s", err)
				return nil, err
			} else if member.Accepted {
				continue
			}
		}

		simpleMember := &models.SimpleUserJSON{
			Id:        int(follower.Id),
			Nickname:  follower.Nickname,
			FirstName: follower.FirstName,
			LastName:  follower.LastName,
			ImagePath: follower.ImagePath,
		}

		simpleMembers[follower.Id] = simpleMember
	}

	following, err := s.UserRepository.GetAllFollowedBy(userId)

	if err != nil {
		s.Logger.Printf("Failed fetching users following: %s", err)
		return nil, err
	}

	for _, followed := range following {

		if simpleMembers[followed.Id] != nil {
			continue
		}

		member, err := s.GroupMemberRepository.GetMemberByGroupId(groupId, followed.Id)

		if err != sql.ErrNoRows {
			if err != nil {
				s.Logger.Printf("Cannot validate user: %s", err)
				return nil, err
			} else if member.Accepted {
				continue
			}
		}

		simpleMember := &models.SimpleUserJSON{
			Id:        int(followed.Id),
			Nickname:  followed.Nickname,
			FirstName: followed.FirstName,
			LastName:  followed.LastName,
			ImagePath: followed.ImagePath,
		}

		simpleMembers[followed.Id] = simpleMember
	}

	simpleMembersArray := make([]*models.SimpleUserJSON, 0, len(simpleMembers))

	for _, simpleMember := range simpleMembers {
		simpleMembersArray = append(simpleMembersArray, simpleMember)
	}

	return simpleMembersArray, nil
}
