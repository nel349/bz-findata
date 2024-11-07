package config

import (
	"context"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Setenv("IS_LOCAL", "true") // Local environment; Prevents from using AWS secrets
    t.Setenv("DB_PASSWORD", "a2s_kjlasjd")
    
    // Set other required env vars
    t.Setenv("EXCHANGE_URL", "wss://ws-feed.exchange.coinbase.com")
    t.Setenv("EXCHANGE_ORIGIN", "https://coinbase.com")
    t.Setenv("EXCHANGE_SYMBOLS", "ETH-BTC,BTC-USD,BTC-EUR")
    t.Setenv("EXCHANGE_CHANNELS", "ticker")
    t.Setenv("DB_HOST", "localhost:3306")
    t.Setenv("DB_USER", "test_mysql")
    t.Setenv("DB_BASE", "test")

	// Logger
	t.Setenv("LOGGER_CALLER", "false")
	t.Setenv("LOGGER_STACKTRACE", "true")
	t.Setenv("LOGGER_LEVEL", "debug")

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{name: testing.CoverMode(), args: args{ctx: context.Background()}, want: &Config{
			Exchange: ExchangeConfig{
				Url:      "wss://ws-feed.exchange.coinbase.com",
				Origin:   "https://coinbase.com",
				Protocol: "",
				Symbols:  []string{"ETH-BTC", "BTC-USD", "BTC-EUR"},
				Channels: []string{"ticker"},
			},
			Database: DatabaseConfig{
				Host:     "localhost:3306",
				User:     "test_mysql",
				Password: "a2s_kjlasjd",
				Base:     "test",
			},
			Logger: LoggerConfig{
				DisableCaller:     false,
				DisableStacktrace: true,
				Level:             "debug",
			},
		}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
