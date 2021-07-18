package bolt

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/tarrsalah/pkt"
	"github.com/tarrsalah/pkt/internal/config"
	"go.etcd.io/bbolt"
	"path/filepath"
)

var bucketName = []byte("pkt")

type DB struct {
	db *bbolt.DB
}

func NewDB() DB {
	path := filepath.Join(config.Dir(), "pkt.bolt")
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	return DB{db}
}

func (b DB) Close() {
	b.db.Close()
}

func (b DB) Get() []pkt.Item {
	var items []pkt.Item

	err := b.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return nil
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var item pkt.Item
			err := json.Unmarshal(v, &item)
			if err != nil {
				return err
			}
			items = append(items, item)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	// Sort
	sort.Slice(items[:], func(i, j int) bool {
		return items[i].AddedAt >= items[j].AddedAt
	})

	return items
}

func (b DB) Put(items []pkt.Item) {
	err := b.db.Update(func(tx *bbolt.Tx) error {
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

	if err != nil {
		panic(err)
	}
}
