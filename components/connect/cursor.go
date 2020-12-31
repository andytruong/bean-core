package connect

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func DecodeCursor(hash string) (*Cursor, error) {
	b, err := base64.StdEncoding.DecodeString(hash)
	if nil != err {
		return nil, err
	}

	chunks := strings.Split(string(b), ":")
	if len(chunks) != 3 {
		return nil, errors.New("invalid format")
	}

	return &Cursor{
		Entity:   chunks[0],
		Property: chunks[1],
		Value:    chunks[2],
	}, nil
}

type Cursor struct {
	Entity   string
	Property string
	Value    string
}

func (cursor Cursor) Encode() string {
	outcome := fmt.Sprintf("%v:%v:%v", cursor.Entity, cursor.Property, cursor.Value)

	return base64.StdEncoding.EncodeToString([]byte(outcome))
}
