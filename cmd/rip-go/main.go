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
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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

	// Настраиваем CORS для работы с Tauri и веб-версией
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Разрешаем все источники (для Tauri и веб)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Регистрируем статические файлы и шаблоны
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "static")

	// Прокси для MinIO изображений
	r.Any("/lab1/*path", func(c *gin.Context) {
		path := c.Param("path")
		minioURL := fmt.Sprintf("http://localhost:9003/lab1%s", path)
		
		// Создаем запрос к MinIO
		req, err := http.NewRequest(c.Request.Method, minioURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}
		
		// Копируем заголовки
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		
		// Выполняем запрос
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to connect to MinIO"})
			return
		}
		defer resp.Body.Close()
		
		// Копируем заголовки ответа
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}
		
		// Устанавливаем статус код
		c.Status(resp.StatusCode)
		
		// Копируем тело ответа
		io.Copy(c.Writer, resp.Body)
	})

	// Регистрируем маршруты
	registerRoutes(r, handler)

	// Запускаем сервер
	serverAddress := fmt.Sprintf("%s:%d", conf.ServiceHost, conf.ServicePort)
	
	if conf.EnableHTTPS {
		logrus.Infof("Starting HTTPS server on %s", serverAddress)
		
		// Создаем HTTP сервер с TLS конфигурацией
		// Настраиваем для работы с самоподписанными сертификатами
		srv := &http.Server{
			Addr:    serverAddress,
			Handler: r,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				// Разрешаем все cipher suites для совместимости с разными клиентами
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				},
				PreferServerCipherSuites: true,
				// Включаем поддержку всех версий TLS для совместимости
				MaxVersion: tls.VersionTLS13,
			},
			// Увеличиваем таймауты для TLS handshake
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		
		// Запускаем HTTPS сервер
		if err := srv.ListenAndServeTLS(conf.CertFile, conf.KeyFile); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		logrus.Infof("Starting HTTP server on %s", serverAddress)
		if err := r.Run(serverAddress); err != nil {
			logrus.Fatal(err)
		}
	}
	logrus.Info("Application terminated")
}

func registerRoutes(r *gin.Engine, handler *handler.Handler) {
	// Маршруты для четырех страниц
	r.GET("/", handler.GetServices)                    // Главная страница со списком услуг
	r.GET("/service/:id", handler.GetService)          // Страница с подробной информацией об услуге
	r.GET("/order", handler.GetLogisticRequestDetails)           // Страница с деталями заявки
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
	r.POST("/api/submitcargoorder", handler.AuthMiddleware.RequireAuth(), handler.SubmitLogisticRequest) // Отправка заявки на грузоперевозку

	// API маршрут для обновления статуса заказа через курсор
	r.PUT("/api/order/:id/status", handler.UpdateLogisticRequestStatus) // Обновление статуса заказа

    // CRUD JSON для услуг
    r.GET("/api/services", handler.GetAllServicesJSON)     // Список всех услуг с фильтрацией
    r.GET("/api/services/:id", handler.GetServiceJSON)
    r.POST("/api/services", handler.CreateService)
    r.PUT("/api/services/:id", handler.UpdateService)
    r.DELETE("/api/services/:id", handler.DeleteService)
    
    // Алиасы для transport-services (для совместимости с фронтендом)
    r.GET("/api/transport-services", handler.GetAllServicesJSON)
    r.GET("/api/transport-services/:id", handler.GetServiceJSON)
    r.POST("/api/transport-services", handler.CreateService)
    r.PUT("/api/transport-services/:id", handler.UpdateService)
    r.DELETE("/api/transport-services/:id", handler.DeleteService)

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
        ordersGroup.GET("", handler.GetLogisticRequests)
        ordersGroup.GET("/:id", handler.GetLogisticRequest)
        ordersGroup.DELETE("/:id", handler.DeleteLogisticRequest)
    }
    
    // Специфичные маршруты с дополнительными сегментами
    r.PUT("/api/orders/:id/form", handler.AuthMiddleware.RequireAuth(), handler.FormLogisticRequest)
    r.PUT("/api/orders/:id/update", handler.AuthMiddleware.RequireAuth(), handler.UpdateLogisticRequest)
    
    // Маршрут модератора для завершения заявки
    moderatorCompleteGroup := r.Group("/api/orders/:id")
    moderatorCompleteGroup.Use(handler.AuthMiddleware.RequireModerator())
    {
        moderatorCompleteGroup.PUT("/complete", handler.CompleteLogisticRequest)
    }

    // Новые маршруты: логистические заявки (основные), алиасы для совместимости
    logisticGroup := r.Group("/api/logistic-requests")
    logisticGroup.Use(handler.AuthMiddleware.RequireAuth())
    {
        logisticGroup.GET("", handler.GetLogisticRequests)
        logisticGroup.GET("/:id", handler.GetLogisticRequest)
        logisticGroup.DELETE("/:id", handler.DeleteLogisticRequest)
        logisticGroup.PUT("/:id/form", handler.FormLogisticRequest)
        logisticGroup.PUT("/:id/update", handler.UpdateLogisticRequest)
        logisticGroup.DELETE("/:id/services/:service_id", handler.RemoveServiceFromLogisticRequest)
        logisticGroup.PUT("/:id/services/:service_id", handler.UpdateLogisticRequestService)
    }
    // Завершение логистической заявки (модератор)
    moderatorLR := r.Group("/api/logistic-requests/:id")
    moderatorLR.Use(handler.AuthMiddleware.RequireModerator())
    {
        moderatorLR.PUT("/complete", handler.CompleteLogisticRequest)
    }

    // Алиас статуса для логистических заявок
    r.PUT("/api/logistic-requests/:id/status", handler.UpdateLogisticRequestStatus)

    // Эндпоинт оформления логистической заявки (новый алиас)
    r.POST("/api/submit-cargo-logistic-request", handler.AuthMiddleware.RequireAuth(), handler.SubmitLogisticRequest)

    // Иконка корзины доступна без авторизации (для фронта)
    r.GET("/api/cart/icon", handler.GetCartIcon)

    // М-М заявка-услуга
    r.DELETE("/api/orders/:id/services/:service_id", handler.RemoveServiceFromLogisticRequest)
    r.PUT("/api/orders/:id/services/:service_id", handler.UpdateLogisticRequestService)

    // Swagger документация
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}