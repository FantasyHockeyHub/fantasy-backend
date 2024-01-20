package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ServiceConfiguration struct {
	PostgresDB `yaml:"db" json:"postgresDB"`
	Api        `yaml:"api" json:"api"`
	User       `yaml:"user" json:"user"`
}

type Api struct {
	PORT string `yaml:"port"`
}

type User struct {
	PasswordSalt string `yaml:"password_salt"`
}

type PostgresDB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func (api *Api) GetAddr() string {
	return fmt.Sprintf("localhost:%s", api.PORT)
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
	return cfg
}
