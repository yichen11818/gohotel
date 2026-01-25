package models

import "time"

type HotelSettings struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	HotelID      int64     `gorm:"not null;uniqueIndex" json:"hotel_id"`
	SettingsJSON string    `gorm:"type:longtext" json:"settings_json"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (HotelSettings) TableName() string {
	return "hotel_settings"
}
