package main

import (
	"log"
	"net/http"
	"os"
	"sbs/handlers"
	"sbs/logger"
	"sbs/repositories"
	"sbs/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Загружаем переменные окружения
	err := godotenv.Load("settings.env")
	if err != nil {
		log.Println("Warning: .env file not found or couldn't load")
	}

	// Подключаемся к базе
	db, err := repositories.NewPostgresDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Создаём репозиторий, сервис и обработчик
	subRepo := repositories.NewSubscriptionRepository(db)
	subService := services.NewSubscriptionService(subRepo)
	subHandler := handlers.NewSubscriptionHandler(subService)

	// Создаём роутер и маршруты
	r := mux.NewRouter()
	r.Use(logger.LoggingMiddleware) // Логируем все запросы
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/subscriptions", subHandler.Create).Methods("POST")
	r.HandleFunc("/subscriptions/{id}", subHandler.GetByID).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", subHandler.Update).Methods("PUT")
	r.HandleFunc("/subscriptions/{id}", subHandler.Delete).Methods("DELETE")
	r.HandleFunc("/subscriptions", subHandler.List).Methods("GET")
	r.HandleFunc("/subscriptions/sum", subHandler.Sum).Methods("GET")

	// Берём порт из env или ставим 8080 по умолчанию
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
