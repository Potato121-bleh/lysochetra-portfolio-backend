package service

import (
	"backend/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// db   *pgxpool.Pool
// repo repository.UserRepoI
func NewService(serviceName string, db *pgxpool.Pool, repo repository.UserRepoI) UserServiceI {
	switch serviceName {
	case "user":
		return &UserService{
			db:   db,
			repo: repo,
		}
	default:
		return nil
	}
}
