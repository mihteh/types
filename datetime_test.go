package types

import (
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
