package scalar

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
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

func (idr *Identifier) MustULID() string {
	val, err := idr.ULID()

	if nil != err {
		panic(err)
	}

	return val
}

func (idr *Identifier) ULID() (string, error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id, err := ulid.New(ulid.Now(), entropy)

	if nil != err {
		return "", err
	}

	return id.String(), nil
}
