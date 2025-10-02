package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rip-go-app/internal/app/ds"
	"rip-go-app/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Service{},
		&ds.Order{},
		&ds.OrderService{},
		&ds.Cart{},
		&ds.CartService{},
	)
	if err != nil {
		panic("cant migrate db")
	}

	// Создаем начальные данные
	services := []ds.Service{
		{
			ID:           1,
			Name:         "Фура",
			Description:  "Полуприцеп для перевозки крупногабаритных грузов. Идеально подходит для перевозки мебели, строительных материалов и других тяжелых грузов.",
			Price:        150.0,
			ImageURL:     "http://localhost:9003/lab1/fura.jpg",
			DeliveryDays: 2,
			MaxWeight:    20000.0,
			MaxVolume:    80.0,
		},
		{
			ID:           2,
			Name:         "Малотоннажный грузовик",
			Description:  "Легкий грузовик для перевозки небольших грузов по городу и между городами. Быстрая доставка с возможностью проезда в центр города.",
			Price:        80.0,
			ImageURL:     "http://localhost:9003/lab1/malotonnazhnyi.jpg",
			DeliveryDays: 1,
			MaxWeight:    3000.0,
			MaxVolume:    15.0,
		},
		{
			ID:           3,
			Name:         "Авиаперевозка",
			Description:  "Быстрая доставка грузов авиатранспортом. Подходит для срочных и ценных грузов. Максимальная скорость доставки.",
			Price:        500.0,
			ImageURL:     "http://localhost:9003/lab1/avia.jpg",
			DeliveryDays: 1,
			MaxWeight:    1000.0,
			MaxVolume:    5.0,
		},
		{
			ID:           4,
			Name:         "Поезд",
			Description:  "Железнодорожные перевозки для крупных партий грузов. Экономичный вариант для больших объемов.",
			Price:        120.0,
			ImageURL:     "http://localhost:9003/lab1/poezd.jpg",
			DeliveryDays: 3,
			MaxWeight:    50000.0,
			MaxVolume:    120.0,
		},
		{
			ID:           5,
			Name:         "Корабль",
			Description:  "Морские перевозки для международной доставки. Подходит для крупных партий и контейнерных перевозок.",
			Price:        200.0,
			ImageURL:     "http://localhost:9003/lab1/korabl.jpg",
			DeliveryDays: 7,
			MaxWeight:    100000.0,
			MaxVolume:    500.0,
		},
		{
			ID:           6,
			Name:         "Мультимодальные",
			Description:  "Комбинированные перевозки с использованием нескольких видов транспорта. Оптимальное решение для сложных маршрутов.",
			Price:        300.0,
			ImageURL:     "http://localhost:9003/lab1/multimodal.jpg",
			DeliveryDays: 5,
			MaxWeight:    30000.0,
			MaxVolume:    100.0,
		},
	}

	// Создаем услуги в БД
	for _, service := range services {
		var existingService ds.Service
		err := db.Where("id = ?", service.ID).First(&existingService).Error
		if err != nil {
			// Услуга не существует, создаем
			db.Create(&service)
		}
	}

	// Создаем пример заявки
	var existingOrder ds.Order
	err = db.Where("id = ?", 1).First(&existingOrder).Error
	if err != nil {
		order := ds.Order{
			ID:        1,
			FromCity:  "Москва",
			ToCity:    "Санкт-Петербург",
			Weight:    500.0,
			Length:    2.0,
			Width:     1.5,
			Height:    1.0,
			TotalCost: 230.0,
			TotalDays: 2,
		}
		db.Create(&order)

		// Создаем услуги в заявке
		orderServices := []ds.OrderService{
			{OrderID: 1, ServiceID: 1, Quantity: 1, Comment: "Основная доставка"},
			{OrderID: 1, ServiceID: 2, Quantity: 1, Comment: "Дополнительная услуга"},
		}
		for _, orderService := range orderServices {
			db.Create(&orderService)
		}
	}

	println("Migration completed successfully!")
}
