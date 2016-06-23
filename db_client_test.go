package hoverfly

import (
	"encoding/json"
	"fmt"
	"github.com/SpectoLabs/hoverfly/cache"
	"github.com/SpectoLabs/hoverfly/core/testutil"
	"net/http"
	"strings"
	"testing"
	"github.com/SpectoLabs/hoverfly/models"
)

func TestSetKey(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	k := []byte("randomkeyhere")
	v := []byte("value")

	err := dbClient.RequestCache.Set(k, v)
	testutil.Expect(t, err, nil)

	value, err := dbClient.RequestCache.Get(k)
	testutil.Expect(t, err, nil)
	testutil.Expect(t, string(value), string(v))
}

func TestPayloadSetGet(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	key := []byte("keySetGetCache")
	resp := models.ResponseDetails{
		Status: 200,
		Body:   "body here",
	}

	payload := models.Payload{Response: resp}
	bts, err := json.Marshal(payload)
	testutil.Expect(t, err, nil)

	err = dbClient.RequestCache.Set(key, bts)
	testutil.Expect(t, err, nil)

	var p models.Payload
	payloadBts, err := dbClient.RequestCache.Get(key)
	err = json.Unmarshal(payloadBts, &p)
	testutil.Expect(t, err, nil)
	testutil.Expect(t, payload.Response.Body, p.Response.Body)

	defer dbClient.RequestCache.DeleteData()
}

func TestGetNonExistingBucket(t *testing.T) {
	cache := cache.NewBoltDBCache(TestDB, []byte("somebucket"))

	_, err := cache.Get([]byte("whatever"))
	testutil.Expect(t, err.Error(), "Bucket \"somebucket\" not found!")
}

func TestDeleteBucket(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	k := []byte("randomkeyhere")
	v := []byte("value")
	// checking whether bucket is okay
	err := dbClient.RequestCache.Set(k, v)
	testutil.Expect(t, err, nil)

	value, err := dbClient.RequestCache.Get(k)
	testutil.Expect(t, err, nil)
	testutil.Expect(t, string(value), string(v))

	// deleting bucket
	err = dbClient.RequestCache.DeleteData()
	testutil.Expect(t, err, nil)

	// deleting it again
	err = dbClient.RequestCache.DeleteData()
	testutil.Refute(t, err, nil)
}

func TestGetAllRequestNoBucket(t *testing.T) {
	cache := cache.NewBoltDBCache(TestDB, []byte("somebucket"))

	cache.CurrentBucket = []byte("no_bucket_for_TestGetAllRequestNoBucket")
	_, err := cache.GetAllValues()
	// expecting nil since this would mean that records were wiped
	testutil.Expect(t, err, nil)
}

func TestCorruptedPayloads(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	k := []byte("randomkeyhere")
	v := []byte("value")

	err := dbClient.RequestCache.Set(k, v)
	testutil.Expect(t, err, nil)

	// corrupted payloads should be just skipped
	payloads, err := dbClient.RequestCache.GetAllValues()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, len(payloads), 1)

}

func TestGetMultipleRecords(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		testutil.Expect(t, err, nil)
		dbClient.captureRequest(req)
	}

	// getting requests
	values, err := dbClient.RequestCache.GetAllValues()
	testutil.Expect(t, err, nil)

	for _, value := range values {
		if payload, err := models.NewPayloadFromBytes(value); err == nil {
			testutil.Expect(t, payload.Request.Method, "GET")
			testutil.Expect(t, payload.Response.Status, 201)
		} else {
			t.Error(err)
		}
	}
}

func TestGetNonExistingKey(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// getting key
	_, err := dbClient.RequestCache.Get([]byte("should not be here"))
	testutil.Refute(t, err, nil)
}

func TestSetGetEmptyValue(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.RequestCache.Set([]byte("shouldbe"), []byte(""))
	testutil.Expect(t, err, nil)
	// getting key
	_, err = dbClient.RequestCache.Get([]byte("shouldbe"))
	testutil.Expect(t, err, nil)
}

func TestGetAllKeys(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		dbClient.RequestCache.Set([]byte(fmt.Sprintf("key%d", i)), []byte("value"))
	}

	keys, err := dbClient.RequestCache.GetAllKeys()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, len(keys), 5)

	for k, v := range keys {
		testutil.Expect(t, strings.HasPrefix(k, "key"), true)
		testutil.Expect(t, v, true)
	}
}

func TestGetAllKeysEmpty(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	keys, err := dbClient.RequestCache.GetAllKeys()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, len(keys), 0)
}
