package service

import (
	"profile-portfolio/internal/db"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/repository"
	"profile-portfolio/internal/util/dbutil"

	"context"
	"fmt"
)

// SignUp(reqUsername string, reqNickname string, reqPassword string) error
type UserServiceI interface {
	Select(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) ([]model.UserData, error)
	Insert(tx db.DatabaseTx, tbName string, colArr []string, valArr []string) error
	Update(tx db.DatabaseTx, tbName string, colArr []string, colVal []string, identifier string, valIdentitier string) error
	Delete(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) error
}

type UserService struct {
	db   db.Database
	repo repository.UserRepository
}

func NewUserService(db db.Database) UserServiceI {
	return &UserService{
		db:   db,
		repo: repository.UserRepository{},
	}
}

func (s *UserService) Select(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) ([]model.UserData, error) {
	if tbName == "" {
		return nil, fmt.Errorf("Please provide all required data")
	}

	cxt := context.Background()
	newTx := dbutil.PrepTx(tx, s.db, cxt)
	queriedData, queriedErr := s.repo.SqlSelect(newTx, tbName, identifier, valIdentifier)
	if queriedErr != nil {
		rollbackErr := newTx.Rollback(context.Background())
		if rollbackErr != nil {
			return nil, fmt.Errorf("failed to perform rollback (%v)", queriedErr.Error())
		}
		return nil, fmt.Errorf("failed to perform sql statement (%v)", queriedErr.Error())
	}
	finalizeFlag := dbutil.FinalizeTx(tx, newTx, cxt)
	if !finalizeFlag {
		return nil, fmt.Errorf("failed to commit transaction by FinalizeTx")
	}
	return queriedData, nil
}

func (s *UserService) Insert(tx db.DatabaseTx, tbName string, colArr []string, valArr []string) error {
	if tbName == "" || len(colArr) == 0 || len(colArr) != len(valArr) {
		return fmt.Errorf("Please provide enough data to perform execution")
	}

	cxt := context.Background()
	newTx := dbutil.PrepTx(tx, s.db, cxt)
	execErr := s.repo.SqlInsert(newTx, tbName, colArr, valArr)
	if execErr != nil {
		newTx.Rollback(context.Background())
		return fmt.Errorf("failed to execute sql statement")
	}

	finalizeFlag := dbutil.FinalizeTx(tx, newTx, cxt)
	if !finalizeFlag {
		return fmt.Errorf("failed to commit the transaction by FinalizeTx")
	}

	return nil
}

// This method allow identifier, If method don't recieved any identifier or valIdentifier, Please put empty string ""
func (s *UserService) Update(tx db.DatabaseTx, tbName string, colArr []string, colVal []string, identifier string, valIdentitier string) error {

	if tbName == "" || len(colArr) == 0 {
		return fmt.Errorf("please provide enough data to perform execution")
	}

	cxt := context.Background()
	newTx := dbutil.PrepTx(tx, s.db, cxt)

	execErr := s.repo.SqlUpdate(newTx, tbName, colArr, colVal, identifier, valIdentitier)
	if execErr != nil {
		rollbackErr := newTx.Rollback(context.Background())
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback: %v", rollbackErr.Error())
		}

		return fmt.Errorf("failed to execute the sql Statement")
	}

	finalizeFlag := dbutil.FinalizeTx(tx, newTx, cxt)
	if !finalizeFlag {
		return fmt.Errorf("failed to commit transaction by FinalizeTx")
	}

	return nil
}

func (s *UserService) Delete(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) error {
	if tbName == "" {
		return fmt.Errorf("Please provide all required data")
	}

	cxt := context.Background()
	newTx := dbutil.PrepTx(tx, s.db, cxt)

	execErr := s.repo.SqlDelete(newTx, tbName, identifier, valIdentifier)
	if execErr != nil {
		rollbackErr := newTx.Rollback(context.Background())
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback")
		}
		return fmt.Errorf("failed to perform sql statement: %v", execErr.Error())
	}

	finalizeFlag := dbutil.FinalizeTx(tx, newTx, cxt)
	if !finalizeFlag {
		return fmt.Errorf("failed to commit the transaction by FinalizeTx")
	}

	return nil
}
