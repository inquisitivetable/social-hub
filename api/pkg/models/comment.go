package models

import (
	"database/sql"
	"log"
	"os"
	"time"
)

const CommentLimit = 5

type Comment struct {
	Id        int64
	PostId    int64
	UserId    int64
	Content   string
	ImagePath string
	CreatedAt time.Time
}

type PostComment struct {
	Id           int64
	UserId       int64
	UserName     string
	Content      string
	ImagePath    string
	CreatedAt    time.Time
	CommentCount int
}

type ICommentRepository interface {
	GetAllByPostId(postId int64, offset int64) ([]*PostComment, error)
	GetAllByUserId(userId int64) ([]*Comment, error)
	GetById(id int64) (*Comment, error)
	Insert(comment *Comment) (int64, error)
	Update(comment *Comment) error
	InsertSeedComment(comment *Comment) (int64, error)
}

type CommentRepository struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		DB:     db,
	}
}

func (repo CommentRepository) Insert(comment *Comment) (int64, error) {
	query := `INSERT INTO comments (post_id, user_id, content, image_path, created_at)
	VALUES(?, ?, ?, ?, ?)`

	args := []interface{}{
		comment.PostId,
		comment.UserId,
		comment.Content,
		comment.ImagePath,
		time.Now(),
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	repo.Logger.Printf("Inserted comment by user %d for post %d (last insert ID: %d)", comment.UserId, comment.PostId, lastId)

	return lastId, nil
}

func (repo CommentRepository) Update(comment *Comment) error {

	//TODO
	//update comment in database
	return nil
}

func (repo CommentRepository) GetById(id int64) (*Comment, error) {
	query := `SELECT id, post_id, user_id, content,  image_path, created_at FROM comments WHERE id = ?`
	row := repo.DB.QueryRow(query, id)
	comment := &Comment{}

	err := row.Scan(&comment.Id, &comment.PostId, &comment.UserId, &comment.Content, &comment.ImagePath, &comment.CreatedAt)

	return comment, err
}

func (repo CommentRepository) GetAllByPostId(postId int64, offset int64) ([]*PostComment, error) {
	query := `SELECT c.id, c.user_id, u.nickname, c.content, c.image_path, c.created_at, cc.comment_count FROM comments c
	LEFT JOIN users u ON c.user_id = u.id
	LEFT JOIN (SELECT COUNT(*) AS comment_count, post_id FROM comments GROUP BY post_id) AS cc ON c.post_id = cc.post_id
	WHERE c.post_id = ? 
	GROUP BY c.id
	ORDER BY c.created_at DESC
	LIMIT ? OFFSET ?`

	args := []interface{}{
		postId,
		CommentLimit,
		(offset * CommentLimit),
	}

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	comments := []*PostComment{}

	for rows.Next() {
		comment := &PostComment{}

		err := rows.Scan(&comment.Id, &comment.UserId, &comment.UserName, &comment.Content, &comment.ImagePath, &comment.CreatedAt, &comment.CommentCount)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo CommentRepository) GetAllByUserId(userId int64) ([]*Comment, error) {
	query := `SELECT id, post_id, user_id, content,  image_path, created_at FROM comments WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := repo.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	comments := []*Comment{}

	for rows.Next() {
		comment := &Comment{}

		err := rows.Scan(&comment.Id, &comment.PostId, &comment.UserId, &comment.Content, &comment.ImagePath, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo CommentRepository) InsertSeedComment(comment *Comment) (int64, error) {
	query := `INSERT INTO comments (post_id, user_id, content, image_path, created_at)
	VALUES(?, ?, ?, ?, ?)`

	args := []interface{}{
		comment.PostId,
		comment.UserId,
		comment.Content,
		comment.ImagePath,
		comment.CreatedAt,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	repo.Logger.Printf("Inserted comment by user %d for post %d (last insert ID: %d)", comment.UserId, comment.PostId, lastId)

	return lastId, nil
}
