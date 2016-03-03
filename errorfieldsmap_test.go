package types

import (
	"reflect"
	"testing"
)

func TestNewErrorFieldsMap(t *testing.T) {
	result := NewErrorFieldsMap()
	if result == nil {
		t.Fatal("Получен nil")
	}
}

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

func TestAddFromMap(t *testing.T) {
	to := ErrorFieldsMap{
		"key1": []string{"value1"},
		"key2": []string{"value2", "value3"},
	}
	from := ErrorFieldsMap{
		"key1": []string{"value4"},
		"key3": []string{"value5"},
		"key4": []string{},
	}
	expectedResult := ErrorFieldsMap{
		"key1": []string{"value1", "value4"},
		"key2": []string{"value2", "value3"},
		"key3": []string{"value5"},
		"key4": []string{},
	}
	to.AddFromMap(from)
	if !reflect.DeepEqual(to, expectedResult) {
		t.Fatalf("Структуры не равны. Ожидалось: \n%v\nПолучено: \n%v\n",
			expectedResult, to)
	}
}
