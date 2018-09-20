package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	lambda.Start(ResetDatabase)
}

// ResetDatabase changes updated to false in database
func ResetDatabase() (bool, error) {
	db, err := gorm.Open("mysql", "Your Connection")
	if err != nil {
		log.Fatalln("Failed to connect database")
	}
	defer db.Close()
	db.SingularTable(true)

	db.Exec("update stock set updated = false where updated = true")
	log.Println("Database Finished Resetting")
	return true, nil
}
