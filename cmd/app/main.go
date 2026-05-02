package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"be-file-uploader/internal/app"
	"be-file-uploader/pkg/geo"

	"github.com/gookit/slog"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	geoService, err := geo.New("data/IP2LOCATION-LITE-DB11.IPV6.BIN")
	if err != nil {
		slog.Fatalf("geo init error: %s", err)
	}
	defer geoService.Close()

	appInstance, db, err := app.CreateApp(geoService)
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

	log.Println("Server is shutting down...")
	if err := appInstance.Shutdown(); err != nil {
		log.Printf("Ошибка при остановке Fiber: %v", err)
	}
}
