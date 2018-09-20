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

func handler(request string) StockPrice {
	urlString := BuildURL(request)
	data := Requesting(urlString)
	results := ParseRecords(data)
	result := results[0]
	resultStock := StockPrice{
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
	return resultStock
}

func handlerConcurrent(symbols []string) []StockPrice {
	var wg sync.WaitGroup
	wg.Add(len(symbols))
	var results = make([]StockPrice, len(symbols))
	for i, symbol := range symbols {
		go func(symbol string, i int) {
			results[i] = handler(symbol)
			wg.Done()
		}(symbol, i)
	}
	wg.Wait()
	return results
}

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
	db, err := gorm.Open("mysql", "Your Connection")
	if err != nil {
		log.Fatalln("Failed to connect database")
	}
	defer db.Close()
	db.SingularTable(true)
	var wg sync.WaitGroup
	var count int
	var stockAmount int

	db.Table("stock").Where("updated = false").Count(&count)
	if count != 0 {
		stockAmount = 5
		if count < 5 {
			stockAmount = count
		}
	} else {
		log.Println("Finished Updating Stocks")
		return true, nil
	}
	var stock []Stock
	var symbols = make([]string, 0, stockAmount)
	db.Where("updated = false").Limit(5).Find(&stock)

	// append symbols with data from stock
	for _, element := range stock {
		symbols = append(symbols, element.Symbol)
	}

	// Create stock_price rows in the database
	var requests = make([]StockPrice, 0, len(symbols))
	requests = handlerConcurrent(symbols)

	wg.Add(len(symbols))
	// Check if the dates are current
	for i, request := range requests {
		go func(i int, request StockPrice) {
			currentDate := time.Now().AddDate(0, 0, -1).Truncate(time.Hour * 24)
			if request.Date.Equal(currentDate) == false {
				log.Printf("%v is no longer current, the date %v was pulled instead", request.Symbol, request.Date)
				// remove index i from requests
				requests[i] = requests[len(requests)-1]
				requests = requests[:len(requests)-1]
			}
			wg.Done()
		}(i, request)
	}
	wg.Wait()

	// Insert stock_prices into database
	insert := insertQuery(requests)
	db.Exec(insert)

	// Change the updated field in stock from false to true
	db.Exec(`UPDATE stock
						SET updated = true
						WHERE symbol in (?)`, symbols)
	return true, nil
}
