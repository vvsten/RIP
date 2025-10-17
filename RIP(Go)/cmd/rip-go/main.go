// @title RIP Go API
// @version 1.0
// @description API for cargo transportation service
// @host localhost:8083
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for authentication. Format: 'Bearer <token>'
package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rip-go-app/internal/app/config"
	"rip-go-app/internal/app/dsn"
	"rip-go-app/internal/app/handler"
	"rip-go-app/internal/app/repository"
	"rip-go-app/internal/app/auth"
	"rip-go-app/internal/app/service"
	"rip-go-app/internal/app/middleware"
	
	// Swagger imports
	_ "rip-go-app/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func main() {
	logrus.Info("Application start up")

	// Загружаем конфигурацию
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	// Получаем строку подключения к БД
	postgresString := dsn.FromEnv()
	fmt.Println("Connecting to database with DSN:", postgresString)

	// Инициализируем репозиторий
	repo, err := repository.New(postgresString)
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}

	// Инициализируем JWT сервис
	jwtService := auth.NewJWTService(
		conf.JWTSecret,
		conf.JWTAccessTokenExpire,
		conf.JWTRefreshTokenExpire,
	)

	// Инициализируем Redis сервис
	redisService := auth.NewRedisService(
		conf.RedisHost,
		conf.RedisPort,
		conf.RedisPassword,
		conf.RedisDB,
	)

	// Проверяем соединение с Redis
	if err := redisService.Ping(); err != nil {
		logrus.Warnf("Redis connection failed: %v", err)
	} else {
		logrus.Info("Redis connected successfully")
	}

	// Инициализируем сервис авторизации
	authService := service.NewAuthService(repo, jwtService, redisService)

	// Инициализируем middleware авторизации
	authMiddleware := middleware.NewAuthMiddleware(jwtService, redisService)

	// Создаем хендлер
	handler := handler.NewHandler(repo, authService, authMiddleware)

	// Создаем роутер
	r := gin.Default()

	// Регистрируем статические файлы и шаблоны
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "static")

	// Регистрируем маршруты
	registerRoutes(r, handler)

	// Запускаем сервер
	serverAddress := fmt.Sprintf("%s:%d", conf.ServiceHost, conf.ServicePort)
	logrus.Infof("Starting server on %s", serverAddress)
	if err := r.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Application terminated")
}

func registerRoutes(r *gin.Engine, handler *handler.Handler) {
	// Маршруты для четырех страниц
	r.GET("/", handler.GetServices)                    // Главная страница со списком услуг
	r.GET("/service/:id", handler.GetService)          // Страница с подробной информацией об услуге
	r.GET("/order", handler.GetOrderDetails)           // Страница с деталями заявки
	r.GET("/calculator", handler.GetCalculator)        // Страница калькулятора
	r.POST("/calculator", handler.PostCalculator)      // Обработка формы калькулятора

	// API маршруты для корзины
	r.POST("/api/cart/add/:id", handler.AddToCart)     // Добавление услуги в корзину
	r.DELETE("/api/cart/:id", handler.RemoveFromCart)  // Удаление всей заявки (корзины) по ID
	r.DELETE("/api/cart/remove/:id", handler.RemoveFromCart) // Старый маршрут для совместимости
	r.GET("/api/cart", handler.GetCart)                // Получение корзины
	r.GET("/api/cart/count", handler.GetCartCount)     // Получение количества в корзине

	// API маршруты для калькулятора (переименованы под грузоперевозки)
	r.POST("/api/searchtrans", handler.SearchTransport) // Поиск транспорта
	r.POST("/api/calculatecargo", handler.CalculateService) // Расчет стоимости грузоперевозки
	r.POST("/api/submitcargoorder", handler.AuthMiddleware.RequireAuth(), handler.SubmitOrder) // Отправка заявки на грузоперевозку

	// API маршрут для обновления статуса заказа через курсор
	r.PUT("/api/order/:id/status", handler.UpdateOrderStatus) // Обновление статуса заказа

    // CRUD JSON для услуг
    r.GET("/api/services/:id", handler.GetServiceJSON)
    r.POST("/api/services", handler.CreateService)
    r.PUT("/api/services/:id", handler.UpdateService)
    r.DELETE("/api/services/:id", handler.DeleteService)

    // Авторизация
    r.POST("/sign_up", handler.RegisterUser)
    r.POST("/login", handler.LoginUser)
    r.POST("/logout", handler.AuthMiddleware.RequireAuth(), handler.LogoutUser)
    r.POST("/refresh", handler.RefreshToken)

    // Пользователи (требуют авторизации)
    authGroup := r.Group("/api/users")
    authGroup.Use(handler.AuthMiddleware.RequireAuth())
    {
        authGroup.GET("/profile", handler.GetUserProfile)
        authGroup.PUT("/profile", handler.UpdateUserProfile)
    }

    // Заявки (требуют авторизации)
    ordersGroup := r.Group("/api/orders")
    ordersGroup.Use(handler.AuthMiddleware.RequireAuth())
    {
        ordersGroup.GET("", handler.GetOrders)
        ordersGroup.GET("/:id", handler.GetOrder)
        ordersGroup.DELETE("/:id", handler.DeleteOrder)
    }
    
    // Специфичные маршруты с дополнительными сегментами
    r.PUT("/api/orders/:id/form", handler.AuthMiddleware.RequireAuth(), handler.FormOrder)
    r.PUT("/api/orders/:id/update", handler.AuthMiddleware.RequireAuth(), handler.UpdateOrder)
    
    // Маршрут модератора для завершения заявки
    moderatorCompleteGroup := r.Group("/api/orders/:id")
    moderatorCompleteGroup.Use(handler.AuthMiddleware.RequireModerator())
    {
        moderatorCompleteGroup.PUT("/complete", handler.CompleteOrder)
    }

    // Корзина (требует авторизации)
    cartGroup := r.Group("/api/cart")
    cartGroup.Use(handler.AuthMiddleware.RequireAuth())
    {
        cartGroup.GET("/icon", handler.GetCartIcon)
    }

    // М-М заявка-услуга
    r.DELETE("/api/orders/:id/services/:service_id", handler.RemoveServiceFromOrder)
    r.PUT("/api/orders/:id/services/:service_id", handler.UpdateOrderService)

    // Swagger документация
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}