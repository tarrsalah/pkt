package internal

import (
	"encoding/json"
	"log"
	"sort"

	"go.etcd.io/bbolt"
	"path/filepath"
)

var bucketName = []byte("pkt")

// DB is a wrapper around bolt database connection
type DB struct {
	db *bbolt.DB
}

// NewDB returns a new database connnection
func NewDB() *DB {
	path := filepath.Join(configDir(), "pkt.bolt")
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{db}
}

// Close closed the database connection
func (b *DB) Close() {
	b.db.Close()
}

// Get gets all stored pocket items
func (b *DB) Get() (Items, error) {
	var items []Item

	err := b.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return nil
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var item Item
			err := json.Unmarshal(v, &item)
			if err != nil {
				return err
			}
			items = append(items, item)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort
	sort.Slice(items[:], func(i, j int) bool {
		return items[i].AddedAt >= items[j].AddedAt
	})

	return items, nil
}

// Put saves a list of pocket items
func (b *DB) Put(items []Item) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		for _, item := range items {
			v, err := json.Marshal(item)
			if err != nil {
				return err
			}
			b.Put([]byte(item.ID), v)
		}

		return nil
	})
}

// Delete deletes a pocket item
func (b *DB) Delete(item Item) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		return b.Delete([]byte(item.ID))
	})
}
