package ds

import "time"

// User - модель пользователя
type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Login     string    `json:"login" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"` // не возвращаем в JSON
	Name      string    `json:"name" gorm:"not null"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role" gorm:"not null;default:'user'"` // user, moderator, admin
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// UserRole - роли пользователей
const (
	RoleUser     = "user"
	RoleModerator = "moderator"
	RoleAdmin    = "admin"
)
