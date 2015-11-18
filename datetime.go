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

type DateTime struct {
	time.Time
	Layout string
}

type Date struct {
	time.Time
	Layout string
}

const (
	DateTimeLayout        = "2006-01-02 15:04:05"
	DateLayout            = "2006-01-02"
	GraphsDateLayout      = "02.01.2006"
	GraphsDateShortLayout = "02.01"
)

var defaultLocation *time.Location

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

func StringToDateTime(s string) (DateTime, error) {
	t, err := time.ParseInLocation(DateTimeLayout, s, defaultLocation)
	if err != nil {
		return DateTime{}, err
	}
	return ToDateTime(t), nil
}

func StringDateToDateTimeHMS(s string, hours int, mins int, secs int) (DateTime, error) {
	t, err := time.ParseInLocation(DateLayout, s, defaultLocation)
	if err != nil {
		return DateTime{}, err
	}
	d := ToDateTime(t)
	d = d.SetHMS(hours, mins, secs)
	return d, nil
}

func StringToDate(s string) (Date, error) {
	t, err := time.ParseInLocation(DateLayout, s, defaultLocation)
	if err != nil {
		return Date{}, err
	}
	return ToDate(t), nil
}

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

func NeverTime() DateTime {
	t, _ := time.ParseInLocation(DateTimeLayout, "1990-01-01 00:00:00", defaultLocation)
	return ToDateTime(t)
}

func (obj *DateTime) setDefaultLayoutIfEmpty() {
	if strings.TrimSpace(obj.Layout) == "" {
		obj.Layout = DateTimeLayout
	}
}

func (obj *Date) setDefaultLayoutIfEmpty() {
	if strings.TrimSpace(obj.Layout) == "" {
		obj.Layout = DateLayout
	}
}

// устанавливает часы, минуты, секунды в объекте типа DateTime
func (obj DateTime) SetHMS(hours int, mins int, secs int) DateTime {
	t := obj.Time
	dS := fmt.Sprintf("%s %02d:%02d:%02d", t.Format(DateLayout), hours, mins, secs)
	d, _ := time.ParseInLocation(DateTimeLayout, dS, defaultLocation)
	obj.Time = d
	return obj
}

func (obj DateTime) ConvertToDate() Date {
	dS := obj.Time.Format(DateLayout)
	t, _ := time.ParseInLocation(DateLayout, dS, defaultLocation)
	return ToDate(t)
}

func (obj Date) ConvertToDateTimeHMS(hours int, mins int, secs int) DateTime {
	d := DateTime{
		Time:   obj.Time,
		Layout: obj.Layout,
	}
	d = d.SetHMS(hours, mins, secs)

	return d
}

func (obj Date) After(d Date) bool {
	return obj.Time.After(d.Time)
}

func (obj Date) Before(d Date) bool {
	return obj.Time.Before(d.Time)
}

func (obj Date) Between(d1, d2 Date) bool {
	return obj.After(d1) && obj.Before(d2)
}

func (obj DateTime) After(d DateTime) bool {
	return obj.Time.After(d.Time)
}

func (obj DateTime) Before(d DateTime) bool {
	return obj.Time.Before(d.Time)
}

func (obj DateTime) Between(d1, d2 DateTime) bool {
	return obj.After(d1) && obj.Before(d2)
}

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

func (obj DateTime) MarshalJSON() ([]byte, error) {
	obj.setDefaultLayoutIfEmpty()
	return []byte(strconv.Quote(obj.String())), nil
}

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

func (obj Date) MarshalJSON() ([]byte, error) {
	obj.setDefaultLayoutIfEmpty()
	return []byte(strconv.Quote(obj.String())), nil
}

func (obj DateTime) String() string {
	obj.setDefaultLayoutIfEmpty()
	return obj.Time.Format(obj.Layout)
}

func (obj Date) String() string {
	obj.setDefaultLayoutIfEmpty()
	return obj.Time.Format(obj.Layout)
}

// Преобразует значение времени в БД к типу DateTime
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

// Преобразует значение типа DateTime к типу в БД
func (dt DateTime) Value() (driver.Value, error) {
	return dt.Time.In(defaultLocation).Format(DateTimeLayout), nil
}

// Преобразует значение времени в БД к типу Date
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

// Преобразует значение типа Date к типу в БД
func (dt Date) Value() (driver.Value, error) {
	return dt.Time.In(defaultLocation).Format(DateLayout), nil
}
