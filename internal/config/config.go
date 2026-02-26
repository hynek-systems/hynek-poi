package config

import (
	"log"
	"os"
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
	OSM        ProviderConfig       `mapstructure:"osm"`
	Google     GoogleProviderConfig `mapstructure:"google"`
	HERE       ProviderConfig       `mapstructure:"here"`
	Foursquare FoursquareProviderConfig `mapstructure:"foursquare"`
}

type ProviderConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	Priority int           `mapstructure:"priority"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Retries  int           `mapstructure:"retries"`
}

type GoogleProviderConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	ApiKey   string        `mapstructure:"api_key"`
	Priority int           `mapstructure:"priority"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Retries  int           `mapstructure:"retries"`
}

type FoursquareProviderConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	ApiKey   string        `mapstructure:"api_key"`
	Priority int           `mapstructure:"priority"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Retries  int           `mapstructure:"retries"`
}

func Load() *Config {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(os.Getenv("HYNEK_POI_CONFIG_FILE"))
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetDefault("server.port", 8080)

	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("providers.osm.enabled", true)
	viper.SetDefault("providers.osm.weight", 10)
	viper.SetDefault("providers.osm.priority", 10)
	viper.SetDefault("providers.osm.timeout", "5s")
	viper.SetDefault("providers.osm.retries", 1)
	viper.SetDefault("providers.google.enabled", false)
	viper.SetDefault("providers.google.weight", 10)
	viper.SetDefault("providers.google.priority", 1)
	viper.SetDefault("providers.google.timeout", "2s")
	viper.SetDefault("providers.google.retries", 2)

	viper.SetDefault("providers.foursquare.enabled", false)
	viper.SetDefault("providers.foursquare.weight", 10)
	viper.SetDefault("providers.foursquare.priority", 5)
	viper.SetDefault("providers.foursquare.timeout", "3s")
	viper.SetDefault("providers.foursquare.retries", 2)

	viper.SetDefault("cache.ttl", "5m")

	viper.SetEnvPrefix("HYNEK_POI")

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
				Timeout:  viper.GetDuration("providers.osm.timeout"),
				Retries:  viper.GetInt("providers.osm.retries"),
			},
			Google: GoogleProviderConfig{
				Enabled:  viper.GetBool("providers.google.enabled"),
				ApiKey:   viper.GetString("providers.google.api_key"),
				Priority: viper.GetInt("providers.google.priority"),
				Timeout:  viper.GetDuration("providers.google.timeout"),
				Retries:  viper.GetInt("providers.google.retries"),
			},
			Foursquare: FoursquareProviderConfig{
				Enabled:  viper.GetBool("providers.foursquare.enabled"),
				ApiKey:   viper.GetString("providers.foursquare.api_key"),
				Priority: viper.GetInt("providers.foursquare.priority"),
				Timeout:  viper.GetDuration("providers.foursquare.timeout"),
				Retries:  viper.GetInt("providers.foursquare.retries"),
			},
		},
	}

	return cfg
}
