package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/nel349/bz-findata/internal/cex-collector/repository"
	"github.com/nel349/bz-findata/pkg/entity"
	"github.com/nel349/bz-findata/pkg/logger"
)

func TestNewExchangeService(t *testing.T) {
	type args struct {
		exchange repository.Exchange
		logger   logger.Logger
	}
	tests := []struct {
		name string
		args args
		want *exchangeService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExchangeService(tt.args.exchange, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExchangeService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_exchangeService_Tick(t *testing.T) {
	type fields struct {
		exchange repository.Exchange
		logger   logger.Logger
	}
	type args struct {
		ctx context.Context
		ch  <-chan entity.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &exchangeService{
				exchange: tt.fields.exchange,
				logger:   tt.fields.logger,
			}
			if err := e.ProcessStream(tt.args.ctx, tt.args.ch); (err != nil) != tt.wantErr {
				t.Errorf("ProcessStream() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
