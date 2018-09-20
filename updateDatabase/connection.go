package main

import (
	"encoding/csv"
	"errors"
	"log"
	"net/http"
	"net/url"
)

const (
	scheme        = "https"
	hostName      = "www.alphavantage.co"
	queryPath     = "query"
	queryFunction = "function"
	querySymbol   = "symbol"
	queryAPIKey   = "apikey"
	queryDataType = "datatype"
)

// BuildURL creates the URL with the host name path, and various queries
func BuildURL(symbol string) *url.URL {

	urlPath := &url.URL{}
	urlPath.Scheme = scheme
	urlPath.Host = hostName
	urlPath.Path = queryPath

	query := urlPath.Query()
	query.Set(queryFunction, "TIME_SERIES_DAILY_ADJUSTED")
	query.Set(querySymbol, symbol)
	query.Set(queryAPIKey, "Your api Key")
	query.Set(queryDataType, "csv")
	urlPath.RawQuery = query.Encode()
	return urlPath
}

// Request data from the Host
func Requesting(URL *url.URL) [][]string {
	var errRequest error
	response, err := http.Get(URL.String())
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	r := csv.NewReader(response.Body)
	records, err := r.ReadAll()
	if err != nil {
		errRequest = errors.New("Too many requests were sent to Alpha Vantage or Alpha Vantage is Down")
		log.Fatalf("%v, %v\n", err, errRequest)
	}

	return records
}
