package supabase

import (
	"context"
	"fmt"
	"os"

	// "time"
	"github.com/nel349/bz-findata/pkg/entity"
	"github.com/supabase-community/supabase-go"
)

type supabaseRepo struct {
	Client *supabase.Client
}

// // NewSupabaseRepository created supabase repository
func NewSupabaseRepository() *supabaseRepo {
	projectURL := os.Getenv("SUPABASE_URL")
	serviceRoleKey := os.Getenv("SERVICE_ROLE_KEY")

	// Lets check if the projectURL and anonKey are set
	if projectURL == "" || serviceRoleKey == "" {
		fmt.Println("SUPABASE_URL or SERVICE_ROLE_KEY is not set")
	}

	client, err := supabase.NewClient(projectURL, serviceRoleKey, &supabase.ClientOptions{
		Headers: map[string]string{
			"Authorization": "Bearer " + serviceRoleKey,
			"apikey":        serviceRoleKey,
		},
	})
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
	return &supabaseRepo{client}
}

// Create the order in supabase
func (s *supabaseRepo) CreateOrder(ctx context.Context, message entity.Message) error {
	fmt.Println("Attempting to insert order to supabase", "order", message.Order)

	var insertedOrder []entity.Order
	if message.Order != nil {
		_, err := s.Client.From("orders").Insert(message.Order, false, "", "", "").ExecuteTo(&insertedOrder)
		if err != nil {
			return err
		}
		fmt.Println("Order inserted successfully", "insertedOrder", insertedOrder)
	}

	return nil
}
