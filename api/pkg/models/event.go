package models

import (
	"database/sql"
	"log"
	"os"
	"time"
)

type Event struct {
	Id           int64
	GroupId      int64
	UserId       int64
	CreatedAt    time.Time
	EventTime    time.Time
	EventEndTime time.Time
	Title        string
	Description  string
}

type CreateGroupEventFormData struct {
	GroupId int `json:"group_id"`
	//UserId       int64  `json:"userId"`
	EventTime    string `json:"startTime"`
	EventEndTime string `json:"endTime"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

type IEventRepository interface {
	GetAllByGroupId(groupId int64) ([]*Event, error)
	GetAllByUserId(userId int64) ([]*Event, error)
	Insert(event *Event) (int64, error)
	InsertSeedEvent(event *Event) (int64, error)
	GetById(id int64) (*Event, error)
}

type EventRepository struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewEventRepo(db *sql.DB) *EventRepository {
	return &EventRepository{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		DB:     db,
	}
}

func (repo EventRepository) Insert(event *Event) (int64, error) {
	query := `INSERT INTO group_events (group_id, user_id, created_at, event_time, event_end_time, title, description)
	VALUES(?, ?, ?, ?, ?, ?, ?)`

	args := []interface{}{
		event.GroupId,
		event.UserId,
		time.Now(),
		event.EventTime,
		event.EventEndTime,
		event.Title,
		event.Description,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	repo.Logger.Printf("Last inserted event '%s' by user %d (last insert ID: %d)", event.Title, event.UserId, lastId)

	return lastId, nil
}

func (repo EventRepository) InsertSeedEvent(event *Event) (int64, error) {
	query := `INSERT INTO group_events (group_id, user_id, created_at, event_time, event_end_time, title, description)
	VALUES(?, ?, ?, ?, ?, ?, ?)`

	args := []interface{}{
		event.GroupId,
		event.UserId,
		event.CreatedAt,
		event.EventTime,
		event.EventEndTime,
		event.Title,
		event.Description,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	repo.Logger.Printf("Last inserted event '%s' by user %d (last insert ID: %d)", event.Title, event.UserId, lastId)

	return lastId, nil
}

func (repo EventRepository) GetAllByGroupId(id int64) ([]*Event, error) {

	query := `SELECT id, group_id, user_id, created_at, event_time, event_end_time, title, description FROM group_events WHERE group_id = ?`

	rows, err := repo.DB.Query(query, id)

	if err != nil {
		return nil, err
	}

	events := []*Event{}

	defer rows.Close()
	for rows.Next() {
		event := &Event{}

		err := rows.Scan(&event.Id, &event.GroupId, &event.UserId, &event.CreatedAt, &event.EventTime, &event.EventEndTime, &event.Title, &event.Description)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	repo.Logger.Printf("Found %d events for group %d", len(events), id)

	return events, err
}

func (repo EventRepository) GetAllByUserId(id int64) ([]*Event, error) {

	query := `SELECT DISTINCT ge.id, group_id, ge.user_id, ge.created_at, ge.event_time, ge.event_end_time, ge.title, ge.description FROM group_events ge
	INNER JOIN group_event_attendance gea
	ON gea.event_id = ge.id
	WHERE gea.user_id = ? AND (gea.is_attending = true OR gea.is_attending IS NULL)`

	rows, err := repo.DB.Query(query, id)

	if err != nil {
		return nil, err
	}

	events := []*Event{}

	defer rows.Close()
	for rows.Next() {
		event := &Event{}

		err := rows.Scan(&event.Id, &event.GroupId, &event.UserId, &event.CreatedAt, &event.EventTime, &event.EventEndTime, &event.Title, &event.Description)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	repo.Logger.Printf("Found %d events for user %d", len(events), id)

	return events, err
}

func (repo EventRepository) GetById(id int64) (*Event, error) {
	query := `SELECT id, group_id, user_id, created_at, event_time, event_end_time, title, description FROM group_events WHERE id = ?`

	row := repo.DB.QueryRow(query, id)

	event := &Event{}

	err := row.Scan(&event.Id, &event.GroupId, &event.UserId, &event.CreatedAt, &event.EventTime, &event.EventEndTime, &event.Title, &event.Description)

	if err != nil {
		return nil, err
	}

	repo.Logger.Printf("Found event %d", event.Id)

	return event, nil
}
