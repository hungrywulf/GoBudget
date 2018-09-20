package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	start := time.Now()

	//users := RandomUserGenerator("first_names.csv", "last_names.csv", 10000)

	db, err := gorm.Open("mysql", "Your Connection")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.SingularTable(true)
	db.LogMode(true)
	//db.Create(&Users{First_name: "Kate", Last_name: "Green", Email: "green@gmail.com"})
	err = gorm.LogErr()
	if err != nil {
		log.Println("Handle this error")
	}

	//increases := []int{0, 1}

	//reoccuring := []string{"false", "daily", "weekly", "monthly", "yearly"}

	// generate items
	//db.Create(&Items{User_ID: uint(userID), Date: date, Amount: itemAmount, Increases: increasesM[increasesR], Reoccuring: reoccuringM[reoccuringR]})

	users := RandomUserGenerator("first_names.csv", "last_names.csv", 100000)

	for _, element := range users {
		fmt.Println(element)
	}
	fmt.Println(len(users))
	fmt.Printf("%v", time.Since(start))
}

func generateSymbols() {
	stocksF, err := ioutil.ReadFile("stocks.csv")
	if err != nil {
		fmt.Println(err)
	}
	stocks := strings.Split(string(stocksF), ",")
	for _, i := range stocks {
		stock := Stock{Symbol: i, Updated: false}
		db.Create(&stock)
	}
}
