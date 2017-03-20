package hoverfly

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func TestSetKey(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	k := []byte("randomkeyhere")
	v := []byte("value")

	err := unit.CacheMatcher.RequestCache.Set(k, v)
	Expect(err).To(BeNil())

	value, err := unit.CacheMatcher.RequestCache.Get(k)
	Expect(err).To(BeNil())
	Expect(value).To(Equal(v))
}

func TestCacheSetGetWithRequestResponsePair(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	key := []byte("keySetGetCache")
	resp := models.ResponseDetails{
		Status: 200,
		Body:   "body here",
	}

	pair := models.RequestResponsePair{Response: resp}
	pairBytes, err := json.Marshal(pair)
	Expect(err).To(BeNil())

	err = unit.CacheMatcher.RequestCache.Set(key, pairBytes)
	Expect(err).To(BeNil())

	var p models.RequestResponsePair
	pairBytes, err = unit.CacheMatcher.RequestCache.Get(key)
	err = json.Unmarshal(pairBytes, &p)
	Expect(err).To(BeNil())
	Expect(pair.Response.Body).To(Equal(p.Response.Body))
}

func TestGetNonExistingBucket(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewBoltDBCache(TestDB, []byte("somebucket"))

	_, err := cache.Get([]byte("whatever"))
	Expect(err).ToNot(BeNil())
	Expect(err).To(MatchError("Bucket \"somebucket\" not found!"))
}

func TestDeleteBucket(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	k := []byte("randomkeyhere")
	v := []byte("value")
	// checking whether bucket is okay
	err := dbClient.CacheMatcher.RequestCache.Set(k, v)
	Expect(err).To(BeNil())

	value, err := dbClient.CacheMatcher.RequestCache.Get(k)
	Expect(err).To(BeNil())
	Expect(value).To(Equal(v))

	// deleting bucket
	err = dbClient.CacheMatcher.RequestCache.DeleteData()
	Expect(err).To(BeNil())

	// deleting it again
	err = dbClient.CacheMatcher.RequestCache.DeleteData()
	Expect(err).ToNot(BeNil())
}

func TestGetAllRequestNoBucket(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewBoltDBCache(TestDB, []byte("somebucket"))

	cache.CurrentBucket = []byte("no_bucket_for_TestGetAllRequestNoBucket")
	_, err := cache.GetAllValues()
	// expecting nil since this would mean that records were wiped
	Expect(err).To(BeNil())
}

func TestCorruptedPairs(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	k := []byte("randomkeyhere")
	v := []byte("value")

	err := unit.CacheMatcher.RequestCache.Set(k, v)
	Expect(err).To(BeNil())

	// corrupted payloads should be just skipped
	pairBytes, err := unit.CacheMatcher.RequestCache.GetAllValues()
	Expect(err).To(BeNil())
	Expect(pairBytes).To(HaveLen(1))
}

func TestGetNonExistingKey(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	// getting key
	_, err := unit.CacheMatcher.RequestCache.Get([]byte("should not be here"))
	Expect(err).ToNot(BeNil())
}

func TestGetAllKeys(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	// inserting some payloads
	for i := 0; i < 5; i++ {
		unit.CacheMatcher.RequestCache.Set([]byte(fmt.Sprintf("key%d", i)), []byte("value"))
	}

	keys, err := unit.CacheMatcher.RequestCache.GetAllKeys()
	Expect(err).To(BeNil())
	Expect(keys).To(HaveLen(5))

	for k, v := range keys {
		Expect(k).To(HavePrefix("key"))
		Expect(v).To(BeTrue())
	}
}

func TestGetAllKeysEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	keys, err := unit.CacheMatcher.RequestCache.GetAllKeys()
	Expect(err).To(BeNil())
	Expect(keys).To(HaveLen(0))
}
