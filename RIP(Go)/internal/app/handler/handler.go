package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rip-go-app/internal/app/repository"
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
	var selectedService repository.Service
	if deliveryType != "" {
		service, err := h.Repository.GetServiceByType(deliveryType)
		if err == nil {
			selectedService = service
		}
	}

	// Рассчитываем стоимость и сроки
	deliveryDays := 0
	totalCost := 0.0

	if selectedService.ID != 0 {
		// Базовые расчеты
		volume := length * width * height
		
		// Проверяем ограничения
		if weight <= selectedService.MaxWeight && volume <= selectedService.MaxVolume {
			// Простая формула расчета расстояния между городами
			distance := calculateDistance(fromCity, toCity)
			
			// Базовые сроки доставки
			deliveryDays = selectedService.DeliveryDays
			
			// Увеличиваем сроки в зависимости от расстояния
			if distance > 1000 {
				deliveryDays += 1 // +1 день за каждые 1000+ км
			}
			if distance > 2000 {
				deliveryDays += 1 // +1 день за каждые 2000+ км
			}
			
			// Базовая стоимость
			totalCost = selectedService.Price
			
			// Учитываем расстояние (рублей за км)
			distanceCost := distance * 0.5 // 50 копеек за км
			totalCost += distanceCost
			
			// Учитываем вес (рублей за кг)
			weightCost := weight * 0.1 // 10 копеек за кг
			totalCost += weightCost
			
			// Учитываем объем (рублей за м³)
			volumeCost := volume * 10 // 10 рублей за м³
			totalCost += volumeCost
			
			// Минимальная стоимость
			if totalCost < selectedService.Price {
				totalCost = selectedService.Price
			}
		} else {
			// Если груз не подходит под ограничения
			deliveryDays = 0
			totalCost = 0
		}
	}

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
