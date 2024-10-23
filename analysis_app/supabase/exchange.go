package supabase

import (
	// "context"
	// "fmt"
	// "os"

	// "time"
	// "github.com/nel349/bz-findata/pkg/entity"
	// "github.com/supabase-community/supabase-go"
)

// type supabaseRepo struct {
// 	supabase *supabase.Client
// }

// // NewSupabaseRepository created supabase repository
// func NewSupabaseRepository() *supabaseRepo {
// 	projectURL := os.Getenv("SUPABASE_URL")
// 	anonKey := os.Getenv("SUPABASE_ANON_KEY")

// 	client, err := supabase.NewClient(projectURL, anonKey, &supabase.ClientOptions{})
// 	if err != nil {
// 		fmt.Println("cannot initalize client", err)
// 	}
// 	return &supabaseRepo{client}
// }

// Create the order in supabase
// func (s *supabaseRepo) CreateOrder(ctx context.Context, message entity.Message) error {
// 	fmt.Println("Attempting to insert order to supabase", "order", message.Order)

// 	var insertedOrder []entity.Order
// 	if message.Order != nil {
// 		_, err := s.supabase.From("orders").Insert(message.Order, false, "", "", "").ExecuteTo(&insertedOrder)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("Order inserted successfully", "insertedOrder", insertedOrder)
// 	}

// 	return nil
// }
