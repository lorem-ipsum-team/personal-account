package rabbit

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	"github.com/google/uuid"
	"github.com/kerilOvs/profile_sevice/internal/config"
	"github.com/kerilOvs/profile_sevice/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	msgTimeout = 5 * time.Second
)

type Repo struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queues  map[string]amqp.Queue // Храним информацию об очередях

	tagsQueueName   string
	photosQueueName string
	anketsQueueName string
}

func New(ctx context.Context, cfg *config.RabbitConfig) (*Repo, error) {
	tagsQueueName := cfg.QueueTagsName
	photosQueueName := cfg.QueuePhotoName
	anketsQueueName := cfg.QueueAnketName

	conn, err := amqp.Dial(cfg.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Rabbit: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Объявляем три разные очереди
	queues := make(map[string]amqp.Queue)

	tagsQueue, err := channel.QueueDeclare(
		tagsQueueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare tags queue: %w", err)
	}
	queues[tagsQueueName] = tagsQueue

	photosQueue, err := channel.QueueDeclare(
		photosQueueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare photos queue: %w", err)
	}
	queues[photosQueueName] = photosQueue

	anketsQueue, err := channel.QueueDeclare(
		anketsQueueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare ankets queue: %w", err)
	}
	queues[anketsQueueName] = anketsQueue

	repo := &Repo{
		conn:            conn,
		channel:         channel,
		queues:          queues,
		tagsQueueName:   tagsQueueName,
		photosQueueName: photosQueueName,
		anketsQueueName: anketsQueueName,
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
	Tags   string
}

func (r *Repo) PublishTags(ctx context.Context, tags Tags) error {
	tagsDto := Tags{
		UserID: tags.UserID,
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
		"",              // exchange
		r.tagsQueueName, // routing key (имя очереди)
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish tags: %w", err)
	}

	return nil
}

type Photo struct {
	ID   uuid.UUID `json:"user_id"`
	Path string    `json:"image_url"`
}

func (r *Repo) PublishPhoto(ctx context.Context, photo Photo) error {
	photoDto := Photo{
		ID:   photo.ID,
		Path: photo.Path,
	}

	body, err := json.Marshal(photoDto)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, msgTimeout)
	defer cancel()

	err = r.channel.PublishWithContext(
		ctx,
		"",                // exchange
		r.photosQueueName, // routing key (имя очереди)
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish photo: %w", err)
	}

	return nil
}

type UserAnket struct {
	ID        uuid.UUID         `json:"user_id"`
	Gender    models.UserGender `json:"gender"`
	BirthDate string            `json:"birth_date"`
	//Description string    `json:"description"`
}

func (r *Repo) PublishAnket(ctx context.Context, anket UserAnket) error {
	anketDto := UserAnket{
		ID:        anket.ID,
		Gender:    anket.Gender,
		BirthDate: anket.BirthDate,
	}

	body, err := json.Marshal(anketDto)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, msgTimeout)
	defer cancel()

	err = r.channel.PublishWithContext(
		ctx,
		"",                // exchange
		r.anketsQueueName, // routing key (имя очереди)
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish anket: %w", err)
	}

	return nil
}
