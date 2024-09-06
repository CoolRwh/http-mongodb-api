package config

import (
	"github.com/spf13/viper"
	"log"
)

func InitConfig() {
	viper.SetConfigFile("config/config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
		return
	}
	log.Println("init config success!")
}
