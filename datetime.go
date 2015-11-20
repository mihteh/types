package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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

type timeModifier interface {
	timeSetter
	timeLayouter
}

type timeSetter interface {
	setTime(time.Time)
}

type timeLayouter interface {
	getLayout() string
	setLayout(layout string)
	fixLayout()
}

// Шаблоны для сериализации
const (
	DateTimeLayout        = "2006-01-02 15:04:05"
	DateLayout            = "2006-01-02"
	GraphsDateLayout      = "02.01.2006"
	GraphsDateShortLayout = "02.01"
)

// defaultLocation хранит значение по умолчанию для Location
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
	dS := t.Format(DateTimeLayout)
	parsedDateTime, _ := time.ParseInLocation(DateTimeLayout, dS, defaultLocation)
	dt := NewDateTime()
	dt.setTime(parsedDateTime)
	return dt
}

// ToDate формирует объект типа Date на основе времени t и шаблона DateLayout
func ToDate(t time.Time) Date {
	dS := fmt.Sprintf("%s %02d:%02d:%02d", t.Format(DateLayout), 0, 0, 0)
	parsedTime, _ := time.ParseInLocation(DateTimeLayout, dS, defaultLocation)
	d := NewDate()
	d.setTime(parsedTime)
	return d
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

// DateNow возвращает объект Date, соответствующий дате сегодня
func DateNow() Date {
	return ToDate(time.Now())
}

// DateTimeNow возвращает объект DateTime, соответствующий дате-времени сейчас
func DateTimeNow() DateTime {
	return ToDateTime(time.Now())
}

// NewDate создаёт новый объект типа Date с шаблоном вывода по умолчанию DateLayout
func NewDate() Date {
	d := Date{}
	d.setLayout(DateLayout)
	return d
}

// NewDateTime создаёт новый объект типа DateTime
// с шаблоном вывода по умолчанию DateTimeLayout
func NewDateTime() DateTime {
	d := DateTime{}
	d.setLayout(DateTimeLayout)
	return d
}

// DateTimeTodayHMS возвращает объект DateTime, соответствующий дате сегодня,
// с установленными значениями часов, минут, секунд согласно заданным параметрам
// hours, mins, secs соответственно.
func DateTimeTodayHMS(hours int, mins int, secs int) DateTime {
	d := ToDateTime(time.Now())
	return d.SetHMS(hours, mins, secs)
}

// NeverDate возвращает объект Date, соответствующий дате в далёком прошлом
func NeverDate() Date {
	t, _ := time.ParseInLocation(DateLayout, "0001-01-01", defaultLocation)
	return ToDate(t)
}

// NeverTime возвращает объект DateTime, соответствующий дате-времени в далёком прошлом
func NeverTime() DateTime {
	t, _ := time.ParseInLocation(DateTimeLayout, "0001-01-01 00:00:00", defaultLocation)
	return ToDateTime(t)
}

// setTime устанавливает время в объекте Date без учёта Location
func (d *Date) setTime(t time.Time) {
	d.Time = t
}

// setTime устанавливает время в объекте DateTime без учёта Location
func (d *DateTime) setTime(t time.Time) {
	d.Time = t
}

// getLayout возвращает строку шаблона вывода в объекте Date
func (d Date) getLayout() string {
	return d.Layout
}

// setLayout устанавливает Layout в объекте Date
func (d *Date) setLayout(layout string) {
	d.Layout = layout
}

// fixLayout устанавливает Layout в объекте Date на DateLayout если он не определён
func (d *Date) fixLayout() {
	if d.getLayout() == "" {
		d.setLayout(DateLayout)
	}
}

// getLayout возвращает строку шаблона вывода в объекте DateTime
func (d DateTime) getLayout() string {
	return d.Layout
}

// setLayout устанавливает Layout в объекте DateTime
func (d *DateTime) setLayout(layout string) {
	d.Layout = layout
}

// fixLayout устанавливает Layout в объекте DateTime на DateTimeLayout если он не определён
func (d *DateTime) fixLayout() {
	if d.getLayout() == "" {
		d.setLayout(DateTimeLayout)
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
	dt := NewDateTime()
	dt.Layout = d.Layout
	dt.setTime(d.Time)
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

// Equal возвращает true если дата d равна дате d1, иначе false
// Сравнение происходит с точностью до дня.
func (d Date) Equal(d1 Date) bool {
	return d.Time.Equal(d1.Time)
}

// After возвращает true если дата-время d позднее d1, иначе false
// Сравнение происходит с точностью до секунд.
func (d DateTime) After(d1 DateTime) bool {
	return d.Time.After(d1.Time)
}

// Before возвращает true если дата-время d ранее d1, иначе false
// Сравнение происходит с точностью до секунд.
func (d DateTime) Before(d1 DateTime) bool {
	return d.Time.Before(d1.Time)
}

// Between возвращает true если дата-время d находится в интервале
// даты-времени (d1; d2), иначе false
// Сравнение происходит с точностью до секунд.
func (d DateTime) Between(d1, d2 DateTime) bool {
	return d.After(d1) && d.Before(d2)
}

// Equal возвращает true если дата-время d равна дате-времени d1, иначе false
// Сравнение происходит с точностью до секунд.
func (d DateTime) Equal(d1 DateTime) bool {
	return d.Time.Equal(d1.Time)
}

// parse устанавливает время в объекте, реализующем интерфейс timeModifier
// на основе строки s и Location defaultLocation
func parse(d timeModifier, s string) error {
	t, err := time.ParseInLocation(d.getLayout(), s, defaultLocation)
	if err == nil {
		d.setTime(t)
	}
	return err
}

func unmarshalJSON(data []byte, to timeModifier) error {
	var s string
	to.fixLayout()
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return parse(to, s)
}

// UnmarshalJSON - реализует интерфейс json.Unmarshaler для объекта DateTime
// десериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d *DateTime) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, d)
}

// MarshalJSON - реализует интерфейс json.Marshaler для объекта DateTime
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d DateTime) MarshalJSON() ([]byte, error) {
	d.fixLayout()
	return []byte(strconv.Quote(d.String())), nil
}

// UnmarshalJSON - реализует интерфейс json.Unmarshaler для объекта Date
// десериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d *Date) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, d)
}

// MarshalJSON - реализует интерфейс json.Marshaler для объекта Date
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d Date) MarshalJSON() ([]byte, error) {
	d.fixLayout()
	return []byte(strconv.Quote(d.String())), nil
}

// String преобразует объект DateTime в строку согласно шаблона в свойстве Layout
func (d DateTime) String() string {
	return d.Time.Format(d.Layout)
}

// String преобразует объект Date в строку согласно шаблона в свойстве Layout
func (d Date) String() string {
	return d.Time.Format(d.Layout)
}

func scan(from interface{}, to timeSetter) error {
	if from == nil {
		return nil
	}

	t := time.Time{}
	t, ok := from.(time.Time)
	if !ok {
		return errors.New("Ошибка преобразования значения к типу time.Time")
	}

	to.setTime(t.In(defaultLocation))
	return nil
}

// Scan преобразует значение времени в БД к типу DateTime
// Реализует интерфейс sql.Scanner
func (d *DateTime) Scan(value interface{}) error {
	d.fixLayout()
	return scan(value, d)
}

// Value преобразует значение типа DateTime к значению в БД
// Реализует интерфейс driver.Valuer
func (d DateTime) Value() (driver.Value, error) {
	return d.Time.In(defaultLocation).Format(DateTimeLayout), nil
}

// Scan преобразует значение времени в БД к типу Date
// Реализует интерфейс sql.Scanner
func (d *Date) Scan(value interface{}) error {
	d.fixLayout()
	return scan(value, d)
}

// Value преобразует значение типа Date к значению в БД
// Реализует интерфейс driver.Valuer
func (d Date) Value() (driver.Value, error) {
	return d.Time.In(defaultLocation).Format(DateLayout), nil
}
