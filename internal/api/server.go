package api

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "log"
  "rip-go-app/internal/app/handler"
  "rip-go-app/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("templates/*.html")
	// добавляем статические файлы (CSS, JS, изображения)
	r.Static("/static", "static")

	// Маршруты для четырех страниц
	r.GET("/", handler.GetServices)                    // Главная страница со списком услуг
	r.GET("/service/:id", handler.GetService)          // Страница с подробной информацией об услуге
	r.GET("/order", handler.GetLogisticRequestDetails)           // Страница с деталями заявки
	r.GET("/calculator", handler.GetCalculator)        // Страница калькулятора
	r.POST("/calculator", handler.PostCalculator)      // Обработка формы калькулятора

	// API маршруты для корзины
	r.POST("/api/cart/add/:id", handler.AddToCart)     // Добавление услуги в корзину
	r.DELETE("/api/cart/remove/:id", handler.RemoveFromCart) // Удаление услуги из корзины
	r.GET("/api/cart", handler.GetCart)                // Получение корзины
	r.GET("/api/cart/count", handler.GetCartCount)     // Получение количества в корзине

	// API маршруты для калькулятора (переименованы под грузоперевозки)
	r.POST("/api/searchtrans", handler.SearchTransport) // Поиск транспорта
	r.POST("/api/calculatecargo", handler.CalculateService) // Расчет стоимости грузоперевозки
	r.POST("/api/submitcargoorder", handler.SubmitLogisticRequest) // Отправка заявки на грузоперевозку

    // API маршруты для заявок — старые алиасы
    r.GET("/api/orders", handler.GetLogisticRequests)
    r.PUT("/api/orders/:id/form", handler.FormLogisticRequest)
    r.PUT("/api/orders/:id/complete", handler.CompleteLogisticRequest)
    r.DELETE("/api/orders/:id/services/:service_id", handler.RemoveServiceFromLogisticRequest)
    r.PUT("/api/orders/:id/services/:service_id", handler.UpdateLogisticRequestService)
    r.GET("/api/orders/:id", handler.GetLogisticRequest)
    r.PUT("/api/orders/:id", handler.UpdateLogisticRequest)
    r.DELETE("/api/orders/:id", handler.DeleteLogisticRequest)

    // Новые маршруты для логистических заявок
    r.GET("/api/logistic-requests", handler.GetLogisticRequests)
    r.PUT("/api/logistic-requests/:id/form", handler.FormLogisticRequest)
    r.PUT("/api/logistic-requests/:id/complete", handler.CompleteLogisticRequest)
    r.DELETE("/api/logistic-requests/:id/services/:service_id", handler.RemoveServiceFromLogisticRequest)
    r.PUT("/api/logistic-requests/:id/services/:service_id", handler.UpdateLogisticRequestService)
    r.GET("/api/logistic-requests/:id", handler.GetLogisticRequest)
    r.PUT("/api/logistic-requests/:id", handler.UpdateLogisticRequest)
    r.DELETE("/api/logistic-requests/:id", handler.DeleteLogisticRequest)

	// API маршруты для пользователей
	r.POST("/api/users/register", handler.RegisterUser)   // Регистрация пользователя
	r.POST("/api/users/login", handler.LoginUser)         // Аутентификация
	r.POST("/api/users/logout", handler.LogoutUser)       // Деавторизация
	r.GET("/api/users/profile", handler.GetUserProfile)   // Получение профиля пользователя
	r.PUT("/api/users/profile", handler.UpdateUserProfile) // Обновление профиля пользователя

	r.Run(":8083") // listen and serve on 0.0.0.0:8083
	log.Println("Server down")
}