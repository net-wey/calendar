package main

import (
	"fmt"
	"log"

	"goproject/internal/storage/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "feedbox"
)

func main() {
	// Формируем строку подключения к БД
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Инициализируем подключение к БД
	db, err := postgres.New(psqlInfo)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Используем db для дальнейшей работы
	_ = db // временно, чтобы не было ошибки неиспользуемой переменной

}
