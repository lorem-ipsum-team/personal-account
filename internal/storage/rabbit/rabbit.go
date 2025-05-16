package rabbit

import (
	"context"
	"encoding/json"
	"fmt"

	//"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/kerilOvs/profile_sevice/internal/config"
	"github.com/kerilOvs/profile_sevice/internal/models"

	//"github.com/kerilOvs/profile_sevice/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	queueName  = "swipes"
	msgTimeout = 5 * time.Second
)

type Repo struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//log     *slog.Logger
}

/*type Swipe struct {
	Init   uuid.UUID `json:"init"`
	Target uuid.UUID `json:"target"`
	Like   bool      `json:"like"`
}*/

func New(ctx context.Context, cfg *config.RabbitConfig) (*Repo, error) {
	//log = log.WithGroup("rabbit_repo")
	//log.Info("connect to rabbit", slog.Any("connection string", logger.Secret(cfg.Url)))

	conn, err := amqp.Dial(cfg.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Rabbit: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	repo := &Repo{
		conn:    conn,
		channel: channel,
		//log:     log,
	}

	go func() {
		<-ctx.Done()

		_ = repo.Close()
	}()

	return repo, nil
}

func (r *Repo) Close() error {
	// r.log.Info("closing rabbit_repo")

	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("failed to close amqp channel: %w", err)
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close amqp connection: %w", err)
	}

	return nil
}

type Tags struct {
	UserID uuid.UUID
	Tags   []string
}

func (r *Repo) PublishTags(ctx context.Context, tags Tags) error {
	tagsDto := Tags{
		UserID: uuid.UUID(tags.UserID),
		Tags:   tags.Tags,
	}

	body, err := json.Marshal(tagsDto)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, msgTimeout)
	defer cancel()

	err = r.channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{ //nolint:exhaustruct
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish swipe: %w", err)
	}

	return nil
}

type Photos struct {
	UserID     uuid.UUID
	Photos_url []string
}

func (r *Repo) PublishPhotos(ctx context.Context, photos Photos) error {
	photosDto := Photos{
		UserID:     uuid.UUID(photos.UserID),
		Photos_url: photos.Photos_url,
	}

	body, err := json.Marshal(photosDto)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, msgTimeout)
	defer cancel()

	err = r.channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{ //nolint:exhaustruct
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish swipe: %w", err)
	}

	return nil
}

type UserAnket struct {
	UserID    uuid.UUID
	BirthDate time.Time
	Gender    models.UserGender
}

func (r *Repo) PublishAnket(ctx context.Context, anket UserAnket) error {
	anketDto := UserAnket{
		UserID:    uuid.UUID(anket.UserID),
		BirthDate: anket.BirthDate,
		Gender:    anket.Gender,
	}

	body, err := json.Marshal(anketDto)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, msgTimeout)
	defer cancel()

	err = r.channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{ //nolint:exhaustruct
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish swipe: %w", err)
	}

	return nil
}
