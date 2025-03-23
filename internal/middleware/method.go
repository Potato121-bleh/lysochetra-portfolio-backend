package middleware

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"profile-portfolio/internal/auth"

	"encoding/json"
	"encoding/pem"
	"fmt"

	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type UserRequestInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func MiddlewareValidateAuth(nextHandler http.HandlerFunc, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqBody UserRequestInfo
		decodeBodyErr := json.NewDecoder(r.Body).Decode(&reqBody)
		if decodeBodyErr != nil {
			http.Error(w, "Failed to validate your request: "+decodeBodyErr.Error(), http.StatusBadRequest)
			return
		}

		userService := service.NewUserService[model.UserData](db)

		queriedData, queriedDataErr := userService.Select(nil, "userauth", "LOWER(username)", strings.ToLower(reqBody.Username))
		if queriedDataErr != nil {
			http.Error(w, "failed to queried user from database: "+queriedDataErr.Error(), http.StatusInternalServerError)
			return
		}

		// Since we use identifier we have to use the first element
		fmt.Println(queriedData)
		fmt.Println("----------")
		fmt.Println(queriedData[0].Password)
		fmt.Println(reqBody.Password)
		if queriedData[0].Password != reqBody.Password {
			http.Error(w, "failed to authenticate user", http.StatusUnauthorized)
			return
		}

		//create context to pass the data to actual handler
		newUser := model.UserData{Id: queriedData[0].Id,
			Username:  queriedData[0].Username,
			Password:  queriedData[0].Password,
			Nickname:  queriedData[0].Nickname,
			SettingId: queriedData[0].SettingId}
		cxtWithData := context.WithValue(r.Context(), auth.NewUserKey, newUser)

		fmt.Println("it passed here")

		nextHandler(w, r.WithContext(cxtWithData))
	}
}

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

func MiddlewareCSRFCheck(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authSvc := service.AuthService{}

		//First we have to check CSRF in header first to check if it not we rejected
		csrfHeaderToken := r.Header.Get("X-CSRF-Token")
		if csrfHeaderToken == "" {
			http.Error(w, "no CSRF token found", http.StatusUnauthorized)
			return
		}

		//Second we have to decode the jwt token to get an actual code inside of it.
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

		//Check whether the two is correct or not
		if csrfHeaderToken != csrfJWTToken {
			http.Error(w, "failed to validate the credential, due to CSRF purposes", http.StatusUnauthorized)
			return
		}

		claimsContext := context.WithValue(r.Context(), auth.JwtClaimsContextKey, jwtClaims)

		nextHandler(w, r.WithContext(claimsContext))

	}
}

func DecodeJWTss(cookie *http.Cookie) (jwt.MapClaims, error) {

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
