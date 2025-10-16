package repository

import (
    "database/sql"
    "fmt"
    "strings"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "rip-go-app/internal/app/ds"
    "rip-go-app/internal/app/calculator"
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

// CRUD для Service
func (r *Repository) CreateService(s *ds.Service) error {
    return r.db.Create(s).Error
}

func (r *Repository) UpdateService(s *ds.Service) error {
    return r.db.Save(s).Error
}

func (r *Repository) DeleteService(id int) error {
    return r.db.Delete(&ds.Service{}, id).Error
}

// CreateCargoOrder создаёт заказ на основе перечня транспортов и параметров груза
type CargoOrderItem struct {
    ServiceID int
    FromCity  string
    ToCity    string
    Length    float64
    Width     float64
    Height    float64
    Weight    float64
}

func (r *Repository) CreateCargoOrder(items []CargoOrderItem) (int, error) {
    if len(items) == 0 {
        return 0, fmt.Errorf("no items provided")
    }

    return r.createCargoOrderTx(items)
}

func (r *Repository) createCargoOrderTx(items []CargoOrderItem) (int, error) {
    calc := calculator.NewDeliveryCalculator()

    returnID := 0
    err := r.db.Transaction(func(tx *gorm.DB) error {
        // Используем параметры первого как общие
        first := items[0]

        order := ds.Order{
            SessionID: "guest",
            IsDraft:   false,
            FromCity:  first.FromCity,
            ToCity:    first.ToCity,
            Weight:    0,
            Length:    0,
            Width:     0,
            Height:    0,
            TotalCost: 0,
            TotalDays: 0,
            Status:    "pending",
        }
        if err := tx.Create(&order).Error; err != nil {
            return err
        }

        // агрегаты
        maxDays := 0
        totalCost := 0.0
        totalWeight := 0.0
        totalLength := 0.0
        totalWidth := 0.0
        totalHeight := 0.0

        for _, it := range items {
            svc, err := r.GetService(it.ServiceID)
            if err != nil {
                return fmt.Errorf("service %d not found", it.ServiceID)
            }
            res := calc.CalculateDelivery(svc, it.FromCity, it.ToCity, it.Length, it.Width, it.Height, it.Weight)
            if !res.IsValid {
                return fmt.Errorf(res.ErrorMessage)
            }

            // создаём строку заказа
            os := ds.OrderService{OrderID: order.ID, ServiceID: it.ServiceID, Quantity: 1}
            if err := tx.Create(&os).Error; err != nil {
                return err
            }

            if res.DeliveryDays > maxDays { maxDays = res.DeliveryDays }
            totalCost += res.TotalCost
            totalWeight += it.Weight
            totalLength += it.Length
            totalWidth += it.Width
            totalHeight += it.Height
        }

        // итоговые поля заказа
        order.TotalDays = maxDays
        order.TotalCost = totalCost
        order.Weight = totalWeight
        order.Length = totalLength
        order.Width = totalWidth
        order.Height = totalHeight

        if err := tx.Save(&order).Error; err != nil {
            return err
        }

        returnID = order.ID
        return nil
    })

    if err != nil {
        return 0, err
    }
    return returnID, nil
}

// GetOrder - получение заявки
func (r *Repository) GetOrder() (ds.Order, error) {
    var order ds.Order
    err := r.db.Where("is_draft = ?", false).Preload("Services").First(&order).Error
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
func (r *Repository) ensureDraftOrder(sessionID string) (int, error) {
    var order ds.Order
    if err := r.db.Where("session_id = ? AND is_draft = ?", sessionID, true).First(&order).Error; err != nil {
        // создаём
        order = ds.Order{SessionID: sessionID, IsDraft: true}
        if err := r.db.Create(&order).Error; err != nil {
            return 0, err
        }
    }
    return order.ID, nil
}

func (r *Repository) AddToCart(serviceID int) error {
    // проверяем услугу
    if _, err := r.GetService(serviceID); err != nil {
        return fmt.Errorf("услуга не найдена")
    }
    // берём черновой заказ как корзину
    orderID, err := r.ensureDraftOrder("guest")
    if err != nil { return err }

    // upsert в order_items
    // используем нативное подключение для ON CONFLICT
    sqlDB, err := r.db.DB(); if err != nil { return err }
    _, err = sqlDB.Exec(`
        INSERT INTO order_items(order_id, service_id, quantity)
        VALUES ($1, $2, 1)
        ON CONFLICT (order_id, service_id)
        DO UPDATE SET quantity = order_items.quantity + 1
    `, orderID, serviceID)
    return err
}

// RemoveFromCart - удаление услуги из корзины
func (r *Repository) RemoveFromCart(serviceID int) error {
    orderID, err := r.ensureDraftOrder("guest")
    if err != nil { return err }

    sqlDB, err := r.db.DB(); if err != nil { return err }
    // уменьшаем qty если >1, иначе удаляем
    var qty int
    err = sqlDB.QueryRow(`SELECT quantity FROM order_items WHERE order_id=$1 AND service_id=$2`, orderID, serviceID).Scan(&qty)
    if err == sql.ErrNoRows { return fmt.Errorf("услуга не найдена в корзине") }
    if err != nil { return err }

    if qty > 1 {
        _, err = sqlDB.Exec(`UPDATE order_items SET quantity = quantity - 1 WHERE order_id=$1 AND service_id=$2`, orderID, serviceID)
    } else {
        _, err = sqlDB.Exec(`DELETE FROM order_items WHERE order_id=$1 AND service_id=$2`, orderID, serviceID)
    }
    return err
}

// GetCart - получение корзины
func (r *Repository) GetCart() (ds.Cart, error) {
    orderID, err := r.ensureDraftOrder("guest")
    if err != nil { return ds.Cart{}, err }
    var items []ds.CartService
    if err := r.db.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
        return ds.Cart{}, err
    }
    return ds.Cart{ID: orderID, SessionID: "guest", IsDraft: true, Services: items}, nil
}

// GetCartServices - получение услуг в корзине с полной информацией
func (r *Repository) GetCartServices() ([]ds.Service, error) {
    orderID, err := r.ensureDraftOrder("guest")
    if err != nil { return nil, err }
    var items []ds.CartService
    if err := r.db.Where("order_id = ?", orderID).Find(&items).Error; err != nil { return nil, err }
    services := make([]ds.Service, 0, len(items))
    for _, it := range items {
        s, err := r.GetService(it.ServiceID)
        if err == nil { services = append(services, s) }
    }
    return services, nil
}

// GetCartCount - получение общего количества услуг в корзине
func (r *Repository) GetCartCount() int {
    orderID, err := r.ensureDraftOrder("guest")
    if err != nil { return 0 }
    sqlDB, err := r.db.DB(); if err != nil { return 0 }
    var count sql.NullInt64
    _ = sqlDB.QueryRow(`SELECT COALESCE(SUM(quantity),0) FROM order_items WHERE order_id=$1`, orderID).Scan(&count)
    if count.Valid { return int(count.Int64) }
    return 0
}

// ClearCart - очистка корзины
func (r *Repository) ClearCart() {
    orderID, err := r.ensureDraftOrder("guest")
    if err != nil { return }
    r.db.Where("order_id = ?", orderID).Delete(&ds.CartService{})
}

// UpdateOrderStatusWithCursor - обновление статуса заказа через курсор без ORM
func (r *Repository) UpdateOrderStatusWithCursor(orderID int, newStatus string) error {
	// Получаем нативное подключение к БД из GORM
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Подготавливаем запрос с курсором
	query := `
		UPDATE orders 
		SET status = $1 
		WHERE id = $2
		RETURNING id, status, from_city, to_city
	`

	// Выполняем запрос через курсор
	rows, err := sqlDB.Query(query, newStatus, orderID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Обрабатываем результат через курсор
	if rows.Next() {
		var id int
		var status, fromCity, toCity string
		
		err := rows.Scan(&id, &status, &fromCity, &toCity)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		
		// Логируем обновление
		fmt.Printf("Order %d status updated to: %s (Route: %s -> %s)\n", id, status, fromCity, toCity)
	} else {
		return fmt.Errorf("order with id %d not found", orderID)
	}

	return nil
}
