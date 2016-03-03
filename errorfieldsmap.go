package types

import (
	"fmt"
	"sort"
	"strings"
)

// Validation реализует отображение имён полей в срезы сообщений об ошибках
// полезно при валидации
type Validation map[string][]string

// NewValidation возвращает пустое отображение
func NewValidation() Validation {
	return make(Validation)
}

// AddError добавляет поле c записью об ошибке
func (e Validation) AddError(name, value string) {
	if _, ok := e[name]; !ok {
		e[name] = []string{value}
		return
	}
	e[name] = append(e[name], value)
}

// AddErros добавляет к e ошибки из валидации from
func (e Validation) AddErrors(from Validation) {
	for name, value := range from {
		if len(value) == 0 {
			e[name] = []string{}
			continue
		}
		for _, element := range value {
			e.AddError(name, element)
		}
	}
}

// HasErrors проверяет, есть ли ошибка. Возвращает true если есть, иначе false
func (e Validation) HasErrors() bool {
	for _, value := range e {
		if len(value) > 0 {
			return true
		}
	}
	return false
}

// String описывает преобразование в строку
func (e Validation) String() string {
	var res string
	for name, value := range e {
		for _, element := range value {
			res += fmt.Sprintf("%s: %s\n", name, element)
		}
	}

	parts := strings.Split(res, "\n")
	sort.Strings(parts)
	return strings.TrimSpace(strings.Join(parts, "\n"))
}
