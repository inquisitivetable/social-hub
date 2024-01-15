package models

import (
	"database/sql"
	"log"
	"os"
	"time"
)

type GroupMember struct {
	UserId   int64
	GroupId  int64
	JoinedAt time.Time
	Accepted bool
}

type GroupMemberJSON struct {
	GroupId int   `json:"groupId"`
	UserIds []int `json:"userIds"`
}

type IGroupMemberRepository interface {
	Insert(groupMember *GroupMember) (int64, error)
	Update(groupMember *GroupMember) error
	Delete(groupMember *GroupMember) error
	GetGroupMembersByGroupId(groupId int64) ([]*GroupMember, error)
	GetMemberByGroupId(groupId int64, userId int64) (*GroupMember, error)
	GetById(id int64) (*GroupMember, error)
}

type GroupMemberRepository struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewGroupMemberRepo(db *sql.DB) *GroupMemberRepository {
	return &GroupMemberRepository{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		DB:     db,
	}
}

func (repo GroupMemberRepository) Insert(groupMember *GroupMember) (int64, error) {
	query := `INSERT INTO user_groups (user_id, group_id, joined_at, accepted)
	VALUES(?, ?, ?, ?)`

	args := []interface{}{
		groupMember.UserId,
		groupMember.GroupId,
		groupMember.JoinedAt,
		groupMember.Accepted,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	repo.Logger.Printf("Last inserted groupuser '%d' for user %d in group %d", lastId, groupMember.UserId, groupMember.GroupId)

	return lastId, nil
}

func (repo GroupMemberRepository) Update(groupMember *GroupMember) error {
	query := `UPDATE user_groups SET joined_at = ?, accepted = ?
	WHERE user_id = ? AND group_id = ?`

	args := []interface{}{
		groupMember.JoinedAt,
		groupMember.Accepted,
		groupMember.UserId,
		groupMember.GroupId,
	}

	_, err := repo.DB.Exec(query, args...)

	return err
}

func (repo GroupMemberRepository) Delete(groupMember *GroupMember) error {
	query := `DELETE FROM user_groups WHERE user_id = ? AND group_id = ?`

	args := []interface{}{
		groupMember.UserId,
		groupMember.GroupId,
	}

	_, err := repo.DB.Exec(query, args...)

	if err != nil {
		return err
	}

	repo.Logger.Printf("Deleted groupuser for user %d in group %d", groupMember.UserId, groupMember.GroupId)

	return nil
}

func (repo GroupMemberRepository) GetGroupMembersByGroupId(groupId int64) ([]*GroupMember, error) {
	query := `SELECT user_id, joined_at, accepted FROM user_groups
	WHERE group_id = ?`

	rows, err := repo.DB.Query(query, groupId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groupMembers := []*GroupMember{}

	for rows.Next() {
		groupMember := &GroupMember{}

		err := rows.Scan(&groupMember.UserId, &groupMember.JoinedAt, &groupMember.Accepted)
		if err != nil {
			return nil, err
		}
		groupMembers = append(groupMembers, groupMember)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groupMembers, nil
}

func (repo GroupMemberRepository) GetMemberByGroupId(groupId int64, userId int64) (*GroupMember, error) {
	query := `SELECT user_id, group_id, joined_at, accepted FROM user_groups
	WHERE user_id = ? AND group_id = ?`

	args := []interface{}{
		userId,
		groupId,
	}

	row := repo.DB.QueryRow(query, args...)

	groupMember := &GroupMember{}

	err := row.Scan(&groupMember.UserId, &groupMember.GroupId, &groupMember.JoinedAt, &groupMember.Accepted)

	if err != nil {
		return nil, err
	}

	return groupMember, nil
}

func (repo GroupMemberRepository) GetById(id int64) (*GroupMember, error) {
	query := `SELECT user_id, group_id, joined_at, accepted FROM user_groups WHERE id = ?`

	row := repo.DB.QueryRow(query, id)

	groupMember := &GroupMember{}

	err := row.Scan(&groupMember.UserId, &groupMember.GroupId, &groupMember.JoinedAt, &groupMember.Accepted)

	if err != nil {
		return nil, err
	}

	return groupMember, nil
}
