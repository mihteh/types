package types

import (
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
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

// Pointer возвращает указатель на объект DateTime
func (d DateTime) Pointer() *DateTime {
	return &d
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

// Pointer возвращает указатель на объект Date
func (d Date) Pointer() *Date {
	return &d
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

// EncodeValues реализует интерфейс query.Encoder для объекта DateTime
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d DateTime) EncodeValues(key string, v *url.Values) error {
	d.fixLayout()
	v.Set(key, d.String())
	return nil
}

// UnmarshalXML реализует интерфейс xml.Unmarshaler для объекта DateTime
// десериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d *DateTime) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := decoder.DecodeElement(&content, &start); err != nil {
		return err
	}
	d.fixLayout()
	return parse(d, content)
}

// MarshalXML реализует интерфейс xml.Marshaler для объекта DateTime
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d DateTime) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	d.fixLayout()
	return encoder.EncodeElement(d.String(), start)
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

// EncodeValues реализует интерфейс query.Encoder для объекта Date
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d Date) EncodeValues(key string, v *url.Values) error {
	d.fixLayout()
	v.Set(key, d.String())
	return nil
}

// UnmarshalXML реализует интерфейс xml.Unmarshaler для объекта Date
// десериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d *Date) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := decoder.DecodeElement(&content, &start); err != nil {
		return err
	}
	d.fixLayout()
	return parse(d, content)
}

// MarshalXML реализует интерфейс xml.Marshaler для объекта Date
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d Date) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	d.fixLayout()
	return encoder.EncodeElement(d.String(), start)
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

// NullDateTime это вспомогательный тип, необходимый для реализации
// интерфейса Valuer на указателе
type NullDateTime struct {
	DateTime
	Valid bool
}

// NullDate это вспомогательный тип, необходимый для реализации
// интерфейса Valuer на указателе
type NullDate struct {
	Date
	Valid bool
}

// Scan преобразует значение времени в БД к типу NullDateTime
// Реализует интерфейс sql.Scanner
func (d *NullDateTime) Scan(value interface{}) error {
	d.fixLayout()
	if value == nil {
		d.Valid = false
		return nil
	}
	d.Valid = true
	return scan(value, d)
}

// Value преобразует значение типа NullDateTime к значению в БД
// Реализует интерфейс driver.Valuer
func (d NullDateTime) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time.In(defaultLocation).Format(DateTimeLayout), nil
}

// Scan преобразует значение времени в БД к типу NullDate
// Реализует интерфейс sql.Scanner
func (d *NullDate) Scan(value interface{}) error {
	d.fixLayout()
	if value == nil {
		d.Valid = false
		return nil
	}
	d.Valid = true
	return scan(value, d)
}

// Value преобразует значение типа NullDate к значению в БД
// Реализует интерфейс driver.Valuer
func (d NullDate) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time.In(defaultLocation).Format(DateLayout), nil
}

// Nullable преобразует тип DateTime в тип NullDateTime
func (dt DateTime) Nullable() NullDateTime {
	ndt := NullDateTime{
		DateTime: dt,
		Valid:    true,
	}
	return ndt
}

// Nullable преобразует тип Date в тип NullDate
func (d Date) Nullable() NullDate {
	nd := NullDate{
		Date:  d,
		Valid: true,
	}
	return nd
}

// isJSONBytesNil проверяет, содержит ли JSON nil значение
func isJSONBytesNil(data []byte) (bool, error) {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return true, err
	}
	switch v.(type) {
	case nil:
		return true, nil
	}
	return false, nil
}

// UnmarshalJSON - реализует интерфейс json.Unmarshaler для объекта NullDateTime
// десериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d *NullDateTime) UnmarshalJSON(data []byte) error {
	isNil, err := isJSONBytesNil(data)
	if err != nil {
		return err
	}
	if isNil {
		d.Valid = false
		return nil
	}

	pobj := &DateTime{}
	err = pobj.UnmarshalJSON(data)
	d.DateTime = *pobj
	d.Valid = err == nil
	return err
}

// UnmarshalJSON - реализует интерфейс json.Unmarshaler для объекта NullDate
// десериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d *NullDate) UnmarshalJSON(data []byte) error {
	isNil, err := isJSONBytesNil(data)
	if err != nil {
		return err
	}
	if isNil {
		d.Valid = false
		return nil
	}

	pobj := &Date{}
	err = pobj.UnmarshalJSON(data)
	d.Date = *pobj
	d.Valid = err == nil
	return err
}

// MarshalJSON - реализует интерфейс json.Marshaler для объекта NullDateTime
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d NullDateTime) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	d.fixLayout()
	return []byte(strconv.Quote(d.String())), nil
}

// String преобразует объект NullDateTime в строку согласно шаблона в свойстве Layout
func (d NullDateTime) String() string {
	if !d.Valid {
		return "null"
	}
	return d.DateTime.String()
}

// String преобразует объект NullDate в строку согласно шаблона в свойстве Layout
func (d NullDate) String() string {
	if !d.Valid {
		return "null"
	}
	return d.Date.String()
}

// UnmarshalXML реализует интерфейс xml.Unmarshaler для объекта NullDateTime
// десериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d *NullDateTime) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := decoder.DecodeElement(&content, &start); err != nil {
		return err
	}
	if content == "" {
		d.Valid = false
		return nil
	}
	pobj := &DateTime{}
	pobj.fixLayout()
	err := parse(pobj, content)
	d.DateTime = *pobj
	d.Valid = err == nil
	return err
}

// MarshalXML реализует интерфейс xml.Marshaler для объекта NullDateTime
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d NullDateTime) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	if !d.Valid {
		return encoder.EncodeElement(nil, start)
	}
	d.fixLayout()
	return encoder.EncodeElement(d.String(), start)
}

// UnmarshalXML реализует интерфейс xml.Unmarshaler для объекта NullDate
// десериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d *NullDate) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := decoder.DecodeElement(&content, &start); err != nil {
		return err
	}
	if content == "" {
		d.Valid = false
		return nil
	}
	pobj := &Date{}
	pobj.fixLayout()
	err := parse(pobj, content)
	d.Date = *pobj
	d.Valid = err == nil
	return err
}

// MarshalXML реализует интерфейс xml.Marshaler для объекта NullDate
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d NullDate) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	if !d.Valid {
		return encoder.EncodeElement(nil, start)
	}
	d.fixLayout()
	return encoder.EncodeElement(d.String(), start)
}

// MarshalJSON - реализует интерфейс json.Marshaler для объекта NullDate
// сериализация происходит с учётом шаблона, заданного в свойстве Layout
func (d NullDate) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	d.fixLayout()
	return []byte(strconv.Quote(d.String())), nil
}

// MakeNullDate возвразает NullDate со значением NULL
func MakeNullDate() NullDate {
	return NullDate{Valid: false}
}

// MakeNullDateTime возвращает NullDateTime со значением NULL
func MakeNullDateTime() NullDateTime {
	return NullDateTime{Valid: false}
}

// OldNeverTime это устаревшая версия метода NeverTime()
func OldNeverTime() DateTime {
	t, _ := time.ParseInLocation(DateTimeLayout, "1990-01-01 00:00:00", defaultLocation)
	return ToDateTime(t)
}

// OldNeverDate это устаревшая версия метода NeverDate()
func OldNeverDate() Date {
	t, _ := time.ParseInLocation(DateLayout, "1990-01-01", defaultLocation)
	return ToDate(t)
}
