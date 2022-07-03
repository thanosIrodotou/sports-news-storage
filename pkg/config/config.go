package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	APP    APP
	Server Server
	API    API
	Mongo  Mongo
	Logger Logger
}

type APP struct {
	Name        string
	Version     string
	Environment string
}

type Server struct {
	Port         int16 `validate:"required,min=80"`
	ReadTimeout  int64 `validate:"required,min=5"`
	WriteTimeout int64 `validate:"required,min=5"`
	IdleTimeout  int64 `validate:"required,min=30"`
}

type API struct {
	GetLatestNewsArticlesUrl     string
	NewsArticlesPerCall          int
	GetArticleDetailsUrl         string
	NewNewsArticlesFetchInterval time.Duration
}

type Mongo struct {
	Host       string
	Port       int16
	Username   string
	Password   string
	Database   string
	Collection string
	TTL        time.Duration
}

type Logger struct {
	LogLevel       uint8
	AppName        string
	AppVersion     string
	AppEnvironment string
}

func New(path ...string) (*Config, error) {
	v := viper.New()
	setDefaults(v)

	v.SetEnvPrefix("app")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return nil, fmt.Errorf("failed to read config, %w", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	c.Logger.AppName = c.APP.Name
	c.Logger.AppEnvironment = c.APP.Environment

	return &c, nil
}

func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "sports-news-storage")
	v.SetDefault("app.version", "v0.0.1")
	v.SetDefault("app.environment", "dev")

	// Logger defaults
	v.SetDefault("logger.loglevel", 6)
	v.SetDefault("logger.appname", "sports-news-storage")
	v.SetDefault("logger.appversion", "v0.0.1")
	v.SetDefault("logger.appenvironment", "dev")

	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.readtimeout", 29)
	v.SetDefault("server.writetimeout", 29)
	v.SetDefault("server.idletimeout", 30)

	// API defaults
	v.SetDefault("api.getLatestNewsArticlesUrl", "https://www.brentfordfc.com/api/incrowd/getnewlistinformation?count=")
	v.SetDefault("api.newsArticlesPerCall", 50)
	v.SetDefault("api.getArticleDetailsUrl", "https://www.brentfordfc.com/api/incrowd/getnewsarticleinformation?id=")
	v.SetDefault("api.newNewsArticlesFetchInterval", "15s")

	// Mongo defaults
	v.SetDefault("mongo.host", "localhost")
	v.SetDefault("mongo.port", 27100)
	v.SetDefault("mongo.username", "root")
	v.SetDefault("mongo.password", "toor")
	v.SetDefault("mongo.database", "news")
	v.SetDefault("mongo.collection", "articles")
	v.SetDefault("mongo.ttl", "168h")
}
