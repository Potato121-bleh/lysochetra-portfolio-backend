package service

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"backend/internal/util/dbutil"

	"github.com/jackc/pgx/v5/pgxpool"
)

// type ServiceI interface {
// 	Select(tx pgx.Tx, tbName string, identifier string, valIdentifier string) ([]repository.UserData, error)
// 	Insert(tx pgx.Tx, tbName string, colArr []string, valArr []string) error
// 	Update(tx pgx.Tx, tbName string, colArr []string, identifier string, valIdentitier string) error
// 	Delete(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error
// }

// db   *pgxpool.Pool
// repo repository.UserRepoI

// This method prepare you a well structured UserService where it use the model as you provide
//
//	Note: "This UserService are using UserRepository"
//
// .
func NewUserService[T dbutil.OnlyStruct](db *pgxpool.Pool) *UserService[T] {
	return &UserService[T]{
		db:   db,
		repo: repository.UserRepository[T]{},
	}
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	// userSvc := NewUserService[model.UserData](db)
	// settingSvc := NewUserService[model.SettingStruct](db)
	return &AuthService{
		db:         db,
		userSvc:    NewUserService[model.UserData](db),
		settingSvc: NewUserService[model.SettingStruct](db),
	}
}
