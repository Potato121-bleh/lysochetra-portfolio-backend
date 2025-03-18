package auth

import (
	"backend/internal/domain/model"
	"backend/internal/domain/service"
	generator "backend/internal/generateKey"
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var NewUserKey ContextKey = "userinfo"

type ContextKey string

type UserRequestInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// type UserData struct {
// 	Id             int       `json:"userId"`
// 	Username       string    `json:"userName"`
// 	Nickname       string    `json:"nickname"`
// 	Password       string    `json:"password"`
// 	RegisteredDate time.Time `json:"registered_date"`
// 	SettingId      int       `json:"settingId"`
// }

type SignUpUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type AuthHandler struct {
	DB         *pgxpool.Pool
	CxtTimeout context.Context
	Svc        *service.AuthService
}

func (s *AuthHandler) Handletesting(w http.ResponseWriter, r *http.Request) {
	userStruct := r.Context().Value(NewUserKey).(model.UserData)
	testprop := "your name is: " + userStruct.Nickname + "and your username is: " + userStruct.Username
	newCookie := http.Cookie{
		Name:     "auth_token",
		Value:    testprop,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Domain:   "localhost",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &newCookie)
	w.Write([]byte("Well, here your cookie"))
}

func (s *AuthHandler) HandleSigningToken(w http.ResponseWriter, r *http.Request) {

	fmt.Println("entered 1")
	userStruct := r.Context().Value(NewUserKey).(model.UserData)
	fmt.Println("entered 2")
	//get key
	loadEnvErrs := godotenv.Load("../.env")
	if loadEnvErrs != nil {
		http.Error(w, "failed to load env file", http.StatusInternalServerError)
		return
	}
	base64PrivateKey := os.Getenv("PRIVATE_KEY")
	if base64PrivateKey == "" {
		http.Error(w, "failed to load env file", http.StatusInternalServerError)
		return
	}

	pemPrivateKey, decodeBase64Err := base64.StdEncoding.DecodeString(base64PrivateKey)
	if decodeBase64Err != nil {
		http.Error(w, "failed to decode base64: "+decodeBase64Err.Error(), http.StatusInternalServerError)
		return
	}

	block, _ := pem.Decode(pemPrivateKey)
	if block.Type != "RSA PRIVATE KEY" {
		http.Error(w, "Sigining key not allowed", http.StatusInternalServerError)
		return
	}

	privateKey, parsePrivateKeyErr := x509.ParsePKCS1PrivateKey(block.Bytes)
	if parsePrivateKeyErr != nil {
		http.Error(w, "failed to parsing key: "+parsePrivateKeyErr.Error(), http.StatusInternalServerError)
		return
	}

	genCSRFKey, genCSRFKeyErr := generator.GenerateCSRFKey()
	if genCSRFKeyErr != nil {
		http.Error(w, "failed to gen key: "+genCSRFKeyErr.Error(), http.StatusInternalServerError)
		return
	}

	jwtClaim := jwt.MapClaims{
		"Id":        userStruct.Id,
		"Username":  userStruct.Username,
		"Password":  userStruct.Password,
		"Nickname":  userStruct.Nickname,
		"SettingId": userStruct.SettingId,
		"CSRFKey":   genCSRFKey,
		"exp":       jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		"iat":       jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaim)
	jwtSignature, signingJwtErr := token.SignedString(privateKey)
	if signingJwtErr != nil {
		http.Error(w, "failed to signing jwt token: "+signingJwtErr.Error(), http.StatusUnauthorized)
		return
	}

	//domainCookie := os.Getenv("COOKIE_DOMAIN")

	newCookie := http.Cookie{
		Name:     "auth_token",
		Value:    jwtSignature,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &newCookie)
	w.Write([]byte("User Authenticated"))
}

func (s *AuthHandler) HandleVerifyToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, retrieveTokenErr := r.Cookie("auth_token")
	if retrieveTokenErr != nil {
		http.Error(w, "failed to retrieve token: "+retrieveTokenErr.Error(), http.StatusBadRequest)
		return
	}

	claims, retrieveClaimsErr := DecodeJWTs(cookie)
	if retrieveClaimsErr != nil {
		http.Error(w, "failed with decodeJWTs func: "+retrieveClaimsErr.Error(), http.StatusUnauthorized)
		return
	}

	responseClient := model.UserData{
		Id:        int(claims["Id"].(float64)),
		Username:  claims["Username"].(string),
		Nickname:  claims["Nickname"].(string),
		Password:  claims["Password"].(string),
		SettingId: int(claims["SettingId"].(float64)),
	}

	fmt.Println("it already passed data response")

	encodeRespErr := json.NewEncoder(w).Encode(responseClient)
	if encodeRespErr != nil {
		http.Error(w, "failed to decode response: "+encodeRespErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *AuthHandler) RetrieveCSRFKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, retrieveCookieErr := r.Cookie("auth_token")
	if retrieveCookieErr != nil {
		http.Error(w, "cookie not found: "+retrieveCookieErr.Error(), http.StatusBadRequest)
		return
	}

	jwtSign := cookie.Value
	if jwtSign == "" {
		http.Error(w, "token not found", http.StatusUnauthorized)
		return
	}

	//Get public key
	loadEnvErr := godotenv.Load("../.env")
	if loadEnvErr != nil {
		http.Error(w, "failed to navigate to env: "+loadEnvErr.Error(), http.StatusInternalServerError)
		return
	}
	base64PublicKey := os.Getenv("PUBLIC_KEY")
	pemPublicKey, convertPemPKErr := base64.StdEncoding.DecodeString(base64PublicKey)
	if convertPemPKErr != nil {
		http.Error(w, "convert base64 failed", http.StatusInternalServerError)
		return
	}

	jwtblock, _ := pem.Decode(pemPublicKey)

	jwtToken, jwtParsingErr := jwt.Parse(jwtSign, func(t *jwt.Token) (interface{}, error) {
		if jwtblock.Type != "PUBLIC KEY" {
			//http.Error("failed to retrieve public key", http.StatusInternalServerError)
			return nil, errors.New("failed to retrieve public key")
		}
		publicKey, x509ParseJwtErr := x509.ParsePKCS1PublicKey(jwtblock.Bytes)
		if x509ParseJwtErr != nil {
			return nil, errors.New("failed to retrieve public key")
		}
		return publicKey, nil
	})
	if jwtParsingErr != nil {
		http.Error(w, "failed to parse jwt: "+jwtParsingErr.Error(), http.StatusUnauthorized)
		return
	}

	jwtClaims, claimsOk := jwtToken.Claims.(jwt.MapClaims)
	if !claimsOk {
		http.Error(w, "failed to parse jwt", http.StatusUnauthorized)
		return
	}

	csrfKey := jwtClaims["CSRFKey"].(string)
	_, writeRespErr := fmt.Fprint(w, csrfKey)
	if writeRespErr != nil {
		http.Error(w, "failed to response", http.StatusInternalServerError)
		return
	}
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

		// userRepo := repository.NewRepository("user", UserData{})
		userService := service.NewUserService[model.UserData](db)
		queriedUser, queriedErr := userService.Select(nil, "userauth", "LOWER(username)", strings.ToLower(reqBody.Username))
		if queriedErr != nil || len(queriedUser) != 1 {
			http.Error(w, "failed to queried user from db", http.StatusUnauthorized)
			return
		}

		if queriedUser[0].Password != reqBody.Password {
			http.Error(w, "failed to authenticate user", http.StatusUnauthorized)
			return
		}

		//create context to pass the data to actual handler
		newUser := model.UserData{Id: queriedUser[0].Id,
			Username:  queriedUser[0].Username,
			Password:  queriedUser[0].Password,
			Nickname:  queriedUser[0].Nickname,
			SettingId: queriedUser[0].SettingId,
		}
		cxtWithData := context.WithValue(r.Context(), NewUserKey, newUser)

		nextHandler(w, r.WithContext(cxtWithData))
	}
}

func (s *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {

	// repoModel := UserData{}
	// userRepo := repository.NewRepository("user", repoModel)
	// userService := service.NewUserService[UserData](s.DB)

	//Validate the CORS
	reqBody := SignUpUser{}
	decodeBodyErr := json.NewDecoder(r.Body).Decode(&reqBody)
	if decodeBodyErr != nil {
		http.Error(w, "failed to read request body: "+decodeBodyErr.Error(), http.StatusBadRequest)
		return
	}

	//Begin transaction
	tx, startTxErr := s.DB.Begin(context.Background())
	if startTxErr != nil {
		http.Error(w, "failed to start transaction: "+startTxErr.Error(), http.StatusInternalServerError)
		return
	}

	if reqBody.Username == "" || reqBody.Password == "" || reqBody.Nickname == "" {
		http.Error(w, "failed to validate user request", http.StatusBadRequest)
		return
	}

	//Check existing user
	validateUsernameRow := s.DB.QueryRow(context.Background(), "SELECT userid FROM userauth WHERE username = $1", reqBody.Username)

	var validateUsernameVar int
	scanUsernameErr := validateUsernameRow.Scan(&validateUsernameVar)
	if scanUsernameErr == nil || validateUsernameVar != 0 {
		http.Error(w, "the user already exist", http.StatusBadRequest)
		return
	}

	// As we want to use Signup which not available in original interface we have to come up with type assertion
	// which allow us to use the extra method outside of assign interface
	signUpErr := s.Svc.SignUp(tx, reqBody.Username, reqBody.Nickname, reqBody.Password)
	if signUpErr != nil {
		rollbackErr := tx.Rollback(context.Background())
		if rollbackErr != nil {
			http.Error(w, fmt.Sprintf("failed to sign user to database & failed to rollback: %v", signUpErr.Error()), http.StatusInternalServerError)
			return
		}

		http.Error(w, fmt.Sprintf("failed to sign user to database: %v", signUpErr.Error()), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User Created Successfully"))
	commitTranErr := tx.Commit(context.Background())
	if commitTranErr != nil {
		tx.Rollback(context.Background())
		http.Error(w, "failed create new user (4): "+commitTranErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Signout Successfully"))
}
