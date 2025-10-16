package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "rip-go-app/internal/app/ds"
    "rip-go-app/internal/app/repository"
    "rip-go-app/internal/app/calculator"
    "net/http"
    "strconv"
    "strings"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// helper для единых ошибок
func fail(ctx *gin.Context, code int, message string) {
    ctx.JSON(code, gin.H{
        "status":  "fail",
        "message": message,
    })
}

// GetServices - главная страница со списком услуг
func (h *Handler) GetServices(ctx *gin.Context) {
	search := ctx.Query("search") // получаем параметр поиска из URL
	
	services, err := h.Repository.GetServices(search)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки услуг",
		})
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"services": services,
		"search":   search, // передаем поисковый запрос для сохранения в поле
	})
}

// GetService - страница с подробной информацией об услуге
func (h *Handler) GetService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Неверный ID услуги",
		})
		return
	}

	service, err := h.Repository.GetService(id)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Услуга не найдена",
		})
		return
	}

	ctx.HTML(http.StatusOK, "service.html", gin.H{
		"service": service,
	})
}

// GetOrderDetails - страница с деталями заявки
func (h *Handler) GetOrderDetails(ctx *gin.Context) {
	order, err := h.Repository.GetOrder()
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки заявки",
		})
		return
	}

	services, err := h.Repository.GetServicesInOrder()
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки услуг в заявке",
		})
		return
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"order":    order,
		"services": services,
	})
}

// GetCalculator - страница калькулятора
func (h *Handler) GetCalculator(ctx *gin.Context) {
	// Получаем услуги из корзины
	cartServices, err := h.Repository.GetCartServices()
	if err != nil {
		logrus.Error(err)
		cartServices = []ds.Service{}
	}

	ctx.HTML(http.StatusOK, "calculator.html", gin.H{
		"FromCity":     "",
		"ToCity":       "",
		"Length":       "",
		"Width":        "",
		"Height":       "",
		"Weight":       "",
		"DeliveryType": "",
		"DeliveryDays": 0,
		"TotalCost":    0,
		"CartServices": cartServices,
	})
}

// PostCalculator - обработка формы калькулятора
func (h *Handler) PostCalculator(ctx *gin.Context) {
	// Получаем данные из формы
	fromCity := ctx.PostForm("from_city")
	toCity := ctx.PostForm("to_city")
	lengthStr := ctx.PostForm("length")
	widthStr := ctx.PostForm("width")
	heightStr := ctx.PostForm("height")
	weightStr := ctx.PostForm("weight")
	deliveryType := ctx.PostForm("delivery_type")

	// Парсим числовые значения
	length, _ := strconv.ParseFloat(lengthStr, 64)
	width, _ := strconv.ParseFloat(widthStr, 64)
	height, _ := strconv.ParseFloat(heightStr, 64)
	weight, _ := strconv.ParseFloat(weightStr, 64)

	// Получаем услугу по типу доставки
	var selectedService ds.Service
	if deliveryType != "" {
		service, err := h.Repository.GetServiceByType(deliveryType)
		if err == nil {
			selectedService = service
		}
	}

	
	deliveryDays, totalCost := selectedService.DeliveryDays + int(weight/1000), selectedService.Price + (length*width*height*50) + (weight*2)

	ctx.HTML(http.StatusOK, "calculator.html", gin.H{
		"FromCity":     fromCity,
		"ToCity":       toCity,
		"Length":       lengthStr,
		"Width":        widthStr,
		"Height":       heightStr,
		"Weight":       weightStr,
		"DeliveryType": deliveryType,
		"DeliveryDays": deliveryDays,
		"TotalCost":    totalCost,
	})
}

// calculateDistance - простая функция расчета расстояния между городами
func calculateDistance(fromCity, toCity string) float64 {
	// Приводим к нижнему регистру для сравнения
	from := strings.ToLower(strings.TrimSpace(fromCity))
	to := strings.ToLower(strings.TrimSpace(toCity))
	
	// Если города одинаковые
	if from == to {
		return 0
	}
	
	// Простая база данных расстояний между основными городами
	distances := map[string]map[string]float64{
		"москва": {
			"санкт-петербург": 635,
			"спб":             635,
			"екатеринбург":    1416,
			"новосибирск":     3354,
			"красноярск":      4205,
			"иркутск":         5152,
			"владивосток":     9100,
			"ростов-на-дону":  1070,
			"сочи":            1360,
			"казань":          820,
			"нижний новгород": 420,
			"самара":          1050,
			"волгоград":       970,
			"воронеж":         520,
			"саратов":         850,
			"пермь":           1380,
			"уфа":             1160,
			"челябинск":       1510,
			"омск":            2550,
			"тюмень":          1720,
		},
		"санкт-петербург": {
			"спб":             0,
			"москва":          635,
			"екатеринбург":    1780,
			"новосибирск":     3720,
			"калининград":     550,
			"мурманск":        1050,
			"архангельск":     1130,
			"петрозаводск":    320,
			"великий новгород": 180,
		},
		"екатеринбург": {
			"москва":          1416,
			"санкт-петербург": 1780,
			"спб":             1780,
			"новосибирск":     1940,
			"челябинск":       200,
			"пермь":           360,
			"тюмень":          320,
			"уфа":             520,
		},
		"новосибирск": {
			"москва":          3354,
			"санкт-петербург": 3720,
			"спб":             3720,
			"екатеринбург":    1940,
			"омск":            650,
			"красноярск":      850,
			"томск":           270,
			"барнаул":         230,
		},
	}
	
	// Ищем расстояние в базе данных
	if cityDistances, exists := distances[from]; exists {
		if distance, found := cityDistances[to]; found {
			return distance
		}
	}
	
	// Если расстояние не найдено, используем примерную оценку
	// Базовое расстояние для неизвестных маршрутов
	return 500.0
}

// AddToCart - добавление услуги в корзину
func (h *Handler) AddToCart(ctx *gin.Context) {
	serviceIDStr := ctx.Param("id")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid service id, must be integer >= 0")
		return
	}

    err = h.Repository.AddToCart(serviceID)
	if err != nil {
        fail(ctx, http.StatusNotFound, err.Error())
		return
	}

	// Возвращаем обновленное количество в корзине
	count := h.Repository.GetCartCount()
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   count,
		"message": "Услуга добавлена в корзину",
	})
}

// RemoveFromCart - удаление услуги из корзины
func (h *Handler) RemoveFromCart(ctx *gin.Context) {
	serviceIDStr := ctx.Param("id")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid service id, must be integer >= 0")
		return
	}

    err = h.Repository.RemoveFromCart(serviceID)
	if err != nil {
        fail(ctx, http.StatusNotFound, err.Error())
		return
	}

	// Возвращаем обновленное количество в корзине
	count := h.Repository.GetCartCount()
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   count,
		"message": "Услуга удалена из корзины",
	})
}

// GetCart - получение корзины
func (h *Handler) GetCart(ctx *gin.Context) {
    cart, err := h.Repository.GetCart()
	if err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to get cart")
		return
	}

    services, err := h.Repository.GetCartServices()
	if err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to get services in cart")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"cart":     cart,
		"services": services,
		"count":    h.Repository.GetCartCount(),
	})
}

// GetCartCount - получение количества услуг в корзине
func (h *Handler) GetCartCount(ctx *gin.Context) {
	count := h.Repository.GetCartCount()
	ctx.JSON(http.StatusOK, gin.H{"count": count})
}

// CalculateService - расчет стоимости грузоперевозки для конкретного типа транспорта
func (h *Handler) CalculateService(ctx *gin.Context) {
	var request struct {
		ServiceID int     `json:"service_id" form:"service_id"`
		FromCity  string  `json:"from_city" form:"from_city"`
		ToCity    string  `json:"to_city" form:"to_city"`
		Length    float64 `json:"length" form:"length"`
		Width     float64 `json:"width" form:"width"`
		Height    float64 `json:"height" form:"height"`
		Weight    float64 `json:"weight" form:"weight"`
	}

	// Пробуем сначала JSON, потом form data
	if err := ctx.ShouldBindJSON(&request); err != nil {
		if err := ctx.ShouldBind(&request); err != nil {
            fail(ctx, http.StatusBadRequest, "invalid request body")
			return
		}
	}

	// Получаем тип транспорта
    service, err := h.Repository.GetService(request.ServiceID)
	if err != nil {
        fail(ctx, http.StatusNotFound, "transport type not found")
		return
	}

    // Используем компонент калькулятора
    calc := calculator.NewDeliveryCalculator()
    res := calc.CalculateDelivery(service, request.FromCity, request.ToCity, request.Length, request.Width, request.Height, request.Weight)

    if !res.IsValid {
        fail(ctx, http.StatusBadRequest, res.ErrorMessage)
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status":        "ok",
        "delivery_days": res.DeliveryDays,
        "total_cost":    res.TotalCost,
        "distance":      res.Distance,
        "volume":        res.Volume,
    })
}

// SubmitOrder - отправка заявки на грузоперевозку
func (h *Handler) SubmitOrder(ctx *gin.Context) {
	var request struct {
		Services []struct {
			ServiceID int     `json:"service_id"`
			FromCity  string  `json:"from_city"`
			ToCity    string  `json:"to_city"`
			Length    float64 `json:"length"`
			Width     float64 `json:"width"`
			Height    float64 `json:"height"`
			Weight    float64 `json:"weight"`
		} `json:"services"`
	}

    if err := ctx.ShouldBindJSON(&request); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

    if len(request.Services) == 0 {
        fail(ctx, http.StatusBadRequest, "no transport types provided")
		return
	}

    // Маппим вход в элементы заказа и сохраняем транзакционно
    items := make([]repository.CargoOrderItem, 0, len(request.Services))
    for _, s := range request.Services {
        items = append(items, repository.CargoOrderItem{
            ServiceID: s.ServiceID,
            FromCity:  s.FromCity,
            ToCity:    s.ToCity,
            Length:    s.Length,
            Width:     s.Width,
            Height:    s.Height,
            Weight:    s.Weight,
        })
    }

    orderID, err := h.Repository.CreateCargoOrder(items)
    if err != nil {
        // Ошибки валидации калькулятора и пр. вернём как 400
        fail(ctx, http.StatusBadRequest, err.Error())
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{
        "status":   "ok",
        "message":  "Заявка на грузоперевозку успешно оформлена",
        "order_id": orderID,
    })
}

// SearchTransport - поиск транспорта (обработка form data)
func (h *Handler) SearchTransport(ctx *gin.Context) {
	// Получаем данные из формы
	searchQuery := ctx.PostForm("search_query")
	transportType := ctx.PostForm("transport_type")
	
	// Если это JSON запрос, обрабатываем по-другому
	if ctx.GetHeader("Content-Type") == "application/json" {
		var request struct {
			SearchQuery   string `json:"search_query"`
			TransportType string `json:"transport_type"`
		}
		
        if err := ctx.ShouldBindJSON(&request); err != nil {
            fail(ctx, http.StatusBadRequest, "invalid request body")
			return
		}
		
		searchQuery = request.SearchQuery
		transportType = request.TransportType
	}
	
	// Поиск транспорта
	services, err := h.Repository.GetServices(searchQuery)
	if err != nil {
        logrus.Error(err)
        if ctx.GetHeader("Content-Type") == "application/json" {
            fail(ctx, http.StatusInternalServerError, "failed to search transports")
        } else {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "Ошибка поиска транспорта",
			})
		}
		return
	}

	// Фильтрация по типу транспорта если указан
	if transportType != "" {
		var filtered []ds.Service
		for _, service := range services {
			if strings.Contains(strings.ToLower(service.Name), strings.ToLower(transportType)) {
				filtered = append(filtered, service)
			}
		}
		services = filtered
	}

	// Возвращаем результат в зависимости от типа запроса
	if ctx.GetHeader("Content-Type") == "application/json" {
        ctx.JSON(http.StatusOK, gin.H{
            "status": "ok",
            "transports": services,
            "count": len(services),
        })
	} else {
		// Возвращаем HTML страницу с результатами
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"services": services,
			"search":   searchQuery,
		})
	}
}

// UpdateOrderStatus - обновление статуса заказа через курсор
func (h *Handler) UpdateOrderStatus(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
    orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid order id, must be integer >= 0")
		return
	}

	// Получаем новый статус из JSON
	var request struct {
		Status string `json:"status" binding:"required"`
	}

    if err := ctx.ShouldBindJSON(&request); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	// Валидация статуса
	validStatuses := []string{"pending", "processing", "shipped", "delivered", "cancelled"}
	isValid := false
	for _, status := range validStatuses {
		if request.Status == status {
			isValid = true
			break
		}
	}

    if !isValid {
        fail(ctx, http.StatusBadRequest, "invalid status. allowed: pending, processing, shipped, delivered, cancelled")
		return
	}

	// Обновляем статус через курсор
    err = h.Repository.UpdateOrderStatusWithCursor(orderID, request.Status)
	if err != nil {
        logrus.Error(err)
        fail(ctx, http.StatusInternalServerError, "failed to update order status")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
        "status": "ok",
        "message": "Статус заказа успешно обновлен",
		"order_id": orderID,
		"new_status": request.Status,
	})
}

// -------------------------
// CRUD JSON для Service
// -------------------------

// CreateService - создание типа транспорта
func (h *Handler) CreateService(ctx *gin.Context) {
    var req ds.Service
    if err := ctx.ShouldBindJSON(&req); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
        return
    }
    if err := h.Repository.CreateService(&req); err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to create service")
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"status": "ok", "service": req})
}

// UpdateService - обновление типа транспорта
func (h *Handler) UpdateService(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid service id")
        return
    }
    var req ds.Service
    if err := ctx.ShouldBindJSON(&req); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
        return
    }
    req.ID = id
    if err := h.Repository.UpdateService(&req); err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to update service")
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "service": req})
}

// DeleteService - удаление типа транспорта
func (h *Handler) DeleteService(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid service id")
        return
    }
    if err := h.Repository.DeleteService(id); err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to delete service")
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetServiceJSON - получение услуги JSON
func (h *Handler) GetServiceJSON(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid service id")
        return
    }
    svc, err := h.Repository.GetService(id)
    if err != nil {
        fail(ctx, http.StatusNotFound, "service not found")
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "service": svc})
}




