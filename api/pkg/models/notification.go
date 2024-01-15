package models

import (
	"database/sql"
	"log"
	"os"
	"time"
)

type Notification struct {
	Id                    int64
	ReceiverId            int64
	NotificationDetailsId int64
	SeenAt                time.Time
	Reaction              sql.NullBool
}

type NotificationDetails struct {
	Id               int64
	SenderId         int64
	NotificationType string
	EntityId         int64
	CreatedAt        time.Time
}

type NotificationJSON struct {
	ReceiverId       int64     `json:"receiver_id"`
	NotificationType string    `json:"notification_type"`
	NotificationId   int64     `json:"notification_id"`
	SenderId         int64     `json:"sender_id"`
	SenderName       string    `json:"sender_name"`
	GroupId          int64     `json:"group_id"`
	GroupName        string    `json:"group_name"`
	EventId          int64     `json:"event_id"`
	EventName        string    `json:"event_name"`
	EventDate        time.Time `json:"event_datetime"`
}

type INotificationRepository interface {
	InsertDetails(notificationDetails *NotificationDetails) (int64, error)
	InsertNotification(notification *Notification) (int64, error)
	Update(notification *Notification) error
	GetById(id int64) (*Notification, error)
	GetDetailsById(id int64) (*NotificationDetails, error)
	GetByReceiverId(receiverId int64) ([]*Notification, error)
	GetByEventAndUserId(eventId int64, userId int64) (*Notification, error)
	GetNotificationType(notificationType string) (int64, error)
}

type NotificationRepository struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		DB:     db,
	}
}

func (repo NotificationRepository) InsertDetails(notificationDetails *NotificationDetails) (int64, error) {

	NotificationTypeID, err := repo.GetNotificationType(notificationDetails.NotificationType)

	if err != nil {
		repo.Logger.Printf("Error getting notification type: %s", err.Error())
		return -1, err
	}

	query := `INSERT INTO notification_details (sender_id, notification_type_id, entity_id, created_at)
	VALUES(?, ?, ?, ?)`

	args := []interface{}{
		notificationDetails.SenderId,
		NotificationTypeID,
		notificationDetails.EntityId,
		notificationDetails.CreatedAt,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		repo.Logger.Printf("Error inserting notification details: %s", err.Error())
		return -1, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		repo.Logger.Printf("Error getting last insert ID: %s", err.Error())
		return -1, err
	}

	repo.Logger.Printf("Inserted notification details (last insert ID: %d)", lastId)

	return lastId, nil
}

func (repo NotificationRepository) InsertNotification(notification *Notification) (int64, error) {

	query := `INSERT INTO notifications (receiver_id, notification_details_id, seen_at, reaction)
	VALUES(?, ?, ?, ?)`

	args := []interface{}{
		notification.ReceiverId,
		notification.NotificationDetailsId,
		notification.SeenAt,
		notification.Reaction,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		repo.Logger.Printf("Error inserting notification: %s", err.Error())
		return -1, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	repo.Logger.Printf("Inserted notification for user %d (last insert ID: %d)", notification.ReceiverId, lastId)

	return lastId, nil

}

func (repo NotificationRepository) Update(notification *Notification) error {
	query := `UPDATE notifications SET seen_at = ?, reaction = ? WHERE id = ?`

	args := []interface{}{
		notification.SeenAt,
		notification.Reaction,
		notification.Id,
	}

	_, err := repo.DB.Exec(query, args...)

	if err != nil {
		repo.Logger.Printf("Error updating notification: %s", err.Error())
		return err
	}

	return nil
}

func (repo NotificationRepository) GetById(id int64) (*Notification, error) {
	query := `SELECT id, receiver_id, notification_details_id, seen_at, reaction FROM notifications
	WHERE id = ?`

	args := []interface{}{
		id,
	}

	notification := &Notification{}

	err := repo.DB.QueryRow(query, args...).Scan(&notification.Id, &notification.ReceiverId, &notification.NotificationDetailsId, &notification.SeenAt, &notification.Reaction)

	if err != nil {
		repo.Logger.Printf("Error getting notification: %s", err.Error())
		return nil, err
	}

	return notification, nil
}

func (repo NotificationRepository) GetDetailsById(id int64) (*NotificationDetails, error) {
	query := `SELECT nd.id, nd.sender_id, nt.name, nd.entity_id, nd.created_at FROM notification_details nd
	JOIN notification_types nt ON nd.notification_type_id = nt.id
	WHERE nd.id = ?`

	args := []interface{}{
		id,
	}

	notificationDetails := &NotificationDetails{}

	err := repo.DB.QueryRow(query, args...).Scan(&notificationDetails.Id, &notificationDetails.SenderId, &notificationDetails.NotificationType, &notificationDetails.EntityId, &notificationDetails.CreatedAt)

	if err != nil {
		repo.Logger.Printf("Error getting notification: %s", err.Error())
		return nil, err
	}

	return notificationDetails, nil
}

func (repo NotificationRepository) GetByReceiverId(userId int64) ([]*Notification, error) {

	query := `SELECT id, seen_at, notification_details_id, reaction FROM notifications
	WHERE receiver_id = ? AND reaction IS NULL`

	args := []interface{}{
		userId,
	}

	rows, err := repo.DB.Query(query, args...)

	if err != nil {
		repo.Logger.Printf("Error getting notifications: %s", err.Error())
		return nil, err
	}

	defer rows.Close()

	notifications := []*Notification{}

	for rows.Next() {
		var notification Notification

		err := rows.Scan(&notification.Id, &notification.SeenAt, &notification.NotificationDetailsId, &notification.Reaction)

		if err != nil {
			repo.Logger.Printf("Error scanning notification: %s", err.Error())
			return nil, err
		}

		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func (repo NotificationRepository) GetByEventAndUserId(eventId int64, userId int64) (*Notification, error) {
	query := `SELECT n.id, n.receiver_id, n.notification_details_id, n.seen_at, n.reaction FROM notifications n
	JOIN notification_details nd ON n.notification_details_id = nd.id
	WHERE nd.entity_id = ? AND n.receiver_id = ?`

	args := []interface{}{
		eventId,
		userId,
	}

	notification := &Notification{}

	err := repo.DB.QueryRow(query, args...).Scan(&notification.Id, &notification.ReceiverId, &notification.NotificationDetailsId, &notification.SeenAt, &notification.Reaction)

	if err != nil {
		repo.Logger.Printf("Error getting notification: %s", err.Error())
		return nil, err
	}

	return notification, nil
}

func (repo NotificationRepository) GetNotificationType(notificationType string) (int64, error) {
	query := `SELECT id FROM notification_types WHERE name = ?`

	args := []interface{}{
		notificationType,
	}

	var id int64

	err := repo.DB.QueryRow(query, args...).Scan(&id)

	if err != nil {
		repo.Logger.Printf("Error getting notification type: %s", err.Error())
		return -1, err
	}

	return id, nil
}
