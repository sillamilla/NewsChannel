package service

import (
	"NewsChanel/internal/models"
	"NewsChanel/internal/storage"
	"context"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	StartProducer() sarama.SyncProducer
	StartConsumer() (*NotificationStore, context.CancelFunc)
	SendNews(producer sarama.SyncProducer, ctx *gin.Context) error

	InsertNotification(ctx context.Context, notification models.Notification) error
	GetNotification(id string, ctx context.Context) ([]models.PieceOfNews, error)
}

type service struct {
	db storage.Storage
}

func New(db storage.Storage) Service {
	return &service{
		db: db,
	}
}

func (s *service) InsertUser(ctx context.Context, user models.User) error {
	userID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	user.ID = userID.String()
	return s.db.InsertUser(ctx, user, time.Now())
}

func (s *service) InsertChannel(ctx context.Context, channel models.Channel) error {
	channelID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	channel.ID = channelID.String()
	return s.db.InsertChannel(ctx, channel, time.Now())
}

func (s *service) InsertNotification(ctx context.Context, notification models.Notification) error {
	dataID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var followers []string
	for _, follower := range notification.From.Followers {
		followers = append(followers, follower.ID)
	}

	channelID := notification.From.ID
	for _, data := range notification.News {
		err = s.db.InsertPieceOfNews(ctx, data, dataID.String(), channelID, followers, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) GetNotification(id string, ctx context.Context) ([]models.PieceOfNews, error) {
	news, err := s.db.GetNotification(id, ctx)
	if err != nil {
		return nil, err
	}

	return news, nil
}
