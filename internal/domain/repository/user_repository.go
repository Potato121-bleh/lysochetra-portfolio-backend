package repository

import (
	"context"
	"fmt"
	sqlbuilder "profile-portfolio/internal/builder/sqlBuilder"
	"profile-portfolio/internal/db"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/util/dbutil"
)

type UserRepository struct {
}

type UserRepositoryI interface {
	SqlSelect(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) ([]model.UserData, error)
	SqlInsert(tx db.DatabaseTx, tbName string, colArr []string, valArr []string) error
	SqlUpdate(tx db.DatabaseTx, tbName string, colArr []string, colVal []string, identifier string, valIdentifier string) error
	SqlDelete(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) error
}

func NewUserRepository() UserRepositoryI {
	return &UserRepository{}
}

func (r *UserRepository) SqlSelect(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) ([]model.UserData, error) {
	builder := sqlbuilder.NewSqlBuilder("select")
	if builder == nil {
		return nil, fmt.Errorf("failed to start the builder")
	}

	sqlStatement := builder.AddTable(tbName)
	if identifier != "" {
		sqlStatement.AddIdentifier(identifier)
	}
	rows, queriedErr := tx.Query(context.Background(), sqlStatement.Build(), valIdentifier)
	if queriedErr != nil {
		return nil, fmt.Errorf("failed to execute transaction")
	}

	responseData := []model.UserData{}
	for rows.Next() {

		newPrepInstance := model.UserData{}
		scanErr := dbutil.ScanRow(rows, &newPrepInstance)
		if scanErr != nil {
			return nil, fmt.Errorf("failed to scan queried row: %v", scanErr.Error())
		}

		responseData = append(responseData, newPrepInstance)

	}

	return responseData, nil

}

func (r *UserRepository) SqlInsert(tx db.DatabaseTx, tbName string, colArr []string, valArr []string) error {
	builder := sqlbuilder.NewSqlBuilder("insert")
	if builder == nil {
		return fmt.Errorf("failed to start the builder")
	}

	sqlStatement := builder.AddTable(tbName).AddColumn(colArr)

	valArrI := make([]interface{}, len(valArr))
	for i, v := range valArr {
		valArrI[i] = v
	}

	insertUserCommandTag, insertUserErr := tx.Exec(context.Background(), sqlStatement.Build(), valArrI...)
	if insertUserErr != nil || insertUserCommandTag.RowsAffected() != 1 {
		errStr := "failed to execute sql statement"
		if insertUserErr != nil {
			errStr += ": " + insertUserErr.Error()
		}
		return fmt.Errorf(errStr)
	}

	return nil
}

func (r *UserRepository) SqlUpdate(tx db.DatabaseTx, tbName string, colArr []string, colVal []string, identifier string, valIdentifier string) error {

	builder := sqlbuilder.NewSqlBuilder("update")
	if builder == nil {
		return fmt.Errorf("failed to start the builder")
	}

	sqlStatement := builder.AddTable(tbName).AddColumn(colArr)
	if identifier != "" {
		sqlStatement.AddIdentifier(identifier)
		colVal = append(colVal, valIdentifier)
	}

	colValI := make([]interface{}, len(colVal))
	for i, v := range colVal {
		colValI[i] = v
	}

	// ColValI already checked: it work well as identifier value stay at the end
	_, updateUseridErr := tx.Exec(context.Background(), sqlStatement.Build(), colValI...)
	if updateUseridErr != nil {
		return fmt.Errorf("failed to perform sql transaction")
	}

	return nil
}

func (r *UserRepository) SqlDelete(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) error {
	builder := sqlbuilder.NewSqlBuilder("delete")
	if builder == nil {
		return fmt.Errorf("failed to start the builder")
	}

	sqlStatement := builder.AddTable(tbName)
	if identifier == "" || valIdentifier == "" {
		return fmt.Errorf("please make sure you input the identifier")
	}
	sqlStatement.AddIdentifier(identifier)

	_, deleteUseridErr := tx.Exec(context.Background(), sqlStatement.Build(), valIdentifier)
	if deleteUseridErr != nil {
		return fmt.Errorf("failed to perform sql transaction")
	}

	return nil
}
