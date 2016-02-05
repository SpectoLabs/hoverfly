package hoverfly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestSetKey(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	k := []byte("randomkeyhere")
	v := []byte("value")

	err := dbClient.Cache.Set(k, v)
	expect(t, err, nil)

	value, err := dbClient.Cache.Get(k)
	expect(t, err, nil)
	expect(t, string(value), string(v))
	dbClient.Cache.DeleteBucket(dbClient.Cache.RequestsBucket)
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

	dbClient.Cache.DeleteBucket(dbClient.Cache.RequestsBucket)
}

func TestGetNonExistingBucket(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.Cache.RequestsBucket = []byte("some_random_bucket")

	_, err := dbClient.Cache.Get([]byte("whatever"))
	expect(t, err.Error(), "Bucket \"some_random_bucket\" not found!")
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
	err = dbClient.Cache.DeleteBucket(dbClient.Cache.RequestsBucket)
	expect(t, err, nil)

	// deleting it again
	err = dbClient.Cache.DeleteBucket(dbClient.Cache.RequestsBucket)
	refute(t, err, nil)
}

func TestGetAllRequestNoBucket(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.Cache.RequestsBucket = []byte("no_bucket_for_TestGetAllRequestNoBucket")
	_, err := dbClient.Cache.GetAllRequests()
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
	payloads, err := dbClient.Cache.GetAllRequests()
	expect(t, err, nil)
	expect(t, len(payloads), 0)

}

func TestGetMultipleRecords(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteBucket(dbClient.Cache.RequestsBucket)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.captureRequest(req)
	}

	// getting requests
	payloads, err := dbClient.Cache.GetAllRequests()
	expect(t, err, nil)

	for _, payload := range payloads {
		expect(t, payload.Request.Method, "GET")
		expect(t, payload.Response.Status, 201)
	}
}

func TestGetNonExistingKey(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteBucket(dbClient.Cache.RequestsBucket)

	// getting key
	_, err := dbClient.Cache.Get([]byte("should not be here"))
	refute(t, err, nil)
}

func TestSetGetEmptyValue(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteBucket(dbClient.Cache.RequestsBucket)

	err := dbClient.Cache.Set([]byte("shouldbe"), []byte(""))
	expect(t, err, nil)
	// getting key
	_, err = dbClient.Cache.Get([]byte("shouldbe"))
	expect(t, err, nil)
}
