package auth

import (
	"backend/internal/domain/repository"
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

type UserData struct {
	Id        int    `json:"userId"`
	Username  string `json:"userName"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	SettingId int    `json:"settingId"`
}

type SignUpUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type AuthHandler struct {
	DB         *pgxpool.Pool
	CxtTimeout context.Context
}

func (s *AuthHandler) Handletesting(w http.ResponseWriter, r *http.Request) {
	userStruct := r.Context().Value(NewUserKey).(UserData)
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
	userStruct := r.Context().Value(NewUserKey).(UserData)
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

	responseClient := UserData{
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

		//fetch user from db
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

		userRepo := repository.NewRepository("user")
		userService := service.NewService("user", db, userRepo)
		queriedUser, queriedErr := userService.Select("userauth", "LOWER(username)", strings.ToLower(reqBody.Username))
		if queriedErr != nil || len(queriedUser) != 1 {
			http.Error(w, "failed to queried user from db", http.StatusUnauthorized)
			return
		}

		if queriedUser[0].Password != reqBody.Password {
			http.Error(w, "failed to authenticate user", http.StatusUnauthorized)
			return
		}

		//create context to pass the data to actual handler
		newUser := UserData{Id: queriedUser[0].Id,
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

	//We insert user data
	insertUserCommandTag, insertUserErr := tx.Exec(
		context.Background(),
		"INSERT INTO userauth (username, nickname, password) VALUES ($1, $2, $3)",
		reqBody.Username,
		reqBody.Nickname,
		reqBody.Password,
	)

	if insertUserErr != nil || insertUserCommandTag.RowsAffected() != 1 {
		tx.Rollback(context.Background())
		http.Error(w, "failed to add user into database (01): "+insertUserErr.Error(), http.StatusInternalServerError)
		return
	}

	//we insert setting for user
	insertSettingCommandTag, insertSettingErr := tx.Exec(context.Background(),
		"INSERT INTO user_setting (darkmode, sound, colorpalettes, font, language) VALUES (0, 0, 0, 1, 1)",
	)
	if insertSettingErr != nil || insertSettingCommandTag.RowsAffected() != 1 {
		tx.Rollback(context.Background())
		http.Error(w, "failed to add user into database", http.StatusInternalServerError)
		return
	}

	//we retrieve userid via username
	userIdRow := tx.QueryRow(context.Background(), "SELECT userid FROM userauth WHERE username = $1", reqBody.Username)

	var recentUserId int
	scanUserIdErr := userIdRow.Scan(&recentUserId)
	if scanUserIdErr != nil {
		tx.Rollback(context.Background())
		http.Error(w, "failed create new user (1): "+scanUserIdErr.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(recentUserId)

	//we retrieve user id from setting and checking for (safety check)
	checkIdRow := tx.QueryRow(context.Background(), "SELECT settingid FROM user_setting WHERE settingid = $1", recentUserId)

	var settingIdValidate int
	idValidationErr := checkIdRow.Scan(&settingIdValidate)
	if idValidationErr != nil {
		tx.Rollback(context.Background())
		http.Error(w, "failed create new user (2): "+idValidationErr.Error(), http.StatusInternalServerError)
		return
	}

	//we can update the references
	updateSettingCommandTag, updateSettingidErr := tx.Exec(context.Background(),
		"UPDATE userauth SET setting_id = $1 WHERE userid = $2",
		recentUserId,
		recentUserId,
	)
	if updateSettingidErr != nil || updateSettingCommandTag.RowsAffected() != 1 {
		tx.Rollback(context.Background())
		http.Error(w, "failed create new user (3): "+updateSettingidErr.Error(), http.StatusInternalServerError)
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
