package backend

import (
	"github.com/ringtail/pansidong/types"
	"github.com/boltdb/bolt"
	"log"
	"encoding/json"
	"os"
	"errors"
	"path/filepath"
)

const (
	defaultBucket = "pansidong"
	defaultDBFile = "pansidong.db"
)

func NewBoltdbBackend(conf *types.BoltDBConfig) types.BackendStore {
	path := conf.Path
	if conf.Path == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal("Failed to getwd in boltdb backend,because of %s", err.Error())
		}
		path = filepath.Join(wd, defaultDBFile)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		newFile, err := os.Create(path)
		if err != nil {
			log.Fatal("Failed to create boltdb pansidong.db,because of %s", err.Error())
		}
		defer newFile.Close()
	}

	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatalf("Your boltdb path %s is not valid,because of %s", path, err.Error())
	}

	bb := &BoltdbBackend{
		db: db,
	}

	return bb
}

type BoltdbBackend struct {
	db *bolt.DB
}

// NotImplement
func (bb *BoltdbBackend) Next(options *types.ListOptions) ([]*types.ProxyIP, error) {
	return nil, nil
}
func (bb *BoltdbBackend) List(options *types.ListOptions) ([]*types.ProxyIP, error) {
	limit := options.Limit
	if limit == 0 {
		limit = 20
	}
	index := 0
	ips := make([]*types.ProxyIP, 0)
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))

		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if index >= limit {
				return nil
			}
			ip := &types.ProxyIP{}
			err := json.Unmarshal(v, ip)
			if err != nil {
				continue
			}
			index = index + 1
			ips = append(ips, ip)
		}
		return nil
	})
	return ips, err
}

// NotImplement
func (bb *BoltdbBackend) Get(key string) (*types.ProxyIP, error) {
	return nil, nil
}

func (bb *BoltdbBackend) Expire(key string) error {
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		v := b.Get([]byte(key))
		if v == nil {
			return errors.New("KeyNotFound")
		}
		err := b.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (bb *BoltdbBackend) Insert(ips []*types.ProxyIP, options *types.InsertOptions) error {
	err := bb.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		if err != nil {
			return err
		}
		for _, ip := range ips {
			ip_b, err := json.Marshal(ip)
			if err != nil {
				continue
			}
			err = b.Put([]byte(ip.IP), ip_b)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}
