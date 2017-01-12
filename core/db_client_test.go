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

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	k := []byte("randomkeyhere")
	v := []byte("value")

	err := dbClient.RequestCache.Set(k, v)
	Expect(err).To(BeNil())

	value, err := dbClient.RequestCache.Get(k)
	Expect(err).To(BeNil())
	Expect(value).To(Equal(v))
}

func TestCacheSetGetWithRequestResponsePair(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	key := []byte("keySetGetCache")
	resp := models.ResponseDetails{
		Status: 200,
		Body:   "body here",
	}

	pair := models.RequestResponsePair{Response: resp}
	pairBytes, err := json.Marshal(pair)
	Expect(err).To(BeNil())

	err = dbClient.RequestCache.Set(key, pairBytes)
	Expect(err).To(BeNil())

	var p models.RequestResponsePair
	pairBytes, err = dbClient.RequestCache.Get(key)
	err = json.Unmarshal(pairBytes, &p)
	Expect(err).To(BeNil())
	Expect(pair.Response.Body).To(Equal(p.Response.Body))

	defer dbClient.RequestCache.DeleteData()
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
	err := dbClient.RequestCache.Set(k, v)
	Expect(err).To(BeNil())

	value, err := dbClient.RequestCache.Get(k)
	Expect(err).To(BeNil())
	Expect(value).To(Equal(v))

	// deleting bucket
	err = dbClient.RequestCache.DeleteData()
	Expect(err).To(BeNil())

	// deleting it again
	err = dbClient.RequestCache.DeleteData()
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

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	k := []byte("randomkeyhere")
	v := []byte("value")

	err := dbClient.RequestCache.Set(k, v)
	Expect(err).To(BeNil())

	// corrupted payloads should be just skipped
	pairBytes, err := dbClient.RequestCache.GetAllValues()
	Expect(err).To(BeNil())
	Expect(pairBytes).To(HaveLen(1))
}

func TestGetMultipleRecords(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		dbClient.Save(&models.RequestDetails{
			Method:      "GET",
			Scheme:      "http",
			Destination: "example.com",
			Query:       fmt.Sprintf("q=%d", i),
		}, &models.ResponseDetails{
			Status: 201,
			Body:   "ok",
		})
	}

	// getting requests
	values, err := dbClient.RequestCache.GetAllValues()
	Expect(err).To(BeNil())

	for _, value := range values {
		if pair, err := models.NewRequestResponsePairFromBytes(value); err == nil {
			Expect(pair.Request.Method).To(Equal("GET"))
			Expect(pair.Response.Status).To(Equal(201))
		} else {
			t.Error(err)
		}
	}
}

func TestGetNonExistingKey(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// getting key
	_, err := dbClient.RequestCache.Get([]byte("should not be here"))
	Expect(err).ToNot(BeNil())
}

func TestSetGetEmptyValue(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.RequestCache.Set([]byte("shouldbe"), []byte(""))
	Expect(err).To(BeNil())
	// getting key
	_, err = dbClient.RequestCache.Get([]byte("shouldbe"))
	Expect(err).To(BeNil())
}

func TestGetAllKeys(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		dbClient.RequestCache.Set([]byte(fmt.Sprintf("key%d", i)), []byte("value"))
	}

	keys, err := dbClient.RequestCache.GetAllKeys()
	Expect(err).To(BeNil())
	Expect(keys).To(HaveLen(5))

	for k, v := range keys {
		Expect(k).To(HavePrefix("key"))
		Expect(v).To(BeTrue())
	}
}

func TestGetAllKeysEmpty(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	keys, err := dbClient.RequestCache.GetAllKeys()
	Expect(err).To(BeNil())
	Expect(keys).To(HaveLen(0))
}
