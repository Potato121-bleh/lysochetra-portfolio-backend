package main

import (
	"profile-portfolio/internal/application"
	// "profile-portfolio/internal/auth"

	// "profile-portfolio/internal/domain/repository"
	"context"
	"log"
	"net/http"
	"os"
	"profile-portfolio/internal/domain/service"
	"profile-portfolio/internal/middleware"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	cxtTimeout, cancelCxt := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCxt()

	//load in our env file
	loadEnvErr := godotenv.Load("../.env")
	if loadEnvErr != nil {
		log.Fatal("Connection failed: " + loadEnvErr.Error())
	}

	//get connection string
	connstr := os.Getenv("CONNECTION_STRING")
	if connstr == "" {
		log.Fatal("No connstr")
	}

	pgxConfig, pgxConfigErr := pgxpool.ParseConfig(connstr)
	if pgxConfigErr != nil {
		log.Fatal("Connection failed: " + pgxConfigErr.Error())
	}

	pgxConfig.MinConns = 20
	pgxConfig.MaxConns = 80
	pgxConfig.HealthCheckPeriod = 5 * time.Minute

	db, dbConErr := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if dbConErr != nil {
		log.Fatal("failed to connect to database")
	}

	muxhandler := mux.NewRouter()

	userSvc := service.NewUserService(db)
	settingSvc := service.NewSettingService(db)
	authSvc := service.NewAuthService(db)

	authHandlers := &application.AuthHandler{DB: db, CxtTimeout: cxtTimeout, AuthSvc: authSvc, UserSvc: userSvc}
	handleQuerys := &application.SettingHandler{DB: db, UserSvc: userSvc, SettingSvc: settingSvc}

	muxhandler.HandleFunc("/user/auth",
		middleware.MiddlewareCORSValidate(middleware.MiddlewareValidateAuth(authHandlers.HandleSigningToken, db))).Methods("OPTIONS", "POST")

	muxhandler.HandleFunc("/user/signup",
		middleware.MiddlewareCORSValidate(authHandlers.HandleSignup)).Methods("POST", "OPTIONS")

	muxhandler.HandleFunc("/user/auth/logout", middleware.MiddlewareCORSValidate(authHandlers.HandleLogout)).Methods("GET", "OPTIONS")

	muxhandler.HandleFunc("/user/verify",
		middleware.MiddlewareCORSValidate(middleware.MiddlewareCSRFCheck(authHandlers.HandleVerifyToken))).Methods("GET", "OPTIONS")

	muxhandler.HandleFunc("/setting/getbyid",
		middleware.MiddlewareCORSValidate(middleware.MiddlewareCSRFCheck(handleQuerys.HandleQuerySetting))).Methods("POST", "OPTIONS")

	muxhandler.HandleFunc("/setting/update",
		middleware.MiddlewareCORSValidate(middleware.MiddlewareCSRFCheck(handleQuerys.HandleUpdateSetting))).Methods("POST", "OPTIONS")

	muxhandler.HandleFunc("/retrieve-CSRFkey",
		middleware.MiddlewareCORSValidate(authHandlers.RetrieveCSRFKey)).Methods("OPTIONS", "GET")

	muxhandler.HandleFunc("/testing-for-cookie/{username}/{password}",
		middleware.MiddlewareCORSValidate(middleware.MiddlewareValidateAuth(authHandlers.Handletesting, db))).Methods("GET", "OPTIONS")

	http.ListenAndServe(":5000", muxhandler)

}
