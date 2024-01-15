package services

import (
	"SocialNetworkRestApi/api/internal/server/utils"
	"SocialNetworkRestApi/api/pkg/enums"
	"SocialNetworkRestApi/api/pkg/models"
	"errors"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type IPostService interface {
	CreatePost(post *models.Post) error
	CreateGroupPost(post *models.Post) error
	GetFeedPosts(userId int64, offset int64) ([]*feedPostJSON, error)
	GetProfilePosts(userId int64, offset int64) ([]*feedPostJSON, error)
	GetGroupPosts(groupId int64, offset int64) ([]*feedPostJSON, error)
	GetUserPosts(userId int64, offset int64, requestingUserId int64) ([]*feedPostJSON, error)
	SavePostImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type PostService struct {
	Logger                *log.Logger
	GroupRepository       models.IGroupRepository
	PostRepository        models.IPostRepository
	AllowedPostRepository models.IAllowedPostRepository
}

func InitPostService(logger *log.Logger, groupRepo *models.GroupRepository, postRepo *models.PostRepository, allowedPostRepo *models.AllowedPostRepository) *PostService {
	return &PostService{
		Logger:                logger,
		GroupRepository:       groupRepo,
		PostRepository:        postRepo,
		AllowedPostRepository: allowedPostRepo,
	}
}

type feedPostJSON struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"userId"`
	UserName     string    `json:"userName"`
	Content      string    `json:"content"`
	ImagePath    string    `json:"imagePath"`
	CommentCount int       `json:"commentCount"`
	CreatedAt    time.Time `json:"createdAt"`
	GroupId      int64     `json:"groupId"`
	GroupName    string    `json:"groupName"`
}

func (s *PostService) CreatePost(post *models.Post) error {

	if len(post.Content) == 0 && len(post.ImagePath) == 0 {
		err := errors.New("content too short")
		log.Printf("CreatePost error: %s", err)
		return err
	}

	postId, err := s.PostRepository.Insert(post)

	if post.PrivacyType == enums.SubPrivate {
		for _, receiver := range post.Receivers {

			receiverId, err := strconv.Atoi(receiver)

			if err != nil {
				s.Logger.Printf("CreatePost atoi parse error: %s", err)
			}

			allowedPost := models.AllowedPost{
				UserId: receiverId,
				PostId: int(postId),
			}

			s.AllowedPostRepository.Insert(&allowedPost)
		}
	}

	if err != nil {
		log.Printf("CreatePost error: %s", err)
	}

	return err
}

func (s *PostService) CreateGroupPost(post *models.Post) error {

	if len(post.Content) == 0 {
		err := errors.New("content too short")
		log.Printf("Create Group Post error: %s", err)
		return err
	}

	postId, err := s.PostRepository.Insert(post)

	if err != nil {
		log.Printf("Create Group Post error: %s", err)
	}

	s.Logger.Printf("Group post inserted: %d", postId)

	return err
}

func (s *PostService) GetFeedPosts(userId int64, offset int64) ([]*feedPostJSON, error) {

	if offset == 0 {
		lastPostId, err := s.PostRepository.GetLastPostId()
		if err != nil {
			s.Logger.Printf("GetFeedPosts error: %s", err)
			return nil, err
		}
		offset = lastPostId + 1
	}

	posts, err := s.PostRepository.GetAllFeedPosts(userId, offset)

	if err != nil {
		s.Logger.Printf("GetFeedPosts error: %s", err)
	}

	feedPosts := []*feedPostJSON{}

	for _, p := range posts {
		// commentCount, err := s.PostRepository.GetCommentCount(p.Id)

		group := &models.Group{}

		if p.GroupId > 0 {
			group, err = s.GroupRepository.GetById(p.GroupId)
			if err != nil {
				s.Logger.Printf("GetFeedPosts error: %s", err)
			}
		}

		if p.Nickname == "" {
			p.Nickname = p.FirstName + " " + p.LastName
		}

		feedPosts = append(feedPosts, &feedPostJSON{
			Id:           p.Id,
			UserId:       p.UserId,
			UserName:     p.Nickname,
			GroupId:      p.GroupId,
			GroupName:    group.Title,
			Content:      p.Content,
			ImagePath:    p.ImagePath,
			CommentCount: p.CommentCount,
			CreatedAt:    p.CreatedAt,
		})
	}

	s.Logger.Printf("Retrived feed posts: %d", len(feedPosts))

	return feedPosts, nil
}

func (s *PostService) GetProfilePosts(userId int64, offset int64) ([]*feedPostJSON, error) {

	if offset == 0 {
		lastPostId, err := s.PostRepository.GetLastPostId()
		if err != nil {
			s.Logger.Printf("GetFeedPosts error: %s", err)
			return nil, err
		}
		offset = lastPostId + 1
	}

	posts, err := s.PostRepository.GetAllByUserId(userId, offset)

	if err != nil {
		s.Logger.Printf("GetFeedPosts error: %s", err)
	}

	feedPosts := []*feedPostJSON{}

	for _, p := range posts {

		group := &models.Group{}

		if p.GroupId > 0 {
			group, err = s.GroupRepository.GetById(p.GroupId)
			if err != nil {
				s.Logger.Printf("GetFeedPosts error: %s", err)
			}
		}

		if p.Nickname == "" {
			p.Nickname = p.FirstName + " " + p.LastName
		}

		feedPosts = append(feedPosts, &feedPostJSON{
			Id:           p.Id,
			UserId:       p.UserId,
			UserName:     p.Nickname,
			GroupId:      p.GroupId,
			GroupName:    group.Title,
			Content:      p.Content,
			ImagePath:    p.ImagePath,
			CommentCount: p.CommentCount,
			CreatedAt:    p.CreatedAt,
		})
	}

	return feedPosts, nil
}

func (s *PostService) GetGroupPosts(groupId int64, offset int64) ([]*feedPostJSON, error) {

	if offset == 0 {
		lastPostId, err := s.PostRepository.GetLastPostId()
		if err != nil {
			s.Logger.Printf("GetFeedPosts error: %s", err)
			return nil, err
		}
		offset = lastPostId + 1
	}

	posts, err := s.PostRepository.GetAllByGroupId(groupId, offset)

	if err != nil {
		s.Logger.Printf("GetFeedPosts error: %s", err)
	}

	group, err := s.GroupRepository.GetById(groupId)

	if err != nil {
		s.Logger.Printf("GetFeedPosts error: %s", err)
	}

	feedPosts := []*feedPostJSON{}

	for _, p := range posts {

		if p.Nickname == "" {
			p.Nickname = p.FirstName + " " + p.LastName
		}

		feedPosts = append(feedPosts, &feedPostJSON{
			p.Id,
			p.UserId,
			p.Nickname,
			p.Content,
			p.ImagePath,
			p.CommentCount,
			p.CreatedAt,
			groupId,
			group.Title,
		})
	}

	return feedPosts, nil
}

func (s *PostService) GetUserPosts(userId int64, offset int64, requestingUserId int64) ([]*feedPostJSON, error) {

	if offset == 0 {
		lastPostId, err := s.PostRepository.GetLastPostId()
		if err != nil {
			s.Logger.Printf("GetFeedPosts error: %s", err)
			return nil, err
		}
		offset = lastPostId + 1
	}

	posts, err := s.PostRepository.GetAllByUserAndRequestingUserIds(userId, offset, requestingUserId)

	if err != nil {
		s.Logger.Printf("GetFeedPosts error: %s", err)
	}

	feedPosts := []*feedPostJSON{}

	for _, p := range posts {

		group := &models.Group{}

		if p.GroupId > 0 {
			group, err = s.GroupRepository.GetById(p.GroupId)
			if err != nil {
				s.Logger.Printf("GetFeedPosts error: %s", err)
			}
		}

		if p.Nickname == "" {
			p.Nickname = p.FirstName + " " + p.LastName
		}

		feedPosts = append(feedPosts, &feedPostJSON{
			Id:           p.Id,
			UserId:       p.UserId,
			UserName:     p.Nickname,
			GroupId:      p.GroupId,
			GroupName:    group.Title,
			Content:      p.Content,
			ImagePath:    p.ImagePath,
			CommentCount: p.CommentCount,
			CreatedAt:    p.CreatedAt,
		})
	}

	return feedPosts, nil
}

func (s *PostService) SavePostImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	// check if file is an image
	if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image") {
		s.Logger.Println("Not an image")
		return "", errors.New("not an image")
	}

	// save image
	imagePath, err := utils.SaveImage(file, fileHeader)
	if err != nil {
		s.Logger.Printf("UpdatePostImage error: %s", err)
	}

	return imagePath, err
}
