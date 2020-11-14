package pndb

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

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

type DB struct {
	db *sqlx.DB
}

type PhoneNumber struct {
	Id     int    `db:"id"`
	Number string `db:"number"`
}

func Open(dbType string, fname string) (DB, error) {
	db, err := sqlx.Connect("sqlite3", "__numbers.db")
	return DB{db}, err
}

func (phonedb DB) Close() {
	phonedb.db.Close()
}

func (pndb DB) SeedPhoneNumberTable() {
	schema := `
DROP TABLE IF EXISTS phone_numbers;
CREATE TABLE phone_numbers (
id     INTEGER PRIMARY KEY,
number VARCHAR(256) DEFAULT ''
);
`
	pndb.db.MustExec(schema)
	tx := pndb.db.MustBegin()
	for _, n := range raw_pns {
		tx.MustExec("INSERT INTO phone_numbers (number) VALUES ($1)", n)
	}
	tx.Commit()
}

func (pndb DB) GetPhoneNumbers() ([]PhoneNumber, error) {
	numbers := []PhoneNumber{}
	if err := pndb.db.Select(&numbers, "SELECT * FROM phone_numbers"); err != nil {
		return nil, err
	}
	return numbers, nil
}

func (pndb DB) FindNumber(number string) bool {
	rows, err := pndb.db.NamedQuery(`
SELECT * FROM phone_numbers
WHERE number=:number`,
		map[string]interface{}{
			"number": number,
		})
	if err != nil {
		return false
	}
	for rows.Next() {
		defer rows.Close()
		return true
	}
	return false
}

func (pndb DB) DeleteNumber(id int) error {
	_, err := pndb.db.NamedExec(`
DELETE FROM phone_numbers WHERE id=:id`,
		map[string]interface{}{
			"id": id,
		})
	return err
}

func (pndb DB) UpdateNumber(n *PhoneNumber) error {
	_, err := pndb.db.NamedExec(`UPDATE phone_numbers SET number=:number WHERE id=:id`,
		map[string]interface{}{
			"id":     n.Id,
			"number": n.Number,
		})
	return err
}
