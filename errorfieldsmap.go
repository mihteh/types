package types

// ErrorFieldsMap реализует отображение имён полей в срезы сообщений об ошибках
// полезно при валидации
type ErrorFieldsMap map[string][]string

// NewErrorFieldsMap возвращает пустое отображение
func NewErrorFieldsMap() ErrorFieldsMap {
	return make(ErrorFieldsMap)
}

// AddError добавляет поле c записью об ошибке
func (e ErrorFieldsMap) AddError(name, value string) {
	if _, ok := e[name]; !ok {
		e[name] = []string{value}
		return
	}
	e[name] = append(e[name], value)
}

// AddFromMap соединяет данные из from с данными из e, результат помещает в e
func (e ErrorFieldsMap) AddFromMap(from ErrorFieldsMap) {
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
