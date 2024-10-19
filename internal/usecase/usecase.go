package usecase

import (
	"context"

	"github.com/nel349/bz-findata/internal/entity"
	"github.com/nel349/bz-findata/internal/repository"
	"github.com/nel349/bz-findata/pkg/logger"
)

// Exchange usecase
type Exchange interface {
	// ProcessStream handles the stream of exchange data
	ProcessStream(ctx context.Context, ch <-chan entity.Message) error
}

// Services struct of usecase layout
type Services struct {
	Exchange
}

// Packages struct of usecase packages
type Packages struct {
	Logger logger.Logger
}

// NewUseCase create usecase layout
func NewUseCase(repos *repository.Repositories, pkg *Packages) *Services {
	return &Services{
		Exchange: NewExchangeService(repos.Exchange, pkg.Logger),
	}
}
