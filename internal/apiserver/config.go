package apiserver

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

const configPath = "configs/config.yml"

type Config struct {
	Port     string `yaml:"port"`
	LogLevel string `yaml:"log_level"`
	DB       struct {
		Username string `yaml:"username"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbname"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"db"`
}

func NewConfig() *Config {
	config := &Config{}
	log.Println("read application configuration")
	if err := cleanenv.ReadConfig(configPath, config); err != nil {
		help, _ := cleanenv.GetDescription(config, nil)
		log.Fatalln(help)
	}
	return config
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("user=%s host=%s dbname=%s password=%s sslmode=%s ",
		c.DB.Username, c.DB.Host, c.DB.DBName, c.DB.Password, c.DB.SSLMode)
}
