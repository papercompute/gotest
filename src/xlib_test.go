package xlib_test

import (
	"../src"
	"encoding/json"
	"fmt"
	"testing"
)

func TestLib1(t *testing.T) {

	testQueries := []string{
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2015-12-31"},
  "reservations": []
}`,
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2015-12-31"},
  "reservations": [
    {"start": "2015-10-01", "end": "2015-12-05"}
  ]
}`,
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2015-12-31"},
  "reservations": [
    {"start": "2015-12-05", "end": "2016-01-05"}
  ]
}`,
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2015-12-31"},
  "reservations": [
    {"start": "2015-12-02", "end": "2015-12-29"}
  ]
}`,
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2016-01-31"},
  "reservations": [
    {"start": "2015-12-02", "end": "2016-01-31"}
  ]
}`,
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2015-12-31"},
  "reservations": [
    {"start": "2015-12-01", "end": "2015-12-01"},
    {"start": "2015-12-31", "end": "2015-12-31"}
  ]
}`,
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2015-12-31"},
  "reservations": [
    {"start": "2015-11-01", "end": "2015-12-01"},
    {"start": "2015-12-17", "end": "2015-12-17"},
    {"start": "2015-12-18", "end": "2015-12-30"}
  ]
}`,
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2016-12-31"},
  "reservations": [
    {"start": "2015-11-01", "end": "2015-11-01"},
    {"start": "2015-12-02", "end": "2015-12-03"},
    {"start": "2015-12-07", "end": "2015-12-13"},
    {"start": "2015-12-18", "end": "2016-11-30"}
  ]
}`}

	testResults := []string{
		`{"start": "2015-12-01", "end": "2015-12-31"}`,
		`{"start": "2015-12-06", "end": "2015-12-31"}`,
		`{"start": "2015-12-01", "end": "2015-12-04"}`,
		`{"start": "2015-12-30", "end": "2015-12-31"}`,
		`{"start": "2015-12-01", "end": "2015-12-01"}`,
		`{"start": "2015-12-02", "end": "2015-12-30"}`,
		`{"start": "2015-12-02", "end": "2015-12-16"}`,
		`{"start": "2016-12-01", "end": "2016-12-31"}`}

	i := 0
	for _, q := range testQueries {
		var d xlib.Data
		fmt.Println(q)
		if err := json.Unmarshal([]byte(q), &d); err != nil {
			t.Fatal(fmt.Sprintf("%v", err))
		}

		p, err := xlib.LongestFreeRange(&d)
		if err != nil {
			t.Fatal(fmt.Sprintf("%v", err))
		}

		if p != testResults[i] {
			t.Fatal(fmt.Sprintf("error %s must be %s", p, testResults[i]))
		}

		i++

	}
}

func TestLib2(t *testing.T) {

	testQueries := []string{
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2015-12-31"},
  "reservations": [
    {"start": "2015-12-17", "end": "2015-12-17"},
    {"start": "2015-12-17", "end": "2015-12-30"}
  ]
}`,
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2015-12-31"},
  "reservations": [
    {"start": "2015-12-17", "end": "2015-12-17"},
    {"start": "2015-12-17", "end": "2015-12-17"}
  ]
}`,
		`{
  "renting_period": {"start": "2015-12-01",  "end": "2015-12-31"},
  "reservations": [
    {"start": "2015-12-31", "end": "2015-12-17"}
  ]
}`}

	for _, q := range testQueries {
		var d xlib.Data
		fmt.Println(q)
		if err := json.Unmarshal([]byte(q), &d); err != nil {
			t.Fatal(fmt.Sprintf("%v", err))
		}

		_, err := xlib.LongestFreeRange(&d)
		if err == nil {
			t.Fatal(fmt.Sprintf("error %v - %v", err, d))
		}
	}

}
