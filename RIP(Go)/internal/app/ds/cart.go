package ds

// Cart - корзина с выбранными услугами
type Cart struct {
	ID       int          `json:"id" gorm:"primaryKey"`
	Services []CartService `json:"services" gorm:"foreignKey:CartID"`
}

// CartService - услуга в корзине
type CartService struct {
	ID        int `json:"id" gorm:"primaryKey"`
	CartID    int `json:"cart_id" gorm:"not null"`
	ServiceID int `json:"service_id" gorm:"not null"`
	Quantity  int `json:"quantity" gorm:"not null"`
}
