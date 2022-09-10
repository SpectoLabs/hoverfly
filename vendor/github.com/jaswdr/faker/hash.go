package faker

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

//Hash is the faker struct for Hashing Functions
type Hash struct {
	Faker *Faker
}

//SHA256 returns a random sha256 based random hashed string
func (hash Hash) SHA256() string {
	hashFunction := sha256.New()
	randomString := hash.Faker.Lorem().Word()
	hashFunction.Write([]byte(randomString))
	return fmt.Sprintf("%x", string(hashFunction.Sum(nil)))
}

//SHA512 returns a random sha512 based random hashed string
func (hash Hash) SHA512() string {
	hashFunction := sha512.New()
	randomString := hash.Faker.Lorem().Word()
	hashFunction.Write([]byte(randomString))
	return fmt.Sprintf("%x", string(hashFunction.Sum(nil)))
}

//MD5 returns a random MD5 based random hashed string
func (hash Hash) MD5() string {
	hashFunction := md5.New()
	randomString := hash.Faker.Lorem().Word()
	hashFunction.Write([]byte(randomString))
	return fmt.Sprintf("%x", string(hashFunction.Sum(nil)))
}
