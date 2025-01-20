package dex

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nel349/bz-findata/pkg/entity"
	"github.com/supabase-community/supabase-go"
)

type Service struct {
	db *sqlx.DB
	supabaseClient *supabase.Client
}

func NewService(db *sqlx.DB, supabaseClient *supabase.Client) *Service {
	return &Service{
		db: db,
		supabaseClient: supabaseClient,
	}
}


// Get the largest swaps in last N hours by Value
func (s *Service) GetLargestSwapsInLastNHours(
	ctx context.Context,
	hours,
	limit int,
)([] entity.SwapTransaction, error) {

	query := `
		SELECT * FROM swap_transactions
		WHERE last_updated > FROM_UNIXTIME(?)
		ORDER BY value DESC
		LIMIT ?
	`
	var swaps []entity.SwapTransaction
	// Log the query parameters for debugging
	log.Printf("Executing query with hours: %d, limit: %d", hours, limit)
	// Convert the timestamp to seconds for FROM_UNIXTIME
	err := s.db.SelectContext(ctx, &swaps, query, time.Now().Add(-time.Duration(hours)*time.Hour).Unix(), limit)
	if err != nil {
		log.Println("error selecting swaps from db", err)
		return nil, err
	}
	// Log the number of swaps found
	log.Printf("Number of swaps found: %d", len(swaps))
	return swaps, nil
}

// Store to supabase StoreLargestSwapsInLastNHours
func (s *Service) StoreLargestSwapsInLastNHours(ctx context.Context, hours, limit int) error {
	swaps, err := s.GetLargestSwapsInLastNHours(ctx, hours, limit)
	if err != nil {
		log.Println("error getting largest swaps", err)
		return err
	}
	return s.storeSwapsInSupabase(swaps, "swap_transactions")
}

func (s *Service) storeSwapsInSupabase(swaps interface{}, tableName string) error {
	_, err := s.supabaseClient.From(tableName).Insert(swaps, false, "", "", "").ExecuteTo(&swaps)
	if err != nil {
		log.Println("error inserting swaps to supabase table", err)
		return err
	}
	return nil
}



