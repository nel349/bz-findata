package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
	awslocal "github.com/nel349/bz-findata/pkg/aws"
)

// authentication for web sockets
type Auth struct {
	Key    string
	Secret string
	Passphrase string
}

func NewAuth() *Auth {
	wsApiKey, wsApiSecret, wsApiPassphrase := GetWSCredentials()
	return &Auth{Key: wsApiKey, Secret: wsApiSecret, Passphrase: wsApiPassphrase}
}
// generate signature
func (a *Auth) GenerateSignature() (string, int64, error) {
	timestamp := time.Now().Unix()
	method := "GET"
	requestPath := "/users/self/verify"
	body := ""

	message := fmt.Sprintf("%d%s%s%s", timestamp, method, requestPath, body)

	key, err := base64.StdEncoding.DecodeString(a.Secret)
    if err != nil {
        return "", 0, fmt.Errorf("failed to decode secret: %w", err)
    }
    
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))


	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// fmt.Printf("Debug: Timestamp: %d\n", timestamp)
	// fmt.Printf("Debug: Message: %s\n", message)
	// fmt.Printf("Debug: Signature: %s\n", signature)
	// fmt.Printf("Debug: API Key: %s\n", a.Key)
	// fmt.Printf("Debug: API Secret: %s\n", a.Secret)
	// fmt.Printf("Debug: API Passphrase: %s\n", a.Passphrase)

	return signature, timestamp, nil
}

func GetWSCredentials() (string, string, string) {
    // wsApiKey := os.Getenv("COINBASE_WS_API_KEY")
    // wsApiSecret := os.Getenv("COINBASE_WS_API_SECRET")
    // wsApiPassphrase := os.Getenv("COINBASE_WS_API_PASSPHRASE")

	secret, err := awslocal.GetAwsSecret()
	if err != nil {
		fmt.Println("Failed to retrieve secret", err)
	}

	wsApiKey := secret.COINBASE_WS_API_KEY
	wsApiSecret := secret.COINBASE_WS_API_SECRET
	wsApiPassphrase := secret.COINBASE_WS_API_PASSPHRASE

    // fmt.Printf("Debug: API Key length: %d\n", len(wsApiKey))
    // fmt.Printf("Debug: API Secret length: %d\n", len(wsApiSecret))
    // fmt.Printf("Debug: API Passphrase length: %d\n", len(wsApiPassphrase))

    if wsApiKey == "" || wsApiSecret == "" || wsApiPassphrase == "" {
        fmt.Println("One or more required environment variables are not set")
    }

	return wsApiKey, wsApiSecret, wsApiPassphrase
}
