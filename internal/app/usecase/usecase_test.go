package usecase

import (
	"reflect"
	"testing"

	"github.com/nel349/bz-findata/internal/app/repository"
)

func TestNewUseCase(t *testing.T) {
	type args struct {
		repos *repository.Repositories
		pkg   *Packages
	}
	tests := []struct {
		name string
		args args
		want *Services
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUseCase(tt.args.repos, tt.args.pkg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
