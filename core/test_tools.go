package hoverfly

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/boltdb/bolt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

const testingDatabaseName = "test.db" // very original

// TestDB - holds connection to database during tests
var testDB *bolt.DB

func testTools(code int, body string) (*httptest.Server, *Hoverfly) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	cfg := InitSettings()
	// disabling auth for testing
	cfg.AuthEnabled = false

	dbClient := GetNewHoverfly(cfg, cache.NewDefaultLRUCache(), nil)

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	dbClient.HTTP = &http.Client{Transport: tr}

	return server, dbClient
}

var src = rand.NewSource(time.Now().UnixNano())

// getRandomName - provides random name for buckets. Each test case gets it's own bucket
func getRandomName(n int) []byte {
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
	testDB = db
}

// teardown does some cleanup after tests
func teardown() {
	testDB.Close()
	os.Remove(testingDatabaseName)
}
