package models

import (
	"database/sql"
	"log"
	"os"
	"time"
)

type Message struct {
	Id          int64
	SenderId    int64
	RecipientId int64
	GroupId     int64
	Content     string
	//ImagePath   string
	SentAt time.Time
	ReadAt time.Time
}

type IMessageRepository interface {
	Insert(event *Message) (int64, error)
	GetMessagesByGroupId(groupId int64, lastMessage int64) ([]*Message, error)
	GetMessagesByUserIds(userId int64, secondUserId int64, lastMessage int64) ([]*Message, error)
	GetChatUsers(id int64) ([]*User, error)
	GetChatGroups(id int64) ([]*Group, error)
	GetLastMessage(userId int64, otherId int64, isGroup bool) (*Message, error)
	GetUnreadCount(userId int64, otherId int64) (int64, error)
	MarkMessagesRead(senderId int64, recipientId int64, messageId int64) error
	GetById(id int64) (*Message, error)
}

type MessageRepository struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewMessageRepo(db *sql.DB) *MessageRepository {
	return &MessageRepository{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		DB:     db,
	}
}

func (repo MessageRepository) Insert(message *Message) (int64, error) {
	query := `INSERT INTO messages (sender_id, recipient_id, group_id, content, sent_at)
	VALUES(?, ?, ?, ?, ?)`

	args := []interface{}{
		message.SenderId,
		message.RecipientId,
		message.GroupId,
		message.Content,
		message.SentAt,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	repo.Logger.Printf("Inserted message by user %d, to group/user %d/%d (last insert ID: %d)", message.SenderId, message.GroupId, message.RecipientId, lastId)

	return lastId, nil
}
func (repo MessageRepository) Update(message *Message) error {
	//TODO
	//Update method needed when readAt is being used
	return nil
}

func (repo MessageRepository) GetMessagesByGroupId(groupId int64, lastMessage int64) ([]*Message, error) {
	query := `SELECT id, sender_id, group_id, content, sent_at FROM messages m
	WHERE group_id = ? AND id < ?
	ORDER BY sent_at DESC LIMIT 10`

	args := []interface{}{
		groupId,
		lastMessage,
	}

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		message := &Message{}

		err := rows.Scan(&message.Id, &message.SenderId, &message.GroupId, &message.Content, &message.SentAt) //, &message.ReadAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (repo MessageRepository) GetMessagesByUserIds(userId int64, secondUserId int64, lastMessage int64) ([]*Message, error) {
	query := `SELECT id, sender_id, recipient_id, group_id, content, sent_at FROM messages m
	WHERE (sender_id = ? AND recipient_id = ? AND id < ?) OR (sender_id = ? AND recipient_id = ? AND id < ?) 
    ORDER BY sent_at DESC LIMIT 10`

	args := []interface{}{
		userId,
		secondUserId,
		lastMessage,
		secondUserId,
		userId,
		lastMessage,
	}

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		message := &Message{}

		err := rows.Scan(&message.Id, &message.SenderId, &message.RecipientId, &message.GroupId, &message.Content, &message.SentAt) //, &message.ReadAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (repo MessageRepository) GetChatUsers(id int64) ([]*User, error) {
	// joint query of all users the user is following + all users the user has sent or received messages from

	query := `
		SELECT u.id, u.forname, u.surname, u.nickname, u.image_path, u.created_at FROM users u 
		JOIN followers f ON u.id = f.following_id WHERE f.follower_id = ? AND f.accepted = 1 GROUP BY u.id
		UNION
		SELECT u.id, u.forname, u.surname, u.nickname, u.image_path, u.created_at FROM users u
		JOIN messages m ON u.id = m.sender_id WHERE m.recipient_id = ? GROUP BY u.id
		UNION
		SELECT u.id, u.forname, u.surname, u.nickname, u.image_path, u.created_at FROM users u
		JOIN messages m ON u.id = m.recipient_id WHERE m.sender_id = ? GROUP BY u.id
		`

	rows, err := repo.DB.Query(query, id, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.ImagePath,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo MessageRepository) GetChatGroups(id int64) ([]*Group, error) {
	//get all groups from user_groups where user is member and title and image_path from groups

	query := `SELECT g.id, g.title, g.created_at, g.image_path FROM groups g 
	JOIN user_groups ug ON g.id = ug.group_id 
	WHERE ug.user_id = ? AND ug.accepted = 1`

	rows, err := repo.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []*Group{}

	for rows.Next() {
		group := &Group{}
		err := rows.Scan(
			&group.Id,
			&group.Title,
			&group.CreatedAt,
			&group.ImagePath,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (repo MessageRepository) GetLastMessage(userId int64, otherId int64, isGroup bool) (*Message, error) {
	//get last message from database
	var query string
	var args []interface{}

	if isGroup {
		query = `SELECT id, sender_id, recipient_id, group_id, content, sent_at FROM messages WHERE group_id = ? ORDER BY sent_at DESC LIMIT 1`
		args = []interface{}{
			otherId,
		}
	} else {
		query = `SELECT id, sender_id, recipient_id, group_id, content, sent_at FROM messages WHERE (sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?) ORDER BY sent_at DESC LIMIT 1`
		args = []interface{}{
			userId,
			otherId,
			otherId,
			userId,
		}
	}

	row := repo.DB.QueryRow(query, args...)

	message := &Message{}

	err := row.Scan(&message.Id, &message.SenderId, &message.RecipientId, &message.GroupId, &message.Content, &message.SentAt) // , &message.ReadAt)

	if err == sql.ErrNoRows {
		return &Message{}, nil
	}

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (repo MessageRepository) GetUnreadCount(userId int64, otherId int64) (int64, error) {
	var query string
	var args []interface{}

	query = `SELECT COUNT(*) FROM messages WHERE (sender_id = ? AND recipient_id = ?) AND read_at IS NULL`
	args = []interface{}{
		otherId,
		userId,
	}

	row := repo.DB.QueryRow(query, args...)

	var count int64

	err := row.Scan(&count)

	if err != nil {
		return -1, err
	}

	return count, nil
}

func (repo MessageRepository) MarkMessagesRead(senderId int64, recipientId int64, messageId int64) error {
	var query string
	var args []interface{}

	query = `UPDATE messages SET read_at = ? WHERE sender_id = ? AND recipient_id = ? AND id <= ? AND read_at IS NULL`
	args = []interface{}{
		time.Now(),
		senderId,
		recipientId,
		messageId,
	}

	_, err := repo.DB.Exec(query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (repo MessageRepository) GetById(id int64) (*Message, error) {
	row := repo.DB.QueryRow("SELECT id, sender_id, recipient_id, group_id, content, sent_at FROM messages WHERE id = ?", id)

	message := &Message{}

	err := row.Scan(&message.Id, &message.SenderId, &message.RecipientId, &message.GroupId, &message.Content, &message.SentAt) // , &message.ReadAt)

	if err == sql.ErrNoRows {
		return &Message{}, nil
	}

	if err != nil {
		return nil, err
	}

	return message, nil
}
