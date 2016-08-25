package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
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

func TestBadStringToDateTime(t *testing.T) {
	if _, err := StringToDateTime("wrong"); err == nil {
		t.Fatal("Ожидалась ошибка")
	}
}

func TestBadStringToDateTimeHMS(t *testing.T) {
	if _, err := StringDateToDateTimeHMS("wrong", 23, 59, 59); err == nil {
		t.Fatal("Ожидалась ошибка")
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

func TestBadStringToDate(t *testing.T) {
	if _, err := StringToDate("wrong"); err == nil {
		t.Fatal("Ожидалась ошибка")
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
	sExpected := "0001-01-01"
	sReceived := nd.String()
	if sExpected != sReceived {
		t.Fatalf("Ожидалось получить строку %s, получена строка %s", sExpected, sReceived)
	}
}

func TestNeverTime(t *testing.T) {
	ndt := NeverTime()
	sExpected := "0001-01-01 00:00:00"
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

func TestDateBeforeAfterBetweenEqual(t *testing.T) {
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

	if dateToday.Equal(dateYesterday) {
		t.Fatal("Сегодня не может быть вчера")
	}
	if dateToday.Equal(dateTomorrow) {
		t.Fatal("Сегодня не может быть завтра")
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
	if !dateToday.Equal(dateToday) {
		t.Fatal("Сегодня должно быть сегодня")
	}
}

func TestDateTimeBeforeAfterBetweenEqual(t *testing.T) {
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

	if dateNow.Equal(dateMinBefore) {
		t.Fatal("Сейчас не может быть через минуту")
	}
	if dateNow.Equal(dateMinAfter) {
		t.Fatal("Сейчас не может быть минуту назад")
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
	if !dateNow.Equal(dateNow) {
		t.Fatal("Сейчас должно быть сейчас")
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

	var dtFromJSON DateTime
	if err := json.Unmarshal(b, &dtFromJSON); err != nil {
		t.Fatal(err)
	}
	dtFromJSONString := dtFromJSON.String()
	if dtFromJSONString != dtString {
		t.Fatalf("Ошибка Unmarshal. Ожидалось получить JSON %s, получен JSON %s", dtString, dtFromJSONString)
	}
}

func TestDateTimeXML(t *testing.T) {
	dtString := "2015-07-30 20:58:59"
	xmlExpected := fmt.Sprintf("<DateTime>%s</DateTime>", dtString)
	dt, err := StringToDateTime(dtString)
	if err != nil {
		t.Fatal(err)
	}
	b, err := xml.Marshal(dt)
	if err != nil {
		t.Fatal(err)
	}
	xmlReceived := string(b)
	if xmlExpected != xmlReceived {
		t.Fatalf("Ошибка Marshal. Ожидалось получить XML %s, получен XML %s", xmlExpected, xmlReceived)
	}

	var dtFromXML DateTime
	if err := xml.Unmarshal(b, &dtFromXML); err != nil {
		t.Fatal(err)
	}
	dtFromXMLString := dtFromXML.String()
	if dtFromXMLString != dtString {
		t.Fatalf("Ошибка Unmarshal. Ожидалось получить XML %s, получен XML %s", dtString, dtFromXMLString)
	}
}

func TestDateTimeUnmarshalBadJSON(t *testing.T) {
	b := []byte(`"wrong"`)
	var dtFromJSON DateTime
	if err := json.Unmarshal(b, &dtFromJSON); err == nil {
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

	var dFromJSON Date
	if err := json.Unmarshal(b, &dFromJSON); err != nil {
		t.Fatal(err)
	}
	dFromJSONString := dFromJSON.String()
	if dFromJSONString != dString {
		t.Fatalf("Ошибка Unmarshal. Ожидалось получить JSON %s, получен JSON %s", dString, dFromJSONString)
	}
}

func TestDateXML(t *testing.T) {
	dString := "2015-07-30"
	xmlExpected := fmt.Sprintf("<Date>%s</Date>", dString)
	d, err := StringToDate(dString)
	if err != nil {
		t.Fatal(err)
	}
	b, err := xml.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	xmlReceived := string(b)
	if xmlExpected != xmlReceived {
		t.Fatalf("Ошибка Marshal. Ожидалось получить XML %s, получен XML %s", xmlExpected, xmlReceived)
	}

	var dFromXML Date
	if err := xml.Unmarshal(b, &dFromXML); err != nil {
		t.Fatal(err)
	}
	dFromXMLString := dFromXML.String()
	if dFromXMLString != dString {
		t.Fatalf("Ошибка Unmarshal. Ожидалось получить XML %s, получен XML %s", dString, dFromXMLString)
	}
}

func TestDateUnmarshalBadJSON(t *testing.T) {
	b := []byte(`"wrong"`)
	var dFromJSON Date
	if err := json.Unmarshal(b, &dFromJSON); err == nil {
		t.Fatal("Ожидалась ошибка")
	}
}

func TestDateTimeScanValueForDB(t *testing.T) {
	dt := DateTimeNow()
	v, err := dt.Value()
	if err != nil {
		t.Fatal(err)
	}
	vString, ok := v.(string)
	if !ok {
		t.Fatal("Ошибка преобразования интерфейса")
	}
	vdt, err := StringToDateTime(vString)
	if err != nil {
		t.Fatal(err)
	}
	var dtFromScan DateTime
	if err := dtFromScan.Scan(vdt.Time); err != nil {
		t.Fatal(err)
	}
	sExpected := dt.String()
	sReceived := dtFromScan.String()
	if sExpected != sReceived {
		t.Fatalf("Ожидалось получить строку %s, получена строка %s", sExpected, sReceived)
	}
}

func TestDateTimeScanIfBadValue(t *testing.T) {
	var dtFromScan DateTime

	if err := dtFromScan.Scan(nil); err != nil {
		t.Fatal(err)
	}

	if err := dtFromScan.Scan("wrong"); err == nil {
		t.Fatal("Ожидалась ошибка")
	}
}

func TestDateScanValueForDB(t *testing.T) {
	d := DateNow()
	v, err := d.Value()
	if err != nil {
		t.Fatal(err)
	}
	vString, ok := v.(string)
	if !ok {
		t.Fatal("Ошибка преобразования интерфейса")
	}
	vd, err := StringToDate(vString)
	if err != nil {
		t.Fatal(err)
	}
	var dFromScan Date
	if err := dFromScan.Scan(vd.Time); err != nil {
		t.Fatal(err)
	}
	sExpected := d.String()
	sReceived := dFromScan.String()
	if sExpected != sReceived {
		t.Fatalf("Ожидалось получить строку %s, получена строка %s", sExpected, sReceived)
	}
}

func TestDateScanIfBadValue(t *testing.T) {
	var dFromScan Date

	if err := dFromScan.Scan(nil); err != nil {
		t.Fatal(err)
	}

	if err := dFromScan.Scan("wrong"); err == nil {
		t.Fatal("Ожидалась ошибка")
	}
}

func TestDateTimeEqualityBug(t *testing.T) {
	var d1, d2 DateTime
	for {
		d1 = DateTimeNow()
		time.Sleep(10 * time.Millisecond)
		d2 = DateTimeNow()
		if d1.Second() == d2.Second() {
			break
		}
	}
	if !d1.Equal(d2) {
		t.Fatalf("Ошибка сравнения DateTime из-за различий в миллисекундах. d1 = %v, d2 = %v", d1, d2)
	}
}

func TestNullableDateTime(t *testing.T) {
	dt := DateTimeNow()
	expected := NullDateTime{
		DateTime: dt,
		Valid:    true,
	}
	received := dt.Nullable()
	if !reflect.DeepEqual(received, expected) {
		t.Fatal("Не равны. Ожидалось: %v, получено: %v", expected, received)
	}
}

func TestNullableDate(t *testing.T) {
	d := DateNow()
	expected := NullDate{Date: d, Valid: true}
	received := d.Nullable()
	if !reflect.DeepEqual(received, expected) {
		t.Fatal("Не равны. Ожидалось: %v, получено: %v", expected, received)
	}
}

func TestDateStringIfNull(t *testing.T) {
	d := MakeNullDate()
	received := fmt.Sprint(d)
	expected := "null"
	if received != expected {
		t.Fatal("Не равны. Ожидалось: %s, получено: %s", expected, received)
	}
}

func TestDateTimeStringIfNull(t *testing.T) {
	dt := MakeNullDateTime()
	received := fmt.Sprint(dt)
	expected := "null"
	if received != expected {
		t.Fatal("Не равны. Ожидалось: %s, получено: %s", expected, received)
	}
}

func TestMarshalUnmarshalJSONForNullDateIfValid(t *testing.T) {
	nd := DateNow().Nullable()
	jsonBytes, err := json.Marshal(nd)
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(`"%s"`, nd)
	received := string(jsonBytes)
	if expected != received {
		t.Fatalf("Неправильное отображение в JSON. Ожидалось: %v, получено: %v", expected, received)
	}

	var fromJSON NullDate
	if err := json.Unmarshal([]byte(received), &fromJSON); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fromJSON, nd) {
		t.Fatalf("Не равны. Получено из JSON: %v, ожидалось: %v", fromJSON, nd)
	}
}

func TestMarshalUnmarshalJSONForNullDateIfNotValid(t *testing.T) {
	nd := MakeNullDate()
	jsonBytes, err := json.Marshal(nd)
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(`null`)
	received := string(jsonBytes)
	if expected != received {
		t.Fatalf("Неправильное отображение в JSON. Ожидалось: %v, получено: %v", expected, received)
	}

	var fromJSON NullDate
	if err := json.Unmarshal([]byte(received), &fromJSON); err != nil {
		t.Fatal(err)
	}
	if fromJSON.Valid {
		t.Fatal("Valid == true, а должно быть false")
	}
}

func TestMarshalUnmarshalJSONForNullDateTimeIfValid(t *testing.T) {
	ndt := DateTimeNow().Nullable()
	jsonBytes, err := json.Marshal(ndt)
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(`"%s"`, ndt)
	received := string(jsonBytes)
	if expected != received {
		t.Fatalf("Неправильное отображение в JSON. Ожидалось: %v, получено: %v", expected, received)
	}

	var fromJSON NullDateTime
	if err := json.Unmarshal([]byte(received), &fromJSON); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fromJSON, ndt) {
		t.Fatalf("Не равны. Получено из JSON: %v, ожидалось: %v", fromJSON, ndt)
	}
}

func TestMarshalUnmarshalJSONForNullDateTimeIfNotValid(t *testing.T) {
	ndt := MakeNullDateTime()
	jsonBytes, err := json.Marshal(ndt)
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(`null`)
	received := string(jsonBytes)
	if expected != received {
		t.Fatalf("Неправильное отображение в JSON. Ожидалось: %v, получено: %v", expected, received)
	}

	var fromJSON NullDateTime
	if err := json.Unmarshal([]byte(received), &fromJSON); err != nil {
		t.Fatal(err)
	}
	if fromJSON.Valid {
		t.Fatal("Valid == true, а должно быть false")
	}
}

func TestMarshalUnmarshalXMLForNullDateTimeIfValid(t *testing.T) {
	ndt := DateTimeNow().Nullable()
	xmlBytes, err := xml.Marshal(ndt)
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(`<NullDateTime>%s</NullDateTime>`, ndt)
	received := string(xmlBytes)
	if expected != received {
		t.Fatalf("Неправильное отображение в XML. Ожидалось: %v, получено: %v", expected, received)
	}

	var fromXML NullDateTime
	if err := xml.Unmarshal([]byte(received), &fromXML); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fromXML, ndt) {
		t.Fatalf("Не равны. Получено из XML: %v, ожидалось: %v", fromXML, ndt)
	}
}

func TestMarshalUnmarshalXMLForNullDateTimeIfNotValid(t *testing.T) {
	ndt := MakeNullDateTime()
	xmlBytes, err := xml.Marshal(ndt)
	if err != nil {
		t.Fatal(err)
	}
	
	//received := string(xmlBytes)
	if xmlBytes != nil {
		t.Fatalf("Неправильное отображение в XML. Ожидалось: nil получено: %v", xmlBytes)
	}

	var fromXML NullDateTime
	expected := fmt.Sprintf("<NullDateTime></NullDateTime>")
	if err := xml.Unmarshal([]byte(expected), &fromXML); err != nil {
		t.Fatal(err)
	}
	if fromXML.Valid {
		t.Fatal("Valid == true, а должно быть false")
	}
}

func TestMarshalUnmarshalXMLForNullDateIfValid(t *testing.T) {
	nd := DateNow().Nullable()
	xmlBytes, err := xml.Marshal(nd)
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf(`<NullDate>%s</NullDate>`, nd)
	received := string(xmlBytes)
	if expected != received {
		t.Fatalf("Неправильное отображение в XML. Ожидалось: %v, получено: %v", expected, received)
	}

	var fromXML NullDate
	if err := xml.Unmarshal([]byte(received), &fromXML); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fromXML, nd) {
		t.Fatalf("Не равны. Получено из XML: %v, ожидалось: %v", fromXML, nd)
	}
}

func TestMarshalUnmarshalXMLForNullDateIfNotValid(t *testing.T) {
	nd := MakeNullDate()
	xmlBytes, err := xml.Marshal(nd)
	if err != nil {
		t.Fatal(err)
	}

	//received := string(xmlBytes)
	if xmlBytes != nil {
		t.Fatalf("Неправильное отображение в XML. Ожидалось: nil, получено: %v", xmlBytes)
	}

	var fromXML NullDate
	expected := fmt.Sprintf("<NullDate></NullDate>")
	if err := xml.Unmarshal([]byte(expected), &fromXML); err != nil {
		t.Fatal(err)
	}
	if fromXML.Valid {
		t.Fatal("Valid == true, а должно быть false")
	}
}
