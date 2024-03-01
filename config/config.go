package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ServiceConfiguration struct {
	PostgresDB `yaml:"postgres_db" json:"postgresDB"`
	RedisDB    `yaml:"redis_db" json:"redisDB"`
	Api        `yaml:"api" json:"api"`
	User       `yaml:"user" json:"user"`
}

type Api struct {
	HOST string
	PORT string `yaml:"port"`
}

type User struct {
	SigningKey           string
	AccessTokenLifetime  int `yaml:"access_token_lifetime"`
	RefreshTokenLifetime int `yaml:"refresh_token_lifetime"`
}

type PostgresDB struct {
	Host     string
	Port     string `yaml:"port"`
	Username string
	Password string
	DBName   string
	SSLMode  string `yaml:"sslmode"`
}

type RedisDB struct {
	Host     string
	Port     string `yaml:"port"`
	Password string
}

func (api *Api) GetAddr() string {
	return fmt.Sprintf("%s:%s", api.HOST, api.PORT)
}

func NewConfig() ServiceConfiguration {
	return Load()
}

func Load() ServiceConfiguration {

	file, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	decoder := yaml.NewDecoder(file)
	var cfg ServiceConfiguration
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	cfg.User.SigningKey = getEnv("SIGNING_KEY")
	cfg.PostgresDB.Password = getEnv("POSTGRES_PASSWORD")
	cfg.PostgresDB.DBName = getEnv("POSTGRES_DB")
	cfg.PostgresDB.Username = getEnv("POSTGRES_USER")
	cfg.PostgresDB.Port = getEnv("POSTGRES_PORT")
	cfg.RedisDB.Password = getEnv("REDIS_PASSWORD")
	cfg.RedisDB.Host = getEnv("REDIS_HOST")
	cfg.PostgresDB.Host = getEnv("POSTGRES_HOST")
	cfg.Api.HOST = getEnv("API_HOST")

	return cfg
}

func getEnv(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("Environment variable %s is not set", key))
	}
	return val
}
