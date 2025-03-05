package service

import (
	"backend/internal/domain/repository"

	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserServiceI interface {
	Select(tbName string, identifier string, valIdentifier string) ([]repository.UserData, error)
	Insert(tbName string, colArr []string, valArr []string) error
	Update(tbName string, colArr []string, identifier string, valIdentitier string) error
	Delete(tbName string, identifier string, valIdentifier string) error
}

type UserService struct {
	db   *pgxpool.Pool
	repo repository.UserRepoI
}

func (s *UserService) Select(tbName string, identifier string, valIdentifier string) ([]repository.UserData, error) {
	if tbName == "" {
		return nil, fmt.Errorf("Please provide all required data")
	}

	tx, startTxErr := s.db.Begin(context.Background())
	if startTxErr != nil {
		return nil, fmt.Errorf("failed to begin the transaction")
	}

	queriedData := s.repo.SqlSelect(tx, tbName, identifier, valIdentifier)
	if queriedData == nil {
		rollbackErr := tx.Rollback(context.Background())
		if rollbackErr != nil {
			return nil, fmt.Errorf("failed to perform rollback")
		}
		return nil, fmt.Errorf("failed to perform sql statement")
	}

	commitErr := tx.Commit(context.Background())
	if commitErr != nil {
		return nil, fmt.Errorf("failed to commit the execution")
	}
	return queriedData, nil
}

func (s *UserService) Insert(tbName string, colArr []string, valArr []string) error {
	if tbName == "" || len(colArr) == 0 || len(colArr) != len(valArr) {
		return fmt.Errorf("Please provide enough data to perform execution")
	}

	tx, startTxErr := s.db.Begin(context.Background())
	if startTxErr != nil {
		return fmt.Errorf("failed to begin transaction")
	}

	execErr := s.repo.SqlInsert(tx, tbName, colArr, valArr)
	if execErr != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("failed to execute sql statement")
	}

	commitErr := tx.Commit(context.Background())
	if commitErr != nil {
		return fmt.Errorf("failed to commit the execution")
	}
	return nil

}

// This method allow identifier, If method don't recieved any identifier or valIdentifier, Please put empty string ""
func (s *UserService) Update(tbName string, colArr []string, identifier string, valIdentitier string) error {

	if tbName == "" || len(colArr) == 0 {
		return fmt.Errorf("please provide enough data to perform execution")
	}

	tx, txBeginErr := s.db.Begin(context.Background())
	if txBeginErr != nil {
		return fmt.Errorf("transaction failed to start: %v", txBeginErr.Error())
	}

	execErr := s.repo.SqlUpdate(tx, tbName, colArr, identifier, valIdentitier)
	if execErr != nil {
		rollbackErr := tx.Rollback(context.Background())
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback: %v", rollbackErr.Error())
		}

		return fmt.Errorf("failed to execute the sql Statement")
	}

	commitErr := tx.Commit(context.Background())
	if commitErr != nil {
		return fmt.Errorf("failed to commit the execution")
	}
	return nil

}

func (s *UserService) Delete(tbName string, identifier string, valIdentifier string) error {
	if tbName == "" {
		return fmt.Errorf("Please provide all required data")
	}

	tx, startTxErr := s.db.Begin(context.Background())
	if startTxErr != nil {
		return fmt.Errorf("failed to begin the transaction")
	}

	execErr := s.repo.SqlDelete(tx, tbName, identifier, valIdentifier)
	if execErr != nil {
		rollbackErr := tx.Rollback(context.Background())
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback")
		}
		return fmt.Errorf("failed to perform sql statement: %v", execErr.Error())
	}

	commitErr := tx.Commit(context.Background())
	if commitErr != nil {
		return fmt.Errorf("failed to commit the execution")
	}

	return nil

}

func (s *UserService) SignUp(reqUsername string) error {
	tx, startTxErr := s.db.Begin(context.Background())
	if startTxErr != nil {
		return fmt.Errorf("failed to start transaction")
	}

	// colArr := []string{"darkmode", "sound", "colorpalettes", "font", "language"}

	s.repo.SqlInsert(
		tx,
		"user_setting",
		[]string{"darkmode", "sound", "colorpalettes", "font", "language"},
		[]string{"0", "0", "0", "1", "1"},
	)

	// tx.Exec(context.Background())

	return nil

}
