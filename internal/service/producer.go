package service

import (
	"NewsChanel/internal/models"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	ProducerPort = ":8080"
	KafkaTopic   = "notifications"
)

func receivePost(ctx *gin.Context) models.Notification {
	var channel models.Notification
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	return channel
}

func (s *service) SendNews(producer sarama.SyncProducer, ctx *gin.Context) error {
	channel := receivePost(ctx)

	for _, user := range channel.From.Followers {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		notification := models.Notification{
			From: channel.From,
			News: channel.News,
		}

		notificationJSON, err := json.Marshal(notification)
		if err != nil {
			return fmt.Errorf("failed to marshal notification: %w", err)
		}
		msg := &sarama.ProducerMessage{
			Topic: KafkaTopic,
			Key:   sarama.StringEncoder(user.ID),
			Value: sarama.StringEncoder(notificationJSON),
		}

		_, _, err = producer.SendMessage(msg)
		if err != nil {
			return err
		}

	}

	return nil
}

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}

func (s *service) StartProducer() sarama.SyncProducer {
	producer, err := setupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}

	fmt.Printf("PRODUCER started at http://localhost%s\n", ProducerPort)

	return producer
}
