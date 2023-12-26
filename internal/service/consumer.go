package service

import (
	"NewsChanel/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"sync"
)

const (
	ConsumerGroup      = "notifications-group"
	ConsumerTopic      = "notifications"
	ConsumerPort       = ":8081"
	KafkaServerAddress = "localhost:9092"
)

type UserNotifications map[string][]models.Notification

type NotificationStore struct {
	data UserNotifications
	mu   sync.RWMutex
}

func (ns *NotificationStore) Add(userID string, notification models.Notification) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.data[userID] = append(ns.data[userID], notification)
}

func (ns *NotificationStore) Get(userID string) []models.Notification {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.data[userID]
}

type Consumer struct {
	store   *NotificationStore
	service *service
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		userID := string(msg.Key)

		var notification models.Notification
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}

		err = consumer.service.InsertNotification(context.Background(), notification)
		if err != nil {
			log.Printf("failed to insert notification: %v", err)
			continue
		}

		consumer.store.Add(userID, notification)
		sess.MarkMessage(msg, "")
	}

	return nil
}

func initializeConsumerGroup() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup([]string{KafkaServerAddress}, ConsumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize consumer group: %w", err)
	}

	return consumerGroup, nil
}

func setupConsumerGroup(ctx context.Context, store *NotificationStore, service *service) {
	consumerGroup, err := initializeConsumerGroup()
	if err != nil {
		log.Printf("initialization error: %v", err)
	}
	defer consumerGroup.Close()

	consumer := &Consumer{
		store:   store,
		service: service,
	}

	for {
		if ctx.Err() != nil {
			return
		}

		err = consumerGroup.Consume(ctx, []string{ConsumerTopic}, consumer)
		if err != nil {
			log.Printf("error from consumer: %v", err)
		}
	}
}

func (s *service) StartConsumer() (*NotificationStore, context.CancelFunc) {
	store := &NotificationStore{
		data: make(UserNotifications),
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		setupConsumerGroup(ctx, store, s)
	}()

	fmt.Printf("CONSUMER (Group: %s) "+" started at http://localhost%s\n", ConsumerGroup, ConsumerPort)

	return store, cancel
}
