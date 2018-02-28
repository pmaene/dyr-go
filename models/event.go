package models

import (
    "encoding/binary"
    "encoding/json"
    "strconv"
    "time"

    "dyr/library"

    "github.com/boltdb/bolt"
)

type Event struct {
    ID           uint64    `json:"id"`
    Accessory    string    `json:"accessory"`
    User         string    `json:"user"`
    CreationTime time.Time `json:"creation_time"`
}

func (e *Event) Bucket() []byte {
    return []byte("events")
}
func (e *Event) Key() []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(e.ID))

    return b
}
func (e *Event) SetKey(k []byte) {
    e.ID = binary.BigEndian.Uint64(k)
}
func (e *Event) EncodeKey() string {
    return strconv.FormatUint(e.ID, 10)
}
func (e *Event) DecodeKey(k string) error {
    id, err := strconv.ParseUint(k, 10, 64)
    if err != nil {
        return err
    }

    e.ID = id
    return nil
}

func (e *Event) Get(db *bolt.DB) error {
    bm := library.BoltModel{e}
    return bm.Get(db)
}

func GetEvents(db *bolt.DB) (map[int]Event, error) {
    r := make(map[int]Event)

    err := db.View(func(tx *bolt.Tx) error {
        bkt := tx.Bucket([]byte("events"))
        if bkt == nil {
            return nil
        }

        i := 0
        bkt.ForEach(func(k, v []byte) error {
            var e Event
            if err := json.Unmarshal(v, &e); err != nil {
                return err
            }

            r[i] = e
            i = i + 1

            return nil
        })

        return nil
    })

    return r, err
}
func CreateEvent(db *bolt.DB, a string, u string) error {
    return db.Update(func(tx *bolt.Tx) error {
        b, err := tx.CreateBucketIfNotExists([]byte("events"))
        if err != nil {
            return err
        }

        id, err := b.NextSequence()
        if err != nil {
            return err
        }

        e := Event{uint64(id), a, u, time.Now()}

        v, err := json.Marshal(e)
        if err != nil {
            return err
        }

        return b.Put(e.Key(), v)
    })
}
