package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
)

func TestNewDexExchangeRepository(t *testing.T) {
	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want *dexExchangeRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDexExchangeRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDexExchangeRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dexExchangeRepo_SaveSwap(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx     context.Context
		tx      *types.Transaction
		version string
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
			e := &dexExchangeRepo{
				db: tt.fields.db,
			}
			if err := e.SaveSwap(tt.args.ctx, tt.args.tx, tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("SaveSwap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
