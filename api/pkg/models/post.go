package models

import (
	"SocialNetworkRestApi/api/pkg/enums"
	"database/sql"
	"log"
	"os"
	"time"
)

type Post struct {
	Id          int64
	UserId      int64
	Content     string
	ImagePath   string
	CreatedAt   time.Time
	PrivacyType enums.PrivacyType
	Receivers   []string
	GroupId     int64
}

type FeedPost struct {
	Id           int64
	UserId       int64
	GroupId      int64
	FirstName    string
	LastName     string
	Nickname     string
	Content      string
	CommentCount int
	ImagePath    string
	CreatedAt    time.Time
	PrivacyType  enums.PrivacyType
}

type IPostRepository interface {
	GetAllByUserId(id int64, offset int64) ([]*FeedPost, error)
	GetAllByGroupId(id int64, offset int64) ([]*FeedPost, error)
	GetAllFeedPosts(currentUserId int64, offset int64) ([]*FeedPost, error)
	GetById(id int64) (*Post, error)
	Insert(post *Post) (int64, error)
	Update(post *Post) error
	GetCommentCount(postId int64) (int, error)
	GetLastPostId() (int64, error)
	GetAllByUserAndRequestingUserIds(userId int64, offset int64, requestingUserId int64) ([]*FeedPost, error)
}

type PostRepository struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepository {
	return &PostRepository{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		DB:     db,
	}
}

const FeedLimit = 10

func (repo PostRepository) Insert(post *Post) (int64, error) {
	query := `INSERT INTO posts (user_id, content, created_at, image_path, privacy_type_id, group_id)
	VALUES(?, ?, ?, ?, ?, ?)`

	args := []interface{}{
		post.UserId,
		post.Content,
		time.Now(),
		post.ImagePath,
		post.PrivacyType,
		post.GroupId,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	repo.Logger.Printf("Inserted post by user %d (last insert ID: %d)", post.UserId, lastId)

	return lastId, nil
}

func (repo PostRepository) Update(post *Post) error {
	query := `UPDATE posts SET content = ?, image_path = ?, privacy_type_id = ? WHERE id = ?`

	args := []interface{}{
		post.Content,
		post.ImagePath,
		post.PrivacyType,
		post.Id,
	}

	_, err := repo.DB.Exec(query, args...)

	return err

}

func (repo PostRepository) GetById(id int64) (*Post, error) {
	query := `SELECT id, user_id, content, created_at, image_path, privacy_type_id, group_id FROM posts WHERE id = ?`
	row := repo.DB.QueryRow(query, id)
	post := &Post{}

	err := row.Scan(&post.Id, &post.UserId, &post.Content, &post.CreatedAt, &post.ImagePath, &post.PrivacyType, &post.GroupId)

	return post, err
}

func (repo PostRepository) GetAllByUserId(id int64, offset int64) ([]*FeedPost, error) {

	stmt := `SELECT p.id, p.user_id, u.forname, u.surname, u.nickname, p.content, p.created_at, p.image_path, p.privacy_type_id, p.group_id, COUNT(DISTINCT c.id) FROM posts p
	LEFT JOIN users u on
	p.user_id = u.id
	LEFT JOIN comments c ON
	p.id = c.post_id
	WHERE p.user_id = ? AND p.id < ?
	GROUP BY p.id
    ORDER BY p.id DESC
	LIMIT ?`

	args := []interface{}{
		id,
		offset,
		FeedLimit,
	}

	rows, err := repo.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*FeedPost{}

	for rows.Next() {
		post := &FeedPost{}

		err := rows.Scan(&post.Id, &post.UserId, &post.FirstName, &post.LastName, &post.Nickname, &post.Content, &post.CreatedAt, &post.ImagePath, &post.PrivacyType, &post.GroupId, &post.CommentCount)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (repo PostRepository) GetAllByGroupId(groupId int64, offset int64) ([]*FeedPost, error) {

	stmt := `SELECT p.id, p.user_id, u.forname, u.surname, u.nickname, p.content, p.created_at, p.image_path, p.privacy_type_id, COUNT(DISTINCT c.id) FROM posts p
	LEFT JOIN users u on
	p.user_id = u.id
	LEFT JOIN comments c ON
	p.id = c.post_id
	WHERE p.group_id = ? AND p.id < ?
	GROUP BY p.id
	ORDER BY p.id DESC
	LIMIT ?`

	args := []interface{}{
		groupId,
		offset,
		FeedLimit,
	}

	rows, err := repo.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*FeedPost{}

	for rows.Next() {
		post := &FeedPost{}

		err := rows.Scan(&post.Id, &post.UserId, &post.FirstName, &post.LastName, &post.Nickname, &post.Content, &post.CreatedAt, &post.ImagePath, &post.PrivacyType, &post.CommentCount)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// Return all posts to the current user by offset
func (m PostRepository) GetAllFeedPosts(currentUserId int64, offset int64) ([]*FeedPost, error) {

	//Change value if needed for testing purposes
	// currentUserId = 11
	// because seeded posts have similar created_at, using p.id as temporary order by

	stmt := `SELECT p.id, p.user_id, u.forname, u.surname, u.nickname, p.content, p.created_at, p.image_path, privacy_type_id, p.group_id, COUNT(DISTINCT c.id) FROM posts p 
	LEFT JOIN users u on
	p.user_id = u.id
	LEFT JOIN  followers f ON  
	p.user_id = f.following_id
	LEFT JOIN allowed_private_posts app ON
	p.id = app.post_id
	LEFT JOIN comments c ON
	p.id = c.post_id
	WHERE (privacy_type_id = 1
	OR p.user_id = ?
	OR (privacy_type_id = 2 AND f.id IS NOT NULL AND f.follower_id = ? AND f.accepted = 1)
	OR (privacy_type_id = 3 AND f.id IS NOT NULL AND f.follower_id = ? AND f.accepted = 1 AND app.id IS NOT NULL AND app.user_id = ?)
	OR p.group_id IN (SELECT group_id FROM user_groups WHERE user_id = ?))
	AND p.id < ?
	GROUP BY p.id
	ORDER BY p.id DESC
	LIMIT ?`

	args := []interface{}{
		currentUserId,
		currentUserId,
		currentUserId,
		currentUserId,
		currentUserId,
		offset,
		FeedLimit,
	}

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*FeedPost{}

	for rows.Next() {
		post := &FeedPost{}

		err := rows.Scan(&post.Id, &post.UserId, &post.FirstName, &post.LastName, &post.Nickname, &post.Content, &post.CreatedAt, &post.ImagePath, &post.PrivacyType, &post.GroupId, &post.CommentCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m PostRepository) GetCommentCount(postId int64) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = ?`
	row := m.DB.QueryRow(query, postId)
	var commentCount int

	err := row.Scan(&commentCount)

	if err != nil {
		return -1, err
	}

	return commentCount, nil
}

func (repo PostRepository) InsertSeedPost(post *Post) (int64, error) {
	query := `INSERT INTO posts (user_id, content, created_at, image_path, privacy_type_id, group_id)
	VALUES(?, ?, ?, ?, ?, ?)`

	args := []interface{}{
		post.UserId,
		post.Content,
		post.CreatedAt,
		post.ImagePath,
		post.PrivacyType,
		post.GroupId,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	repo.Logger.Printf("Inserted post by user %d (last insert ID: %d)", post.UserId, lastId)

	return lastId, nil
}

func (repo PostRepository) GetLastPostId() (int64, error) {
	query := `SELECT id FROM posts ORDER BY id DESC LIMIT 1`
	row := repo.DB.QueryRow(query)
	var id int64

	err := row.Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (repo PostRepository) GetAllByUserAndRequestingUserIds(userId int64, offset int64, requestingUserId int64) ([]*FeedPost, error) {

	stmt := `SELECT p.id, p.user_id, u.forname, u.surname, u.nickname, p.content, p.created_at, p.image_path, p.privacy_type_id, p.group_id, COUNT(DISTINCT c.id) FROM posts p
	LEFT JOIN users u on p.user_id = u.id
	LEFT JOIN followers f ON p.user_id = f.following_id
	LEFT JOIN allowed_private_posts app ON p.id = app.post_id 
	LEFT JOIN comments c ON p.id = c.post_id
	WHERE (p.privacy_type_id = 1 AND p.user_id = ?
	OR (p.privacy_type_id = 2 AND p.user_id = ? AND f.id IS NOT NULL AND f.follower_id = ? AND f.accepted = 1)
	OR (p.privacy_type_id = 3 AND p.user_id = ? AND f.id IS NOT NULL AND f.follower_id = ? AND f.accepted = 1 AND app.id IS NOT NULL AND app.user_id = ?))
	AND p.id < ?
	GROUP BY p.id
	ORDER BY p.id DESC
	LIMIT ?`

	args := []interface{}{
		userId,
		userId,
		requestingUserId,
		userId,
		requestingUserId,
		requestingUserId,
		offset,
		FeedLimit,
	}

	rows, err := repo.DB.Query(stmt, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*FeedPost{}

	for rows.Next() {
		post := &FeedPost{}

		err := rows.Scan(&post.Id, &post.UserId, &post.FirstName, &post.LastName, &post.Nickname, &post.Content, &post.CreatedAt, &post.ImagePath, &post.PrivacyType, &post.GroupId, &post.CommentCount)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil

}
