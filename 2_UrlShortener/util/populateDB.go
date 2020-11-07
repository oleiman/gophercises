package main

import (
	"github.com/boltdb/bolt"
	"log"
)

func main() {
	db, err := bolt.Open("redirect.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("Mappings"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("/gatech"), []byte("https://omscs.gatech.edu"))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
