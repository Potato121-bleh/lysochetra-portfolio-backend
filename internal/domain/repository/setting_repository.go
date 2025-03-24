package repository

import (
	"context"
	"fmt"

	sqlbuilder "profile-portfolio/internal/builder/sqlBuilder"
	"profile-portfolio/internal/domain/model"
	"profile-portfolio/internal/util/dbutil"

	"github.com/jackc/pgx/v5"
)

type SettingRepository struct {
}

type SettingRepositoryI interface {
	SqlSelect(tx pgx.Tx, tbName string, identifier string, valIdentifier string) ([]model.SettingStruct, error)
	SqlInsert(tx pgx.Tx, tbName string, colArr []string, valArr []string) error
	SqlUpdate(tx pgx.Tx, tbName string, colArr []string, colVal []string, identifier string, valIdentifier string) error
	SqlDelete(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error
}

func NewSettingRepository() SettingRepositoryI {
	return &SettingRepository{}
}

func (r *SettingRepository) SqlSelect(tx pgx.Tx, tbName string, identifier string, valIdentifier string) ([]model.SettingStruct, error) {
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

		// var useridTem int
		// var usernameTem string
		// var userpasswordTem string
		// var usernicknameTem string
		// var usersettingidTem int
		// scannedRowErr := rows.Scan(&useridTem, &usernameTem, &usernicknameTem, &userpasswordTem, nil, &usersettingidTem)
		// if scannedRowErr != nil {
		// 	return nil
		// }

		// //create context to pass the data to actual handler
		// newUser := UserData{
		// 	Id:        useridTem,
		// 	Username:  usernameTem,
		// 	Password:  userpasswordTem,
		// 	Nickname:  usernicknameTem,
		// 	SettingId: usersettingidTem,
		// }
		// responseData = append(responseData, newUser)
	}

	// instead of returning the queried data we returning error
	// so that the service know there is no error upon query the data, the service can call getData
	// r.queriedData = responseData
	return responseData, nil

}

func (r *SettingRepository) SqlInsert(tx pgx.Tx, tbName string, colArr []string, valArr []string) error {
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
		return fmt.Errorf("failed to execute sql statement: %v", insertUserErr.Error())
	}

	return nil
}

func (r *SettingRepository) SqlUpdate(tx pgx.Tx, tbName string, colArr []string, colVal []string, identifier string, valIdentifier string) error {

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

	updateUserCommandTag, updateUseridErr := tx.Exec(context.Background(), sqlStatement.Build(), colValI...)
	if updateUseridErr != nil || updateUserCommandTag.RowsAffected() != 1 {
		return fmt.Errorf("failed to perform sql transaction")
	}

	return nil
}

func (r *SettingRepository) SqlDelete(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error {
	builder := sqlbuilder.NewSqlBuilder("delete")
	if builder == nil {
		return fmt.Errorf("failed to start the builder")
	}

	sqlStatement := builder.AddTable(tbName)
	if identifier == "" || valIdentifier == "" {
		return fmt.Errorf("please make sure you input the identifier")
	}
	sqlStatement.AddIdentifier(identifier)

	deleteUserCommandTag, deleteUseridErr := tx.Exec(context.Background(), sqlStatement.Build(), valIdentifier)
	if deleteUseridErr != nil || deleteUserCommandTag.RowsAffected() != 1 {
		return fmt.Errorf("failed to perform sql transaction")
	}

	return nil
}
