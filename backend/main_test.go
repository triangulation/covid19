package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mustTime(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestParse(t *testing.T) {
	s := `
Date,      Country, Province, Lat, Long, Confirmed, Recovered, Deaths
2020-01-01,    Abc,         ,   2,    1,         3,         0,      1
2020-01-02,    Abc,         ,   2,    1,        10,         2,      6
2020-01-03,    Abc,         ,   2,    1,        12,         4,      5
`
	b := bytes.NewBufferString(strings.TrimSpace(s))
	got, err := parse(b)
	require.NoError(t, err)

	want := &Data{
		Dates: []time.Time{mustTime("2020-01-01"), mustTime("2020-01-02"), mustTime("2020-01-03")},
		Countries: map[string]*Country{
			"Abc": {
				Name:      "Abc",
				Lat:       2,
				Lng:       1,
				Confirmed: []int{3, 7, 2},
				Recovered: []int{0, 2, 2},
				Dead:      []int{1, 5, -1},
			},
		},
	}
	require.Equal(t, want, got)
}

func TestParseProvince(t *testing.T) {
	s := `
Date,      Country, Province, Lat, Long, Confirmed, Recovered, Deaths
2020-01-01,    Abc,    prov1,   2,    1,         3,         0,      1
2020-01-02,    Abc,    prov1,   2,    1,         5,         2,      2
2020-01-01,    Abc,    prov2,   3,    1,         2,         0,      1
2020-01-02,    Abc,    prov2,   3,    1,         6,         2,      3
`
	b := bytes.NewBufferString(strings.TrimSpace(s))
	got, err := parse(b)
	require.NoError(t, err)

	want := &Data{
		Dates: []time.Time{mustTime("2020-01-01"), mustTime("2020-01-02")},
		Countries: map[string]*Country{
			"Abc": {
				Name:      "Abc",
				Lat:       2,
				Lng:       1,
				Confirmed: []int{5, 6},
				Recovered: []int{0, 4},
				Dead:      []int{2, 3},
			},
		},
	}
	require.Equal(t, want, got)
}

func TestParseFile(t *testing.T) {
	f, err := os.Open("testdata/data.csv")
	require.NoError(t, err)
	data, err := parse(f)
	require.NoError(t, err)

	got, err := json.Marshal(data)
	require.NoError(t, err)
	want, err := ioutil.ReadFile("testdata/want.json")
	require.NoError(t, err)

	require.JSONEq(t, string(want), string(got))
}
