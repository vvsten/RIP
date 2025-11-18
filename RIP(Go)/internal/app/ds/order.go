package ds

import "time"

// Order - модель заявки
type Order struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	SessionID string `json:"session_id"`
	IsDraft   bool   `json:"is_draft" gorm:"not null;default:true"`
	FromCity  string `json:"from_city"`
	ToCity    string `json:"to_city"`
	// Параметры груза
	Weight    float64        `json:"weight" gorm:"not null;default:0"`
	Length    float64        `json:"length" gorm:"not null;default:0"`
	Width     float64        `json:"width" gorm:"not null;default:0"`
	Height    float64        `json:"height" gorm:"not null;default:0"`
	Services  []OrderService `json:"TransportService" gorm:"foreignKey:OrderID"`
	TotalCost float64        `json:"total_cost"`
	TotalDays int            `json:"total_days"`
	Status    string         `json:"status" gorm:"type:varchar(32);not null;default:'draft'"`

	// Системные поля
	CreatorID   int        `json:"creator_id" gorm:"not null"`
	ModeratorID *int       `json:"moderator_id"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	FormedAt    *time.Time `json:"formed_at"`
	CompletedAt *time.Time `json:"completed_at"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"-" gorm:"index"`

	// Связи
	Creator   User  `json:"creator" gorm:"foreignKey:CreatorID"`
	Moderator *User `json:"moderator" gorm:"foreignKey:ModeratorID"`
}

// OrderStatus - статусы заявки
const (
	StatusDraft     = "draft"     // черновик
	StatusFormed    = "formed"    // сформирован
	StatusCompleted = "completed" // завершён
	StatusRejected  = "rejected"  // отклонён
	StatusDeleted   = "deleted"   // удалён
)

// OrderService - услуга в заявке (м-м)
type OrderService struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	OrderID   int    `json:"order_id" gorm:"not null"`
	ServiceID int    `json:"service_id" gorm:"not null"`
	Quantity  int    `json:"quantity" gorm:"not null;default:1"`
	Comment   string `json:"comment" gorm:"type:text"`
	Order     int    `json:"order" gorm:"not null;default:0"` // порядок в заявке

	// Связи
	Service Service `json:"service" gorm:"foreignKey:ServiceID"`
}

// TableName переопределяет имя таблицы для Order
func (Order) TableName() string {
	return "LogisticRequest"
}

// TableName переопределяет имя таблицы для OrderService
func (OrderService) TableName() string {
	return "order_TransportService"
}
