package hoverfly

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/boltdb/bolt"
	"net/url"
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

// TestDB - holds connection to database during tests
var TestDB *bolt.DB

func testTools(code int, body string) (*httptest.Server, *Hoverfly) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	// creating random buckets for everyone!
	bucket := GetRandomName(10)
	metaBucket := GetRandomName(10)

	requestCache := cache.NewBoltDBCache(TestDB, bucket)
	metaCache := cache.NewBoltDBCache(TestDB, metaBucket)

	cfg := InitSettings()
	// disabling auth for testing
	cfg.AuthEnabled = false

	dbClient := GetNewHoverfly(cfg, requestCache, metaCache, nil)

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	dbClient.HTTP = &http.Client{Transport: tr}

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
	db := cache.GetDB(testingDatabaseName)
	TestDB = db
}

// teardown does some cleanup after tests
func teardown() {
	TestDB.Close()
	os.Remove(testingDatabaseName)
}
