package dbutil

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OnlyStruct interface{}

// This function are built for specific use with all services where it need to make sure the transaction are good to go.
// This function will take in your pgx.Tx, if your transaction is nil,
// then the PrepTx will create new tx for later execution
func PrepTx(tx pgx.Tx, db *pgxpool.Pool, cxt context.Context) pgx.Tx {
	if tx != nil {
		return tx
	}

	newTx, beginTxErr := db.Begin(cxt)
	if beginTxErr != nil {
		return nil
	}

	return newTx

}

// This function are built for specific use with all services where it need to make sure whether to commit the transaction or not.
//
// This function will take in your original transaction and check.
//
//   - if it nil, it mean that the handler want us to do method of: one-time use only which we need to commit the transaction as we done.
//
//   - if not nil, it mean that the handler want us to just execute the query and not commit it yet.
//
// @return:
//
//	True: "indicate that all expected execution are good"
//	False: "indicate that all expected execution are failed, (This might due to commit or rollback failed)""
//
// .
func FinalizeTx(originTx pgx.Tx, currentTx pgx.Tx, cxt context.Context) bool {
	if originTx == nil {
		commitErr := currentTx.Commit(cxt)
		if commitErr != nil {
			rollbackErr := currentTx.Rollback(cxt)
			if rollbackErr != nil {
				return false
			}
			return false
		}
	}

	return true

}

// This method will perform operation to retrieve data from provided row into the provided struct
func ScanRow(row pgx.Rows, dest interface{}) error {
	// we retrieve the value of the dest
	fmt.Println("about to ENTER an error")
	val := reflect.ValueOf(dest).Elem()
	fmt.Println("about to PASS an error")

	values := make([]interface{}, val.NumField())

	for i := range values {
		values[i] = val.Field(i).Addr().Interface()
	}

	scanErr := row.Scan(values...)
	if scanErr != nil {
		return fmt.Errorf(scanErr.Error())
	}

	return nil
}

// This method will accept you unknown generic struct and convert it into usable struct
// This method will take your Generic struct and make a copy into new struct for later usage
func GenericStructConversion[T OnlyStruct]() T {
	t := reflect.TypeOf(new(T)).Elem()
	tPtr := reflect.New(t).Elem()
	newStructInstance := tPtr.Interface().(T)
	return newStructInstance
}
