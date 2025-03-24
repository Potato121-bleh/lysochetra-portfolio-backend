package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/service"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var NewUserKey ContextKey = "userinfo"

type ContextKey string

type UserRequestInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type AuthHandler struct {
	DB         *pgxpool.Pool
	CxtTimeout context.Context
	AuthSvc    service.AuthServiceI
	UserSvc    service.UserServiceI
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

	userStruct := r.Context().Value(NewUserKey).(model.UserData)
	jwtSignature, signingJwtErr := s.AuthSvc.SigningToken(userStruct)
	if signingJwtErr != nil {
		http.Error(w, "failed to signing jwt token: "+signingJwtErr.Error(), http.StatusUnauthorized)
		return
	}

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

	claims, retrieveClaimsErr := s.AuthSvc.ParseJwt(cookie)
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

	jwtClaims, parseJwtErr := s.AuthSvc.ParseJwt(cookie)
	if parseJwtErr != nil {
		http.Error(w, "failed to parse jwt: "+parseJwtErr.Error(), http.StatusUnauthorized)
		return
	}

	csrfKey := jwtClaims["CSRFKey"].(string)
	_, writeRespErr := fmt.Fprint(w, csrfKey)
	if writeRespErr != nil {
		http.Error(w, "failed to response", http.StatusInternalServerError)
		return
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

	// As we want to use Signup which not available in original interface we have to come up with type assertion
	// which allow us to use the extra method outside of assign interface
	signUpErr := s.AuthSvc.SignUp(tx, s.UserSvc, reqBody.Username, reqBody.Nickname, reqBody.Password)
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
