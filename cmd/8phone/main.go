package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"regexp"
)

var schema = `
DROP TABLE IF EXISTS phone_numbers;
CREATE TABLE phone_numbers (
id     INTEGER PRIMARY KEY,
number VARCHAR(256) DEFAULT ''
);
`

var raw_pns = []string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

type PhoneNumber struct {
	Id     int    `db:"id"`
	Number string `db:"number"`
}

func normalize(phoneNumber string) string {
	reg, _ := regexp.Compile("[^0-9]+")
	return reg.ReplaceAllString(phoneNumber, "")
}

func addNumbers(db *sqlx.DB) error {
	db.MustExec(schema)

	tx := db.MustBegin()
	for _, n := range raw_pns {
		tx.MustExec("INSERT INTO phone_numbers (number) VALUES ($1)", n)
	}
	tx.Commit()
	return nil
}

func getNumbers(db *sqlx.DB) ([]PhoneNumber, error) {
	numbers := []PhoneNumber{}
	if err := db.Select(&numbers, "SELECT * FROM phone_numbers"); err != nil {
		return nil, err
	}
	return numbers, nil
}

func deleteIfDupe(db *sqlx.DB, n *PhoneNumber) (bool, error) {
	rows, err := db.NamedQuery(`SELECT * FROM phone_numbers WHERE number=:number AND id!=:id`,
		map[string]interface{}{
			"number": n.Number,
			"id":     n.Id,
		})
	if err != nil {
		return false, err
	}
	delete := false
	for rows.Next() {
		delete = true
		rows.Close()
	}

	if delete {
		_, err := db.NamedExec(`DELETE FROM phone_numbers WHERE id=:id`,
			map[string]interface{}{
				"id": n.Id,
			})
		if err != nil {
			return false, err
		}
	}
	return delete, nil
}

func updateNumber(db *sqlx.DB, n *PhoneNumber) error {
	_, err := db.NamedExec(`UPDATE phone_numbers SET number=:number WHERE id=:id`,
		map[string]interface{}{
			"id":     n.Id,
			"number": n.Number,
		})
	return err
}

func main() {
	db, err := sqlx.Connect("sqlite3", "__numbers.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := addNumbers(db); err != nil {
		log.Fatal(err)
	}

	numbers, err := getNumbers(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range numbers {
		n.Number = normalize(n.Number)
		deleted, err := deleteIfDupe(db, &n)
		if err != nil {
			log.Fatal(err)
		} else if deleted {
			continue
		}
		if err := updateNumber(db, &n); err != nil {
			log.Fatal(err)
		}
	}

	numbers, err = getNumbers(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range numbers {
		fmt.Println(n)
	}
}
