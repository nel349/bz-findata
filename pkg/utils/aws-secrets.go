package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AwsSecret struct {
	COINBASE_WS_API_KEY string
	COINBASE_WS_API_SECRET string
	COINBASE_WS_API_PASSPHRASE string
	SUPABASE_URL string
	SERVICE_ROLE_KEY string
}

func GetAwsSecret() (AwsSecret, error) {
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