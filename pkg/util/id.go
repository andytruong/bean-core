package util

import (
	"time"

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

func (this *Identifier) UUID() (string, error) {
	panic("wip")
}
