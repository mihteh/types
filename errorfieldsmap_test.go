package types

import (
	"reflect"
	"testing"
)

func TestAddError(t *testing.T) {
	expectedErr := ErrorFieldsMap{
		"key1": []string{"value1"},
		"key2": []string{"value2", "value3"},
	}
	err := make(ErrorFieldsMap)
	err.AddError("key1", "value1")
	err.AddError("key2", "value2")
	err.AddError("key2", "value3")
	if !reflect.DeepEqual(expectedErr, err) {
		t.Fatalf("Структуры не равны. Ожидалось: \n%v\nПолучено: \n%v\n",
			expectedErr, err)
	}
}
