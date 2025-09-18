package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
	services []Service
	order    Order
}

func NewRepository() (*Repository, error) {
	repo := &Repository{}
	repo.initData()
	return repo, nil
}

// Service - модель услуги (тип транспорта)
type Service struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	DeliveryDays int    `json:"delivery_days"`
	MaxWeight   float64 `json:"max_weight"`
	MaxVolume   float64 `json:"max_volume"`
}

// Order - модель заявки
type Order struct {
	ID           int               `json:"id"`
	FromCity     string            `json:"from_city"`
	ToCity       string            `json:"to_city"`
	Weight       float64           `json:"weight"`
	Length       float64           `json:"length"`
	Width        float64           `json:"width"`
	Height       float64           `json:"height"`
	Services     []OrderService    `json:"services"`
	TotalCost    float64           `json:"total_cost"`
	TotalDays    int               `json:"total_days"`
}

// OrderService - услуга в заявке
type OrderService struct {
	ServiceID int     `json:"service_id"`
	Quantity  int     `json:"quantity"`
	Comment   string  `json:"comment"`
}

func (r *Repository) initData() {
	// Инициализируем услуги
	r.services = []Service{
		{
			ID:           1,
			Name:         "Фура",
			Description:  "Полуприцеп для перевозки крупногабаритных грузов. Идеально подходит для перевозки мебели, строительных материалов и других тяжелых грузов.",
			Price:        150.0,
			ImageURL:     "http://localhost:9000/lab1/fura.jpg",
			DeliveryDays: 2,
			MaxWeight:    20000.0,
			MaxVolume:    80.0,
		},
		{
			ID:           2,
			Name:         "Малотоннажный грузовик",
			Description:  "Легкий грузовик для перевозки небольших грузов по городу и между городами. Быстрая доставка с возможностью проезда в центр города.",
			Price:        80.0,
			ImageURL:     "http://localhost:9000/lab1/malotonnazhnyi.jpg",
			DeliveryDays: 1,
			MaxWeight:    3000.0,
			MaxVolume:    15.0,
		},
		{
			ID:           3,
			Name:         "Авиаперевозка",
			Description:  "Быстрая доставка грузов авиатранспортом. Подходит для срочных и ценных грузов. Максимальная скорость доставки.",
			Price:        500.0,
			ImageURL:     "http://localhost:9000/lab1/avia.jpg",
			DeliveryDays: 1,
			MaxWeight:    1000.0,
			MaxVolume:    5.0,
		},
		{
			ID:           4,
			Name:         "Поезд",
			Description:  "Железнодорожные перевозки для крупных партий грузов. Экономичный вариант для больших объемов.",
			Price:        120.0,
			ImageURL:     "http://localhost:9000/lab1/poezd.jpg",
			DeliveryDays: 3,
			MaxWeight:    50000.0,
			MaxVolume:    120.0,
		},
		{
			ID:           5,
			Name:         "Корабль",
			Description:  "Морские перевозки для международной доставки. Подходит для крупных партий и контейнерных перевозок.",
			Price:        200.0,
			ImageURL:     "http://localhost:9000/lab1/korabl.jpg",
			DeliveryDays: 7,
			MaxWeight:    100000.0,
			MaxVolume:    500.0,
		},
		{
			ID:           6,
			Name:         "Мультимодальные",
			Description:  "Комбинированные перевозки с использованием нескольких видов транспорта. Оптимальное решение для сложных маршрутов.",
			Price:        300.0,
			ImageURL:     "http://localhost:9000/lab1/multimodal.jpg",
			DeliveryDays: 5,
			MaxWeight:    30000.0,
			MaxVolume:    100.0,
		},
	}

	// Инициализируем заявку
	r.order = Order{
		ID:        1,
		FromCity:  "Москва",
		ToCity:    "Санкт-Петербург",
		Weight:    500.0,
		Length:    2.0,
		Width:     1.5,
		Height:    1.0,
		Services: []OrderService{
			{ServiceID: 1, Quantity: 1, Comment: "Основная доставка"},
			{ServiceID: 2, Quantity: 1, Comment: "Дополнительная услуга"},
		},
		TotalCost: 230.0,
		TotalDays: 2,
	}
}

// GetServices - получение всех услуг с возможностью фильтрации
func (r *Repository) GetServices(search string) ([]Service, error) {
	if search == "" {
		return r.services, nil
	}

	var filtered []Service
	searchLower := strings.ToLower(search)
	
	for _, service := range r.services {
		if strings.Contains(strings.ToLower(service.Name), searchLower) ||
		   strings.Contains(strings.ToLower(service.Description), searchLower) {
			filtered = append(filtered, service)
		}
	}

	return filtered, nil
}

// GetService - получение услуги по ID
func (r *Repository) GetService(id int) (Service, error) {
	for _, service := range r.services {
		if service.ID == id {
			return service, nil
		}
	}
	return Service{}, fmt.Errorf("услуга не найдена")
}

// GetOrder - получение заявки
func (r *Repository) GetOrder() (Order, error) {
	return r.order, nil
}

// GetServicesInOrder - получение услуг в заявке
func (r *Repository) GetServicesInOrder() ([]Service, error) {
	var services []Service
	for _, orderService := range r.order.Services {
		service, err := r.GetService(orderService.ServiceID)
		if err != nil {
			continue
		}
		services = append(services, service)
	}
	return services, nil
}

// GetServiceByType - получение услуги по типу доставки
func (r *Repository) GetServiceByType(deliveryType string) (Service, error) {
	typeMap := map[string]int{
		"fura":           1,
		"malotonnazhnyi": 2,
		"avia":           3,
		"poezd":          4,
		"korabl":         5,
		"multimodal":     6,
	}
	
	if id, exists := typeMap[deliveryType]; exists {
		return r.GetService(id)
	}
	
	return Service{}, fmt.Errorf("тип доставки не найден")
}