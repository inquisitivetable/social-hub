package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

type Group struct {
	Id          int64
	CreatorId   int64
	Title       string
	Description string
	ImagePath   string
	CreatedAt   time.Time
}

type UserGroup struct {
	Id    int64  `json:"groupId"`
	Title string `json:"groupName"`
}

type GroupJSON struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImagePath   string `json:"imagePath"`
	IsMember    bool   `json:"isMember"`
	IsCreator   bool   `json:"isCreator"`
}

type SearchResult struct {
	GroupId   int64  `json:"groupId"`
	UserId    int64  `json:"userId"`
	Name      string `json:"name"`
	ImagePath string `json:"imagePath"`
}

type IGroupRepository interface {
	GetAllByCreatorId(userId int64) ([]*Group, error)
	GetAllByMemberId(userId int64) ([]*Group, error)
	GetById(id int64) (*Group, error)
	Insert(group *Group) (int64, error)
	SearchGroupsAndUsersByString(userId int64, searchString string) ([]*SearchResult, error)
	UpdateImagePath(groupId int64, imagePath string) error
}

type GroupRepository struct {
	Logger *log.Logger
	DB     *sql.DB
}

func NewGroupRepo(db *sql.DB) *GroupRepository {
	return &GroupRepository{
		Logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
		DB:     db,
	}
}

func (repo GroupRepository) Insert(group *Group) (int64, error) {

	query := `INSERT INTO groups (creator_id, title, description, created_at, image_path)
	VALUES(?, ?, ?, ?, ?)`

	args := []interface{}{
		group.CreatorId,
		group.Title,
		group.Description,
		time.Now(),
		group.ImagePath,
	}

	result, err := repo.DB.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	repo.Logger.Printf("Last inserted group '%s' by user %d (last insert ID: %d)", group.Title, group.CreatorId, lastId)

	return lastId, nil
}

func (p GroupRepository) GetById(id int64) (*Group, error) {
	query := `SELECT id, creator_id, title, description, created_at, image_path FROM groups WHERE id = ?`
	row := p.DB.QueryRow(query, id)
	group := &Group{}

	err := row.Scan(&group.Id, &group.CreatorId, &group.Title, &group.Description, &group.CreatedAt, &group.ImagePath)

	return group, err
}

func (repo GroupRepository) GetAllByCreatorId(userId int64) ([]*Group, error) {

	stmt := `SELECT id, creator_id,  title, description, created_at, image_path FROM groups
	WHERE creator_id = ?
    ORDER BY title ASC`

	rows, err := repo.DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := []*Group{}

	for rows.Next() {
		group := &Group{}

		err := rows.Scan(&group.Id, &group.CreatorId, &group.Title, &group.Description, &group.CreatedAt, &group.ImagePath)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (repo GroupRepository) GetAllByMemberId(userId int64) ([]*Group, error) {

	stmt := `SELECT DISTINCT g.id, g.creator_id,  g.title, g.description, g.created_at, g.image_path FROM groups g
	INNER JOIN user_groups ug ON
	g.id = ug.group_id
	WHERE ug.user_id = ? AND ug.accepted = TRUE
    ORDER BY title ASC`

	rows, err := repo.DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := []*Group{}

	for rows.Next() {
		group := &Group{}

		err := rows.Scan(&group.Id, &group.CreatorId, &group.Title, &group.Description, &group.CreatedAt, &group.ImagePath)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (repo GroupRepository) SearchGroupsAndUsersByString(userId int64, searchString string) ([]*SearchResult, error) {

	formattedSearchString := fmt.Sprintf("%%%s%%", searchString)

	//repo.Logger.Println(formattedSearchString)

	stmt := `SELECT * FROM(SELECT 0 as UserId, g.Id as GroupId, g.Title as Name, g.image_path as ImagePath FROM groups g
		UNION
		SELECT u.Id as UserId, 0 as GroupId, u.forname ||  " " || u.nickname || " " || u.surname as Name, u.image_path as ImagePath FROM users u)
	WHERE Name LIKE ? AND UserId != ?`

	rows, err := repo.DB.Query(stmt, formattedSearchString, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := []*SearchResult{}

	for rows.Next() {
		group := &SearchResult{}

		err := rows.Scan(&group.UserId, &group.GroupId, &group.Name, &group.ImagePath)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (repo GroupRepository) UpdateImagePath(groupId int64, imagePath string) error {

	stmt := `UPDATE groups SET image_path = ? WHERE id = ?`

	_, err := repo.DB.Exec(stmt, imagePath, groupId)

	if err != nil {
		return err
	}

	return nil
}
