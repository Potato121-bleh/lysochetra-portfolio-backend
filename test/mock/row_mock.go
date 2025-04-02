package mockUtil

import (
	"fmt"
	// "strings"
	"time"

	"github.com/stretchr/testify/mock"
)

func NewMockRow(dataset [][]interface{}) *MockRow {
	return &MockRow{
		currentIndex: -1,
		data:         dataset,
	}
}

func NewMockQueryRow(dataset []interface{}) *MockQueryRow {
	return &MockQueryRow{
		data: dataset,
	}
}

type MockRow struct {
	mock.Mock
	currentIndex int
	data         [][]interface{}
}

func (r *MockRow) Next() bool {
	r.currentIndex++
	return r.currentIndex < len(r.data)
}

func (r *MockRow) Scan(dest ...interface{}) error {
	mockArgs := r.Called(dest...)
	if r.currentIndex >= len(r.data) {
		return fmt.Errorf("no row exist")
	}

	for i, ele := range dest {
		switch valType := ele.(type) {
		case *int:
			*valType = r.data[r.currentIndex][i].(int)
		case *string:
			*valType = r.data[r.currentIndex][i].(string)
		case *time.Time:
			*valType = r.data[r.currentIndex][i].(time.Time)
		default:
			return fmt.Errorf("unsupported type")
		}
	}

	return mockArgs.Error(0)
}

func (r *MockRow) Close() {
	return
}

type MockQueryRow struct {
	mock.Mock
	data []interface{}
}

func (r *MockQueryRow) Scan(dest ...interface{}) error {
	mockArgs := r.Called(dest...)

	if len(r.data) == 0 {
		return mockArgs.Error(0)
	}

	for i, ele := range dest {
		switch valType := ele.(type) {
		case *int:
			*valType = r.data[i].(int)
		case *string:
			*valType = r.data[i].(string)
		case *time.Time:
			*valType = r.data[i].(time.Time)
		default:
			return fmt.Errorf("unsupported type")
		}
	}

	return mockArgs.Error(0)
}
