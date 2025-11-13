package ds

import "time"

// Service - модель услуги (тип транспорта)
type Service struct {
	ID           int     `json:"id" gorm:"primaryKey"`
	Name         string  `json:"name" gorm:"not null"`
	Description  string  `json:"description" gorm:"type:text"`
	Price        float64 `json:"price" gorm:"not null"`
	ImageURL     string  `json:"image_url" gorm:"type:varchar(500)"`
	DeliveryDays int     `json:"delivery_days" gorm:"not null"`
	MaxWeight    float64 `json:"max_weight" gorm:"not null"`
	MaxVolume    float64 `json:"max_volume" gorm:"not null"`
	
	// Системные поля
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}
