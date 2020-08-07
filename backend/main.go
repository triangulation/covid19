package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"foxygo.at/s/errs"
	"github.com/tidwall/pretty"
)

// ErrInvalidInput is a sentinel error for invalid input in CSV file.
var ErrInvalidInput = fmt.Errorf("invalid input")

// Data is the top level structure holding all covid time series data
// per country. For memory efficiency dates are only listed once.
type Data struct {
	Dates     []time.Time         `json:"dates"`
	Countries map[string]*Country `json:"countries"`
}

// Country is a sub-struct of Data holding all time series data for a
// specific country some meta-data such as population.
type Country struct {
	Name       string  `json:"name"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Population int64   `json:"population"`
	// per date new, not cumulative, confirmed cases/recovered/dead:
	Confirmed []int `json:"confirmed"`
	Recovered []int `json:"recovered"`
	Dead      []int `json:"dead"`
}

func main() {
	data, err := parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	opts := &pretty.Options{Width: 10000, Prefix: "", Indent: "  "}
	fmt.Println(string(pretty.PrettyOptions(b, opts)))
}

func parse(r io.Reader) (*Data, error) {
	csvr := csv.NewReader(r)
	csvr.FieldsPerRecord = 8
	csvr.ReuseRecord = true
	// Date,       Country/Region,Province/State,Lat,  Long, Confirmed,Recovered,Deaths
	// 2020-01-22, Afghanistan,                 ,33.0, 65.0, 0,        0,        0
	rows, err := csvr.ReadAll()
	if err != nil {
		return nil, err
	}
	rows = rows[1:] // skip header
	dates, err := parseDates(rows)
	if err != nil {
		return nil, err
	}
	countries, err := parseCountries(rows, dates)
	if err != nil {
		return nil, err
	}
	return &Data{Dates: dates, Countries: countries}, nil
}

func parseDates(rows [][]string) ([]time.Time, error) {
	d := map[string]bool{}
	dates := []time.Time{}
	for _, row := range rows {
		if d[row[0]] {
			break
		}
		d[row[0]] = true
		date, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}
	return dates, nil
}

func parseCountries(rows [][]string, dates []time.Time) (map[string]*Country, error) {
	countries := map[string]*Country{}
	dateIdx := 0
	dateCnt := len(dates)
	for _, row := range rows {
		country := strings.TrimSpace(row[1])
		if countries[country] == nil {
			c, err := newCountry(row, dateCnt)
			if err != nil {
				return nil, err
			}
			countries[country] = c
		}
		if err := updateCountry(countries[country], row, dates[dateIdx], dateIdx); err != nil {
			return nil, err
		}
		dateIdx = (dateIdx + 1) % dateCnt
	}
	// Store deltas (i.e. daily new cases/recoveries/deaths) rather than
	// cumulative numbers:
	for _, country := range countries {
		for i := dateCnt - 1; i > 0; i-- {
			country.Confirmed[i] -= country.Confirmed[i-1]
			country.Recovered[i] -= country.Recovered[i-1]
			country.Dead[i] -= country.Dead[i-1]
		}
	}
	return countries, nil
}

func updateCountry(country *Country, row []string, wantDate time.Time, dateIdx int) error {
	confirmed, err := atoi0(row[5])
	if err != nil {
		return errs.Errorf("%v: %s: %v", ErrInvalidInput, row, err)
	}
	recovered, err := atoi0(row[6])
	if err != nil {
		return errs.Errorf("%v: %s: %v", ErrInvalidInput, row, err)
	}
	dead, err := atoi0(row[7])
	if err != nil {
		return errs.Errorf("%v: %s: %v", ErrInvalidInput, row, err)
	}
	country.Confirmed[dateIdx] += confirmed
	country.Recovered[dateIdx] += recovered
	country.Dead[dateIdx] += dead
	date, err := time.Parse("2006-01-02", row[0])
	if err != nil {
		return errs.Errorf("%v: %s: %v", ErrInvalidInput, row, err)
	}
	if wantDate != date {
		return errs.Errorf("%v: %s: want date %s got %s", ErrInvalidInput, row, wantDate, date)
	}
	return nil
}

func newCountry(row []string, dateCnt int) (*Country, error) {
	lat, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
	if err != nil {
		return nil, errs.Errorf("%v: %s: %v", ErrInvalidInput, row, err)
	}
	lng, err := strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
	if err != nil {
		return nil, errs.Errorf("%v: %s: %v", ErrInvalidInput, row, err)
	}
	country := &Country{
		Name:      strings.TrimSpace(row[1]),
		Lat:       lat,
		Lng:       lng,
		Confirmed: make([]int, dateCnt),
		Recovered: make([]int, dateCnt),
		Dead:      make([]int, dateCnt),
	}
	return country, nil
}

func atoi0(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}
