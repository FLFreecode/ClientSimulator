package config

import (
	"bytes"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const Namespace = "Virtual-Clients"

const Default = `
log:
  level: 0

clients:
  filename: "ClientsData.txt"
  numberofclients: 1
  resetuuid: false
  numberofrequest: 1000
  dilay: 1000  			# miliseconds 
  qoute: " Qoute Example : One swallow doesn't make summer "
  Postduplicateqoute: false

server:
  addr: "http://localhost:8690"
  `

type (
	Config struct {
		Log     Log     `mapstructure:"log" validate:"required"`
		Clients Clients `mapstructure:"clients" validate:"required"`
		Server  Server  `mapstructure:"server" validate:"required"`
	}

	Log struct {
		Level int `mapstructure:"level"`
	}

	Clients struct {
		FileName           string `mapstructure:"filename" validate:"required"`
		NumberOfClients    int    `mapstructure:"numberofclients" validate:"required"`
		ResetUUID          bool   `mapstructure:"resetuuid" validate:"required"`
		NumberOfRequest    int    `mapstructure:"numberofrequest" validate:"required"`
		Dilay              int    `mapstructure:"dilay" validate:"required"`
		Qoute              string `mapstructure:"qoute" validate:"required"`
		PostDuplicateQoute bool   `mapstructure:"postduplicateqoute" validate:"required"`
	}

	Server struct {
		Addr string `mapstructure:"addr" validate:"required"`
	}
)

var (
	config Config
	logger = log.With().Str("Service", "Visrtual Clients").Logger()
)

func Get() *Config {
	return &config
}

func (c Config) Validate() error {
	return validator.New().Struct(c)
}

func Load(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetEnvPrefix(Namespace)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadConfig(bytes.NewReader([]byte(Default))); err != nil {
		logger.Error().Err(err).Msgf("error loading default configs")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Info().Msgf("Config file changed %s", file)
		reload(e.Name)
	})

	return reload(file)
}

func reload(file string) bool {
	err := viper.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Error().Err(err).Msgf("config file not found %s", file)
		} else {
			logger.Error().Err(err).Msgf("config file read failed %s", file)
		}
		return false
	}

	err = viper.GetViper().UnmarshalExact(&config)
	if err != nil {
		logger.Error().Err(err).Msgf("config file loaded failed %s", file)
		return false
	}

	if err = config.Validate(); err != nil {
		logger.Error().Err(err).Msgf("invalid configuration %s", file)
	}

	logger.Info().Msgf("Config file loaded %s", file)
	return true
}
