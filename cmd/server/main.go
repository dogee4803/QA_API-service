package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"qna-api/internal/config"
	"qna-api/internal/database"
	"qna-api/internal/handler"
	"qna-api/internal/repository"
	"qna-api/internal/service"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := config.Load()
	log.Printf("Конфигурация загружена: порт %s, БД %s", cfg.ServerPort, cfg.DBName)

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии подключения к БД: %v", err)
		} else {
			log.Println("Подключение к базе данных закрыто")
		}
	}()

	questionRepo := repository.NewQuestionRepository(db.DB)
	answerRepo := repository.NewAnswerRepository(db.DB)

	questionService := service.NewQuestionService(questionRepo, answerRepo)
	answerService := service.NewAnswerService(answerRepo, questionRepo)

	router := handler.SetupRoutes(questionService, answerService)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
		
		ReadTimeout: 15 * time.Second,
		
		WriteTimeout: 15 * time.Second,
		
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		log.Printf("HTTP сервер запускается на порту %s", cfg.ServerPort)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	
	signal.Notify(quit, 
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	
	sig := <-quit

	log.Println("Завершение работы сервера...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Принудительное завершение сервера: %v", err)
	}

	log.Println("Сервер корректно остановлен")
}