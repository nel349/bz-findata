package config

import (
	"context"
	"log"
	"os"

	awslocal "github.com/nel349/bz-findata/pkg/aws"
	"github.com/sethvargo/go-envconfig"
)

// Config default config structure
type Config struct {
	Exchange ExchangeConfig `env:",prefix=EXCHANGE_,required"`
	Database DatabaseConfig `env:",prefix=DB_,required"`
	Logger   LoggerConfig   `env:",prefix=LOGGER_"`
}

// AnalysisConfig for analysis configuration
type AnalysisConfig struct {
	Database DatabaseConfig `env:",prefix=DB_,required"`
}

// DexConfig for dex configuration
type DexConfig struct {
	Database DatabaseConfig `env:",prefix=DB_,required"`
}

// LoggerConfig for logger configuration
type LoggerConfig struct {
	DisableCaller     bool   `env:"CALLER,default=false"`
	DisableStacktrace bool   `env:"STACKTRACE,default=false"`
	Level             string `env:"LEVEL,default=debug"`
}

// ExchangeConfig for exchange configuration
type ExchangeConfig struct {
	Url      string   `env:"URL,required"`
	Origin   string   `env:"ORIGIN,required"`
	Protocol string   `env:"PROTOCOL,default="`
	Symbols  []string `env:"SYMBOLS,required"`
	Channels []string `env:"CHANNELS,required"`
}

// DatabaseConfig for db config
type DatabaseConfig struct {
	Host     string `env:"HOST,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD"`
	Base     string `env:"BASE"`
}

// NewConfig init default config for application
func NewConfig(ctx context.Context) (*Config, error) {
	var cfg Config

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}
	setDBPassword(&cfg)
	return &cfg, nil
}

func NewAnalysisConfig(ctx context.Context) (*AnalysisConfig, error) {
	var cfg AnalysisConfig

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}
	setDBPassword(&cfg)
	return &cfg, nil
}

func NewDexConfig(ctx context.Context) (*DexConfig, error) {
	var cfg DexConfig

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}
	setDBPassword(&cfg)
	return &cfg, nil
}

func setDBPassword(cfg interface{}) {
	var dbConfig *DatabaseConfig
	switch c := cfg.(type) {
	case *Config:
		dbConfig = &c.Database
	case *AnalysisConfig:
		dbConfig = &c.Database
	case *DexConfig:
		dbConfig = &c.Database
	default:
		log.Fatal("unsupported config type")
	}

	if os.Getenv("IS_LOCAL") == "true" {
		dbConfig.Password = os.Getenv("DB_PASSWORD")
	} else {
		dbSecret, err := awslocal.GetDefaultDBSecret()
		if err != nil {
			log.Fatalf("failed to retrieve DB secret: %v", err)
		}
		dbConfig.Password = dbSecret.DB_PASSWORD
	}
}
