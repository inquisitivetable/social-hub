package models

import (
	"database/sql"
	"log"
	"os"
)

type AllowedPost struct {
	UserId int
	PostId int
}

type IAllowedPostRepository interface {
	Insert(allowedPost *AllowedPost) (int64, error)
}

type AllowedPostRepository struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewAllowedPostRepo(db *sql.DB) *AllowedPostRepository {
	return &AllowedPostRepository{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		DB:     db,
	}
}

func (repo AllowedPostRepository) Insert(allowedPost *AllowedPost) (int64, error) {
	query := `INSERT INTO allowed_private_posts (post_id, user_id)
	VALUES(?, ?)`

	args := []interface{}{
		allowedPost.PostId,
		allowedPost.UserId,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, nil
}
