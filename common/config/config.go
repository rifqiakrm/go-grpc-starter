// Package config parses environment variable into usable structs
package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// Config holds configuration for the project.
type Config struct {
	Env          string `env:"ENV,default=development"`
	ServiceName  string `env:"SERVICE_NAME,default=grpc-starter"`
	Port         Port
	HashID       HashID
	Google       Google
	Postgres     Postgres
	Redis        Redis
	Jaeger       Jaeger
	JWTConfig    JWTConfig
	SMTP         SMTP
	Mailgun      Mailgun
	Sendgrid     Sendgrid
	CloudStorage CloudStorage
}

// Port holds configuration for project's port.
type Port struct {
	GRPC string `env:"PORT_GRPC,default=8081"`
	REST string `env:"PORT,default=8080"`
}

// Google holds configuration for the Google.
type Google struct {
	ProjectID          string `env:"GOOGLE_APPLICATION_PROJECT_ID,required"`
	ServiceAccountFile string `env:"GOOGLE_APPLICATION_CREDENTIALS,required"`
}

// HashID holds configuration for HashID.
type HashID struct {
	Salt      string `env:"HASHID_SALT"`
	MinLength int    `env:"HASHID_MIN_LENGTH,default=10"`
}

// Postgres holds all configuration for PostgreSQL.
type Postgres struct {
	Host            string `env:"POSTGRES_HOST,default=localhost"`
	Port            string `env:"POSTGRES_PORT,default=5432"`
	User            string `env:"POSTGRES_USER,required"`
	Password        string `env:"POSTGRES_PASSWORD,required"`
	Name            string `env:"POSTGRES_NAME,required"`
	MaxOpenConns    string `env:"POSTGRES_MAX_OPEN_CONNS,default=5"`
	MaxConnLifetime string `env:"POSTGRES_MAX_CONN_LIFETIME,default=10m"`
	MaxIdleLifetime string `env:"POSTGRES_MAX_IDLE_LIFETIME,default=5m"`
}

// Redis holds configuration for the Redis.
type Redis struct {
	Address  string `env:"REDIS_ADDRESS,required"`
	Password string `env:"REDIS_PASSWORD"`
}

// Jaeger holds configuration for the Jaeger.
type Jaeger struct {
	Address string `env:"JAEGER_ADDRESS"`
}

// JWTConfig holds configuration for jwt.
type JWTConfig struct {
	SecretKey string `env:"JWT_SECRET_KEY"`
}

// SMTP holds configuration for smtp email.
type SMTP struct {
	Host      string `env:"SMTP_HOST"`
	Port      int    `env:"SMTP_PORT,default=587"`
	User      string `env:"SMTP_USER"`
	Pass      string `env:"SMTP_PASS"`
	FromName  string `env:"SMTP_FROM_NAME"`
	FromEmail string `env:"SMTP_FROM_EMAIL"`
}

// Mailgun holds configuration for mailgun service.
type Mailgun struct {
	APIKey string `env:"MAILGUN_API_KEY"`
	Domain string `env:"MAILGUN_DOMAIN"`
}

// Sendgrid holds configuration for sendgrid service.
type Sendgrid struct {
	APIKey string `env:"SENDGRID_API_KEY"`
}

// CloudStorage holds configuration for file service.
type CloudStorage struct {
	AssetURL string `env:"ASSET_URL"`
}

// NewConfig creates an instance of Config.
// It needs the path of the env file to be used.
func NewConfig(env string) (*Config, error) {
	// just skip loading env files if it is not exists, env files only used in local dev
	_ = godotenv.Load(env)

	var config Config
	if err := envdecode.Decode(&config); err != nil {
		return nil, errors.Wrap(err, "[NewConfig] error decoding env")
	}

	return &config, nil
}
