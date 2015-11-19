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
// Если d было вчера, endDate - сегодня, то возвращается 1
// Если endDate было раньше чем d, то возвращается отрицательное число.
func (d Date) DaysBefore(endDate Date) int {
	return int(endDate.Time.Sub(d.Time).Hours() / 24)
}

// StringToDateTime преобразует строку по стандартному шаблону даты-времени в дату-время
func StringToDateTime(s string) (DateTime, error) {
	t, err := time.ParseInLocation(DateTimeLayout, s, defaultLocation)
	if err != nil {
		return DateTime{}, err
	}
	return ToDateTime(t), nil
}

// StringDateToDateTimeHMS преобразует строку по стандартному шаблону
// даты в дату-время с заданным значением часов, минут и секунд
func StringDateToDateTimeHMS(s string, hours int, mins int, secs int) (DateTime, error) {
	t, err := time.ParseInLocation(DateLayout, s, defaultLocation)
	if err != nil {
		return DateTime{}, err
	}
	d := ToDateTime(t)
	d = d.SetHMS(hours, mins, secs)
	return d, nil
}

// StringToDate преобразует строку по стандартному шаблону даты в дату
func StringToDate(s string) (Date, error) {
	t, err := time.ParseInLocation(DateLayout, s, defaultLocation)
	if err != nil {
		return Date{}, err
	}
	return ToDate(t), nil
}

// NeverDate возвращает дату в далёком прошлом
func NeverDate() Date {
	t, _ := time.ParseInLocation(DateLayout, "1990-01-01", defaultLocation)
	return ToDate(t)
}

// DateNow возвращает дату сегодня
func DateNow() Date {
	return ToDate(time.Now())
}

// DateTimeNow возвращает дату-время сейчас
func DateTimeNow() DateTime {
	return ToDateTime(time.Now())
}

// DateTimeTodayHMS возвращает дату-время сегодня в заданными значениями
// часов, минут, секунд
func DateTimeTodayHMS(hours int, mins int, secs int) DateTime {
	d := ToDateTime(time.Now())
	return d.SetHMS(hours, mins, secs)
}

// NeverTime возвращает дату-время в далёком прошлом
func NeverTime() DateTime {
	t, _ := time.ParseInLocation(DateTimeLayout, "1990-01-01 00:00:00", defaultLocation)
	return ToDateTime(t)
}

// setDefaultLayoutIfEmpty устанавливает шаблон вывода даты-времени
// по умолчанию, если шаблон не установлен
func (d *DateTime) setDefaultLayoutIfEmpty() {
	if strings.TrimSpace(d.Layout) == "" {
		d.Layout = DateTimeLayout
	}
}

// setDefaultLayoutIfEmpty устанавливает шаблон вывода даты по умолчанию, если шаблон не установлен
func (d *Date) setDefaultLayoutIfEmpty() {
	if strings.TrimSpace(d.Layout) == "" {
		d.Layout = DateLayout
	}
}

// SetHMS устанавливает значения часов, минут и секунд
func (d DateTime) SetHMS(hours int, mins int, secs int) DateTime {
	t := d.Time
	dS := fmt.Sprintf("%s %02d:%02d:%02d", t.Format(DateLayout), hours, mins, secs)
	dt, _ := time.ParseInLocation(DateTimeLayout, dS, defaultLocation)
	d.Time = dt
	return d
}

// ConvertToDate преобразует дату-время в дату
func (d DateTime) ConvertToDate() Date {
	dS := d.Time.Format(DateLayout)
	t, _ := time.ParseInLocation(DateLayout, dS, defaultLocation)
	return ToDate(t)
}

// ConvertToDateTimeHMS преобразует дату в дату-время с заданными значениями часов, минут и секунд
func (d Date) ConvertToDateTimeHMS(hours int, mins int, secs int) DateTime {
	dt := DateTime{
		Time:   d.Time,
		Layout: d.Layout,
	}
	dt = dt.SetHMS(hours, mins, secs)

	return dt
}

// After возвращает true если дата d позднее d1, иначе false
// Сравнение с точностью до дня.
func (d Date) After(d1 Date) bool {
	return d.Time.After(d1.Time)
}

// Before возвращает true если дата d ранее d1, иначе false
// Сравнение с точностью до дня.
func (d Date) Before(d1 Date) bool {
	return d.Time.Before(d1.Time)
}

// Between возвращает true если дата d находится в интервале дат (d1; d2), иначе false
// Сравнение с точностью до дня.
func (d Date) Between(d1, d2 Date) bool {
	return d.After(d1) && d.Before(d2)
}

// After возвращает true если дата-время obj позднее d1, иначе false
func (d DateTime) After(d1 DateTime) bool {
	return d.Time.After(d1.Time)
}

// Before возвращает true если дата-время obj ранее d1, иначе false
func (d DateTime) Before(d1 DateTime) bool {
	return d.Time.Before(d1.Time)
}

// Between возвращает true если дата-время d находится в интервале
// даты-времени (d1; d2), иначе false
func (d DateTime) Between(d1, d2 DateTime) bool {
	return d.After(d1) && d.Before(d2)
}

// UnmarshalJSON - правило преобразования поля JSON в объект DateTime
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

// MarshalJSON - правило преобразования объекта DateTime в поле JSON
func (d DateTime) MarshalJSON() ([]byte, error) {
	d.setDefaultLayoutIfEmpty()
	return []byte(strconv.Quote(d.String())), nil
}

// UnmarshalJSON - правило преобразования поля JSON в объект Date
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

// MarshalJSON - правило преобразования объекта Date в поле JSON
func (d Date) MarshalJSON() ([]byte, error) {
	d.setDefaultLayoutIfEmpty()
	return []byte(strconv.Quote(d.String())), nil
}

// String преобразует объект DateTime в строку согласно заданного шаблона
func (d DateTime) String() string {
	d.setDefaultLayoutIfEmpty()
	return d.Time.Format(d.Layout)
}

// String преобразует объект Date в строку согласно заданного шаблона
func (d Date) String() string {
	d.setDefaultLayoutIfEmpty()
	return d.Time.Format(d.Layout)
}

func scanInternal(value interface{}) (time.Time, error) {
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
	t, err := scanInternal(value)
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
