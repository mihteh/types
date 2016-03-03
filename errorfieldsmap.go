package types

// ErrorFieldsMap реализует отображение имён полей в срезы сообщений об ошибках
// полезно при валидации
type ErrorFieldsMap map[string][]string

func (e ErrorFieldsMap) AddError(name, value string) {
	if _, ok := e[name]; !ok {
		e[name] = []string{value}
		return
	}
	e[name] = append(e[name], value)
}
