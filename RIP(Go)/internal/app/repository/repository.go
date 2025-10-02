package repository

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rip-go-app/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db: db,
	}, nil
}

// GetServices - получение всех услуг с возможностью фильтрации
func (r *Repository) GetServices(search string) ([]ds.Service, error) {
	var services []ds.Service
	
	if search == "" {
		err := r.db.Find(&services).Error
		if err != nil {
			return nil, err
		}
		return services, nil
	}

	searchLower := strings.ToLower(search)
	err := r.db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", 
		"%"+searchLower+"%", "%"+searchLower+"%").Find(&services).Error
	if err != nil {
		return nil, err
	}

	return services, nil
}

// GetService - получение услуги по ID
func (r *Repository) GetService(id int) (ds.Service, error) {
	var service ds.Service
	err := r.db.Where("id = ?", id).First(&service).Error
	if err != nil {
		return ds.Service{}, fmt.Errorf("услуга не найдена")
	}
	return service, nil
}

// GetServiceByType - получение услуги по типу доставки
func (r *Repository) GetServiceByType(deliveryType string) (ds.Service, error) {
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
	
	return ds.Service{}, fmt.Errorf("тип доставки не найден")
}

// GetOrder - получение заявки
func (r *Repository) GetOrder() (ds.Order, error) {
	var order ds.Order
	err := r.db.Preload("Services").First(&order).Error
	if err != nil {
		return ds.Order{}, err
	}
	return order, nil
}

// GetServicesInOrder - получение услуг в заявке
func (r *Repository) GetServicesInOrder() ([]ds.Service, error) {
	var order ds.Order
	err := r.db.Preload("Services").First(&order).Error
	if err != nil {
		return nil, err
	}

	var services []ds.Service
	for _, orderService := range order.Services {
		service, err := r.GetService(orderService.ServiceID)
		if err != nil {
			continue
		}
		services = append(services, service)
	}
	return services, nil
}

// AddToCart - добавление услуги в корзину
func (r *Repository) AddToCart(serviceID int) error {
	// Проверяем, существует ли услуга
	_, err := r.GetService(serviceID)
	if err != nil {
		return fmt.Errorf("услуга не найдена")
	}

	// Получаем или создаем корзину
	var cart ds.Cart
	err = r.db.First(&cart).Error
	if err != nil {
		// Создаем новую корзину
		cart = ds.Cart{}
		err = r.db.Create(&cart).Error
		if err != nil {
			return err
		}
	}

	// Проверяем, есть ли уже такая услуга в корзине
	var cartService ds.CartService
	err = r.db.Where("cart_id = ? AND service_id = ?", cart.ID, serviceID).First(&cartService).Error
	if err == nil {
		// Увеличиваем количество
		cartService.Quantity++
		err = r.db.Save(&cartService).Error
		if err != nil {
			return err
		}
	} else {
		// Добавляем новую услугу в корзину
		cartService = ds.CartService{
			CartID:    cart.ID,
			ServiceID: serviceID,
			Quantity:  1,
		}
		err = r.db.Create(&cartService).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveFromCart - удаление услуги из корзины
func (r *Repository) RemoveFromCart(serviceID int) error {
	var cart ds.Cart
	err := r.db.First(&cart).Error
	if err != nil {
		return fmt.Errorf("корзина не найдена")
	}

	var cartService ds.CartService
	err = r.db.Where("cart_id = ? AND service_id = ?", cart.ID, serviceID).First(&cartService).Error
	if err != nil {
		return fmt.Errorf("услуга не найдена в корзине")
	}

	if cartService.Quantity > 1 {
		cartService.Quantity--
		err = r.db.Save(&cartService).Error
		if err != nil {
			return err
		}
	} else {
		err = r.db.Delete(&cartService).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCart - получение корзины
func (r *Repository) GetCart() (ds.Cart, error) {
	var cart ds.Cart
	err := r.db.Preload("Services").First(&cart).Error
	if err != nil {
		return ds.Cart{}, err
	}
	return cart, nil
}

// GetCartServices - получение услуг в корзине с полной информацией
func (r *Repository) GetCartServices() ([]ds.Service, error) {
	var cart ds.Cart
	err := r.db.Preload("Services").First(&cart).Error
	if err != nil {
		// Если корзина не найдена, возвращаем пустой массив
		return []ds.Service{}, nil
	}

	var services []ds.Service
	for _, cartService := range cart.Services {
		service, err := r.GetService(cartService.ServiceID)
		if err != nil {
			continue
		}
		services = append(services, service)
	}
	return services, nil
}

// GetCartCount - получение общего количества услуг в корзине
func (r *Repository) GetCartCount() int {
	var cart ds.Cart
	err := r.db.Preload("Services").First(&cart).Error
	if err != nil {
		return 0
	}

	count := 0
	for _, cartService := range cart.Services {
		count += cartService.Quantity
	}
	return count
}

// ClearCart - очистка корзины
func (r *Repository) ClearCart() {
	var cart ds.Cart
	err := r.db.First(&cart).Error
	if err != nil {
		return
	}
	
	r.db.Where("cart_id = ?", cart.ID).Delete(&ds.CartService{})
}
