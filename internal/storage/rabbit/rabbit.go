package rabbit

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	"github.com/google/uuid"
	"github.com/kerilOvs/profile_sevice/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	queueName  = "swipes"
	msgTimeout = 5 * time.Second
)

type Repo struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(ctx context.Context, cfg *config.RabbitConfig) (*Repo, error) {
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
	}

	go func() {
		<-ctx.Done()

		_ = repo.Close()
	}()

	return repo, nil
}

func (r *Repo) Close() error {
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

type Photo struct {
	ID   uuid.UUID `json:"id"`
	Path string    `json:"image_url"`
}

func (r *Repo) PublishPhoto(ctx context.Context, photos Photo) error {
	photosDto := Photo{
		ID:   uuid.UUID(photos.ID),
		Path: photos.Path,
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
	ID          uuid.UUID `json:"user_id"`
	Description string    `json:"description"`
}

func (r *Repo) PublishAnket(ctx context.Context, anket UserAnket) error {
	anketDto := UserAnket{
		ID:          uuid.UUID(anket.ID),
		Description: anket.Description,
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
