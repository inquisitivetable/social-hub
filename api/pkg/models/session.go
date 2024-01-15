package models

import (
	"database/sql"
	"log"
	"os"
	"time"
)

type Session struct {
	Id        int64
	UserId    int64
	Token     string
	CreatedAt time.Time
}

type ISessionRepository interface {
	DeleteByToken(token string) error
	GetByToken(token string) (*Session, error)
	GetUserSessions(id int) ([]*Session, error)
	Insert(session *Session) (int64, error)
}

type SessionRepository struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		DB:     db,
	}
}

// Inserts a new user session to database
func (repo SessionRepository) Insert(session *Session) (int64, error) {
	query := `INSERT INTO user_sessions (user_id, token, created_at)
	VALUES(?, ?, ?)`

	args := []interface{}{
		session.UserId,
		session.Token,
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

	repo.Logger.Printf("Inserted session for user %d (last insert ID: %d)", session.UserId, lastId)

	return lastId, nil
}

// Returns a session by token
func (repo SessionRepository) GetByToken(token string) (*Session, error) {

	query := `SELECT id, user_id, token, created_at FROM user_sessions WHERE token = ?`
	row := repo.DB.QueryRow(query, token)
	session := &Session{}

	err := row.Scan(&session.Id, &session.UserId, &session.Token, &session.CreatedAt)

	return session, err
}

// Returns all user sessions
func (repo SessionRepository) GetUserSessions(id int) ([]*Session, error) {

	//TODO
	//Not sure if needed
	return nil, nil
}

// remove a session by token
func (repo SessionRepository) DeleteByToken(token string) error {
	query := `DELETE FROM user_sessions WHERE token = ?`

	result, err := repo.DB.Exec(query, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	repo.Logger.Printf("Deleted %d session(s) with token %s", rowsAffected, token)

	return err
}
