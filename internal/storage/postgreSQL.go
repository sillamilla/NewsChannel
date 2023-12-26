package storage

import (
	"NewsChanel/internal/models"
	"context"
	"database/sql"
	"github.com/lib/pq"
	"log"
	"time"
)

type Storage interface {
	InsertUser(ctx context.Context, user models.User, date time.Time) error
	InsertChannel(ctx context.Context, channel models.Channel, date time.Time) error
	InsertPieceOfNews(ctx context.Context, oneNews models.PieceOfNews, dataID, channelID string, followers []string, date time.Time) error
	GetNotification(id string, ctx context.Context) ([]models.PieceOfNews, error)
}

type postgreSQL struct {
	db *sql.DB
}

func New(db *sql.DB) Storage {
	return &postgreSQL{
		db: db,
	}
}

func (p *postgreSQL) InsertUser(ctx context.Context, user models.User, date time.Time) error {
	query := `INSERT INTO users PREMARY KEY ID (id, name, date) VALUES ($1, $2, $3)`
	_, err := p.db.ExecContext(ctx, query, user.ID, user.Name, date)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgreSQL) InsertChannel(ctx context.Context, channel models.Channel, date time.Time) error {
	query := `INSERT INTO channels (id, name, image, followers, date) VALUES ($1, $2, $3, $4, $5)`
	_, err := p.db.ExecContext(ctx, query, channel.ID, channel.Name, channel.Image, channel.Followers, date)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgreSQL) InsertPieceOfNews(ctx context.Context, data models.PieceOfNews, dataID, channelID string, followers []string, date time.Time) error {
	query := `INSERT INTO news (id, followers, channel_id, image, title, content, date) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := p.db.ExecContext(ctx, query, dataID, pq.Array(followers), channelID, data.Image, data.Title, data.Content, date)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgreSQL) GetNotification(userID string, ctx context.Context) ([]models.PieceOfNews, error) {
	query := `SELECT * FROM news WHERE user_id`
	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Printf("failed to close the rows in GetNotification: %v", err)
		}
	}()

	var allNews []models.PieceOfNews
	for rows.Next() {
		var news models.PieceOfNews
		err = rows.Scan(&news.Image, &news.Title, &news.Content, &news.Date)
		if err != nil {
			return nil, err
		}
		allNews = append(allNews, news)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allNews, nil
}
