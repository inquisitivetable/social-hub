package models

import (
	"database/sql"
	"log"
	"os"
)

type AttendeeJSON struct {
	Id          int    `json:"id"`
	Nickname    string `json:"nickname"`
	ImagePath   string `json:"imagePath"`
	IsAttending bool   `json:"isAttending"`
}

type EventAttendance struct {
	UserId      int64 `json:"userId"`
	EventId     int64 `json:"eventId"`
	IsAttending bool  `json:"isAttending"`
}

type IEventAttendanceRepository interface {
	Insert(attendance *EventAttendance) (int64, error)
	Update(attendance *EventAttendance) (int64, error)
	GetAttendeesByEventId(eventId int64) ([]*EventAttendance, error)
	GetAttendee(eventId int64, userId int64) (*EventAttendance, error)
}

type EventAttendanceRepository struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewEventAttendanceRepo(db *sql.DB) *EventAttendanceRepository {
	return &EventAttendanceRepository{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		DB:     db,
	}
}

func (repo EventAttendanceRepository) Insert(attendance *EventAttendance) (int64, error) {
	query := `INSERT INTO group_event_attendance (user_id, event_id, is_attending)
	VALUES(?, ?, ?)`

	args := []interface{}{
		attendance.UserId,
		attendance.EventId,
		attendance.IsAttending,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return -1, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	repo.Logger.Printf("User %d added to attend event %d", attendance.UserId, attendance.EventId)

	return lastId, nil
}

func (repo EventAttendanceRepository) Update(attendance *EventAttendance) (int64, error) {
	query := `UPDATE group_event_attendance SET is_attending = ? WHERE user_id = ? AND event_id = ?`

	args := []interface{}{
		attendance.IsAttending,
		attendance.UserId,
		attendance.EventId,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return -1, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return -1, err
	}

	repo.Logger.Printf("User %d updated to attend event %d", attendance.UserId, attendance.EventId)

	return rowsAffected, nil
}

func (repo EventAttendanceRepository) GetAttendeesByEventId(eventId int64) ([]*EventAttendance, error) {
	query := `SELECT user_id, event_id, is_attending FROM group_event_attendance WHERE event_id = ?`

	rows, err := repo.DB.Query(query, eventId)

	if err != nil {
		return nil, err
	}

	attendees := []*EventAttendance{}

	for rows.Next() {
		attendee := &EventAttendance{}
		err := rows.Scan(&attendee.UserId, &attendee.EventId, &attendee.IsAttending)

		if err != nil {
			return nil, err
		}

		attendees = append(attendees, attendee)
	}

	repo.Logger.Printf("Fetched %d attendees for event %d", len(attendees), eventId)

	return attendees, nil
}

func (repo EventAttendanceRepository) GetAttendee(eventId int64, userId int64) (*EventAttendance, error) {
	query := `SELECT user_id, event_id, is_attending FROM group_event_attendance 
	WHERE event_id = ? AND user_id = ?`

	row := repo.DB.QueryRow(query, eventId, userId)

	attendee := &EventAttendance{}

	err := row.Scan(&attendee.UserId, &attendee.EventId, &attendee.IsAttending)

	return attendee, err

}
