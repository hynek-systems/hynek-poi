package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Redis     RedisConfig
	Cache     CacheConfig
	Providers ProvidersConfig
}

type ServerConfig struct {
	Port int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type CacheConfig struct {
	TTL time.Duration
}

type ProvidersConfig struct {
	OSM    ProviderConfig       `mapstructure:"osm"`
	Google GoogleProviderConfig `mapstructure:"google"`
	HERE   ProviderConfig       `mapstructure:"here"`
}

type ProviderConfig struct {
	Enabled  bool `mapstructure:"enabled"`
	Priority int  `mapstructure:"priority"`
}

type GoogleProviderConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	ApiKey   string `mapstructure:"api_key"`
	Priority int    `mapstructure:"priority"`
}

func Load() *Config {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetDefault("server.port", 8080)

	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("providers.osm.enabled", true)
	viper.SetDefault("providers.osm.weight", 10)
	viper.SetDefault("providers.osm.priority", 10)
	viper.SetDefault("providers.google.enabled", false)
	viper.SetDefault("providers.google.weight", 10)
	viper.SetDefault("providers.google.priority", 1)

	viper.SetDefault("cache.ttl", "5m")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		log.Println("No config file found, using defaults")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: viper.GetInt("server.port"),
		},

		Redis: RedisConfig{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},

		Cache: CacheConfig{
			TTL: viper.GetDuration("cache.ttl"),
		},

		Providers: ProvidersConfig{
			OSM: ProviderConfig{
				Enabled:  viper.GetBool("providers.osm.enabled"),
				Priority: viper.GetInt("providers.osm.priority"),
			},
			Google: GoogleProviderConfig{
				Enabled:  viper.GetBool("providers.google.enabled"),
				ApiKey:   viper.GetString("providers.google.api_key"),
				Priority: viper.GetInt("providers.google.priority"),
			},
		},
	}

	return cfg
}
