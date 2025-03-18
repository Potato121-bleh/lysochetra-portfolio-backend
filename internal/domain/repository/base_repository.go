package repository

import "backend/internal/util/dbutil"

// type RepoI interface {
// 	SqlSelect(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error
// 	SqlInsert(tx pgx.Tx, tbName string, colArr []string, valArr []string) error
// 	SqlUpdate(tx pgx.Tx, tbName string, colArr []string, identifier string, valIdentifier string) error
// 	SqlDelete(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error
// }

// This factory are returning UserRepository and using your provided Generic as model.
func NewUserRepository[T dbutil.OnlyStruct]() *UserRepository[T] {
	return &UserRepository[T]{}
}
