package config

import (
	"log/slog"
	"os"

	errors "github.com/kerilOvs/profile_sevice/internal/errorsExt"
	"github.com/kerilOvs/profile_sevice/pkg/logger"

	yaml "gopkg.in/yaml.v3"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint" env:"MINIO_ENDPOINT"`
	AccessKey string `yaml:"access_key" env:"MINIO_ACCESS_KEY"`
	SecretKey string `yaml:"secret_key" env:"MINIO_SECRET_KEY"`
	Bucket    string `yaml:"bucket" env:"MINIO_BUCKET"`
	UseSSL    bool   `yaml:"use_ssl" env:"MINIO_USE_SSL"`
}

type RabbitConfig struct {
	Url            string `yaml:"url" env:"RABBIT_URL"`
	QueuePhotoName string `yaml:"queue_photo_name" env:"RABBIT_PHOTO_NAME"`
	QueueTagsName  string `yaml:"queue_tags_name" env:"RABBIT_TAGS_NAME"`
	QueueAnketName string `yaml:"queue_anket_name" env:"RABBIT_ANKET_NAME"`
}

type LogConfig struct {
	LogLevel  string `yaml:"endpoint" env:"MINIO_ENDPOINT"`
	LogFormat string `yaml:"queue_photo_name" env:"RABBIT_PHOTO_NAME"`
}

type Config struct {
	Database DBConfig     `yaml:"database"`
	Server   ServerConfig `yaml:"server"`
	Minio    MinioConfig  `yaml:"minio"`
	Rabbit   RabbitConfig `yaml:"rabbit"`
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Group("database",
			slog.String("host", c.Database.Host),
			slog.Int("port", c.Database.Port),
			slog.String("user", c.Database.User),
			slog.Any("password", logger.Secret(c.Database.Password)),
			slog.String("name", c.Database.User),
		),
		slog.Group("server",
			slog.Int("port", c.Server.Port),
		),
		slog.Group("minio",
			slog.String("endpoint", c.Minio.Endpoint),
			slog.Any("access_key", logger.Secret(c.Minio.AccessKey)),
			slog.Any("secret_key", logger.Secret(c.Minio.SecretKey)),
			slog.String("bucket", c.Minio.Bucket),
			slog.Bool("use_ssl", c.Minio.UseSSL),
		),
		slog.Group("rabbit",
			slog.String("url", c.Rabbit.Url),
			slog.String("queue_photo_name", c.Rabbit.QueuePhotoName),
			slog.String("queue_tags_name", c.Rabbit.QueueTagsName),
			slog.String("queue_anket_name", c.Rabbit.QueueAnketName),
		),
	)
}
func ReadConfig() (Config, error) {

	var config Config

	data, err := os.ReadFile(fileName)
	if err != nil {
		return config, errors.ErrorLocate(err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, errors.ErrorLocate(err)
	}
	return config, nil
}
