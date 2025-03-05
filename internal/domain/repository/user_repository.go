package repository

import (
	sqlbuilder "backend/internal/builder/sqlBuilder"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type UserRepoI interface {
	SqlSelect(tx pgx.Tx, tbName string, identifier string, valIdentifier string) []UserData
	SqlInsert(tx pgx.Tx, tbName string, colArr []string, valArr []string) error
	SqlUpdate(tx pgx.Tx, tbName string, colArr []string, identifier string, valIdentifier string) error
	SqlDelete(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error
}

type UserRepository struct {
	//db *pgxpool.Pool
}

// This type are for temporary use
type UserData struct {
	Id        int    `json:"userId"`
	Username  string `json:"userName"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	SettingId int    `json:"settingId"`
}

// tbName, colArr, valueArr, identifier?

func (r *UserRepository) SqlSelect(tx pgx.Tx, tbName string, identifier string, valIdentifier string) []UserData {
	builder := sqlbuilder.NewSqlBuilder("select")
	if builder == nil {
		return nil
	}
	sqlStatement := builder.AddTable(tbName)
	if identifier != "" {
		sqlStatement.AddIdentifier(identifier)
	}

	// format to []interface{}
	rows, queriedErr := tx.Query(context.Background(), sqlStatement.Build(), valIdentifier)
	if queriedErr != nil {
		return nil
	}

	responseData := []UserData{}

	for rows.Next() {
		var useridTem int
		var usernameTem string
		var userpasswordTem string
		var usernicknameTem string
		var usersettingidTem int
		scannedRowErr := rows.Scan(&useridTem, &usernameTem, &usernicknameTem, &userpasswordTem, nil, &usersettingidTem)
		if scannedRowErr != nil {
			return nil
		}

		//create context to pass the data to actual handler
		newUser := UserData{
			Id:        useridTem,
			Username:  usernameTem,
			Password:  userpasswordTem,
			Nickname:  usernicknameTem,
			SettingId: usersettingidTem,
		}
		responseData = append(responseData, newUser)
	}

	return responseData

}

func (r *UserRepository) SqlInsert(tx pgx.Tx, tbName string, colArr []string, valArr []string) error {
	builder := sqlbuilder.NewSqlBuilder("insert")
	if builder == nil {
		return fmt.Errorf("failed to start the builder")
	}

	sqlStatement := builder.AddTable(tbName).AddColumn(colArr)

	valArrI := make([]interface{}, len(valArr))
	for i, v := range colArr {
		valArrI[i] = v
	}

	insertUserCommandTag, insertUserErr := tx.Exec(context.Background(), sqlStatement.Build(), valArrI...)
	if insertUserErr != nil || insertUserCommandTag.RowsAffected() != 1 {
		return fmt.Errorf("failed to execute sql statement")
	}

	// commitTranErr := tx.Commit(context.Background())
	// if commitTranErr != nil {
	// 	tx.Rollback(context.Background())
	// 	return false
	// }

	return nil
}

func (r *UserRepository) SqlUpdate(tx pgx.Tx, tbName string, colArr []string, identifier string, valIdentifier string) error {

	builder := sqlbuilder.NewSqlBuilder("update")
	if builder != nil {
		return fmt.Errorf("failed to start the builder")
	}

	colArrI := make([]interface{}, len(colArr))
	for i, v := range colArr {
		colArrI[i] = v
	}

	sqlStatement := builder.AddTable(tbName).AddColumn(colArr)
	if identifier != "" {
		sqlStatement.AddIdentifier(identifier)
		colArrI[len(colArrI)] = valIdentifier
	}

	updateUserCommandTag, updateUseridErr := tx.Exec(context.Background(), sqlStatement.Build(), colArrI...)
	if updateUseridErr != nil || updateUserCommandTag.RowsAffected() != 1 {
		return fmt.Errorf("failed to perform sql transaction")
	}

	// commitTranErr := tx.Commit(context.Background())
	// if commitTranErr != nil {
	// 	tx.Rollback(context.Background())
	// 	return false
	// }

	return nil
}

func (r *UserRepository) SqlDelete(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error {
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
