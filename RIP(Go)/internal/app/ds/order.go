package ds

// Order - модель заявки
type Order struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	FromCity  string         `json:"from_city" gorm:"not null"`
	ToCity    string         `json:"to_city" gorm:"not null"`
	Weight    float64        `json:"weight" gorm:"not null"`
	Length    float64        `json:"length" gorm:"not null"`
	Width     float64        `json:"width" gorm:"not null"`
	Height    float64        `json:"height" gorm:"not null"`
	Services  []OrderService `json:"services" gorm:"foreignKey:OrderID"`
	TotalCost float64        `json:"total_cost" gorm:"not null"`
	TotalDays int            `json:"total_days" gorm:"not null"`
}

// OrderService - услуга в заявке
type OrderService struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	OrderID   int    `json:"order_id" gorm:"not null"`
	ServiceID int    `json:"service_id" gorm:"not null"`
	Quantity  int    `json:"quantity" gorm:"not null"`
	Comment   string `json:"comment" gorm:"type:text"`
}
