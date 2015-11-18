package types

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

/*
	Проверяет функции ToDateTime(), String() и StringToDateTime()
*/
func TestDateTimeStringConversion(t *testing.T) {
	timeMonthAgo := time.Now().AddDate(0, -1, 0)
	dateTime := ToDateTime(timeMonthAgo)
	dateTimeString := dateTime.String()

	dateTimeFromString, err := StringToDateTime(dateTimeString)
	if err != nil {
		t.Fatal(err)
	}

	expectedString := dateTimeFromString.String()
	if expectedString != dateTimeString {
		t.Fatalf("Ошибка %v. Ожидалось %v, получено %v", err, expectedString, dateTimeString)
	}
}

/*
	Проверяет функции ToDate(), String() и StringToDate()
*/
func TestDateStringConversion(t *testing.T) {
	timeMonthAgo := time.Now().AddDate(0, -1, 0)
	date := ToDate(timeMonthAgo)
	dateString := date.String()

	dateFromString, err := StringToDate(dateString)
	if err != nil {
		t.Fatal(err)
	}

	expectedString := dateFromString.String()
	if expectedString != dateString {
		t.Fatalf("Ошибка %v. Ожидалось %v, получено %v", err, expectedString, dateString)
	}
}

func TestDaysBefore(t *testing.T) {
	dateToday := DateNow()
	dateYesterday := ToDate(time.Now().AddDate(0, 0, -1))
	dateTomorrow := ToDate(time.Now().AddDate(0, 0, 1))

	if d := dateToday.DaysBefore(dateToday); d != 0 {
		t.Fatalf("Ожидалось 0, получено %d дней", d)
	}

	if d := dateToday.DaysBefore(dateYesterday); d != -1 {
		t.Fatalf("Ожидалось -1, получено %d дней", d)
	}

	if d := dateToday.DaysBefore(dateTomorrow); d != 1 {
		t.Fatalf("Ожидалось 1, получено %d дней", d)
	}

	if d := dateYesterday.DaysBefore(dateTomorrow); d != 2 {
		t.Fatalf("Ожидалось 2, получено %d дней", d)
	}
}

func TestStringToDateTimeHMS(t *testing.T) {
	s := "2015-07-30"
	dt, err := StringDateToDateTimeHMS(s, 5, 6, 7)
	if err != nil {
		t.Fatal(err)
	}
	sExpected := "2015-07-30 05:06:07"
	sReceived := dt.String()
	if sExpected != sReceived {
		t.Fatalf("Ожидалось получить строку %s, получена строка %s", sExpected, sReceived)
	}
}

func TestNeverDate(t *testing.T) {
	nd := NeverDate()
	sExpected := "1990-01-01"
	sReceived := nd.String()
	if sExpected != sReceived {
		t.Fatalf("Ожидалось получить строку %s, получена строка %s", sExpected, sReceived)
	}
}

func TestNeverTime(t *testing.T) {
	ndt := NeverTime()
	sExpected := "1990-01-01 00:00:00"
	sReceived := ndt.String()
	if sExpected != sReceived {
		t.Fatalf("Ожидалось получить строку %s, получена строка %s", sExpected, sReceived)
	}
}

/*
	Проверяет возврат текущих даты, даты-времени и
	преобразование Date -> DateTime с учётом HMS
*/
func TestNow(t *testing.T) {
	dateTimeNow := DateTimeNow()
	dateNow := dateTimeNow.ConvertToDate()
	dateTimeAfterConversion := dateNow.ConvertToDateTimeHMS(
		dateTimeNow.Hour(), dateTimeNow.Minute(), dateTimeNow.Second())
	dateTimeAfterConversion.Layout = DateTimeLayout

	sExpected := dateTimeNow.String()
	sReceived := dateTimeAfterConversion.String()
	if sExpected != sReceived {
		t.Fatalf("Ожидалось получить строку %s, получена строка %s", sExpected, sReceived)
	}
}

func TestDateTimeTodayHMS(t *testing.T) {
	dt := DateTimeTodayHMS(23, 59, 59)
	sExpected := DateNow().String() + " 23:59:59"
	sReceived := dt.String()
	if sExpected != sReceived {
		t.Fatalf("Ожидалось получить строку %s, получена строка %s", sExpected, sReceived)
	}
}

func TestDateSetDefaultLayoutIfEmpty(t *testing.T) {
	d := DateNow()
	d.Layout = ""
	d.setDefaultLayoutIfEmpty()
	if d.Layout != DateLayout {
		t.Fatal("Установлен неправильный Layout")
	}
}

func TestDateTimeSetDefaultLayoutIfEmpty(t *testing.T) {
	dt := DateTimeNow()
	dt.Layout = ""
	dt.setDefaultLayoutIfEmpty()
	if dt.Layout != DateTimeLayout {
		t.Fatal("Установлен неправильный Layout")
	}
}

func TestDateBeforeAfterBetween(t *testing.T) {
	dateToday := DateNow()
	dateYesterday := ToDate(time.Now().AddDate(0, 0, -1))
	dateTomorrow := ToDate(time.Now().AddDate(0, 0, 1))

	if dateToday.Before(dateToday) {
		t.Fatal("Сегодня не может быть раньше чем сегодня")
	}
	if dateToday.After(dateToday) {
		t.Fatal("Сегодня не может быть позже чем сегодня")
	}
	if dateToday.Between(dateToday, dateTomorrow) {
		t.Fatal("Сегодня не может быть между сегодня и завтра")
	}
	if dateToday.Between(dateYesterday, dateToday) {
		t.Fatal("Сегодня не может быть между вчера и сегодня")
	}

	if dateToday.Before(dateYesterday) {
		t.Fatal("Сегодня не может быть раньше чем вчера")
	}
	if dateYesterday.After(dateToday) {
		t.Fatal("Вчера не может быть позже чем сегодня")
	}
	if dateYesterday.Between(dateToday, dateTomorrow) {
		t.Fatal("Вчера не может быть между сегодня и завтра")
	}

	if !dateToday.Before(dateTomorrow) {
		t.Fatal("Сегодня должно быть раньше чем завтра")
	}
	if !dateToday.After(dateYesterday) {
		t.Fatal("Сегодня должно быть позже чем вчера")
	}
	if !dateToday.Between(dateYesterday, dateTomorrow) {
		t.Fatal("Сегодня должно быть между вчера и завтра")
	}
}

func TestDateTimeBeforeAfterBetween(t *testing.T) {
	dateNow := DateTimeNow()
	dateMinBefore := ToDateTime(time.Now().Add(-time.Minute))
	dateMinAfter := ToDateTime(time.Now().Add(time.Minute))

	if dateNow.Before(dateNow) {
		t.Fatal("Сейчас не может быть раньше чем сейчас")
	}
	if dateNow.After(dateNow) {
		t.Fatal("Сейчас не может быть позже чем сейчас")
	}
	if dateNow.Between(dateNow, dateMinAfter) {
		t.Fatal("Сейчас не может быть между сейчас и через минуту")
	}
	if dateNow.Between(dateMinBefore, dateNow) {
		t.Fatal("Сейчас не может быть между минуту назад и сейчас")
	}

	if dateNow.Before(dateMinBefore) {
		t.Fatal("Сейчас не может быть раньше чем минуту назад")
	}
	if dateMinBefore.After(dateNow) {
		t.Fatal("Минуту назад не может быть позже чем сейчас")
	}
	if dateMinBefore.Between(dateNow, dateMinAfter) {
		t.Fatal("Минуту назад не может быть между сейчас и через минуту")
	}

	if !dateNow.Before(dateMinAfter) {
		t.Fatal("Сейчас должно быть раньше чем через минуту")
	}
	if !dateNow.After(dateMinBefore) {
		t.Fatal("Сейчас должно быть позже чем минуту назад")
	}
	if !dateNow.Between(dateMinBefore, dateMinAfter) {
		t.Fatal("Сейчас должно быть между минуту назад и через минуту")
	}
}

func TestDateTimeJSON(t *testing.T) {
	dtString := "2015-07-30 20:58:59"
	jsonExpected := fmt.Sprintf(`"%s"`, dtString)
	dt, err := StringToDateTime(dtString)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(dt)
	if err != nil {
		t.Fatal(err)
	}
	jsonReceived := string(b)
	if jsonExpected != jsonReceived {
		t.Fatalf("Ошибка Marshal. Ожидалось получить JSON %s, получен JSON %s", jsonExpected, jsonReceived)
	}

	var dtFromJson DateTime
	if err := json.Unmarshal(b, &dtFromJson); err != nil {
		t.Fatal(err)
	}
	dtFromJsonString := dtFromJson.String()
	if dtFromJsonString != dtString {
		t.Fatalf("Ошибка Unmarshal. Ожидалось получить JSON %s, получен JSON %s", dtString, dtFromJsonString)
	}
}

func TestDateTimeUnmarshalBadJSON(t *testing.T) {
	b := []byte(`"wrong"`)
	var dtFromJson DateTime
	if err := json.Unmarshal(b, &dtFromJson); err == nil {
		t.Fatal("Ожидалась ошибка")
	}
}

func TestDateJSON(t *testing.T) {
	dString := "2015-07-30"
	jsonExpected := fmt.Sprintf(`"%s"`, dString)
	d, err := StringToDate(dString)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	jsonReceived := string(b)
	if jsonExpected != jsonReceived {
		t.Fatalf("Ошибка Marshal. Ожидалось получить JSON %s, получен JSON %s", jsonExpected, jsonReceived)
	}

	var dFromJson Date
	if err := json.Unmarshal(b, &dFromJson); err != nil {
		t.Fatal(err)
	}
	dFromJsonString := dFromJson.String()
	if dFromJsonString != dString {
		t.Fatalf("Ошибка Unmarshal. Ожидалось получить JSON %s, получен JSON %s", dString, dFromJsonString)
	}
}

func TestDateUnmarshalBadJSON(t *testing.T) {
	b := []byte(`"wrong"`)
	var dFromJson Date
	if err := json.Unmarshal(b, &dFromJson); err == nil {
		t.Fatal("Ожидалась ошибка")
	}
}
