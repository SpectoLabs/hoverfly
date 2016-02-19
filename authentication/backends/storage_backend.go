package backends

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
)

type AuthBackend interface {
	SetValue(key, value []byte) error
	GetValue(key []byte) ([]byte, error)
	Delete(key []byte) error
}
