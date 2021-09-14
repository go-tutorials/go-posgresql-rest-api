package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	. "go-service/internal/models"
	"log"
	"strings"
)

type SqlUserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *SqlUserService {
	return &SqlUserService{DB: db}
}

func (m *SqlUserService) GetAll(ctx context.Context) (*[]User, error) {
	query := "select id, username, email, phone, date_of_birth, interests, skills, achievements, settings from users"
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	var res []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Phone, &user.Email, &user.DateOfBirth, pq.Array(&user.Interests), pq.Array(&user.Skills), pq.Array(&user.Achievements), &user.Settings)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	return &res, nil
}

func (m *SqlUserService) Load(ctx context.Context, id string) (*User, error) {
	var user User
	query := "select id, username, email, phone, date_of_birth, interests, skills, achievements, settings from users WHERE ID = $1"
	row,err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	for row.Next(){
		err = row.Scan(&user.Id, &user.Username, &user.Phone, &user.Email, &user.DateOfBirth, pq.Array(&user.Interests), pq.Array(&user.Skills), pq.Array(&user.Achievements), &user.Settings)
	}
	if err != nil {
		errMsg := err.Error()
		if strings.Compare(fmt.Sprintf(errMsg), "0 row(s) returned") == 0 {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func (m *SqlUserService) Insert(ctx context.Context, user *User) (int64, error) {
	query := "insert into users (id, username, email, phone, date_of_birth, interests, skills, achievements, settings) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	result,err := m.DB.Exec(query, user.Id, user.Username, user.Email, user.Phone, user.DateOfBirth, pq.Array(user.Interests), pq.Array(user.Skills), pq.Array(user.Achievements), user.Settings)
	if err != nil {
		log.Println(err)
		return -1, nil
	}
	return result.RowsAffected()
}

func (m *SqlUserService) Update(ctx context.Context, user *User) (int64, error) {
	query := "update users set username = $2, email = $3, phone = $4, date_of_birth = $5, interests = $6, skills = $7, achievements = $8, settings = $9 where id = $1"
	result, err := m.DB.Exec(query, user.Id, user.Username, user.Email, user.Phone, user.DateOfBirth, pq.Array(user.Interests), pq.Array(user.Skills), pq.Array(user.Achievements), user.Settings)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (m *SqlUserService) Delete(ctx context.Context, id string) (int64, error) {
	query := "delete from users where id = $1"
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return -1, err
	}
	rowAffect, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	fmt.Printf("Total rows/records has been deleted %v", rowAffect)
	return rowAffect, nil
}
