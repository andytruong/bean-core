package util

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/speps/go-hashids"
)

type Identifier struct {
}

func (this *Identifier) Encode(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

func (this *Identifier) HashInt64(salt string, current time.Time) (string, error) {
	data := hashids.NewData()
	data.Salt = salt
	data.MinLength = 7
	hash, _ := hashids.NewWithData(data)

	return hash.EncodeInt64([]int64{current.Unix()})
}

func (this *Identifier) MustULID() string {
	val, err := this.ULID()

	if nil != err {
		panic(err)
	}

	return val
}

func (this *Identifier) ULID() (string, error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id, err := ulid.New(ulid.Now(), entropy)

	if nil != err {
		return "", err
	}

	return id.String(), nil
}

func (this *Identifier) UUID() (string, error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	uid := ulid.MustNew(ulid.Now(), entropy)
	id, err := uuid.FromBytes(uid[:])

	if nil != err {
		return "", err
	}

	return id.String(), nil
}
