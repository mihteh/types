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

/*
	Тип для хранения даты-времени
*/
type DateTime struct {
	time.Time
	Layout string
}

/*
	Тип для хранения даты
*/
type Date struct {
	time.Time
	Layout string
}

/*
	Шаблоны вывода
*/
const (
	DateTimeLayout        = "2006-01-02 15:04:05"
	DateLayout            = "2006-01-02"
	GraphsDateLayout      = "02.01.2006"
	GraphsDateShortLayout = "02.01"
)

var defaultLocation *time.Location

/*
	Задаёт часовой пояс по умолчанию
*/
func init() {
	var err error
	defaultLocation, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Ошибка time.LoadLocation")
	}
}

/*
	Возвращает время типом DateTime по стандартному шаблону
*/
func ToDateTime(t time.Time) DateTime {
	return DateTime{t.In(defaultLocation), DateTimeLayout}
}

/*
	Возвращает время типом Date по стандартному шаблону
*/
func ToDate(t time.Time) Date {
	dS := fmt.Sprintf("%s %02d:%02d:%02d", t.Format(DateLayout), 0, 0, 0)
	d, _ := time.ParseInLocation(DateTimeLayout, dS, defaultLocation)
	return Date{d.In(defaultLocation), DateLayout}
}

/*
	Получает количество полных дней, прошедших от a до b
	Если a было вчера, b - сегодня, то возвращается 1
	Если b было раньше чем a, то возвращается отрицательное число.
*/
func (a Date) DaysBefore(b Date) int {
	return int(b.Time.Sub(a.Time).Hours() / 24)
}

/*
	Преобразует строку по стандартному шаблону даты-времени в дату-время
*/
func StringToDateTime(s string) (DateTime, error) {
	t, err := time.ParseInLocation(DateTimeLayout, s, defaultLocation)
	if err != nil {
		return DateTime{}, err
	}
	return ToDateTime(t), nil
}

/*
	Преобразует строку по стандартному шаблону даты в дату-время с заданным значением
	часов, минут и секунд
*/
func StringDateToDateTimeHMS(s string, hours int, mins int, secs int) (DateTime, error) {
	t, err := time.ParseInLocation(DateLayout, s, defaultLocation)
	if err != nil {
		return DateTime{}, err
	}
	d := ToDateTime(t)
	d = d.SetHMS(hours, mins, secs)
	return d, nil
}

/*
	Преобразует строку по стандартному шаблону даты в дату
*/
func StringToDate(s string) (Date, error) {
	t, err := time.ParseInLocation(DateLayout, s, defaultLocation)
	if err != nil {
		return Date{}, err
	}
	return ToDate(t), nil
}

/*
	Возвращает дату в далёком прошлом
*/
func NeverDate() Date {
	t, _ := time.ParseInLocation(DateLayout, "1990-01-01", defaultLocation)
	return ToDate(t)
}

/*
	Возвращает дату сегодня
*/
func DateNow() Date {
	return ToDate(time.Now())
}

/*
	Возвращает дату-время сейчас
*/
func DateTimeNow() DateTime {
	return ToDateTime(time.Now())
}

/*
	Возвращает дату-время сегодня в заданными значениями
	часов, минут, секунд
*/
func DateTimeTodayHMS(hours int, mins int, secs int) DateTime {
	d := ToDateTime(time.Now())
	return d.SetHMS(hours, mins, secs)
}

/*
	Возвращает дату-время в далёком прошлом
*/
func NeverTime() DateTime {
	t, _ := time.ParseInLocation(DateTimeLayout, "1990-01-01 00:00:00", defaultLocation)
	return ToDateTime(t)
}

/*
	Устанавливает шаблон вывода даты-времени по умолчанию, если шаблон не установлен
*/
func (obj *DateTime) setDefaultLayoutIfEmpty() {
	if strings.TrimSpace(obj.Layout) == "" {
		obj.Layout = DateTimeLayout
	}
}

/*
	Устанавливает шаблон вывода даты по умолчанию, если шаблон не установлен
*/
func (obj *Date) setDefaultLayoutIfEmpty() {
	if strings.TrimSpace(obj.Layout) == "" {
		obj.Layout = DateLayout
	}
}

/*
	Устанавливает значения часов, минут и секунд
*/
func (obj DateTime) SetHMS(hours int, mins int, secs int) DateTime {
	t := obj.Time
	dS := fmt.Sprintf("%s %02d:%02d:%02d", t.Format(DateLayout), hours, mins, secs)
	d, _ := time.ParseInLocation(DateTimeLayout, dS, defaultLocation)
	obj.Time = d
	return obj
}

/*
	Преобразует дату-время в дату
*/
func (obj DateTime) ConvertToDate() Date {
	dS := obj.Time.Format(DateLayout)
	t, _ := time.ParseInLocation(DateLayout, dS, defaultLocation)
	return ToDate(t)
}

/*
	Преобразует дату в дату-время с заданными значениями часов, минут и секунд
*/
func (obj Date) ConvertToDateTimeHMS(hours int, mins int, secs int) DateTime {
	d := DateTime{
		Time:   obj.Time,
		Layout: obj.Layout,
	}
	d = d.SetHMS(hours, mins, secs)

	return d
}

/*
	Возвращает true если дата obj позднее d, иначе false
	Сравнение с точностью до дня.
*/
func (obj Date) After(d Date) bool {
	return obj.Time.After(d.Time)
}

/*
	Возвращает true если дата obj ранее d, иначе false
	Сравнение с точностью до дня.
*/
func (obj Date) Before(d Date) bool {
	return obj.Time.Before(d.Time)
}

/*
	Возвращает true если дата obj находится в интервале дат (d1; d2), иначе false
	Сравнение с точностью до дня.
*/
func (obj Date) Between(d1, d2 Date) bool {
	return obj.After(d1) && obj.Before(d2)
}

/*
	Возвращает true если дата-время obj позднее d, иначе false
*/
func (obj DateTime) After(d DateTime) bool {
	return obj.Time.After(d.Time)
}

/*
	Возвращает true если дата-время obj ранее d, иначе false
*/
func (obj DateTime) Before(d DateTime) bool {
	return obj.Time.Before(d.Time)
}

/*
	Возвращает true если дата-время obj находится в интервале
	даты-времени (d1; d2), иначе false
*/
func (obj DateTime) Between(d1, d2 DateTime) bool {
	return obj.After(d1) && obj.Before(d2)
}

/*
	Правило преобразования поля JSON в объект DateTime
*/
func (obj *DateTime) UnmarshalJSON(data []byte) error {
	obj.setDefaultLayoutIfEmpty()
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.ParseInLocation(obj.Layout, s, defaultLocation)
	if err != nil {
		return err
	}
	obj.Time = t
	return nil
}

/*
	Правило преобразования объекта DateTime в поле JSON
*/
func (obj DateTime) MarshalJSON() ([]byte, error) {
	obj.setDefaultLayoutIfEmpty()
	return []byte(strconv.Quote(obj.String())), nil
}

/*
	Правило преобразования поля JSON в объект Date
*/
func (obj *Date) UnmarshalJSON(data []byte) error {
	obj.setDefaultLayoutIfEmpty()
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.ParseInLocation(obj.Layout, s, defaultLocation)
	if err != nil {
		return err
	}
	obj.Time = t
	return nil
}

/*
	Правило преобразования объекта Date в поле JSON
*/
func (obj Date) MarshalJSON() ([]byte, error) {
	obj.setDefaultLayoutIfEmpty()
	return []byte(strconv.Quote(obj.String())), nil
}

/*
	Преобразует объект DateTime в строку согласно заданного шаблона
*/
func (obj DateTime) String() string {
	obj.setDefaultLayoutIfEmpty()
	return obj.Time.Format(obj.Layout)
}

/*
	Преобразует объект Date в строку согласно заданного шаблона
*/
func (obj Date) String() string {
	obj.setDefaultLayoutIfEmpty()
	return obj.Time.Format(obj.Layout)
}

/*
	Преобразует значение времени в БД к типу DateTime
*/
func (dt *DateTime) Scan(value interface{}) error {
	dt.setDefaultLayoutIfEmpty()
	if value == nil {
		return nil
	}

	time, ok := value.(time.Time)
	if !ok {
		return errors.New("Ошибка преобразования значения к типу time.Time")
	}

	dt.Time = time.In(defaultLocation)
	return nil
}

/*
	Преобразует значение типа DateTime к значению в БД
*/
func (dt DateTime) Value() (driver.Value, error) {
	return dt.Time.In(defaultLocation).Format(DateTimeLayout), nil
}

/*
	Преобразует значение времени в БД к типу Date
*/
func (dt *Date) Scan(value interface{}) error {
	dt.setDefaultLayoutIfEmpty()
	if value == nil {
		return nil
	}

	time, ok := value.(time.Time)
	if !ok {
		return errors.New("Ошибка преобразования значения к типу time.Time")
	}

	dt.Time = time.In(defaultLocation)
	return nil
}

/*
	Преобразует значение типа Date к значению в БД
*/
func (dt Date) Value() (driver.Value, error) {
	return dt.Time.In(defaultLocation).Format(DateLayout), nil
}
