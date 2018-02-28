package models

import (
    "dyr/library"

    "github.com/boltdb/bolt"
    "github.com/lhecker/argon2"
)

type Password struct {
    User string `json:"user"`
    Hash []byte `json:"hash"`
}
func (p Password) Verify(pwd string) (bool, error) {
    h, err := argon2.Decode(p.Hash)
    if err != nil {
        return false, err
    }

    return h.Verify([]byte(pwd))
}

func (p *Password) Bucket() []byte {
    return []byte("passwords")
}
func (p *Password) Key() []byte {
    return []byte(p.User)
}
func (p *Password) SetKey(k []byte) {
    p.User = string(k)
}
func (p *Password) EncodeKey() string {
    return string(p.Key())
}
func (p *Password) DecodeKey(k string) error {
    p.User = k
    return nil
}

func (p Password) Put(db *bolt.DB) error {
    bm := library.BoltModel{&p}
    return bm.Put(db)
}
func (p *Password) Get(db *bolt.DB) error {
    bm := library.BoltModel{p}
    return bm.Get(db)
}
func (p Password) Delete(db *bolt.DB) error {
    bm := library.BoltModel{&p}
    return bm.Delete(db)
}

func NewPassword(db *bolt.DB, u string, pwd string) (Password, error) {
    cfg := argon2.DefaultConfig()

    h, err := cfg.Hash([]byte(pwd), nil)
    if err != nil {
        return Password{}, err
    }

    return Password{u, h.Encode()}, nil
}
