package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id          string        `json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"required,max=40"`
	Username    string        `json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty" validate:"required,username,max=100"`
	Email       string        `json:"email2,omitempty" gorm:"column:email" bson:"email3,omitempty" dynamodbav:"email,omitempty" firestore:"email,omitempty" validate:"email,max=100"`
	Phone       string        `json:"phone,omitempty" gorm:"column:phone" bson:"phone,omitempty" dynamodbav:"phone,omitempty" firestore:"required,phone,omitempty" validate:"required,phone,max=18"`
	DateOfBirth *time.Time    `json:"dateOfBirth,omitempty" gorm:"column:date_of_birth" bson:"dateOfBirth,omitempty" dynamodbav:"dateOfBirth,omitempty" firestore:"dateOfBirth,omitempty"`
	Interests   []string      `json:"interests,omitempty" gorm:"-" bson:"interests,omitempty" firestore:"interests" mar:"InterestsB"`
	Skills      Skills        `json:"skills,omitempty" gorm:"type:jsonb;column:skills" bson:"skills,omitempty" firestore:"skills"`
	Settings    *UserSettings `json:"settings,omitempty" gorm:"type:settings;column:skills" bson:"settings,omitempty" firestore:"settings"`
}

type UserSettings struct {
	UserId         string `bson:"_id" json:"-"`
	Language       string `bson:"language" json:"language"`
	DateFormat     string `bson:"dateFormat" json:"dateFormat"`
	DateTimeFormat string `bson:"dateTimeFormat" json:"dateTimeFormat"`
	TimeFormat     string `bson:"timeFormat" json:"timeFormat"`
	Notification   bool   `bson:"notification" json:"notification"`
}

type Skill struct {
	Skill   string `json:"skill" gorm:"primary_key;" bson:"skill"`
	Hirable bool   `json:"hirable" bson:"hirable" `
}

type Skills []Skill

func (o Skills) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *Skills) Scan(input interface{}) error {
	bytes, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", input))
	}
	return json.Unmarshal(bytes, o)
}
