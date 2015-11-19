package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DateTime хранит дату-время и шаблон для преобразования при сериализации
type DateTime struct {
	time.Time
	Layout string
}

// Date хранит дату и шаблон для преобразования при сериализации
type Date struct {
	time.Time
	Layout string
}

// Шаблоны для сериализации
const (
	DateTimeLayout        = "2006-01-02 15:04:05"
	DateLayout            = "2006-01-02"
	GraphsDateLayout      = "02.01.2006"
	GraphsDateShortLayout = "02.01"
)

var defaultLocation *time.Location

// Задаёт часовой пояс по умолчанию
func init() {
	var err error
	defaultLocation, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Ошибка time.LoadLocation")
	}
}

// ToDateTime формирует объект типа DateTime на основе времени t и шаблона DateTimeLayout
func ToDateTime(t time.Time) DateTime {
	return DateTime{t.In(defaultLocation), DateTimeLayout}
}

// ToDate формирует объект типа Date на основе времени t и шаблона DateLayout
func ToDate(t time.Time) Date {
	dS := fmt.Sprintf("%s %02d:%02d:%02d", t.Format(DateLayout), 0, 0, 0)
	d, _ := time.ParseInLocation(DateTimeLayout, dS, defaultLocation)
	return Date{d.In(defaultLocation), DateLayout}
}

// DaysBefore возвращает количество полных дней, прошедших от d до endDate
// Например, если d было вчера, endDate - сегодня, то возвращается 1
// Если endDate было раньше чем d, то возвращается отрицательное число.
func (d Date) DaysBefore(endDate Date) int {
	return int(endDate.Time.Sub(d.Time).Hours() / 24)
}

// StringToDateTime формирует объект типа DateTime на основе строки s,
// заданной по шаблону DateTimeLayout
func StringToDateTime(s string) (DateTime, error) {
	t, err := time.ParseInLocation(DateTimeLayout, s, defaultLocation)
	if err != nil {
		return DateTime{}, err
	}
	return ToDateTime(t), nil
}

// StringDateToDateTimeHMS формирует объект типа DateTime на основе строки s,
// заданной по шаблону DateTimeLayout, и значений часов, минут и секунд,
// заданных параметрами hours, mins, secs соответственно.
func StringDateToDateTimeHMS(s string, hours int, mins int, secs int) (DateTime, error) {
	t, err := time.ParseInLocation(DateLayout, s, defaultLocation)
	if err != nil {
		return DateTime{}, err
	}
	d := ToDateTime(t)
	d = d.SetHMS(hours, mins, secs)
	return d, nil
}

// StringToDate формирует объект типа Date на основе строки s, заданной по шаблону DateLayout
func StringToDate(s string) (Date, error) {
	t, err := time.ParseInLocation(DateLayout, s, defaultLocation)
	if err != nil {
		return Date{}, err
	}
	return ToDate(t), nil
}

// NeverDate возвращает объект Date, соответствующий дате в далёком прошлом
func NeverDate() Date {
	t, _ := time.ParseInLocation(DateLayout, "1990-01-01", defaultLocation)
	return ToDate(t)
}

// DateNow возвращает объект Date, соответствующий дате сегодня
func DateNow() Date {
	return ToDate(time.Now())
}

// DateTimeNow возвращает объект DateTime, соответствующий дате-времени сейчас
func DateTimeNow() DateTime {
	return ToDateTime(time.Now())
}

// DateTimeTodayHMS возвращает объект DateTime, соответствующий дате сегодня,
// с установленными значениями часов, минут, секунд согласно заданным параметрам
// hours, mins, secs соответственно.
func DateTimeTodayHMS(hours int, mins int, secs int) DateTime {
	d := ToDateTime(time.Now())
	return d.SetHMS(hours, mins, secs)
}

// NeverTime возвращает объект DateTime, соответствующий дате-времени в далёком прошлом
func NeverTime() DateTime {
	t, _ := time.ParseInLocation(DateTimeLayout, "1990-01-01 00:00:00", defaultLocation)
	return ToDateTime(t)
}

// setDefaultLayoutIfEmpty устанавливает в объекте DateTime шаблон вывода даты-времени
// по умолчанию DateTimeLayout, если шаблон не установлен
func (d *DateTime) setDefaultLayoutIfEmpty() {
	if strings.TrimSpace(d.Layout) == "" {
		d.Layout = DateTimeLayout
	}
}

// setDefaultLayoutIfEmpty устанавливает в объекте Date шаблон вывода даты
// по умолчанию DateLayout, если шаблон не установлен
func (d *Date) setDefaultLayoutIfEmpty() {
	if strings.TrimSpace(d.Layout) == "" {
		d.Layout = DateLayout
	}
}

// SetHMS возвращает новый объект DateTime на основе объекта d,
// с заданными значениями часов, минут, секунд в параметрах
// hours, mins, secs соответственно.
func (d DateTime) SetHMS(hours int, mins int, secs int) DateTime {
	t := d.Time
	dS := fmt.Sprintf("%s %02d:%02d:%02d", t.Format(DateLayout), hours, mins, secs)
	dt, _ := time.ParseInLocation(DateTimeLayout, dS, defaultLocation)
	d.Time = dt
	return d
}

// ConvertToDate преобразует объект DateTime в объект Date
func (d DateTime) ConvertToDate() Date {
	dS := d.Time.Format(DateLayout)
	t, _ := time.ParseInLocation(DateLayout, dS, defaultLocation)
	return ToDate(t)
}

// ConvertToDateTimeHMS преобразует объект Date в объект DateTime
// с учётом заданных часов, минут, секунд в параметрах hours, mind, secs соответственно.
func (d Date) ConvertToDateTimeHMS(hours int, mins int, secs int) DateTime {
	dt := DateTime{
		Time:   d.Time,
		Layout: d.Layout,
	}
	dt = dt.SetHMS(hours, mins, secs)

	return dt
}

// After возвращает true если дата d позднее d1, иначе false
// Сравнение происходит с точностью до дня.
func (d Date) After(d1 Date) bool {
	return d.Time.After(d1.Time)
}

// Before возвращает true если дата d ранее d1, иначе false
// Сравнение происходит с точностью до дня.
func (d Date) Before(d1 Date) bool {
	return d.Time.Before(d1.Time)
}

// Between возвращает true если дата d находится в интервале дат (d1; d2), иначе false
// Сравнение происходит с точностью до дня.
func (d Date) Between(d1, d2 Date) bool {
	return d.After(d1) && d.Before(d2)
}

// After возвращает true если дата-время d позднее d1, иначе false
func (d DateTime) After(d1 DateTime) bool {
	return d.Time.After(d1.Time)
}

// Before возвращает true если дата-время d ранее d1, иначе false
func (d DateTime) Before(d1 DateTime) bool {
	return d.Time.Before(d1.Time)
}

// Between возвращает true если дата-время d находится в интервале
// даты-времени (d1; d2), иначе false
func (d DateTime) Between(d1, d2 DateTime) bool {
	return d.After(d1) && d.Before(d2)
}

// UnmarshalJSON - реализует интерфейс json.Unmarshaler для объекта DateTime
// десериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d *DateTime) UnmarshalJSON(data []byte) error {
	d.setDefaultLayoutIfEmpty()
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.ParseInLocation(d.Layout, s, defaultLocation)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON - реализует интерфейс json.Marshaler для объекта DateTime
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d DateTime) MarshalJSON() ([]byte, error) {
	d.setDefaultLayoutIfEmpty()
	return []byte(strconv.Quote(d.String())), nil
}

// UnmarshalJSON - реализует интерфейс json.Unmarshaler для объекта Date
// десериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d *Date) UnmarshalJSON(data []byte) error {
	d.setDefaultLayoutIfEmpty()
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.ParseInLocation(d.Layout, s, defaultLocation)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON - реализует интерфейс json.Marshaler для объекта Date
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d Date) MarshalJSON() ([]byte, error) {
	d.setDefaultLayoutIfEmpty()
	return []byte(strconv.Quote(d.String())), nil
}

// String преобразует объект DateTime в строку согласно шаблона в свойстве Layout
func (d DateTime) String() string {
	d.setDefaultLayoutIfEmpty()
	return d.Time.Format(d.Layout)
}

// String преобразует объект Date в строку согласно шаблона в свойстве Layout
func (d Date) String() string {
	d.setDefaultLayoutIfEmpty()
	return d.Time.Format(d.Layout)
}

func scan(value interface{}) (time.Time, error) {
	t := time.Time{}
	if value == nil {
		return t, nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return t, errors.New("Ошибка преобразования значения к типу time.Time")
	}
	return t.In(defaultLocation), nil
}

// Scan преобразует значение времени в БД к типу DateTime
func (d *DateTime) Scan(value interface{}) error {
	d.setDefaultLayoutIfEmpty()
	t, err := scan(value)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// Value преобразует значение типа DateTime к значению в БД
func (d DateTime) Value() (driver.Value, error) {
	return d.Time.In(defaultLocation).Format(DateTimeLayout), nil
}

// Scan преобразует значение времени в БД к типу Date
func (d *Date) Scan(value interface{}) error {
	d.setDefaultLayoutIfEmpty()
	t, err := scanInternal(value)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// Value преобразует значение типа Date к значению в БД
func (d Date) Value() (driver.Value, error) {
	return d.Time.In(defaultLocation).Format(DateLayout), nil
}
