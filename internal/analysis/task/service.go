package task

import (
	"context"
	"log"

	"github.com/nel349/bz-findata/internal/analysis/orders"
)

type Service struct {
	analysisService *analysis.Service
}

func NewService(service *analysis.Service) *Service {
	return &Service{
		analysisService: service,
	}
}

func (s *Service) StoreMatchOrders(ctx context.Context, hours, limit int) error {
	log.Printf("Starting task: StoreMatchOrders with hours=%d, limit=%d", hours, limit)
	err := s.analysisService.StoreMatchOrdersInSupabase(ctx, hours, limit)
	if err != nil {
		log.Printf("Error executing StoreMatchOrders task: %v", err)
		return err
	}
	log.Printf("Successfully completed StoreMatchOrders task")
	return nil
}
