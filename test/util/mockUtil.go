package testUtilTool

import (
	"reflect"

	"github.com/stretchr/testify/mock"
)

// As we don't know what field it need but we know it base on model
// we have to take the model get the field and get the total val of it - 1
//   - This method will accept your model and convert into an array of mock.Anything
func CountMockAnything(model interface{}) []interface{} {
	fieldCount := reflect.TypeOf(model).NumField()
	mockScanArg := make([]interface{}, fieldCount)
	for i, _ := range mockScanArg {
		mockScanArg[i] = mock.Anything
	}
	return mockScanArg
}
