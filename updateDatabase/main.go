package main

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	lambda.Start(UpdateDatabase)
}

// handler sends a request to Alpha Vantage
func handler(request string) (StockPrice, error) {
	var resultStock StockPrice
	urlString := BuildURL(request)
	data, err := Requesting(urlString)
	if err != nil {
		return resultStock, err
	}
	results := ParseRecords(data)
	result := results[0]
	resultStock = StockPrice{
		Date:          result.Date,
		Symbol:        request,
		Open:          result.Open,
		High:          result.High,
		Low:           result.Low,
		Close:         result.Close,
		AdjustedClose: result.AdjustedClose,
		Volume:        result.Volume,
		Dividend:      result.Dividend,
		Split:         result.Split,
	}
	return resultStock, nil
}

func handlerConcurrent(symbols []string) ([]StockPrice, error) {
	errors := make(chan error, 4)
	var wg sync.WaitGroup
	wg.Add(len(symbols))
	var results = make([]StockPrice, len(symbols))
	for i, symbol := range symbols {
		go func(symbol string, i int) {
			var err error
			results[i], err = handler(symbol)
			if err != nil {
				errors <- err
			}
			wg.Done()
		}(symbol, i)
	}
	wg.Wait()
	if len(errors) > 0 {
		return results, <-errors
	}
	return results, nil
}

// updateQuery changes updated from false to true in database for the symbols in the array
func updateQuery(symbols []string) string {
	valuesQuery := "UPDATE stock SET updated = true WHERE symbol IN("
	for _, symbol := range symbols {

		valuesQuery += "'" + symbol + "'" + ","
	}
	// Trim , at the end
	valuesQuery = valuesQuery[0 : len(valuesQuery)-1]
	valuesQuery += ")"
	return valuesQuery
}

// insertQuery creates a string that combines values into an insert query
func insertQuery(requests []StockPrice) string {
	valuesQuery := "INSERT INTO stock_price VALUES "
	for _, row := range requests {

		date := row.Date.String()
		symbol := row.Symbol
		open := strconv.FormatFloat(row.Open, 'f', -1, 64)
		high := strconv.FormatFloat(row.High, 'f', -1, 64)
		low := strconv.FormatFloat(row.Low, 'f', -1, 64)
		close := strconv.FormatFloat(row.Close, 'f', -1, 64)
		adjustedClose := strconv.FormatFloat(row.AdjustedClose, 'f', -1, 64)
		volume := strconv.Itoa(int(row.Volume))
		dividend := strconv.FormatFloat(row.Dividend, 'f', -1, 64)
		split := strconv.FormatFloat(row.Split, 'f', -1, 64)
		data := "\"" + date + "\"" + "," + "'" + symbol + "'" + "," + open + "," + high + "," + low + "," + close + ","
		data += adjustedClose + "," + volume + "," + dividend + "," + split

		valuesQuery += "(" + data + "),"
	}
	// Trim , at the end
	valuesQuery = valuesQuery[0 : len(valuesQuery)-1]
	return valuesQuery
}

// UpdateDatabase updates database
func UpdateDatabase() (bool, error) {
	db, err := gorm.Open("mysql", "SQL Connection")
	if err != nil {
		log.Fatalln("Failed to connect database")
	}
	defer db.Close()
	db.SingularTable(true)
	//db.LogMode(true)
	var count int
	var stockAmount int

	db.Table("stock").Where("updated = false").Count(&count)
	if count != 0 {
		stockAmount = 4
		if count < 4 {
			stockAmount = count
		}
	} else {
		log.Println("Finished Updating Stocks")
	}
	var stock []Stock
	var symbols = make([]string, 0, stockAmount)
	db.Where("updated = false").Limit(4).Find(&stock)

	// append symbols with data from stock
	for _, element := range stock {
		symbols = append(symbols, element.Symbol)
	}

	// Create stock_price rows in the database
	var requests = make([]StockPrice, 0, len(symbols))
	requests, err = handlerConcurrent(symbols)
	if err != nil {
		log.Println(err)
		return false, nil
	}
	// Check if the dates are current
	requests = checkDates(requests)
	//requests = checkDuplicate(requests)
	// Insert stock_prices into database
	insert := insertQuery(requests)

	// check if requests is empty
	if len(requests) > 0 {
		db.Exec(insert)
	}

	// Change the updated field in stock from false to true
	db.Exec(`UPDATE stock
						SET updated = true
						WHERE symbol in (?)`, symbols)
	return true, nil
}

// checkDates checks if all the dates in the request are current
func checkDates(requests []StockPrice) []StockPrice {
	for i := len(requests) - 1; i >= 0; i-- {
		currentDate := time.Now().AddDate(0, 0, -1).Truncate(time.Hour * 24)
		if requests[i].Date.Equal(currentDate) == false {
			log.Printf("%v is no longer current, the date %v was pulled instead", requests[i].Symbol, requests[i].Date)
			// remove index i from requests
			requests[len(requests)-1], requests[i] = requests[i], requests[len(requests)-1]
			requests = requests[:len(requests)-1]
		}
	}
	return requests
}

// checkDuplicate checks if a request contains duplicates already in the database
func checkDuplicate(requests []StockPrice) []StockPrice {
	db, err := gorm.Open("mysql", "SQL Connection")
	if err != nil {
		log.Fatalln("Failed to connect database")
	}
	defer db.Close()
	db.SingularTable(true)

	for i := len(requests) - 1; i >= 0; i-- {
		var count int
		db.Table("stock_price").Where("symbol = ? AND date = ?", requests[i].Symbol, requests[i].Date).Count(&count)
		if count != 0 {
			log.Printf("%v, %v is a duplicate", requests[i].Date, requests[i].Symbol)
			// remove duplicate
			requests[len(requests)-1], requests[i] = requests[i], requests[len(requests)-1]
			requests = requests[:len(requests)-1]
		}
	}
	return requests
}
