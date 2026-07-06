package env

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Env struct {
	IsProduction bool
	Port         string
	Host         string
	DatabaseUrl  string

	AccessTokenExpiresAt time.Duration
	AccessTokenSecret    string

	RefreshTokenExpiresAt time.Duration
	RefreshTokenSecret    string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	DomainFE string
	DomainBE string

	CloudinaryURL string

	RedisAddr string
	RedisPass string

	ElasticAddr            string
	ElasticUser            string
	ElasticPass            string
	ElasticCertFingerprint string

	RabbitMQURL string
}

func New() *Env {
	godotenv.Load()

	isProduction := os.Getenv("IS_PRODUCTION") == "true"
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	DatabaseUrl := os.Getenv("DATABASE_URL")

	AccessTokenExpiresAtString := os.Getenv("ACCESS_TOKEN_EXPIRES_AT")
	AccessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	AccessTokenExpiresAt := getDuration(AccessTokenExpiresAtString)

	RefreshTokenExpiresAtString := os.Getenv("REFRESH_TOKEN_EXPIRES_AT")
	RefreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	RefreshTokenExpiresAt := getDuration(RefreshTokenExpiresAtString)

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	domainFE := os.Getenv("DOMAIN_FE")
	domainBE := os.Getenv("DOMAIN_BE")
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")

	RedisAddr := os.Getenv("REDIS_ADDR")
	RedisPass := os.Getenv("REDIS_PASS")

	elasticAddr := os.Getenv("ELASTIC_ADDRS")
	elasticUser := os.Getenv("ELASTIC_USER")
	elasticPass := os.Getenv("ELASTIC_PASSWORD")
	elasticCertFingerprint := os.Getenv("ELASTIC_CERT_FINGERPRINT")

	rabbitMQURL := os.Getenv("RABBIT_MQ_URL")

	return &Env{
		IsProduction:           isProduction,
		Port:                   port,
		Host:                   host,
		DatabaseUrl:            DatabaseUrl,
		AccessTokenExpiresAt:   AccessTokenExpiresAt,
		AccessTokenSecret:      AccessTokenSecret,
		RefreshTokenExpiresAt:  RefreshTokenExpiresAt,
		RefreshTokenSecret:     RefreshTokenSecret,
		GoogleClientID:         googleClientID,
		GoogleClientSecret:     googleClientSecret,
		GoogleRedirectURL:      googleRedirectURL,
		DomainFE:               domainFE,
		DomainBE:               domainBE,
		CloudinaryURL:          cloudinaryURL,
		RedisAddr:              RedisAddr,
		RedisPass:              RedisPass,
		ElasticAddr:            elasticAddr,
		ElasticUser:            elasticUser,
		ElasticPass:            elasticPass,
		ElasticCertFingerprint: elasticCertFingerprint,
		RabbitMQURL:            rabbitMQURL,
	}
}

func getDuration(durationString string) time.Duration {
	durationTime, err := time.ParseDuration(durationString)
	if err != nil {
		log.Fatal("Error Parser")
	}

	return durationTime
}
