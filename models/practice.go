package models

import (

	"time"
)

type Latihan struct {
	Id       uint   `json:"id,omitempty" bson:"_id,omitempty"`
	Judul string `json:"judul" bson:"judul"`
	Deskripsi string `json:"deskripsi" bson:"deskripsi"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time
}



