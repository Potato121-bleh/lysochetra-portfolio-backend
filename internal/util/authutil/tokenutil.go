package authutil

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetJWTBlock(key string, keyType string) (*pem.Block, error) {
	loadEnvErrs := godotenv.Load("../.env")
	if loadEnvErrs != nil {
		return nil, fmt.Errorf("failed to load env file")
	}
	base64JWTkey := os.Getenv(key)
	if base64JWTkey == "" {
		return nil, fmt.Errorf("failed to load env file")
	}

	pemJWTkey, decodeBase64Err := base64.StdEncoding.DecodeString(base64JWTkey)
	if decodeBase64Err != nil {
		return nil, fmt.Errorf("failed to decode base64: " + decodeBase64Err.Error())
	}

	block, _ := pem.Decode(pemJWTkey)
	if block.Type != keyType {
		return nil, fmt.Errorf("Sigining key not allowed")
	}

	return block, nil
}
