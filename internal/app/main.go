package main

import (
	"NewsChanel/internal/config"
	"NewsChanel/internal/handler"
	"NewsChanel/internal/service"
	"NewsChanel/internal/storage"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

const (
	ProducerPort = ":8080"
	ConsumerPort = ":8081"
)

func main() {
	cfg := config.GetConfig()

	dsn := fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Printf("failed to close the database: %v", err)
		}
	}()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	storage := storage.New(db)
	service := service.New(storage)
	handler := handler.New(service)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	go func() {
		producer := service.StartProducer()
		defer func() {
			err = producer.Close()
			if err != nil {
				log.Printf("failed to close the producer: %v", err)
			}
		}()

		router.POST("/send", handler.SendNewsHandler(producer))
		if err = router.Run(ProducerPort); err != nil {
			log.Printf("failed to run the server: %v", err)
		}

	}()

	store, cancel := service.StartConsumer()
	defer cancel()
	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		handler.HandleNotifications(ctx, store)
	})
	if err = router.Run(ConsumerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}

}
