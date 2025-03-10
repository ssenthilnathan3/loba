package config

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Target struct {
	Host_name        string
	Host_url         string
	Capacity         string
	Consumption_rate string
	Hash             string
}

type Configuration struct {
	Server struct {
		Listen_port string
	}
	Targets []Target
}

var configuration *Configuration

func Hash(url string) (uint32, error) {
	hostURL := strings.Split(strings.SplitAfter(`//`, url)[1], ":")[0]
	if hostURL == "" {
		return 0, fmt.Errorf("error reading host url %s", "")
	}

	h := sha256.New()

	h.Write([]byte(hostURL))

	hashedURL := h.Sum(nil)

	return binary.BigEndian.Uint32(hashedURL[:4]), nil
}

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
