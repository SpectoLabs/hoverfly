package hoverfly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestSetKey(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	k := []byte("randomkeyhere")
	v := []byte("value")

	err := dbClient.Cache.Set(k, v)
	expect(t, err, nil)

	value, err := dbClient.Cache.Get(k)
	expect(t, err, nil)
	expect(t, string(value), string(v))
}

func TestPayloadSetGet(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	key := []byte("keySetGetCache")
	resp := ResponseDetails{
		Status: 200,
		Body:   "body here",
	}

	payload := Payload{Response: resp}
	bts, err := json.Marshal(payload)
	expect(t, err, nil)

	err = dbClient.Cache.Set(key, bts)
	expect(t, err, nil)

	var p Payload
	payloadBts, err := dbClient.Cache.Get(key)
	err = json.Unmarshal(payloadBts, &p)
	expect(t, err, nil)
	expect(t, payload.Response.Body, p.Response.Body)

	defer dbClient.Cache.DeleteData()
}

func TestGetNonExistingBucket(t *testing.T) {
	cache := NewBoltDBCache(TestDB, []byte("somebucket"))

	_, err := cache.Get([]byte("whatever"))
	expect(t, err.Error(), "Bucket \"somebucket\" not found!")
}

func TestDeleteBucket(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	k := []byte("randomkeyhere")
	v := []byte("value")
	// checking whether bucket is okay
	err := dbClient.Cache.Set(k, v)
	expect(t, err, nil)

	value, err := dbClient.Cache.Get(k)
	expect(t, err, nil)
	expect(t, string(value), string(v))

	// deleting bucket
	err = dbClient.Cache.DeleteData()
	expect(t, err, nil)

	// deleting it again
	err = dbClient.Cache.DeleteData()
	refute(t, err, nil)
}

func TestGetAllRequestNoBucket(t *testing.T) {
	cache := NewBoltDBCache(TestDB, []byte("somebucket"))

	cache.RequestsBucket = []byte("no_bucket_for_TestGetAllRequestNoBucket")
	_, err := cache.GetAllValues()
	// expecting nil since this would mean that records were wiped
	expect(t, err, nil)
}

func TestCorruptedPayloads(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	k := []byte("randomkeyhere")
	v := []byte("value")

	err := dbClient.Cache.Set(k, v)
	expect(t, err, nil)

	// corrupted payloads should be just skipped
	payloads, err := dbClient.Cache.GetAllValues()
	expect(t, err, nil)
	expect(t, len(payloads), 1)

}

func TestGetMultipleRecords(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.captureRequest(req)
	}

	// getting requests
	values, err := dbClient.Cache.GetAllValues()
	expect(t, err, nil)

	for _, value := range values {
		if payload, err := decodePayload(value); err == nil {
			expect(t, payload.Request.Method, "GET")
			expect(t, payload.Response.Status, 201)
		} else {
			t.Error(err)
		}
	}
}

func TestGetNonExistingKey(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	// getting key
	_, err := dbClient.Cache.Get([]byte("should not be here"))
	refute(t, err, nil)
}

func TestSetGetEmptyValue(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	err := dbClient.Cache.Set([]byte("shouldbe"), []byte(""))
	expect(t, err, nil)
	// getting key
	_, err = dbClient.Cache.Get([]byte("shouldbe"))
	expect(t, err, nil)
}

func TestGetAllKeys(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		dbClient.Cache.Set([]byte(fmt.Sprintf("key%d", i)), []byte("value"))
	}

	keys, err := dbClient.Cache.GetAllKeys()
	expect(t, err, nil)
	expect(t, len(keys), 5)

	for k, v := range keys {
		expect(t, strings.HasPrefix(k, "key"), true)
		expect(t, v, true)
	}
}

func TestGetAllKeysEmpty(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	keys, err := dbClient.Cache.GetAllKeys()
	expect(t, err, nil)
	expect(t, len(keys), 0)
}
