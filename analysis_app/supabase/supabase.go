package supabase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/nel349/bz-findata/pkg/entity"
	"github.com/supabase-community/supabase-go"
)

type supabaseRepo struct {
	Client *supabase.Client
}

type AwsSecret struct {
	COINBASE_WS_API_KEY string
	COINBASE_WS_API_SECRET string
	COINBASE_WS_API_PASSPHRASE string
	SUPABASE_URL string
	SERVICE_ROLE_KEY string
}

func getAwsSecret() (AwsSecret, error) {
	secretName := "prod/supabase/coinbase"
	region := "us-east-2"

	// Add timeout to context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if we have environment credentials
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		fmt.Println("AWS credentials not found in environment")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion(region),
        // config.WithClientLogMode(aws.LogSigning|aws.LogRetries),
    )

	if err != nil {
		return AwsSecret{}, fmt.Errorf("unable to load SDK config: %w", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return AwsSecret{}, err
	}

	// fmt.Println("Secret retrieved successfully", "secret", *result.SecretString)

	var secret AwsSecret
	err = json.Unmarshal([]byte(*result.SecretString), &secret)
	if err != nil {
		return AwsSecret{}, err
	}

	return secret, nil
}

// // NewSupabaseRepository created supabase repository
func NewSupabaseRepository() *supabaseRepo {
	// projectURL := os.Getenv("SUPABASE_URL")
	// serviceRoleKey := os.Getenv("SERVICE_ROLE_KEY")

	secret, err := getAwsSecret()
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
