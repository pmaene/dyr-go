package models

import (
    "encoding/json"
    "time"

    "dyr/library"

    "github.com/boltdb/bolt"
    jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
    Username string `json:"username"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Admin    bool   `json:"admin"`
}
func (u User) GetTokenString(e time.Time, s library.Scope) (string, error) {
    t := jwt.New(jwt.SigningMethodHS256)

	claims := t.Claims.(jwt.MapClaims)
	claims["sub"] = u.Username
    claims["iat"] = time.Now().Unix()
    claims["exp"] = e.Unix()

    claims["admin"] = u.Admin
    claims["scope"] = s

	token, err := t.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

    return token, nil
}

func (u *User) Bucket() []byte {
    return []byte("users")
}
func (u *User) Key() []byte {
    return []byte(u.Username)
}
func (u *User) SetKey(k []byte) {
    u.Username = string(k)
}
func (u *User) EncodeKey() string {
    return string(u.Key())
}
func (u *User) DecodeKey(k string) error {
    u.Username = k
    return nil
}

func (u User) Put(db *bolt.DB) error {
    bm := library.BoltModel{&u}
    return bm.Put(db)
}
func (u *User) Get(db *bolt.DB) error {
    bm := library.BoltModel{u}
    return bm.Get(db)
}
func (u User) Delete(db *bolt.DB) error {
    var p Password

    if err := p.DecodeKey(u.Username); err != nil {
        return err
    }

    if err := p.Delete(db); err != nil {
        return err
    }

    bm := library.BoltModel{&u}
    return bm.Delete(db)
}

func GetUsers(db *bolt.DB) (map[int]User, error) {
    r := make(map[int]User)

    err := db.View(func(tx *bolt.Tx) error {
        bkt := tx.Bucket([]byte("users"))
        if bkt == nil {
            return nil
        }

        i := 0
        bkt.ForEach(func(k, v []byte) error {
            var u User
            if err := json.Unmarshal(v, &u); err != nil {
                return err
            }

            r[i] = u
            i = i + 1

            return nil
        })

        return nil
    })

    return r, err
}
func GetTokenUser(db *bolt.DB, c jwt.MapClaims) (User, error) {
    var u User
    if err := u.DecodeKey(c["sub"].(string)); err != nil {
        return User{}, err
    }

    if err := u.Get(db); err != nil {
        return User{}, err
    }

    return u, nil
}
