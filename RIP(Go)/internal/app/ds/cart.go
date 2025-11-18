package ds

// Cart — представление корзины через заказ-черновик (LogisticRequest)
type Cart struct {
    ID        int          `json:"id" gorm:"primaryKey"`    // это id заказа в таблице LogisticRequest
    SessionID string       `json:"session_id"`
    IsDraft   bool         `json:"is_draft" gorm:"not null;default:true"`
    Services  []CartService `json:"TransportService" gorm:"foreignKey:OrderID"`
}

func (Cart) TableName() string { return "LogisticRequest" }

// CartService — строка корзины хранится в order_TransportService
type CartService struct {
    ID        int `json:"id" gorm:"primaryKey"`
    OrderID   int `json:"order_id" gorm:"not null"`
    ServiceID int `json:"service_id" gorm:"not null"`
    Quantity  int `json:"quantity" gorm:"not null"`
}

func (CartService) TableName() string { return "order_TransportService" }
