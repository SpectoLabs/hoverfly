package hoverfly

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

const testingDatabaseName = "test.db" // very original

// Client structure to be injected into functions to perform HTTP calls
type Client struct {
	HTTPClient *http.Client
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

// TestDB - holds connection to database during tests
var TestDB *bolt.DB

func testTools(code int, body string) (*httptest.Server, *DBClient) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	// creating random buckets for everyone!
	bucket := GetRandomName(10)
	metaBucket := GetRandomName(10)

	cache := NewBoltDBCache(TestDB, bucket)
	md := NewBoltDBMetadata(TestDB, metaBucket)

	cfg := InitSettings()
	counter := NewModeCounter()
	// preparing client
	dbClient := &DBClient{
		HTTP:    &http.Client{Transport: tr},
		Cache:   cache,
		Cfg:     cfg,
		Counter: counter,
		MD:      md,
	}
	return server, dbClient
}

var src = rand.NewSource(time.Now().UnixNano())

// GetRandomName - provides random name for buckets. Each test case gets it's own bucket
func GetRandomName(n int) []byte {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func setup() {
	// we don't really want to see what's happening
	log.SetLevel(log.FatalLevel)
	db := GetDB(testingDatabaseName)
	TestDB = db
}

// teardown does some cleanup after tests
func teardown() {
	TestDB.Close()
	os.Remove(testingDatabaseName)
}
