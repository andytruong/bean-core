package util

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/speps/go-hashids"
)

type Identifier struct {
}

func (this *Identifier) Hash(entityType string, current time.Time) (string, error) {
	data := hashids.NewData()
	data.Salt = "QEjpevoA7yN1V:" + entityType
	data.MinLength = 7
	hash, _ := hashids.NewWithData(data)

	return hash.EncodeInt64([]int64{current.Unix()})
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
