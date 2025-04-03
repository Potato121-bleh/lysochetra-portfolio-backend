package repository

import (
	"context"
	"fmt"

	sqlbuilder "profile-portfolio/internal/builder/sqlBuilder"
	"profile-portfolio/internal/db"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/util/dbutil"
)

type SettingRepository struct {
}

type SettingRepositoryI interface {
	SqlSelect(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) ([]model.SettingStruct, error)
	SqlInsert(tx db.DatabaseTx, tbName string, colArr []string, valArr []string) error
	SqlUpdate(tx db.DatabaseTx, tbName string, colArr []string, colVal []string, identifier string, valIdentifier string) error
	SqlDelete(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) error
}

func NewSettingRepository() SettingRepositoryI {
	return &SettingRepository{}
}

func (r *SettingRepository) SqlSelect(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) ([]model.SettingStruct, error) {
	builder := sqlbuilder.NewSqlBuilder("select")
	if builder == nil {
		return nil, fmt.Errorf("failed to start the builder")
	}

	sqlStatement := builder.AddTable(tbName)
	if identifier != "" {
		sqlStatement.AddIdentifier(identifier)
	}

	// format to []interface{}
	rows, queriedErr := tx.Query(context.Background(), sqlStatement.Build(), valIdentifier)
	if queriedErr != nil {
		return nil, fmt.Errorf("failed to execute transaction")
	}

	responseData := []model.SettingStruct{}

	for rows.Next() {

		// prepStruct := T{}

		// As we got the T of struct our sqlSelect doesn't know what the exact type of T
		// By this we have to reflect the T type into actual struct so we can use it later
		// t := reflect.TypeOf(new(T)).Elem()
		// tPtr := reflect.New(t).Elem()
		// newStructInstance := tPtr.Interface().(T)
		newPrepInstance := model.SettingStruct{}

		scanErr := dbutil.ScanRow(rows, &newPrepInstance)
		if scanErr != nil {
			return nil, fmt.Errorf("failed to scan queried row: %v", scanErr.Error())
		}

		responseData = append(responseData, newPrepInstance)

	}

	return responseData, nil
}

func (r *SettingRepository) SqlInsert(tx db.DatabaseTx, tbName string, colArr []string, valArr []string) error {
	builder := sqlbuilder.NewSqlBuilder("insert")
	if builder == nil {
		return fmt.Errorf("failed to start the builder")
	}

	sqlStatement := builder.AddTable(tbName).AddColumn(colArr)

	valArrI := make([]interface{}, len(valArr))
	for i, v := range valArr {
		valArrI[i] = v
	}

	insertSettingCommandTag, insertSettingErr := tx.Exec(context.Background(), sqlStatement.Build(), valArrI...)

	rowAffected := insertSettingCommandTag.RowsAffected()
	if insertSettingErr != nil || rowAffected != 1 {
		errStr := "failed to execute sql statement"
		if insertSettingErr != nil {
			errStr += ": " + insertSettingErr.Error()
		}
		return fmt.Errorf(errStr)
	}
	return nil
}

func (r *SettingRepository) SqlUpdate(tx db.DatabaseTx, tbName string, colArr []string, colVal []string, identifier string, valIdentifier string) error {

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

	_, updateUseridErr := tx.Exec(context.Background(), sqlStatement.Build(), colValI...)
	if updateUseridErr != nil {
		return fmt.Errorf("failed to perform sql transaction")
	}

	return nil
}

func (r *SettingRepository) SqlDelete(tx db.DatabaseTx, tbName string, identifier string, valIdentifier string) error {
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
