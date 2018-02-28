package library

import (
    "errors"
    "encoding/json"

    "github.com/boltdb/bolt"
)

var (
    BucketDoesNotExistError = errors.New("bucket does not exist")
    KeyNotFoundError = errors.New("key not found")
)

type BoltValue interface {
    Bucket()          []byte

    Key()             []byte
    SetKey([]byte)
    EncodeKey()       string
    DecodeKey(string) error
}

type BoltModel struct {
    Value BoltValue
}
func (bm BoltModel) Put(db *bolt.DB) error {
    return db.Update(func(tx *bolt.Tx) error {
        b, err := tx.CreateBucketIfNotExists(bm.Value.Bucket())
        if err != nil {
            return err
        }

        v, err := json.Marshal(bm.Value)
        if err != nil {
            return err
        }

        return b.Put(bm.Value.Key(), v)
    })
}
func (bm *BoltModel) Get(db *bolt.DB) error {
    return db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(bm.Value.Bucket())
        if b == nil {
            return BucketDoesNotExistError
        }

        v := b.Get(bm.Value.Key())
        if v == nil {
            return KeyNotFoundError
        }

        if err := json.Unmarshal(v, bm.Value); err != nil {
            return err
        }

        return nil
    })
}
func (bm BoltModel) Delete(db *bolt.DB) error {
    return db.Update(func(tx *bolt.Tx) error {
        b, err := tx.CreateBucketIfNotExists(bm.Value.Bucket())
        if err != nil {
            return err
        }

        return b.Delete(bm.Value.Key())
    })
}
