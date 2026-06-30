package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-default:"development"`
	Database   Database   `yaml:"database"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Auth       Auth       `yaml:"auth"`
	RabbitMQ   RabbitMQ   `yaml:"rabbitmq"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"go_postgres"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	DBName   string `yaml:"dbname" env-default:"go_app_db"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Auth struct {
	PrivateKeyPath string        `yaml:"private_key_path" env:"JWT_PRIVATE_KEY_PATH" env-default:"./internal/security/keys/private.pem"`
	PublicKeyPath  string        `yaml:"public_key_path" env:"JWT_PUBLIC_KEY_PATH" env-default:"./internal/security/keys/public.pem"`
	TokenTTL       time.Duration `yaml:"token_ttl" env:"JWT_TOKEN_TTL" env-default:"24h"`
}

type RabbitMQ struct {
	Host          string `yaml:"host" env:"host" env-default:"rabbitmq"`
	Port          string `yaml:"port" env:"port" env-default:"5672"`
	User          string `yaml:"user" env:"user" env-default:"guest"`
	Password      string `yaml:"password" env:"password" env-default:"guest"`
	RequestQueue  string `yaml:"request_queue" env:"request_queue" env-default:"api.requests"`
	ResponseQueue string `yaml:"response_queue" env:"response_queue" env-default:"api.responses"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
