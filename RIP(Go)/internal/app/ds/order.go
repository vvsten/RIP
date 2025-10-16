package ds

// Order - модель заявки
type Order struct {
    ID        int            `json:"id" gorm:"primaryKey"`
    SessionID string         `json:"session_id"`
    IsDraft   bool           `json:"is_draft" gorm:"not null;default:true"`
    FromCity  string         `json:"from_city"`
    ToCity    string         `json:"to_city"`
    // Параметры груза
    Weight    float64        `json:"weight" gorm:"not null;default:0"`
    Length    float64        `json:"length" gorm:"not null;default:0"`
    Width     float64        `json:"width" gorm:"not null;default:0"`
    Height    float64        `json:"height" gorm:"not null;default:0"`
    Services  []OrderService `json:"services" gorm:"foreignKey:OrderID"`
    TotalCost float64        `json:"total_cost"`
    TotalDays int            `json:"total_days"`
    Status    string         `json:"status" gorm:"type:varchar(32);not null;default:'pending'"`
}

// OrderService - услуга в заявке
type OrderService struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	OrderID   int    `json:"order_id" gorm:"not null"`
	ServiceID int    `json:"service_id" gorm:"not null"`
	Quantity  int    `json:"quantity" gorm:"not null"`
	Comment   string `json:"comment" gorm:"type:text"`
}
