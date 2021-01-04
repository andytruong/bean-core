package scalar

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/speps/go-hashids"
)

type Identifier struct {
}

func (idr *Identifier) Encode(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

func (idr *Identifier) HashInt64(salt string, current time.Time) (string, error) {
	data := hashids.NewData()
	data.Salt = salt
	data.MinLength = 7
	hash, _ := hashids.NewWithData(data)

	return hash.EncodeInt64([]int64{current.Unix()})
}

func (idr *Identifier) ULID() string {
	val, _ := idr.ulid()

	return val
}

func (idr *Identifier) ulid() (string, error) {
	seed := time.Now().UnixNano()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(seed)), 0)
	val, err := ulid.New(ulid.Now(), entropy)

	if nil != err {
		panic(err)
	}

	return val.String(), nil
}
