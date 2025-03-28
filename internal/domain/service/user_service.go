package service

import (
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/domain/repository"
	"profile-portfolio/internal/util/dbutil"

	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SignUp(reqUsername string, reqNickname string, reqPassword string) error
type UserServiceI interface {
	Select(tx pgx.Tx, tbName string, identifier string, valIdentifier string) ([]model.UserData, error)
	Insert(tx pgx.Tx, tbName string, colArr []string, valArr []string) error
	Update(tx pgx.Tx, tbName string, colArr []string, colVal []string, identifier string, valIdentitier string) error
	Delete(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error
}

type UserService struct {
	db   *pgxpool.Pool
	repo repository.UserRepository
}

func NewUserService(db *pgxpool.Pool) UserServiceI {
	return &UserService{
		db:   db,
		repo: repository.UserRepository{},
	}
}

func (s *UserService) Select(tx pgx.Tx, tbName string, identifier string, valIdentifier string) ([]model.UserData, error) {
	if tbName == "" {
		return nil, fmt.Errorf("Please provide all required data")
	}

	cxt := context.Background()

	newTx := dbutil.PrepTx(tx, s.db, cxt)

	// tx, startTxErr := s.db.Begin(context.Background())
	// if startTxErr != nil {
	// 	return nil, fmt.Errorf("failed to begin the transaction")
	// }

	queriedData, queriedErr := s.repo.SqlSelect(newTx, tbName, identifier, valIdentifier)
	if queriedErr != nil {
		rollbackErr := newTx.Rollback(context.Background())
		if rollbackErr != nil {
			return nil, fmt.Errorf("failed to perform rollback (%v)", queriedErr.Error())
		}
		return nil, fmt.Errorf("failed to perform sql statement (%v)", queriedErr.Error())
	}

	// As now we can retrieve data from getData in repository
	// as this method not available in current interface we have to do type assertion
	// if getQueriedData, ok := s.repo.()

	// if queriedData == nil {
	// 	rollbackErr := newTx.Rollback(context.Background())
	// 	if rollbackErr != nil {
	// 		return nil, fmt.Errorf("failed to perform rollback")
	// 	}
	// 	return nil, fmt.Errorf("failed to perform sql statement")
	// }

	finalizeFlag := dbutil.FinalizeTx(tx, newTx, cxt)
	if !finalizeFlag {
		return nil, fmt.Errorf("failed to commit transaction by FinalizeTx")
	}

	// commitErr := tx.Commit(context.Background())
	// if commitErr != nil {
	// 	return nil, fmt.Errorf("failed to commit the execution")
	// }
	return queriedData, nil
}

func (s *UserService) Insert(tx pgx.Tx, tbName string, colArr []string, valArr []string) error {
	if tbName == "" || len(colArr) == 0 || len(colArr) != len(valArr) {
		return fmt.Errorf("Please provide enough data to perform execution")
	}

	cxt := context.Background()
	newTx := dbutil.PrepTx(tx, s.db, cxt)

	// tx, startTxErr := s.db.Begin(context.Background())
	// if startTxErr != nil {
	// 	return fmt.Errorf("failed to begin transaction")
	// }

	execErr := s.repo.SqlInsert(newTx, tbName, colArr, valArr)
	if execErr != nil {
		newTx.Rollback(context.Background())
		return fmt.Errorf("failed to execute sql statement")
	}

	finalizeFlag := dbutil.FinalizeTx(tx, newTx, cxt)
	if !finalizeFlag {
		return fmt.Errorf("failed to commit the transaction by FinalizeTx")
	}

	// commitErr := newTx.Commit(context.Background())
	// if commitErr != nil {
	// 	return fmt.Errorf("failed to commit the execution")
	// }
	return nil

}

// This method allow identifier, If method don't recieved any identifier or valIdentifier, Please put empty string ""
func (s *UserService) Update(tx pgx.Tx, tbName string, colArr []string, colVal []string, identifier string, valIdentitier string) error {

	fmt.Println("update has triggered")

	if tbName == "" || len(colArr) == 0 {
		return fmt.Errorf("please provide enough data to perform execution")
	}

	cxt := context.Background()
	newTx := dbutil.PrepTx(tx, s.db, cxt)

	// tx, txBeginErr := s.db.Begin(context.Background())
	// if txBeginErr != nil {
	// 	return fmt.Errorf("transaction failed to start: %v", txBeginErr.Error())
	// }

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

	// commitErr := tx.Commit(context.Background())
	// if commitErr != nil {
	// 	return fmt.Errorf("failed to commit the execution")
	// }

	return nil

}

func (s *UserService) Delete(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error {
	if tbName == "" {
		return fmt.Errorf("Please provide all required data")
	}

	cxt := context.Background()
	newTx := dbutil.PrepTx(tx, s.db, cxt)

	// tx, startTxErr := s.db.Begin(context.Background())
	// if startTxErr != nil {
	// 	return fmt.Errorf("failed to begin the transaction")
	// }

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

	// commitErr := tx.Commit(context.Background())
	// if commitErr != nil {
	// 	return fmt.Errorf("failed to commit the execution")
	// }

	return nil

}
