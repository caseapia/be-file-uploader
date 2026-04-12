package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"be-file-uploader/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	appInstance, db, err := app.CreateApp()
	if err != nil {
		log.Fatalf("Error creating app instance: %s", err)
	}
	defer db.Close()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		if err := appInstance.Listen(":" + port); err != nil {
			log.Fatalf("Error starting app instance: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	log.Println("Завершение работы сервера...")
	if err := appInstance.Shutdown(); err != nil {
		log.Printf("Ошибка при остановке Fiber: %v", err)
	}
}
