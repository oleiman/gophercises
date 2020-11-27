package database

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	// "fmt"
	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
)

var HOMEDIR, _ = homedir.Dir()
var DB_FNAME = filepath.Join(HOMEDIR, ".gophertask.db")

// Get tasks, but only those that satisfy the provided predicate
func GetTasks(pred func(Task) bool) ([]Task, error) {
	db, err := bolt.Open(DB_FNAME, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var tasks []Task
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		if b == nil {
			return errors.New("Database not initialized")
		}
		cursor := b.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var item Task
			err := json.Unmarshal(v, &item)
			if err != nil || !pred(item) {
				continue
			}
			tasks = append(tasks, item)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func UpdateTask(task *Task) error {
	err := DeleteTask(task.ID)
	if err != nil {
		return err
	}
	db, err := bolt.Open(DB_FNAME, 0600, nil)
	if err != nil {
		return nil
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		if b == nil {
			return errors.New("Database not initialized")
		}
		binkey := make([]byte, 4)
		binary.BigEndian.PutUint32(binkey, uint32(task.ID))
		buf, err := json.Marshal(*task)
		if err != nil {
			return err
		}
		if err := b.Put(binkey, buf); err != nil {
			return err
		}
		return nil
	})
}

func DeleteTask(key int) error {
	db, err := bolt.Open(DB_FNAME, 0600, nil)
	if err != nil {
		return nil
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		if b == nil {
			return errors.New("Database not initialized")
		}
		binkey := make([]byte, 4)
		binary.BigEndian.PutUint32(binkey, uint32(key))
		if err := b.Delete(binkey); err != nil {
			return err
		}
		return nil
	})
}

func AddTasks(descriptions []string) error {
	db, err := bolt.Open(DB_FNAME, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		for _, description := range descriptions {
			b := tx.Bucket([]byte("tasks"))
			if b == nil {
				return errors.New("Database not initialized")
			}
			id, _ := b.NextSequence()
			item := Task{
				ID:      int(id),
				Content: description,
				Done:    false,
			}
			buf, err := json.Marshal(item)
			if err != nil {
				return err
			}
			key := make([]byte, 4)
			binary.BigEndian.PutUint32(key, uint32(item.ID))
			b.Put(key, buf)
		}
		return nil
	})
}

func InitBolt() error {
	db, err := bolt.Open(DB_FNAME, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tasks"))
		return err
	})

}
