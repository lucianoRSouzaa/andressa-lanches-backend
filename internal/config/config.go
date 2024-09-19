package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	JWTSecret     string
	DatabaseURL   string
	ServerAddress string
	AuthUser      string
	AuthPassword  string
)

type Config struct {
	DatabaseURL   string
	JWTSecret     string
	ServerAddress string
	AuthUser      string
	AuthPassword  string
}

func LoadConfig() Config {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Nenhum arquivo de configuração encontrado. Usando variáveis de ambiente.")
	}

	config := Config{
		DatabaseURL:   viper.GetString("DATABASE_URL"),
		JWTSecret:     viper.GetString("JWT_SECRET"),
		ServerAddress: viper.GetString("SERVER_ADDRESS"),
		AuthUser:      viper.GetString("AUTH_USER"),
		AuthPassword:  viper.GetString("AUTH_PASSWORD"),
	}

	if config.DatabaseURL == "" || config.JWTSecret == "" || config.ServerAddress == "" || config.AuthUser == "" || config.AuthPassword == "" {
		log.Fatal("Variáveis de ambiente DATABASE_URL, JWT_SECRET, SERVER_ADDRESS, AUTH_USER e AUTH_PASSWORD são obrigatórias")
	}

	JWTSecret = config.JWTSecret
	DatabaseURL = config.DatabaseURL
	ServerAddress = config.ServerAddress
	AuthUser = config.AuthUser
	AuthPassword = config.AuthPassword

	return config
}
