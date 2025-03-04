package viper

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Viper *viper.Viper
}

func Init(configName string) Config {
	config := Config{Viper: viper.New()}

	vp := config.Viper

	// config type
	vp.SetConfigType("yml")

	// config name
	vp.SetConfigName(configName)

	// config path
	vp.AddConfigPath("./config")
	vp.AddConfigPath("../config")
	vp.AddConfigPath("../../config")

	// read config file
	if err := vp.ReadInConfig(); err != nil {
		log.Fatalf("Read config %s error, content is %+v", configName, err)
	}

	return config
}
