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
	COINBASE_WS_API_KEY        string
	COINBASE_WS_API_SECRET     string
	COINBASE_WS_API_PASSPHRASE string
	SUPABASE_URL               string
	SERVICE_ROLE_KEY           string
}

type AwsDBSecret struct {
	DB_HOST     string
	DB_USER     string `json:"username"`
	DB_PASSWORD string `json:"password"`
	DB_BASE     string
}

func GetAwsSecret[T any](secretName string, region string) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		fmt.Println("AWS credentials not found in environment")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	if err != nil {
		var empty T
		return empty, fmt.Errorf("unable to load SDK config: %w", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		var empty T
		return empty, err
	}

	var secret T
	err = json.Unmarshal([]byte(*result.SecretString), &secret)
	if err != nil {
		var empty T
		return empty, err
	}

	return secret, nil
}

func GetDefaultCoinbaseSecret() (AwsSecret, error) {
	secretName := "prod/supabase/coinbase"
	region := "us-east-2"
	return GetAwsSecret[AwsSecret](secretName, region)
}

func GetDefaultDBSecret() (AwsDBSecret, error) {
	secretName := "rds!db-78e7999e-5cdb-40a8-9a31-f5b15afbc492"
	region := "us-east-2"
	return GetAwsSecret[AwsDBSecret](secretName, region)
}
