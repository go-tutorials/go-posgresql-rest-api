package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	. "go-service/internal/models"
	"strings"
	"time"
)

type SqlUserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *SqlUserService {
	return &SqlUserService{DB: db}
}

func (m *SqlUserService) GetAll(ctx context.Context) (*[]User, error) {
	query := "select id, username, email, phone, date_of_birth, interests, skills, settings from users"
	var user User
	var bufferSkills []byte
	var bufferSettings []byte
	var interests pq.StringArray
	var skills []*Skill
	var settings UserSettings
	var result []User
	user.Interests = interests

	rows, err := m.DB.Query(query)
	if err != nil {
		fmt.Printf("The is an error when reading the user list: %v", err)
	}

	for rows.Next() {
		user := User{}

		err = rows.Scan(&user.Id, &user.Username, &user.Phone, &user.Email, &user.DateOfBirth, &interests, &bufferSkills, &bufferSettings)
		if err != nil {
			fmt.Printf("The is an error when reading the user list: %v", err)
		}

		user.Interests = interests
		err = json.Unmarshal(bufferSkills, &skills)
		if err != nil {
			fmt.Printf("The is an error when reading the skill of user: %v", err)
		}

		for _, v := range skills {
			user.Skills = append(user.Skills, Skill{Skill: v.Skill, Hirable: v.Hirable})
		}

		err = json.Unmarshal(bufferSettings, &settings)
		if err != nil {
			fmt.Printf("The is an error when reading the setting of user: %v", err)
		}

		user.Settings = &settings
		result = append(result, user)
	}
	return &result, nil
}

func (m *SqlUserService) Load(ctx context.Context, id string) (*User, error) {
	var user User
	var bufferSkills []byte
	var bufferSettings []byte
	var interests pq.StringArray
	var skills []*Skill
	var settings UserSettings

	query := "select id, username, email, phone, date_of_birth, interests, skills, settings from users WHERE ID = '$1'"
	query = strings.ReplaceAll(query, "$1", id)

	user.Interests = interests

	row := m.DB.QueryRow(query)

	err := row.Scan(&user.Id, &user.Username, &user.Phone, &user.Email, &user.DateOfBirth, &interests, &bufferSkills, &bufferSettings)
	if err != nil {
		errMsg := err.Error()
		if strings.Compare(fmt.Sprintf(errMsg), "0 row(s) returned") == 0 {
			return nil, nil
		} else {
			return nil, err
		}
	}

	user.Interests = interests

	err = json.Unmarshal(bufferSkills, &skills)
	err = json.Unmarshal(bufferSkills, &skills)
	if err != nil {
		fmt.Printf("The is an error when reading the skill of user: %v", err)
	}

	for _, v := range skills {
		user.Skills = append(user.Skills, Skill{Skill: v.Skill, Hirable: v.Hirable})
	}

	err = json.Unmarshal(bufferSettings, &settings)
	if err != nil {
		fmt.Printf("The is an error when reading the setting of user: %v", err)
	}

	user.Settings = &settings

	return &user, nil
}

func (m *SqlUserService) Insert(ctx context.Context, user *User) (int64, error) {
	query := "insert into users (id, username, email, phone, date_of_birth, interests, skills, settings) values ('$1', '$2', '$3', '$4', '$5', '$6', '$7', '$8') RETURNING id"
	query = strings.ReplaceAll(query, "$1", user.Id)
	query = strings.ReplaceAll(query, "$2", user.Username)
	query = strings.ReplaceAll(query, "$3", user.Email)
	query = strings.ReplaceAll(query, "$4", user.Phone)
	query = strings.ReplaceAll(query, "$5", user.DateOfBirth.Format(time.RFC3339))
	query = strings.ReplaceAll(query, "$6", formatStringArrayToPostGres(user.Interests))
	query = strings.ReplaceAll(query, "$7", formatStructToJSON(user.Skills))
	query = strings.ReplaceAll(query, "$8", formatStructToJSON(user.Settings))

	var id int64
	err := m.DB.QueryRow(query).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return -1, nil
	}

	fmt.Printf("Inserted a new user succesfully with id:= %v", id)
	return id, nil
}

func (m *SqlUserService) Update(ctx context.Context, user *User) (int64, error) {
	query := "update users set username = '$2', email = '$3', phone = '$4', date_of_birth = '$5', interests = '$6', skills = '$7', settings = '$8' where id = '$1'  RETURNING id"
	query = strings.ReplaceAll(query, "$1", user.Id)
	query = strings.ReplaceAll(query, "$2", user.Username)
	query = strings.ReplaceAll(query, "$3", user.Email)
	query = strings.ReplaceAll(query, "$4", user.Phone)
	query = strings.ReplaceAll(query, "$5", user.DateOfBirth.Format(time.RFC3339))
	query = strings.ReplaceAll(query, "$6", formatStringArrayToPostGres(user.Interests))
	query = strings.ReplaceAll(query, "$7", formatStructToJSON(user.Skills))
	query = strings.ReplaceAll(query, "$8", formatStructToJSON(user.Settings))

	result, err := m.DB.Exec(query)
	if err != nil {
		return -1, err
	}
	rowAffect, err := result.RowsAffected()
	fmt.Printf("Total records has been updated %v", rowAffect)
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

func formatStructToJSON(object interface{}) string {
	jsonData, err := json.Marshal(object)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return string(jsonData)
}

func formatStringArrayToPostGres(object interface{}) string  {
	result := formatStructToJSON(object)
	result = strings.ReplaceAll(result, "[", "{")
	result = strings.ReplaceAll(result, "]", "}")
	return result
}
