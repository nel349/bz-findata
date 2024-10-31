package supabase

import (
	"context"
	"fmt"

	// "github.com/aws/aws-sdk-go-v2/aws"
	awslocal "github.com/nel349/bz-findata/pkg/aws"
	"github.com/nel349/bz-findata/pkg/entity"
	"github.com/supabase-community/supabase-go"
)

type supabaseRepo struct {
	Client *supabase.Client
}

// // NewSupabaseRepository created supabase repository
func NewSupabaseRepository() *supabaseRepo {
	// projectURL := os.Getenv("SUPABASE_URL")
	// serviceRoleKey := os.Getenv("SERVICE_ROLE_KEY")

	secret, err := awslocal.GetAwsSecret()
	if err != nil {
		fmt.Println("Failed to retrieve secret", err)
	}

	fmt.Println("Secret retrieved successfully: ")
	// fmt.Println("COINBASE_WS_API_KEY", secret.COINBASE_WS_API_KEY)
	// fmt.Println("COINBASE_WS_API_SECRET", secret.COINBASE_WS_API_SECRET)
	// fmt.Println("COINBASE_WS_API_PASSPHRASE", secret.COINBASE_WS_API_PASSPHRASE)
	// fmt.Println("SUPABASE_URL", secret.SUPABASE_URL)
	// fmt.Println("SERVICE_ROLE_KEY", secret.SERVICE_ROLE_KEY)


	// Lets check if the projectURL and anonKey are set
	if secret.SUPABASE_URL == "" || secret.SERVICE_ROLE_KEY == "" {
		fmt.Println("SUPABASE_URL or SERVICE_ROLE_KEY is not set")
	}

	client, err := supabase.NewClient(secret.SUPABASE_URL, secret.SERVICE_ROLE_KEY, &supabase.ClientOptions{
		Headers: map[string]string{
			"Authorization": "Bearer " + secret.SERVICE_ROLE_KEY,
			"apikey":        secret.SERVICE_ROLE_KEY,
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
