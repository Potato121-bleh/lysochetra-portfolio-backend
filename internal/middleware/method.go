package middleware

import (
	"backend/internal/auth"
	"context"
	"crypto/x509"
	"encoding/base64"

	"encoding/json"
	"encoding/pem"
	"fmt"

	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	//"github.com/gorilla/mux"
	"backend/internal/domain/repository"
	"backend/internal/domain/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type UserRequestInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
type userData struct {
	Id        int    `json:"userId"`
	Username  string `json:"userName"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	SettingId int    `json:"settingId"`
}
*/

func MiddlewareValidateAuth(nextHandler http.HandlerFunc, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		/*
			urlVars := mux.Vars(r)
			urlUsername := urlVars["username"]
			urlPassword := urlVars["password"]
		*/

		var reqBody UserRequestInfo
		decodeBodyErr := json.NewDecoder(r.Body).Decode(&reqBody)
		if decodeBodyErr != nil {
			http.Error(w, "Failed to validate your request: "+decodeBodyErr.Error(), http.StatusBadRequest)
			return
		}

		//fetch user from db

		//fetchedRow := db.QueryRow(context.Background(), "SELECT * FROM userauth WHERE LOWER(username) = $1", strings.ToLower(reqBody.Username))

		userRepository := repository.NewRepository("user")
		userService := service.NewService("user", db, userRepository)

		// fetchedRow := db.QueryRow(context.Background(), "SELECT * FROM userauth WHERE LOWER(username) = $1", strings.ToLower(reqBody.Username))

		// var useridTem int
		// var usernameTem string
		// var userpasswordTem string
		// var usernicknameTem string
		// var usersettingidTem int
		// scannedRowErr := fetchedRow.Scan(&useridTem, &usernameTem, &usernicknameTem, &userpasswordTem, nil, &usersettingidTem)
		// if scannedRowErr != nil {
		// 	http.Error(w, "User not found: "+scannedRowErr.Error(), http.StatusUnauthorized)
		// 	return
		// }

		queriedData, queriedDataErr := userService.Select("userauth", "LOWER(username)", strings.ToLower(reqBody.Username))
		if queriedDataErr != nil {
			http.Error(w, "failed to queried user from database", http.StatusInternalServerError)
			return
		}

		// Since we use identifier we have to use the first element
		if queriedData[0].Password != reqBody.Password {
			http.Error(w, "failed to authenticate user", http.StatusUnauthorized)
			return
		}

		//create context to pass the data to actual handler
		newUser := auth.UserData{Id: queriedData[0].Id,
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

		//Get origin from env

		loadingEnvErr := godotenv.Load("../.env")
		if loadingEnvErr != nil {
			http.Error(w, "failed to load env: "+loadingEnvErr.Error(), http.StatusInternalServerError)
			return
		}
		//checking for preflight

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

// X-CSRF-Token
func MiddlewareCSRFCheck(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		jwtClaims, decodeJwtErr := auth.DecodeJWTs(cookie)
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
