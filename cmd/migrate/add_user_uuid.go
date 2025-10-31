package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

func main() {
	// Подключение к БД
	dsn := "host=localhost port=5433 user=myuser password=mypassword dbname=mydb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Добавляем колонку UUID если её нет
	err = db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS uuid VARCHAR(36)").Error
	if err != nil {
		log.Fatal("Failed to add uuid column:", err)
	}

	// Добавляем уникальный индекс для UUID
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_uuid ON users(uuid)").Error
	if err != nil {
		log.Fatal("Failed to create uuid index:", err)
	}

	// Обновляем существующих пользователей, добавляя UUID
	var users []struct {
		ID int
	}
	
	err = db.Raw("SELECT id FROM users WHERE uuid IS NULL OR uuid = ''").Scan(&users).Error
	if err != nil {
		log.Fatal("Failed to get users without UUID:", err)
	}

	for _, user := range users {
		newUUID := uuid.New().String()
		err = db.Exec("UPDATE users SET uuid = ? WHERE id = ?", newUUID, user.ID).Error
		if err != nil {
			log.Printf("Failed to update user %d with UUID: %v", user.ID, err)
		} else {
			fmt.Printf("Updated user %d with UUID: %s\n", user.ID, newUUID)
		}
	}

	// Делаем колонку UUID обязательной
	err = db.Exec("ALTER TABLE users ALTER COLUMN uuid SET NOT NULL").Error
	if err != nil {
		log.Fatal("Failed to make uuid column NOT NULL:", err)
	}

	fmt.Println("Migration completed successfully!")
}

