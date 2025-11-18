package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "rip-go-app/internal/app/ds"
    "rip-go-app/internal/app/repository"
    "rip-go-app/internal/app/calculator"
    "rip-go-app/internal/app/service"
    "rip-go-app/internal/app/middleware"
    "net/http"
    "strconv"
    "strings"
    "time"
    "golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Repository   *repository.Repository
	AuthService  *service.AuthService
	AuthMiddleware *middleware.AuthMiddleware
}

func NewHandler(r *repository.Repository, authService *service.AuthService, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		Repository:     r,
		AuthService:    authService,
		AuthMiddleware: authMiddleware,
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
	
	TransportService, err := h.Repository.GetServices(search)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки услуг",
		})
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"TransportService": TransportService,
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
	// Получаем первую сформированную заявку для демонстрации
	LogisticRequest, err := h.Repository.GetOrders("formed", nil, nil)
	if err != nil || len(LogisticRequest) == 0 {
		logrus.Error(err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Ошибка загрузки заявки",
		})
		return
	}

	order := LogisticRequest[0]
	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"order":    order,
		"TransportService": order.Services,
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

// RemoveFromCart - удаление всей заявки (корзины) по ID
func (h *Handler) RemoveFromCart(ctx *gin.Context) {
    // Очищаем текущую гостевую корзину независимо от order_id в URL
    h.Repository.ClearCart()

    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "count":   0,
        "message": "Корзина очищена",
    })
}

// GetCart - получение корзины
func (h *Handler) GetCart(ctx *gin.Context) {
    cart, err := h.Repository.GetCart()
	if err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to get cart")
		return
	}

    TransportService, err := h.Repository.GetCartServices()
	if err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to get TransportService in cart")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"cart":     cart,
		"TransportService": TransportService,
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

// FormOrder - формирование заявки создателем (дата формирования)
func (h *Handler) FormOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fail(ctx, http.StatusBadRequest, "invalid order id")
		return
	}

	var request struct {
		FromCity string  `json:"from_city"`
		ToCity   string  `json:"to_city"`
		Weight   float64 `json:"weight"`
		Length   float64 `json:"length"`
		Width    float64 `json:"width"`
		Height   float64 `json:"height"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		fail(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	// Проверяем обязательные поля
	if request.FromCity == "" || request.ToCity == "" {
		fail(ctx, http.StatusBadRequest, "from_city and to_city are required")
		return
	}

	if request.Weight <= 0 || request.Length <= 0 || request.Width <= 0 || request.Height <= 0 {
		fail(ctx, http.StatusBadRequest, "weight, length, width, height must be greater than 0")
		return
	}

	err = h.Repository.FormOrder(id, request.FromCity, request.ToCity, request.Weight, request.Length, request.Width, request.Height)
	if err != nil {
		fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Заявка успешно сформирована",
	})
}

// SubmitOrder - отправка логистической заявки на грузоперевозку
// @Summary Submit cargo logistic request
// @Description Submit a new cargo transportation logistic request
// @Tags logistic-requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Logistic request data with TransportService"
// @Success 201 {object} map[string]interface{} "Logistic request submitted successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/submit-cargo-logistic-request [post]
func (h *Handler) SubmitOrder(ctx *gin.Context) {
	userUUID, exists := middleware.GetUserUUID(ctx)
	if !exists {
		fail(ctx, http.StatusUnauthorized, "authentication required")
		return
	}

	// Получаем пользователя для creatorID
	user, err := h.Repository.GetUserByUUID(userUUID)
	if err != nil {
		fail(ctx, http.StatusInternalServerError, "failed to get user")
		return
	}

	var request struct {
		Services []struct {
			ServiceID int     `json:"service_id"`
			FromCity  string  `json:"from_city"`
			ToCity    string  `json:"to_city"`
			Length    float64 `json:"length"`
			Width     float64 `json:"width"`
			Height    float64 `json:"height"`
			Weight    float64 `json:"weight"`
		} `json:"TransportService"`
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

    orderID, err := h.Repository.CreateCargoOrder(items, user.ID)
    if err != nil {
        // Ошибки валидации калькулятора и пр. вернём как 400
        fail(ctx, http.StatusBadRequest, err.Error())
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{
        "status":   "success",
        "message":  "Заявка на грузоперевозку успешно оформлена",
        "order_id": orderID,
        "creator_id": user.ID,
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
	TransportService, err := h.Repository.GetServices(searchQuery)
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
		for _, service := range TransportService {
			if strings.Contains(strings.ToLower(service.Name), strings.ToLower(transportType)) {
				filtered = append(filtered, service)
			}
		}
		TransportService = filtered
	}

	// Возвращаем результат в зависимости от типа запроса
	if ctx.GetHeader("Content-Type") == "application/json" {
        ctx.JSON(http.StatusOK, gin.H{
            "status": "ok",
            "transports": TransportService,
            "count": len(TransportService),
        })
	} else {
		// Возвращаем HTML страницу с результатами
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"TransportService": TransportService,
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

// GetAllServicesJSON - получение всех услуг в JSON с фильтрацией
// Поддерживает фильтрацию по search, minPrice, maxPrice, dateFrom, dateTo
func (h *Handler) GetAllServicesJSON(ctx *gin.Context) {
    // Получаем параметры запроса из URL
    search := ctx.Query("search")
    
    // Обработка minPrice
    var minPrice *float64
    if minPriceStr := ctx.Query("minPrice"); minPriceStr != "" {
        if parsed, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
            minPrice = &parsed
        }
    }
    
    // Обработка maxPrice
    var maxPrice *float64
    if maxPriceStr := ctx.Query("maxPrice"); maxPriceStr != "" {
        if parsed, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
            maxPrice = &parsed
        }
    }
    
    // Обработка dateFrom
    var dateFrom *time.Time
    if dateFromStr := ctx.Query("dateFrom"); dateFromStr != "" {
        if parsed, err := time.Parse("2006-01-02", dateFromStr); err == nil {
            dateFrom = &parsed
        }
    }
    
    // Обработка dateTo
    var dateTo *time.Time
    if dateToStr := ctx.Query("dateTo"); dateToStr != "" {
        if parsed, err := time.Parse("2006-01-02", dateToStr); err == nil {
            dateTo = &parsed
        }
    }
    
    // Получаем отфильтрованные услуги из репозитория
    TransportService, err := h.Repository.GetServicesWithFilters(search, minPrice, maxPrice, dateFrom, dateTo)
    if err != nil {
        logrus.Error("Error getting TransportService:", err)
        fail(ctx, http.StatusInternalServerError, "failed to get TransportService")
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "TransportService": TransportService})
}

// ==================== ПОЛЬЗОВАТЕЛИ ====================

// RegisterUser - регистрация пользователя
// @Summary Register new user
// @Description Register a new user with login, email, password and other details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "Registration data"
// @Success 201 {object} service.AuthResponse "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 409 {object} map[string]string "User already exists"
// @Router /sign_up [post]
func (h *Handler) RegisterUser(ctx *gin.Context) {
    var req service.RegisterRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
        return
    }

    // Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to hash password")
        return
    }

    // Заменяем пароль на хеш
    req.Password = string(hashedPassword)

    response, err := h.AuthService.Register(req)
    if err != nil {
        if err.Error() == "user with this login already exists" {
            fail(ctx, http.StatusConflict, "user with this login already exists")
            return
        }
        fail(ctx, http.StatusInternalServerError, err.Error())
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{
        "status":         "success",
        "message":        "User registered successfully",
        "access_token":   response.AccessToken,
        "refresh_token":  response.RefreshToken,
        "user":           response.User,
        "expires_at":     response.ExpiresAt,
    })
}

// GetUserProfile - получение профиля пользователя
func (h *Handler) GetUserProfile(ctx *gin.Context) {
    userID := ds.GetCreatorID() // используем системного пользователя
    
    user, err := h.Repository.GetUser(userID)
    if err != nil {
        fail(ctx, http.StatusNotFound, "user not found")
        return
    }

    user.Password = ""
    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "user": user})
}

// UpdateUserProfile - обновление профиля пользователя
func (h *Handler) UpdateUserProfile(ctx *gin.Context) {
    userID := ds.GetCreatorID()
    
    var req struct {
        Name  string `json:"name"`
        Phone string `json:"phone"`
        Email string `json:"email" binding:"email"`
    }

    if err := ctx.ShouldBindJSON(&req); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
        return
    }

    user, err := h.Repository.GetUser(userID)
    if err != nil {
        fail(ctx, http.StatusNotFound, "user not found")
        return
    }

    if req.Name != "" {
        user.Name = req.Name
    }
    if req.Phone != "" {
        user.Phone = req.Phone
    }
    if req.Email != "" {
        user.Email = req.Email
    }

    if err := h.Repository.UpdateUser(&user); err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to update user")
        return
    }

    user.Password = ""
    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "user": user})
}

// LoginUser - аутентификация
// @Summary User login
// @Description Authenticate user with login and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "Login credentials"
// @Success 200 {object} service.AuthResponse "Login successful"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /login [post]
func (h *Handler) LoginUser(ctx *gin.Context) {
    var req service.LoginRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
        return
    }

    // Получаем пользователя
    user, err := h.Repository.GetUserByLogin(req.Login)
    if err != nil {
        fail(ctx, http.StatusUnauthorized, "invalid credentials")
        return
    }

    // Проверяем пароль
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        fail(ctx, http.StatusUnauthorized, "invalid credentials")
        return
    }

    // Используем сервис авторизации для входа
    response, err := h.AuthService.Login(req, user.Password)
    if err != nil {
        fail(ctx, http.StatusUnauthorized, err.Error())
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// LogoutUser - деавторизация
// @Summary User logout
// @Description Logout user and invalidate tokens
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Logout successful"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /logout [post]
func (h *Handler) LogoutUser(ctx *gin.Context) {
    userUUID, exists := middleware.GetUserUUID(ctx)
    if !exists {
        fail(ctx, http.StatusUnauthorized, "user not authenticated")
        return
    }

    // Извлекаем токен из заголовка
    authHeader := ctx.GetHeader("Authorization")
    if authHeader == "" {
        fail(ctx, http.StatusUnauthorized, "authorization header required")
        return
    }

    token := strings.TrimPrefix(authHeader, "Bearer ")
    if token == authHeader {
        fail(ctx, http.StatusUnauthorized, "invalid authorization header format")
        return
    }

    err := h.AuthService.Logout(userUUID, token)
    if err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to logout")
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Logged out successfully",
    })
}

// RefreshToken - обновление токенов
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} service.AuthResponse "Tokens refreshed successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid refresh token"
// @Router /refresh [post]
func (h *Handler) RefreshToken(ctx *gin.Context) {
    var req struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }

    if err := ctx.ShouldBindJSON(&req); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
        return
    }

    response, err := h.AuthService.RefreshTokens(req.RefreshToken)
    if err != nil {
        fail(ctx, http.StatusUnauthorized, err.Error())
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status":         "success",
        "message":        "Tokens refreshed successfully",
        "access_token":   response.AccessToken,
        "refresh_token":  response.RefreshToken,
        "user":           response.User,
        "expires_at":     response.ExpiresAt,
    })
}

// ==================== ЛОГИСТИЧЕСКИЕ ЗАЯВКИ ====================

// GetOrders - получение списка логистических заявок с фильтрацией
// @Summary Get logistic requests list
// @Description Get logistic requests list with filtering by status and date range
// @Tags logistic-requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Logistic request status filter"
// @Param date_from query string false "Date from (YYYY-MM-DD)"
// @Param date_to query string false "Date to (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Logistic requests retrieved successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/logistic-requests [get]
func (h *Handler) GetOrders(ctx *gin.Context) {
    userUUID, exists := middleware.GetUserUUID(ctx)
    if !exists {
        fail(ctx, http.StatusUnauthorized, "authentication required")
        return
    }

    userRole, _ := middleware.GetUserRole(ctx)
    status := ctx.Query("status")
    dateFromStr := ctx.Query("date_from")
    dateToStr := ctx.Query("date_to")

    var dateFrom, dateTo *time.Time
    if dateFromStr != "" {
        if t, err := time.Parse("2006-01-02", dateFromStr); err == nil {
            dateFrom = &t
        }
    }
    if dateToStr != "" {
        if t, err := time.Parse("2006-01-02", dateToStr); err == nil {
            dateTo = &t
        }
    }

    LogisticRequest, err := h.Repository.GetOrders(status, dateFrom, dateTo)
    if err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to get LogisticRequest")
        return
    }

    // Фильтрация по ролям
    if userRole == ds.RoleBuyer {
        // Buyer видит только свои заявки
        user, err := h.Repository.GetUserByUUID(userUUID)
        if err != nil {
            fail(ctx, http.StatusInternalServerError, "failed to get user")
            return
        }
        
        var userOrders []ds.Order
        for _, order := range LogisticRequest {
            if order.CreatorID == user.ID {
                userOrders = append(userOrders, order)
            }
        }
        LogisticRequest = userOrders
    }
    // Manager и Admin видят все заявки

    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "LogisticRequest": LogisticRequest})
}

// GetOrder - получение заявки по ID
func (h *Handler) GetOrder(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid order id")
        return
    }

    order, err := h.Repository.GetOrder(id)
    if err != nil {
        fail(ctx, http.StatusNotFound, "order not found")
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "order": order})
}

// UpdateOrder - обновление заявки
func (h *Handler) UpdateOrder(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid order id")
        return
    }

    var req struct {
        FromCity string  `json:"from_city"`
        ToCity   string  `json:"to_city"`
        Weight   float64 `json:"weight"`
        Length   float64 `json:"length"`
        Width    float64 `json:"width"`
        Height   float64 `json:"height"`
    }

    if err := ctx.ShouldBindJSON(&req); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
        return
    }

    order, err := h.Repository.GetOrder(id)
    if err != nil {
        fail(ctx, http.StatusNotFound, "order not found")
        return
    }

    if order.Status != ds.StatusDraft {
        fail(ctx, http.StatusBadRequest, "can only update draft LogisticRequest")
        return
    }

    if req.FromCity != "" {
        order.FromCity = req.FromCity
    }
    if req.ToCity != "" {
        order.ToCity = req.ToCity
    }
    if req.Weight > 0 {
        order.Weight = req.Weight
    }
    if req.Length > 0 {
        order.Length = req.Length
    }
    if req.Width > 0 {
        order.Width = req.Width
    }
    if req.Height > 0 {
        order.Height = req.Height
    }

    if err := h.Repository.UpdateOrder(&order); err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to update order")
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "order": order})
}


// CompleteOrder - завершение/отклонение логистической заявки модератором
// @Summary Complete or reject logistic request
// @Description Complete or reject logistic request by moderator
// @Tags logistic-requests
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Logistic request ID"
// @Param request body map[string]string true "Logistic request status (completed/rejected)"
// @Success 200 {object} map[string]string "Logistic request completed successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/logistic-requests/{id}/complete [put]
func (h *Handler) CompleteOrder(ctx *gin.Context) {
    // Middleware уже проверил авторизацию и роль модератора
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid order id")
        return
    }

    var req struct {
        Status string `json:"status" binding:"required"`
    }

    if err := ctx.ShouldBindJSON(&req); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
        return
    }

    // Валидация статуса
    if req.Status != ds.StatusCompleted && req.Status != ds.StatusRejected {
        fail(ctx, http.StatusBadRequest, "invalid status. allowed: completed, rejected")
        return
    }

    // Получаем пользователя для moderatorID
    userUUID, _ := middleware.GetUserUUID(ctx)
    user, err := h.Repository.GetUserByUUID(userUUID)
    if err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to get user")
        return
    }

    err = h.Repository.CompleteOrder(id, req.Status, user.ID)
    if err != nil {
        fail(ctx, http.StatusBadRequest, err.Error())
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Order completed successfully",
    })
}

// DeleteOrder - удаление заявки
func (h *Handler) DeleteOrder(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid order id")
        return
    }

    err = h.Repository.DeleteOrder(id)
    if err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to delete order")
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "order deleted successfully"})
}

// GetCartIcon - получение иконки корзины
func (h *Handler) GetCartIcon(ctx *gin.Context) {
    creatorID := ds.GetCreatorID()
    orderID, count, err := h.Repository.GetCartIcon(creatorID)
    if err != nil {
        fail(ctx, http.StatusInternalServerError, "failed to get cart")
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status":   "ok",
        "order_id": orderID,
        "count":    count,
    })
}

// ==================== М-М ЗАЯВКА-УСЛУГА ====================

// AddServiceToOrder - добавление услуги в заявку
func (h *Handler) AddServiceToOrder(ctx *gin.Context) {
    var req struct {
        OrderID   int `json:"order_id" binding:"required"`
        ServiceID int `json:"service_id" binding:"required"`
    }

    if err := ctx.ShouldBindJSON(&req); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
        return
    }

    err := h.Repository.AddServiceToOrder(req.OrderID, req.ServiceID)
    if err != nil {
        fail(ctx, http.StatusBadRequest, err.Error())
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "service added to order"})
}

// RemoveServiceFromOrder - удаление услуги из заявки
func (h *Handler) RemoveServiceFromOrder(ctx *gin.Context) {
    orderIDStr := ctx.Param("id")
    serviceIDStr := ctx.Param("service_id")
    
    orderID, err := strconv.Atoi(orderIDStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid order id")
        return
    }
    
    serviceID, err := strconv.Atoi(serviceIDStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid service id")
        return
    }

    err = h.Repository.RemoveServiceFromOrder(orderID, serviceID)
    if err != nil {
        fail(ctx, http.StatusBadRequest, err.Error())
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "service removed from order"})
}

// UpdateOrderService - обновление м-м
func (h *Handler) UpdateOrderService(ctx *gin.Context) {
    orderIDStr := ctx.Param("id")
    serviceIDStr := ctx.Param("service_id")
    
    orderID, err := strconv.Atoi(orderIDStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid order id")
        return
    }
    
    serviceID, err := strconv.Atoi(serviceIDStr)
    if err != nil {
        fail(ctx, http.StatusBadRequest, "invalid service id")
        return
    }

    var req struct {
        Quantity int    `json:"quantity" binding:"required,min=1"`
        Order    int    `json:"order"`
        Comment  string `json:"comment"`
    }

    if err := ctx.ShouldBindJSON(&req); err != nil {
        fail(ctx, http.StatusBadRequest, "invalid request body")
        return
    }

    err = h.Repository.UpdateOrderService(orderID, serviceID, req.Quantity, req.Order, req.Comment)
    if err != nil {
        fail(ctx, http.StatusBadRequest, err.Error())
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "order service updated"})
}
