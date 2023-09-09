package models

import (

	"time"

)

type User struct {
	IdUser       uint   `json:"id,omitempty" bson:"_id,omitempty"`
	Nama      string `json:"nama" bson:"nama"`
	Username  string `json:"username" bson:"username"`
	Password  string `json:"password" bson:"password"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time
}