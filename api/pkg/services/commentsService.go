package services

import (
	"SocialNetworkRestApi/api/pkg/models"
	"errors"
	"log"
	"time"
)

type ICommentService interface {
	GetPostComments(userId int64, offset int64) ([]*CommentJSON, error)
	CreateComment(comment *models.Comment) error
}

// Controller contains the service, which contains database-related logic, as an injectable dependency, allowing us to decouple business logic from db logic.
type CommentService struct {
	Logger            *log.Logger
	CommentRepository models.ICommentRepository
	UserRepository    models.IUserRepository
}

func InitCommentService(logger *log.Logger, commentRepo *models.CommentRepository, userRepo *models.UserRepository) *CommentService {
	return &CommentService{
		Logger:            logger,
		CommentRepository: commentRepo,
		UserRepository:    userRepo,
	}
}

type CommentJSON struct {
	Id           int       `json:"id"`
	UserId       int       `json:"userId"`
	UserName     string    `json:"userName"`
	Content      string    `json:"content"`
	ImagePath    string    `json:"imagePath"`
	CreatedAt    time.Time `json:"createdAt"`
	CommentCount int       `json:"commentCount"`
}

func (s *CommentService) GetPostComments(postId int64, offset int64) ([]*CommentJSON, error) {

	result, err := s.CommentRepository.GetAllByPostId(postId, offset)

	if err != nil {
		s.Logger.Printf("Failed fetching comments: %s", err)
	}

	comments := []*CommentJSON{}

	for _, p := range result {

		user, err := s.UserRepository.GetById(p.UserId)
		if err != nil {
			s.Logger.Printf("Failed fetching user: %s", err)
		}

		if user.Nickname == "" {
			user.Nickname = user.FirstName + " " + user.LastName
		}

		comments = append(comments, &CommentJSON{
			int(p.Id),
			int(p.UserId),
			user.Nickname,
			p.Content,
			p.ImagePath,
			p.CreatedAt,
			p.CommentCount,
		})
	}

	return comments, nil
}

func (s *CommentService) CreateComment(comment *models.Comment) error {

	if len(comment.Content) == 0 && len(comment.ImagePath) == 0 {
		err := errors.New("comment content too short")
		s.Logger.Printf("CreateComment error: %s", err)
		return err
	}

	_, err := s.CommentRepository.Insert(comment)

	if err != nil {
		log.Printf("CreateComment error: %s", err)
	}

	return err
}
