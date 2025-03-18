package repository

import (
	sqlbuilder "backend/internal/builder/sqlBuilder"
	"backend/internal/util/dbutil"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepository[T dbutil.OnlyStruct] struct {
	//db *pgxpool.Pool
}

// This type are for temporary use
type UserData struct {
	Id             int       `json:"userId"`
	Username       string    `json:"userName"`
	Password       string    `json:"password"`
	Nickname       string    `json:"nickname"`
	RegisteredDate time.Time `json:"registered_date"`
	SettingId      int       `json:"settingId"`
}

// tbName, colArr, valueArr, identifier?

func (r *UserRepository[T]) SqlSelect(tx pgx.Tx, tbName string, identifier string, valIdentifier string) ([]T, error) {
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

	responseData := []T{}

	for rows.Next() {

		// prepStruct := T{}

		// As we got the T of struct our sqlSelect doesn't know what the exact type of T
		// By this we have to reflect the T type into actual struct so we can use it later
		// t := reflect.TypeOf(new(T)).Elem()
		// tPtr := reflect.New(t).Elem()
		// newStructInstance := tPtr.Interface().(T)
		newPrepInstance := dbutil.GenericStructConversion[T]()

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

func (r *UserRepository[T]) SqlInsert(tx pgx.Tx, tbName string, colArr []string, valArr []string) error {
	builder := sqlbuilder.NewSqlBuilder("insert")
	if builder == nil {
		return fmt.Errorf("failed to start the builder")
	}

	sqlStatement := builder.AddTable(tbName).AddColumn(colArr)
	fmt.Println(valArr)

	valArrI := make([]interface{}, len(valArr))
	for i, v := range valArr {
		valArrI[i] = v
	}

	fmt.Println(sqlStatement.Build())
	fmt.Println("the value of field: ")
	fmt.Println(colArr)
	fmt.Println(valArrI)

	insertUserCommandTag, insertUserErr := tx.Exec(context.Background(), sqlStatement.Build(), valArrI...)
	if insertUserErr != nil || insertUserCommandTag.RowsAffected() != 1 {
		return fmt.Errorf("failed to execute sql statement: %v", insertUserErr.Error())
	}

	// commitTranErr := tx.Commit(context.Background())
	// if commitTranErr != nil {
	// 	tx.Rollback(context.Background())
	// 	return false
	// }

	return nil
}

func (r *UserRepository[T]) SqlUpdate(tx pgx.Tx, tbName string, colArr []string, colVal []string, identifier string, valIdentifier string) error {

	fmt.Println("sql update has triggered")

	// unknown error?
	builder := sqlbuilder.NewSqlBuilder("update")
	fmt.Println("passed level 1")
	fmt.Println(builder)
	if builder == nil {
		fmt.Println("it triggered error")
		return fmt.Errorf("failed to start the builder")
	}

	// colValI := make([]interface{}, len(colVal))
	// fmt.Println("passed level 2")
	// for i, v := range colVal {
	// 	colValI[i] = v
	// }

	fmt.Println("passed level 3")
	fmt.Println(colArr)
	sqlStatement := builder.AddTable(tbName).AddColumn(colArr)
	if identifier != "" {
		sqlStatement.AddIdentifier(identifier)
		colVal = append(colVal, valIdentifier)
	}

	colValI := make([]interface{}, len(colVal))
	fmt.Println("passed level 2")
	for i, v := range colVal {
		colValI[i] = v
	}
	// else {
	// 	return fmt.Errorf("failed to perform sql transaction, Please input the required column value")
	// }

	fmt.Println(sqlStatement.Build())
	fmt.Println(colValI)

	updateUserCommandTag, updateUseridErr := tx.Exec(context.Background(), sqlStatement.Build(), colValI...)
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

func (r *UserRepository[T]) SqlDelete(tx pgx.Tx, tbName string, identifier string, valIdentifier string) error {
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

// func (r *UserRepository[T]) getData() []T {
// 	return r.queriedData
// }
