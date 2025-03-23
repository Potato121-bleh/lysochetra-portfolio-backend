package auth

import (
	"crypto/x509"
	"encoding/base64"

	"encoding/pem"

	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// Not use anymore we use service instead
func DecodeJWTs(cookie *http.Cookie) (jwt.MapClaims, error) {

	//retrieve the public key
	loadEnvErr := godotenv.Load("../.env")
	if loadEnvErr != nil {
		return nil, errors.New("failed to load env file")
	}

	pemPKformat, pemPKformatErr := base64.StdEncoding.DecodeString(os.Getenv("PUBLIC_KEY"))
	if pemPKformatErr != nil {
		return nil, errors.New("failed to decode base64 cookie")
	}

	pemBlock, _ := pem.Decode(pemPKformat)
	publicKey, parsingPKErr := x509.ParsePKCS1PublicKey(pemBlock.Bytes)
	if parsingPKErr != nil {
		return nil, errors.New("failed to parse PK")
	}

	jwtToken, parseJwtErr := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if pemBlock.Type != "PUBLIC KEY" {
			return nil, errors.New("RSA key type not allowed")
		}
		if token.Method.(*jwt.SigningMethodRSA) != jwt.SigningMethodRS256 {
			return nil, errors.New("algorithm of jwt is not allowed")
		}
		return publicKey, nil
	})
	if parseJwtErr != nil {
		return nil, errors.New(parseJwtErr.Error())
	}

	if !jwtToken.Valid {
		return nil, errors.New("token invalid, can be expired")
	}

	jwtClaim := jwtToken.Claims.(jwt.MapClaims)

	return jwtClaim, nil
}
