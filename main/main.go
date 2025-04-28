package main

import (
	"context"
	"demo-prismao-apicodegen/prisma/db"
	"demo-prismao-apicodegen/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Khởi tạo Prisma Client
	client := db.NewClient()

	// Kết nối database
	if err := client.Prisma.Connect(); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			log.Printf("failed to disconnect database: %v", err)
		}
	}()

	e := server.SetupEcho(client)

	go func() {
		if err := e.Start(":8080"); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
