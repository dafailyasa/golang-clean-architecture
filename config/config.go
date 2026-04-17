package config

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Keycloak KeycloakConfig
}

type AppConfig struct {
	Name        string
	Version     string
	Schema      string
	Host        string
	Environment string
}

type ServerConfig struct {
	Port     string
	Debug    bool
	TimeZone string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	DBName   string
	UserName string
	Password string
	Debug    bool
	Pool     PoolConfig
}

type PoolConfig struct {
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifetime time.Duration
}

type KeycloakConfig struct {
	KeycloakBaseURL string
	Realm           string
	ClientID        string
	AdminKeycloak   AdminKeycloak
	ClientSecret    string
}

type AdminKeycloak struct {
	Email    string
	Password string
}

func LoadConfigPath() (Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./../")
	v.AddConfigPath("./")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return Config{}, errors.New("config file not found")
		}
		return Config{}, err
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return Config{}, err
	}

	return c, nil
}
