package models

import "time"

type Hotel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null;size:100" json:"name"`
	Status    string    `gorm:"not null;default:'active';size:20" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Hotel) TableName() string {
	return "hotels"
}
