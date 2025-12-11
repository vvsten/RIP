package ds

// Cart — представление корзины через логистическую заявку-черновик
type Cart struct {
    ID        int          `json:"id" gorm:"primaryKey"`    // это id заявки в таблице logistic_requests
    SessionID string       `json:"session_id"`
    IsDraft   bool         `json:"is_draft" gorm:"not null;default:true"`
    Services  []CartService `json:"services" gorm:"foreignKey:LogisticRequestID"`
}

func (Cart) TableName() string { return "logistic_requests" }

// CartService — строка корзины хранится в logistic_request_services
type CartService struct {
    ID                 int `json:"id" gorm:"primaryKey"`
    LogisticRequestID  int `json:"logistic_request_id" gorm:"not null"`
    TransportServiceID int `json:"transport_service_id" gorm:"not null"`
    Quantity           int `json:"quantity" gorm:"not null"`
}

func (CartService) TableName() string { return "logistic_request_services" }
