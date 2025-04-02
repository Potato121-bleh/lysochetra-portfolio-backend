package db

type RowsScanner interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}
