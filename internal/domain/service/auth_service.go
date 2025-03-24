package service

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"profile-portfolio/internal/domain/model"
	generator "profile-portfolio/internal/generateKey"
	"profile-portfolio/internal/util/authutil"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthServiceI interface {
	SignUp(tx pgx.Tx, userSvc UserServiceI, reqUsername string, reqNickname string, reqPassword string) error
	SigningToken(userStruct model.UserData) (string, error)
	ParseJwt(cookie *http.Cookie) (jwt.MapClaims, error)
}

type AuthService struct {
	db *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) AuthServiceI {
	return &AuthService{
		db: db,
	}
}

func (s *AuthService) SignUp(tx pgx.Tx, userSvc UserServiceI, reqUsername string, reqNickname string, reqPassword string) error {

	// we insert setting, and query id from setting
	newSettingRow := tx.QueryRow(context.Background(), "INSERT INTO user_setting (darkmode, sound, colorpalettes, font, language) VALUES (0, 0, 0, 1, 1)	RETURNING settingid")
	var latestSettingId int
	queriedSettingIdErr := newSettingRow.Scan(&latestSettingId)
	if queriedSettingIdErr != nil {
		return fmt.Errorf("failed to query settingid")
	}

	// we insert new user with setting id we currently
	fmt.Println([]string{reqUsername, reqNickname, reqPassword, strconv.Itoa(latestSettingId)})
	// insertUserErr := s.userSvc.SqlInsert(
	// 	tx,
	// 	"userauth",
	// 	[]string{"username", "nickname", "password", "setting_id"},
	// 	[]string{reqUsername, reqNickname, reqPassword, strconv.Itoa(latestSettingId)},
	// )
	insertUserErr := userSvc.Insert(
		tx,
		"userauth",
		[]string{"username", "nickname", "password", "setting_id"},
		[]string{reqUsername, reqNickname, reqPassword, strconv.Itoa(latestSettingId)},
	)
	if insertUserErr != nil {
		return fmt.Errorf("failed to insert userauth (%v)", insertUserErr.Error())
	}

	return nil

}

func (s *AuthService) SigningToken(userStruct model.UserData) (string, error) {
	block, parseBlockErr := authutil.GetJWTBlock("PRIVATE_KEY", "RSA PRIVATE KEY")
	if parseBlockErr != nil {
		// http.Error(w, "failed to parsing key: "+parsePrivateKeyErr.Error(), http.StatusInternalServerError)
		return "", fmt.Errorf("failed to parsing key: %v", parseBlockErr.Error())
	}

	privateKey, parseJWTkeyErr := x509.ParsePKCS1PrivateKey(block.Bytes)
	if parseJWTkeyErr != nil {
		return "", fmt.Errorf("failed to parsing key: " + parseJWTkeyErr.Error())
	}

	genCSRFKey, genCSRFKeyErr := generator.GenerateCSRFKey()
	if genCSRFKeyErr != nil {
		// http.Error(w, "failed to gen key: "+genCSRFKeyErr.Error(), http.StatusInternalServerError)
		return "", fmt.Errorf("failed to gen key: %v", genCSRFKeyErr.Error())
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
		return "", fmt.Errorf("failed to signing jwt token: %v", signingJwtErr.Error())
	}

	return jwtSignature, nil
}

func (s *AuthService) ParseJwt(cookie *http.Cookie) (jwt.MapClaims, error) {
	jwtSign := cookie.Value
	if jwtSign == "" {
		// http.Error(w, "token not found", http.StatusUnauthorized)
		return nil, fmt.Errorf("token not found")
	}

	jwtblock, parseJwtBlockErr := authutil.GetJWTBlock("PUBLIC_KEY", "PUBLIC KEY")
	if parseJwtBlockErr != nil {
		// http.Error(w, "failed to retrieve public key: "+parseJwtBlockErr.Error(), http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to retrieve public key: " + parseJwtBlockErr.Error())
	}

	publicKey, x509ParseJwtErr := x509.ParsePKCS1PublicKey(jwtblock.Bytes)
	if x509ParseJwtErr != nil {
		// http.Error(w, "failed to retrieve public key", http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to retrieve public key: " + x509ParseJwtErr.Error())
	}

	jwtToken, jwtParsingErr := jwt.Parse(jwtSign, func(t *jwt.Token) (interface{}, error) {
		if t.Method.(*jwt.SigningMethodRSA) != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("algorithm of jwt is not allowed")
		}
		return publicKey, nil
	})
	if jwtParsingErr != nil {
		// http.Error(w, "failed to parse jwt: "+jwtParsingErr.Error(), http.StatusUnauthorized)
		return nil, fmt.Errorf("failed to parse jwt: " + jwtParsingErr.Error())
	}

	if !jwtToken.Valid {
		return nil, fmt.Errorf("token invalid, can be expired")
	}

	jwtClaims, claimsOk := jwtToken.Claims.(jwt.MapClaims)
	if !claimsOk {
		// http.Error(w, "failed to parse jwt", http.StatusUnauthorized)
		return nil, fmt.Errorf("failed to parse jwt claims")
	}

	return jwtClaims, nil

}
