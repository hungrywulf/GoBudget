package main

import (
	"log"
	"strconv"
	"time"
)

// Represents the items data as integers
const (
	Date = iota
	Open
	High
	Low
	Close
	AdjustedClose
	Volume
	Dividend
	Split
)

// parseFloat parses input string as a float
func parseFloat(str string) float64 {

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Println(err)
	}
	return f
}

// parseInt parses input string into an int
func parseInt(str string) uint {

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Println(err)
	}
	return uint(i)
}

// parseDate parses input string into time.Time
func parseDate(str string) time.Time {

	format := "2006-01-02"
	time, err := time.Parse(format, str)

	if err != nil {
		log.Println(err)
	}
	return time
}

// parseRecords parses the input into the TimeSeriesAD struct
func ParseRecords(records [][]string) []*TimeSeriesAD {
	records = records[1:][:]
	values := &TimeSeriesAD{}
	collection := []*TimeSeriesAD{}
	for i := 0; i < len(records); i++ {
		values = &TimeSeriesAD{
			Date:          parseDate(records[i][Date]),
			Open:          parseFloat(records[i][Open]),
			High:          parseFloat(records[i][High]),
			Low:           parseFloat(records[i][Low]),
			Close:         parseFloat(records[i][Close]),
			AdjustedClose: parseFloat(records[i][AdjustedClose]),
			Volume:        parseInt(records[i][Volume]),
			Dividend:      parseFloat(records[i][Dividend]),
			Split:         parseFloat(records[i][Split]),
		}
		collection = append(collection, values)
	}

	return collection
}
