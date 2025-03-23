package service

import (
	"profile-portfolio/internal/domain/repository"
	"profile-portfolio/internal/util/dbutil"

	"github.com/jackc/pgx/v5/pgxpool"
)

// type ServiceI interface {
// 	Select(tx pgx.Tx, tbName string, identifier string, valIdentifier string) ([]repository.UserData, error)
// 	Insert(tx pgx.Tx, tbName string, colArr []string, valArr []string) error
// 	Update(tx pgx.Tx, tbName string, colArr []string, identifier string, valIdentitier string) error
// 	Delete(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error
// }

func NewUserService[T dbutil.OnlyStruct](db *pgxpool.Pool) *UserService[T] {
	return &UserService[T]{
		db:   db,
		repo: repository.UserRepository[T]{},
	}
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{
		db: db,
	}
}
