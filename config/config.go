package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type Config struct {
	Env            string
	Port           int
	Log            Log
	Postgres       Postgres
	JWT            JWT
	ExchangeClient ExchangeClient
	HTTP           HTTP
	Redis          Redis
}

type Log struct {
	Enable         bool
	Level          string
	TimeLayout     string
	Caller         bool
	Trace          bool
	FilePath       string
	FileMaxSize    int
	FileMaxBackups int
	FileMaxAge     int
	FileCompress   bool
}

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWT struct {
	SecretKey         string
	TokenExpiration   time.Duration
	RefreshExpiration time.Duration
	Issuer            string
}

type ExchangeClient struct {
	Timeout time.Duration
}

type HTTP struct {
	Port  int
	Mode  string
	Guard string
}

func NewConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	overrideWithEnv(v, &cfg)

	return &cfg, nil
}

func overrideWithEnv(v *viper.Viper, cfg *Config) {
	if env := os.Getenv("ENV"); env != "" {
		cfg.Env = env
	}
	if port := os.Getenv("PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Port)
	}

	// اورراید سایر مقادیر به همین شکل...

	// مثلا برای لاگ
	if val := os.Getenv("LOG_ENABLE"); val != "" {
		cfg.Log.Enable = strings.ToLower(val) == "true"
	}
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}

	// سایر متغیرهای ENV رو به صورت مشابه اضافه کن
}
