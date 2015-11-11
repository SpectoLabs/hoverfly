package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

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

	// getting redis configuration
	redisAddress := os.Getenv("RedisAddress")
	if redisAddress == "" {
		redisAddress = ":6379"
	}
	AppConfig.redisAddress = redisAddress

	redisPool := getRedisPool()
	//	defer redisPool.Close()

	cache := Cache{pool: redisPool}

	// preparing client
	dbClient := &DBClient{
		http:  &http.Client{Transport: tr},
		cache: cache,
	}

	//	client := &Client{httpClient}
	return server, dbClient
}
