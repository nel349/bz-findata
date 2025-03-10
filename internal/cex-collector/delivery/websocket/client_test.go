package websocket

import (
	"reflect"
	"testing"

	"github.com/nel349/bz-findata/config"
	"github.com/nel349/bz-findata/internal/cex-collector/usecase"
	"github.com/nel349/bz-findata/pkg/exchange"
	"github.com/nel349/bz-findata/pkg/logger"
)

func TestNewSocketClient(t *testing.T) {
	type args struct {
		conn   exchange.Manager
		uc     *usecase.Services
		logger logger.Logger
		cfg    config.ExchangeConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *client
		wantErr bool
	}{
		{name: testing.CoverMode(), args: args{conn: nil, uc: nil, logger: nil, cfg: config.ExchangeConfig{}}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSocketClient(tt.args.conn, tt.args.uc, tt.args.logger, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSocketClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSocketClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}
