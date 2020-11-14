package main

import (
	"fmt"
	"gophercises/pkg/pndb"
	"log"
	"regexp"
)

func normalize(phoneNumber string) string {
	reg, _ := regexp.Compile("[^0-9]+")
	return reg.ReplaceAllString(phoneNumber, "")
}

func main() {

	db, err := pndb.Open("sqlite3", "__numbers.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SeedPhoneNumberTable()
	numbers, err := db.GetPhoneNumbers()
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range numbers {
		normalized := normalize(n.Number)
		if normalized == n.Number {
			continue
		}

		if db.FindNumber(normalized) {
			if err := db.DeleteNumber(n.Id); err != nil {
				log.Fatal(err)
			}
		} else {
			n.Number = normalized
			if err := db.UpdateNumber(&n); err != nil {
				log.Fatal(err)
			}
		}
	}

	numbers, err = db.GetPhoneNumbers()
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range numbers {
		fmt.Println(n)
	}
}
