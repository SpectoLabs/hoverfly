package backends

import (
	"bytes"
	"encoding/json"

	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
)

type User struct {
	UUID     string `json:"uuid" form:"-"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	IsAdmin  bool   `json:"is_admin" form:"is_admin"`
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

func DecodeUser(user []byte) (*User, error) {
	var u *User
	buf := bytes.NewBuffer(user)
	dec := json.NewDecoder(buf)
	err := dec.Decode(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Authentication - generic interface for authentication backend
type Authentication interface {
	AddUser(username, password string, admin bool) (err error)
	AddUserHashedPassword(username, passwordHash string, admin bool) (err error)
	GetUser(username string) (user *User, err error)
	GetAllUsers() (users []User, err error)
	InvalidateToken(token string) (err error)
	IsTokenBlacklisted(token string) (blacklisted bool, err error)
}

// NewCacheBasedAuthBackend - takes two caches - one for token and one for users
func NewCacheBasedAuthBackend(tokenCache, userCache cache.Cache) *CacheAuthBackend {
	return &CacheAuthBackend{
		TokenCache: tokenCache,
		userCache:  userCache,
	}
}

// UserBucketName - default name for BoltDB bucket that stores user info
const UserBucketName = "authbucket"

// TokenBucketName
const TokenBucketName = "tokenbucket"

// CacheAuthBackend - container to implement Cache instance with i.e. BoltDB backend for storage
type CacheAuthBackend struct {
	TokenCache cache.Cache
	userCache  cache.Cache
}

// AddUser - adds user with provided username, password and admin parameters
func (b *CacheAuthBackend) AddUser(username, password string, admin bool) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := User{
		UUID:     uuid.New(),
		Username: username,
		Password: string(hashedPassword),
		IsAdmin:  admin,
	}
	userBytes, err := user.Encode()
	if err != nil {
		logUserError(err, username)
		return err
	}
	err = b.userCache.Set([]byte(username), userBytes)
	return err
}

// AddUserHashedPassword - adds user with provided username, hashed password and admin parameters
func (b *CacheAuthBackend) AddUserHashedPassword(username, hashedPassword string, admin bool) error {
	user := User{
		UUID:     uuid.New(),
		Username: username,
		Password: hashedPassword,
		IsAdmin:  admin,
	}
	userBytes, err := user.Encode()
	if err != nil {
		logUserError(err, username)
		return err
	}
	err = b.userCache.Set([]byte(username), userBytes)
	return err
}

func (b *CacheAuthBackend) GetUser(username string) (user *User, err error) {
	userBytes, err := b.userCache.Get([]byte(username))

	if err != nil {
		logUserError(err, username)
		return
	}

	user, err = DecodeUser(userBytes)

	if err != nil {
		logUserError(err, username)
		return
	}
	return
}

func (b *CacheAuthBackend) InvalidateToken(token string) error {
	return b.TokenCache.Set([]byte(token), []byte("whentoexpire"))
}

// IsTokenBlacklisted - checks if token is blacklisted.
func (b *CacheAuthBackend) IsTokenBlacklisted(token string) (bool, error) {
	blacklistedToken, err := b.TokenCache.Get([]byte(token))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"token": token,
		}).Debug("got error while looking for blacklisted token")
		if blacklistedToken != nil {
			return true, err
		}
		return false, err
	}
	if blacklistedToken == nil {
		return false, nil
	}
	return true, nil

}

func (b *CacheAuthBackend) GetAllUsers() (users []User, err error) {
	values, _ := b.userCache.GetAllValues()
	users = make([]User, len(values), len(values))
	for i, user := range values {
		decodedUser, err := DecodeUser(user)
		users[i] = *decodedUser
		return users, err
	}
	return users, err
}

func logUserError(err error, username string) {
	log.WithFields(log.Fields{
		"error":    err.Error(),
		"username": username,
	})
}
