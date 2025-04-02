package middleware

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"encoding/json"
	"encoding/pem"

	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"profile-portfolio/internal/application"
	"profile-portfolio/internal/db"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/service"

	"github.com/joho/godotenv"
)

func MiddlewareCORSValidate(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		loadingEnvErr := godotenv.Load("../.env")
		if loadingEnvErr != nil {
			http.Error(w, "failed to load env: "+loadingEnvErr.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGIN"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		nextHandler(w, r)
	}
}

func MiddlewareValidateAuth(nextHandler http.HandlerFunc, db db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqBody model.UserData
		decodeBodyErr := json.NewDecoder(r.Body).Decode(&reqBody)
		if decodeBodyErr != nil {
			http.Error(w, "Failed to validate your request: "+decodeBodyErr.Error(), http.StatusBadRequest)
			return
		}

		userService := service.NewUserService(db)
		queriedUser, queriedErr := userService.Select(nil, "userauth", "LOWER(username)", strings.ToLower(reqBody.Username))
		if queriedErr != nil || len(queriedUser) != 1 {
			http.Error(w, "failed to queried user from db", http.StatusUnauthorized)
			return
		}

		if queriedUser[0].Password != reqBody.Password {
			http.Error(w, "failed to authenticate user", http.StatusUnauthorized)
			return
		}

		newUser := model.UserData{Id: queriedUser[0].Id,
			Username:  queriedUser[0].Username,
			Password:  queriedUser[0].Password,
			Nickname:  queriedUser[0].Nickname,
			SettingId: queriedUser[0].SettingId,
		}
		cxtWithData := context.WithValue(r.Context(), application.NewUserKey, newUser)

		nextHandler(w, r.WithContext(cxtWithData))
	}
}

func MiddlewareCSRFCheck(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authSvc := service.AuthService{}

		csrfHeaderToken := r.Header.Get("X-CSRF-Token")
		if csrfHeaderToken == "" {
			http.Error(w, "no CSRF token found", http.StatusUnauthorized)
			return
		}

		cookie, retrieveTokenErr := r.Cookie("auth_token")
		if retrieveTokenErr != nil {
			http.Error(w, "failed to retrieve the cookie", http.StatusUnauthorized)
			return
		}

		jwtClaims, decodeJwtErr := authSvc.ParseJwt(cookie)
		if decodeJwtErr != nil {
			http.Error(w, "failed to decode jwt token: "+decodeJwtErr.Error(), http.StatusUnauthorized)
			return
		}

		csrfJWTToken := jwtClaims["CSRFKey"].(string)

		if csrfHeaderToken != csrfJWTToken {
			http.Error(w, "failed to validate the credential, due to CSRF purposes", http.StatusUnauthorized)
			return
		}

		claimsContext := context.WithValue(r.Context(), application.JwtClaimsContextKey, jwtClaims)

		nextHandler(w, r.WithContext(claimsContext))

	}
}

func DecodeJWTss(cookie *http.Cookie) (jwt.MapClaims, error) {

	loadEnvErr := godotenv.Load("../.env")
	if loadEnvErr != nil {
		return nil, fmt.Errorf("failed to load env file")
	}

	pemPKformat, pemPKformatErr := base64.StdEncoding.DecodeString(os.Getenv("PUBLIC_KEY"))
	if pemPKformatErr != nil {
		return nil, fmt.Errorf("failed to decode base64 cookie")
	}

	pemBlock, _ := pem.Decode(pemPKformat)
	publicKey, parsingPKErr := x509.ParsePKCS1PublicKey(pemBlock.Bytes)
	if parsingPKErr != nil {
		return nil, fmt.Errorf("failed to parse PK")
	}

	jwtToken, parseJwtErr := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if pemBlock.Type != "PUBLIC KEY" {
			return nil, fmt.Errorf("RSA key type not allowed")
		}
		if token.Method.(*jwt.SigningMethodRSA) != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("algorithm of jwt is not allowed")
		}
		return publicKey, nil
	})
	if parseJwtErr != nil {
		return nil, fmt.Errorf(parseJwtErr.Error())
	}

	if !jwtToken.Valid {
		return nil, fmt.Errorf("token invalid, can be expired")
	}

	jwtClaim := jwtToken.Claims.(jwt.MapClaims)

	return jwtClaim, nil
}
