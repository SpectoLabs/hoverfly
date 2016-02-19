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

func NewBoltDBAuthBackend(db *bolt.DB, tokenBucket, userBucket []byte) *BoltAuth {
	return &BoltAuth{
		DS:          db,
		TokenBucket: []byte(tokenBucket),
		UserBucket:  []byte(userBucket),
	}
}

// UserBucketName - default name for BoltDB bucket that stores user info
const UserBucketName = "authbucket"

// TokenBucketName
const TokenBucketName = "tokenbucket"

// BoltCache - container to implement Cache instance with BoltDB backend for storage
type BoltAuth struct {
	DS          *bolt.DB
	TokenBucket []byte
	UserBucket  []byte
}

