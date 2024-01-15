package services

import (
	"SocialNetworkRestApi/api/internal/server/utils"
	"SocialNetworkRestApi/api/pkg/models"
	"errors"
	"log"
	"mime/multipart"
	"strings"
	"time"
)

type IGroupService interface {
	GetUserGroups(userId int64) ([]*models.UserGroup, error)
	GetUserCreatedGroups(userId int64) ([]*models.UserGroup, error)
	GetGroupById(groupId int64) (models.GroupJSON, error)
	SearchGroupsAndUsers(userId int64, searchString string) ([]*models.SearchResult, error)
	CreateGroup(groupFormData *models.GroupJSON, userId int64) (int64, error)
	UpdateGroupImage(userId int64, groupId int64, imageFile multipart.File, header *multipart.FileHeader) error
	GetGroupCreator(groupId int64) (*models.User, error)
}

type GroupService struct {
	Logger          *log.Logger
	GroupRepository models.IGroupRepository
	GroupMemberRepo models.IGroupMemberRepository
	UserRepository  models.IUserRepository
}

func InitGroupService(
	logger *log.Logger,
	groupRepo *models.GroupRepository,
	groupMemberRepo *models.GroupMemberRepository,
	userRepo *models.UserRepository,
) *GroupService {
	return &GroupService{
		Logger:          logger,
		GroupRepository: groupRepo,
		GroupMemberRepo: groupMemberRepo,
		UserRepository:  userRepo,
	}
}

func (s *GroupService) GetUserGroups(userId int64) ([]*models.UserGroup, error) {

	result, err := s.GroupRepository.GetAllByMemberId(userId)

	if err != nil {
		s.Logger.Printf("Failed fetching groups: %s", err)
	}

	groups := []*models.UserGroup{}

	for _, p := range result {
		groups = append(groups, &models.UserGroup{
			Id:    p.Id,
			Title: p.Title,
		})
	}

	return groups, nil
}

func (s *GroupService) GetUserCreatedGroups(userId int64) ([]*models.UserGroup, error) {

	result, err := s.GroupRepository.GetAllByCreatorId(userId)

	if err != nil {
		s.Logger.Printf("Failed fetching groups: %s", err)
	}

	groups := []*models.UserGroup{}

	for _, p := range result {
		groups = append(groups, &models.UserGroup{
			Id:    p.Id,
			Title: p.Title,
		})
	}

	return groups, nil
}

func (s *GroupService) GetGroupById(groupId int64) (models.GroupJSON, error) {
	result, err := s.GroupRepository.GetById(groupId)

	group := models.GroupJSON{
		Title:       result.Title,
		Description: result.Description,
		ImagePath:   result.ImagePath,
	}

	if err != nil {
		s.Logger.Printf("Failed fetching groups: %s", err)
	}

	return group, err
}

func (s *GroupService) SearchGroupsAndUsers(userId int64, searchString string) ([]*models.SearchResult, error) {

	result, err := s.GroupRepository.SearchGroupsAndUsersByString(userId, searchString)

	if err != nil {
		s.Logger.Printf("Failed searching groups: %s", err)
	}

	return result, err
}

func (s *GroupService) CreateGroup(groupFormData *models.GroupJSON, userId int64) (int64, error) {
	group := &models.Group{
		CreatorId:   userId,
		ImagePath:   groupFormData.ImagePath,
		Title:       groupFormData.Title,
		Description: groupFormData.Description,
	}

	result, err := s.GroupRepository.Insert(group)

	if err != nil {
		s.Logger.Printf("Failed inserting group: %s", err)
		return -1, err
	}

	creator := &models.GroupMember{
		UserId:   userId,
		GroupId:  result,
		JoinedAt: time.Now(),
		Accepted: true,
	}

	_, err = s.GroupMemberRepo.Insert(creator)

	if err != nil {
		s.Logger.Printf("Failed inserting group member: %s", err)
		return -1, err
	}

	return result, err
}

func (s *GroupService) UpdateGroupImage(userId int64, groupId int64, imageFile multipart.File, header *multipart.FileHeader) error {

	// check if group exists
	_, err := s.GroupRepository.GetById(groupId)
	if err != nil {
		s.Logger.Printf("User not found: %s", err)
		return err
	}

	// check if user is creator of group
	creatorUser, err := s.GetGroupCreator(groupId)
	if err != nil {
		s.Logger.Printf("Error with checking group ownership: %s", err)
		return err
	}

	if userId != creatorUser.Id {
		return errors.New("user is not creator of group")
	}

	// check if file is an image
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image") {
		s.Logger.Println("Not an image")
		return errors.New("not an image")
	}

	// save image
	imagePath, err := utils.SaveImage(imageFile, header)
	if err != nil {
		s.Logger.Printf("Cannot save image: %s", err)
		return err
	}

	// update user image path
	err = s.GroupRepository.UpdateImagePath(groupId, imagePath)
	if err != nil {
		return err
	}

	return nil
}

func (s *GroupService) GetGroupCreator(groupId int64) (*models.User, error) {
	group, err := s.GroupRepository.GetById(groupId)
	if err != nil {
		s.Logger.Printf("Group not found: %s", err)
		return nil, err
	}

	userData, err := s.UserRepository.GetById(group.CreatorId)
	if err != nil {
		s.Logger.Printf("User not found: %s", err)
		return nil, err
	}

	return userData, nil
}
