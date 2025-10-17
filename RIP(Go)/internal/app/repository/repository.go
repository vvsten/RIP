package repository

import (
    "database/sql"
    "fmt"
    "strings"
    "time"

    "github.com/google/uuid"
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

// GetServices - получение всех услуг с возможностью фильтрации (исключая удалённые)
func (r *Repository) GetServices(search string) ([]ds.Service, error) {
	var services []ds.Service
	
	query := r.db.Where("deleted_at IS NULL")
	
	if search != "" {
		searchLower := strings.ToLower(search)
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", 
			"%"+searchLower+"%", "%"+searchLower+"%")
	}
	
	err := query.Find(&services).Error
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

func (r *Repository) CreateCargoOrder(items []CargoOrderItem, creatorID int) (int, error) {
    if len(items) == 0 {
        return 0, fmt.Errorf("no items provided")
    }

    return r.createCargoOrderTx(items, creatorID)
}

func (r *Repository) createCargoOrderTx(items []CargoOrderItem, creatorID int) (int, error) {
    calc := calculator.NewDeliveryCalculator()

    returnID := 0
    err := r.db.Transaction(func(tx *gorm.DB) error {
        // Используем параметры первого как общие
        first := items[0]

        order := ds.Order{
            SessionID: "guest",
            IsDraft:   true,
            FromCity:  first.FromCity,
            ToCity:    first.ToCity,
            Weight:    0,
            Length:    0,
            Width:     0,
            Height:    0,
            TotalCost: 0,
            TotalDays: 0,
            Status:    ds.StatusDraft,
            CreatorID: creatorID, // используем переданный creatorID
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

// ==================== ПОЛЬЗОВАТЕЛИ ====================

// CreateUser - создание пользователя
func (r *Repository) CreateUser(user *ds.User) error {
    // Генерируем UUID если он не задан
    if user.UUID == "" {
        user.UUID = uuid.New().String()
    }
    return r.db.Create(user).Error
}

// GetUserByLogin - получение пользователя по логину
func (r *Repository) GetUserByLogin(login string) (ds.User, error) {
    var user ds.User
    err := r.db.Where("login = ?", login).First(&user).Error
    if err != nil {
        return ds.User{}, fmt.Errorf("пользователь не найден")
    }
    return user, nil
}

// GetUser - получение пользователя по ID
func (r *Repository) GetUser(id int) (ds.User, error) {
    var user ds.User
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        return ds.User{}, fmt.Errorf("пользователь не найден")
    }
    return user, nil
}

// GetUserByUUID - получение пользователя по UUID
func (r *Repository) GetUserByUUID(userUUID string) (ds.User, error) {
    var user ds.User
    err := r.db.Where("uuid = ?", userUUID).First(&user).Error
    if err != nil {
        return ds.User{}, fmt.Errorf("пользователь не найден")
    }
    return user, nil
}

// UpdateUser - обновление пользователя
func (r *Repository) UpdateUser(user *ds.User) error {
    return r.db.Save(user).Error
}

// ==================== ЗАЯВКИ ====================

// GetOrders - получение списка заявок с фильтрацией (исключая удалённые и черновики)
func (r *Repository) GetOrders(status string, dateFrom, dateTo *time.Time) ([]ds.Order, error) {
    var orders []ds.Order
    
    query := r.db.Preload("Creator").Preload("Moderator").
        Where("deleted_at IS NULL AND status != ?", ds.StatusDraft)
    
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    if dateFrom != nil {
        query = query.Where("formed_at >= ?", *dateFrom)
    }
    
    if dateTo != nil {
        query = query.Where("formed_at <= ?", *dateTo)
    }
    
    err := query.Order("created_at DESC").Find(&orders).Error
    return orders, err
}

// GetOrder - получение заявки по ID с услугами
func (r *Repository) GetOrder(id int) (ds.Order, error) {
    var order ds.Order
    err := r.db.Preload("Services.Service").Preload("Creator").Preload("Moderator").
        Where("id = ? AND deleted_at IS NULL", id).First(&order).Error
    if err != nil {
        return ds.Order{}, fmt.Errorf("заявка не найдена")
    }
    return order, nil
}

// GetDraftOrder - получение черновика заявки пользователя
func (r *Repository) GetDraftOrder(creatorID int) (ds.Order, error) {
    var order ds.Order
    err := r.db.Preload("Services.Service").
        Where("creator_id = ? AND status = ? AND deleted_at IS NULL", creatorID, ds.StatusDraft).
        First(&order).Error
    if err != nil {
        return ds.Order{}, fmt.Errorf("черновик не найден")
    }
    return order, nil
}

// CreateDraftOrder - создание черновика заявки
func (r *Repository) CreateDraftOrder(creatorID int) (ds.Order, error) {
    order := ds.Order{
        CreatorID: creatorID,
        Status:    ds.StatusDraft,
        IsDraft:   true,
    }
    err := r.db.Create(&order).Error
    return order, err
}

// UpdateOrder - обновление заявки
func (r *Repository) UpdateOrder(order *ds.Order) error {
    return r.db.Save(order).Error
}

// FormOrder - формирование заявки создателем (проверка обязательных полей)
func (r *Repository) FormOrder(orderID int, fromCity, toCity string, weight, length, width, height float64) error {
    var order ds.Order
    err := r.db.Preload("Services").Where("id = ?", orderID).First(&order).Error
    if err != nil {
        return fmt.Errorf("заявка не найдена")
    }
    
    // Проверяем, что заявка в статусе draft
    if order.Status != ds.StatusDraft {
        return fmt.Errorf("можно формировать только черновики")
    }
    
    // Проверяем обязательные поля
    if fromCity == "" || toCity == "" || weight <= 0 || length <= 0 || width <= 0 || height <= 0 {
        return fmt.Errorf("не заполнены обязательные поля: города и параметры груза")
    }
    
    if len(order.Services) == 0 {
        return fmt.Errorf("в заявке нет услуг")
    }
    
    // Обновляем заявку
    now := time.Now()
    order.FromCity = fromCity
    order.ToCity = toCity
    order.Weight = weight
    order.Length = length
    order.Width = width
    order.Height = height
    order.Status = ds.StatusFormed
    order.FormedAt = &now
    order.IsDraft = false
    
    return r.db.Save(&order).Error
}

// CompleteOrder - завершение/отклонение заявки модератором
func (r *Repository) CompleteOrder(orderID int, status string, moderatorID int) error {
    if status != ds.StatusCompleted && status != ds.StatusRejected {
        return fmt.Errorf("неверный статус для завершения")
    }
    
    var order ds.Order
    err := r.db.Preload("Services.Service").Where("id = ?", orderID).First(&order).Error
    if err != nil {
        return fmt.Errorf("заявка не найдена")
    }
    
    if order.Status != ds.StatusFormed {
        return fmt.Errorf("можно завершать только сформированные заявки")
    }
    
    // Рассчитываем стоимость и сроки при завершении
    if status == ds.StatusCompleted {
        calc := calculator.NewDeliveryCalculator()
        totalCost := 0.0
        maxDays := 0
        
        for _, orderService := range order.Services {
            res := calc.CalculateDelivery(orderService.Service, order.FromCity, order.ToCity, 
                order.Length, order.Width, order.Height, order.Weight)
            if res.IsValid {
                totalCost += res.TotalCost
                if res.DeliveryDays > maxDays {
                    maxDays = res.DeliveryDays
                }
            }
        }
        
        order.TotalCost = totalCost
        order.TotalDays = maxDays
    }
    
    now := time.Now()
    order.Status = status
    order.ModeratorID = &moderatorID
    order.CompletedAt = &now
    
    return r.db.Save(&order).Error
}

// DeleteOrder - удаление заявки (мягкое удаление)

func (r *Repository) DeleteOrder(orderID int) error {
    // Каскадное удаление автоматически удалит связанные записи в order_services
    return r.db.Where("id = ?", orderID).Delete(&ds.Order{}).Error
}

// GetCartIcon - получение иконки корзины (количество услуг в черновике)
func (r *Repository) GetCartIcon(creatorID int) (int, int, error) {
    var order ds.Order
    err := r.db.Preload("Services").Where("creator_id = ? AND status = ? AND deleted_at IS NULL", 
        creatorID, ds.StatusDraft).First(&order).Error
    if err != nil {
        // Создаём черновик если нет
        order, err = r.CreateDraftOrder(creatorID)
        if err != nil {
            return 0, 0, err
        }
    }
    
    count := len(order.Services)
    return order.ID, count, nil
}

// ==================== М-М ЗАЯВКА-УСЛУГА ====================

// AddServiceToOrder - добавление услуги в заявку-черновик
func (r *Repository) AddServiceToOrder(orderID, serviceID int) error {
    // Проверяем что заявка - черновик
    var order ds.Order
    err := r.db.Where("id = ? AND status = ?", orderID, ds.StatusDraft).First(&order).Error
    if err != nil {
        return fmt.Errorf("заявка не найдена или не является черновиком")
    }
    
    // Проверяем услугу
    _, err = r.GetService(serviceID)
    if err != nil {
        return fmt.Errorf("услуга не найдена")
    }
    
    // Проверяем не добавлена ли уже
    var existing ds.OrderService
    err = r.db.Where("order_id = ? AND service_id = ?", orderID, serviceID).First(&existing).Error
    if err == nil {
        // Увеличиваем количество
        existing.Quantity++
        return r.db.Save(&existing).Error
    }
    
    // Добавляем новую
    orderService := ds.OrderService{
        OrderID:   orderID,
        ServiceID: serviceID,
        Quantity:  1,
    }
    return r.db.Create(&orderService).Error
}

// RemoveServiceFromOrder - удаление услуги из заявки
func (r *Repository) RemoveServiceFromOrder(orderID, serviceID int) error {
    var orderService ds.OrderService
    err := r.db.Where("order_id = ? AND service_id = ?", orderID, serviceID).First(&orderService).Error
    if err != nil {
        return fmt.Errorf("услуга не найдена в заявке")
    }
    
    return r.db.Delete(&orderService).Error
}

// UpdateOrderService - обновление количества/порядка в м-м
func (r *Repository) UpdateOrderService(orderID, serviceID int, quantity, orderNum int, comment string) error {
    var orderService ds.OrderService
    err := r.db.Where("order_id = ? AND service_id = ?", orderID, serviceID).First(&orderService).Error
    if err != nil {
        return fmt.Errorf("услуга не найдена в заявке")
    }
    
    orderService.Quantity = quantity
    orderService.Order = orderNum
    orderService.Comment = comment
    
    return r.db.Save(&orderService).Error
}


// AddToCart - добавление услуги в корзину
func (r *Repository) ensureDraftOrder(sessionID string) (int, error) {
    var order ds.Order
    if err := r.db.Where("session_id = ? AND is_draft = ? AND deleted_at IS NULL", sessionID, true).First(&order).Error; err != nil {
        // создаём с системным создателем
        order = ds.Order{
            SessionID: sessionID, 
            IsDraft: true,
            CreatorID: ds.GetCreatorID(),
            Status: ds.StatusDraft,
        }
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

    // upsert в order_services
    // используем нативное подключение для ON CONFLICT
    sqlDB, err := r.db.DB(); if err != nil { return err }
    _, err = sqlDB.Exec(`
        INSERT INTO order_services(order_id, service_id, quantity)
        VALUES ($1, $2, 1)
        ON CONFLICT (order_id, service_id)
        DO UPDATE SET quantity = order_services.quantity + 1
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
    err = sqlDB.QueryRow(`SELECT quantity FROM order_services WHERE order_id=$1 AND service_id=$2`, orderID, serviceID).Scan(&qty)
    if err == sql.ErrNoRows { return fmt.Errorf("услуга не найдена в корзине") }
    if err != nil { return err }

    if qty > 1 {
        _, err = sqlDB.Exec(`UPDATE order_services SET quantity = quantity - 1 WHERE order_id=$1 AND service_id=$2`, orderID, serviceID)
    } else {
        _, err = sqlDB.Exec(`DELETE FROM order_services WHERE order_id=$1 AND service_id=$2`, orderID, serviceID)
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
    _ = sqlDB.QueryRow(`SELECT COALESCE(SUM(quantity),0) FROM order_services WHERE order_id=$1`, orderID).Scan(&count)
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
