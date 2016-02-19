package backends

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"

	log "github.com/Sirupsen/logrus"
)

type User struct {
	UUID     string `json:"uuid" form:"-"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (u *User) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(u)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeUser(user []bytes) (*User, error) {
	var u *User
	buf := bytes.NewBuffer(user)
	dec := json.NewDecoder(buf)
	err := dec.Decode(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

type AuthBackend interface {
	SetValue(key, value []byte) error
	GetValue(key []byte) ([]byte, error)

	AddUser(username, password []byte) error
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

func (b *BoltAuth) AddUser(username, password []byte) error {
	err := b.DS.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(b.UserBucket)
		if err != nil {
			return err
		}
		hashedPassword, _ := bcrypt.GenerateFromPassword(password, 10)
		u := User{
			UUID:     uuid.New(),
			Username: username,
			Password: string(hashedPassword),
		}
		bts, err := u.Encode()
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err.Error(),
				"username": username,
			})
			return err
		}
		err = bucket.Put(username, bts)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (b *BoltAuth) SetValue(key, value []byte) error {
	err := b.DS.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(b.TokenBucket)
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (b *BoltAuth) Delete(key []byte) error {
	err := b.DS.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(b.TokenBucket)
		if err != nil {
			return err
		}
		err = bucket.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (b *BoltAuth) GetValue(key []byte) (value []byte, err error) {

	err = b.DS.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.TokenBucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", b.TokenBucket)
		}
		// "Byte slices returned from Bolt are only valid during a transaction."
		var buffer bytes.Buffer
		val := bucket.Get(key)

		// If it doesn't exist then it will return nil
		if val == nil {
			return fmt.Errorf("key %q not found \n", key)
		}

		buffer.Write(val)
		value = buffer.Bytes()
		return nil
	})

	return
}
