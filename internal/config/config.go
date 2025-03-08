package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Target struct {
	Host_name        string
	Host_url         string
	Capacity         string
	Consumption_rate string
}

type Configuration struct {
	Server struct {
		Listen_port string
	}
	Targets []Target
}

var configuration *Configuration

func NewConfiguration() (*Configuration, error) {
	viper.AddConfigPath("internal")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()

	if err != nil {
		return nil, fmt.Errorf("error occurred when reading config %s", err)
	}

	err = viper.Unmarshal(&configuration)

	if err != nil {
		return nil, fmt.Errorf("error occurred when unmarshaling %s", err)
	}
	return configuration, nil
}
