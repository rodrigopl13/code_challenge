package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configs struct {
	//CrtFile        string   `mapstructure:"crt-file"`
	//Keyfile        string   `mapstructure:"key-file"`
	StooqURLString string   `mapstructure:"stooq-url"`
	JwtSecretKey   string   `mapstructure:"SECRET_KEY"`
	Repo           Database `mapstructure:"database"`
	RabbitMQ       Broker   `mapstructure:"broker"`
}

type Database struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"DB_PASSWORD"`
}

type Broker struct {
	AmqpUrl        string `mapstructure:"amqp-url"`
	StockQueueName string `mapstructure:"stock-queue-name"`
}

func Load() *Configs {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	c := Configs{}
	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("Fatal error unmarshal config file: %w \n", err))
	}
	c.Repo.Password = viper.GetString("DB_PASSWORD")
	return &c
}
