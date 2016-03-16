package boltdb

import (
	"os"
	"reflect"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

// TestDB - holds connection to database during tests
var TestDB *bolt.DB

var testingDatabaseName = "bolt_test.db"

func TestSetGetValue(t *testing.T) {
	db := NewBoltDBCache(TestDB, []byte("bucket1"))

	err := db.Set([]byte("foo"), []byte("bar"))
	expect(t, err, nil)

	val, err := db.Get([]byte("foo"))
	expect(t, string(val), "bar")
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

// TestMain prepares database for testing and then performs a cleanup
func TestMain(m *testing.M) {

	setup()

	retCode := m.Run()

	// delete test database
	teardown()

	// call with result of m.Run()
	os.Exit(retCode)
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
