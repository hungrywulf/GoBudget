package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
)

func randomNameGenerator(firstFile string, lastFile string, count int) []string {

	first, err := ioutil.ReadFile(firstFile)
	if err != nil {
		fmt.Println(err)
	}
	last, err := ioutil.ReadFile(lastFile)
	if err != nil {
		fmt.Println(err)
	}
	firstS := strings.TrimSpace(string(first))
	lastS := strings.TrimSpace(string(last))
	firstNames := strings.Split(firstS, ",")
	lastNames := strings.Split(lastS, ",")
	names := make([]string, count)

	for i := range names {
		firstRandom := rand.Intn(5494)
		lastRandom := rand.Intn(88799)
		names[i] = firstNames[firstRandom] + " " + lastNames[lastRandom]
	}
	return names
}

func randomEmailGenerator(firstFile string, lastFile string, count int) ([]string, []string) {

	names := randomNameGenerator("first_names.csv", "last_names.csv", count)
	emails := make([]string, count)

	domains := []string{"aol.com", "att.net", "comcast.net", "facebook.com", "gmail.com", "gmx.com",
		"googlemail.com", "google.com", "hotmail.com", "hotmail.co.uk", "mac.com", "me.com", "mail.com",
		"msn.com", "live.com", "sbcglobal.net", "verizon.net", "yahoo.com", "yahoo.co.uk"}

	for i, name := range names {
		name = strings.Replace(name, " ", "", -1)
		domain := domains[rand.Intn(19)]
		firstNumber := strconv.Itoa(rand.Intn(10))
		secondNumber := strconv.Itoa(rand.Intn(10))
		emails[i] = name + firstNumber + secondNumber + "@" + domain
	}
	return names, emails
}

// RandomUserGenerator returns a random user
func RandomUserGenerator(firstFile string, lastFile string, count int) [][]string {

	names, emails := randomEmailGenerator("first_names.csv", "last_names.csv", count)
	users := make([][]string, count)
	for i := range users {
		users[i] = make([]string, 3)
	}

	passwords := make([]string, len(names), len(names))
	specialCharacter := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "`", "[", "]", ",", "{", "}", "?", "<", ">"}

	for i, name := range names {
		name = strings.Replace(name, " ", "", -1)

		firstNumber := strconv.Itoa(rand.Intn(10))
		secondNumber := strconv.Itoa(rand.Intn(10))
		thirdNumber := strconv.Itoa(rand.Intn(10))
		randomSpecial := specialCharacter[rand.Intn(20)]
		passwords[i] = name + firstNumber + secondNumber + thirdNumber + randomSpecial
		passwords[i] = strings.Replace(passwords[i], " ", "", -1)
	}
	for i := range users {
		users[i][0] = names[i]
		users[i][1] = emails[i]
		users[i][2] = passwords[i]
	}
	return users
}
