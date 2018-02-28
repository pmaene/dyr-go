package models

import (
    "bufio"
    "encoding/binary"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net"
    "strconv"
    "strings"
    "time"

    "dyr/library"

    "github.com/ankitkalbande/simonspeck"
    "github.com/boltdb/bolt"
)

const (
    ArduinoTimeout = time.Duration(5)*time.Second
)

var (
    ArduinoError = errors.New("arduino")
)

type Door struct {
    Accessory

    Host        string  `json:"host"`
    Port        uint32  `json:"port"`

    Latitude    float64 `json:"latitude"`
    Longitude   float64 `json:"longitude"`
    MaxDistance float64 `json:"max_distance"`

    Nonce       uint32  `json:"nonce"`
}
func (d *Door) IncrementNonce(db *bolt.DB) error {
    d.Nonce = d.Nonce + 1

    if err := d.Put(db); err != nil {
        return err
    }

    return nil
}
func (d Door) Switch(db *bolt.DB) error {
    if err := d.IncrementNonce(db); err != nil {
        return err
    }

    // speck implementation expects data in reverse, i.e., first word in hex
    // string is used as the first key word instead of the last (notation used
    // the paper)
    //
    // unsigned long key[4] = {0x75b7a326, 0x38aed491, 0x735e4aa9, 0x2e83e923};
    k, err := hex.DecodeString("75b7a32638aed491735e4aa92e83e923")
    if err != nil {
        return err
    }

    s := simonspeck.NewSpeck64(convertBytes(k))

    h, err := strconv.Atoi(strings.Replace(d.Host, ".", "", -1))
    if err != nil {
        return err
    }

    src := make([]byte, 8)
    binary.BigEndian.PutUint32(src[0:4], d.Nonce)
    binary.BigEndian.PutUint32(src[4:8], uint32(h))

    dst := make([]byte, 8)
    s.Encrypt(dst, convertBytes(src))
    t := hex.EncodeToString(convertBytes(dst))
    m := fmt.Sprintf("switch/nonce/%08x/token/%s", d.Nonce, t)

    var r string
    r, err = send(d, m)
    if err != nil {
        return err
    }

    if r == "status/error" {
        return ArduinoError
    }

    return nil
}

func (d *Door) Bucket() []byte {
    return []byte("doors")
}
func (d *Door) Key() []byte {
    return []byte(d.Name)
}
func (d *Door) SetKey(k []byte) {
    d.Name = string(k)
}
func (d *Door) EncodeKey() string {
    return string(d.Key())
}
func (d *Door) DecodeKey(k string) error {
    d.Name = k
    return nil
}

func (d Door) Put(db *bolt.DB) error {
    bm := library.BoltModel{&d}
    return bm.Put(db)
}
func (d *Door) Get(db *bolt.DB) error {
    bm := library.BoltModel{d}
    return bm.Get(db)
}
func (d Door) Delete(db *bolt.DB) error {
    bm := library.BoltModel{&d}
    return bm.Delete(db)
}

func GetDoors(db *bolt.DB) (map[int]Door, error) {
    r := make(map[int]Door)

    err := db.View(func(tx *bolt.Tx) error {
        bkt := tx.Bucket([]byte("doors"))
        if bkt == nil {
            return nil
        }

        i := 0
        bkt.ForEach(func(k, v []byte) error {
            var d Door
            if err := json.Unmarshal(v, &d); err != nil {
                return err
            }

            r[i] = d
            i = i + 1

            return nil
        })

        return nil
    })

    return r, err
}

func send(d Door, msg string) (string, error) {
    c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", d.Host, d.Port), ArduinoTimeout)
    defer c.Close()

    if err != nil {
        return "", err
    }

    fmt.Fprintf(c, msg)

    var r string
    r, err = bufio.NewReader(c).ReadString('\n')
    if err != nil {
        if err != io.EOF {
            return "", err
        }
    }

    return r, nil
}

// endianness conversion
func convertBytes(b []byte) []byte {
    var n = len(b)
	for i := 0; i < n; i += 4 {
		b[i], b[i+1], b[i+2], b[i+3] = b[i+3], b[i+2], b[i+1], b[i]
	}
	return b
}
