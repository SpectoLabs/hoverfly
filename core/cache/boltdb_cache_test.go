package cache

import (
	. "github.com/onsi/gomega"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

// TestDB - holds connection to database during tests
var TestDB *bolt.DB

var testingDatabaseName = "bolt_test.db"

func TestSetGetValue(t *testing.T) {
	RegisterTestingT(t)
	db := NewBoltDBCache(TestDB, []byte("bucket1"))

	err := db.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	val, err := db.Get([]byte("foo"))
	Expect(string(val)).To(Equal("bar"))
}

func TestRecordsCountZero(t *testing.T) {
	RegisterTestingT(t)

	db := NewBoltDBCache(TestDB, []byte("bucketRecordsCountZero"))

	ct, err := db.RecordsCount()
	Expect(err).To(BeNil())
	Expect(ct).To(Equal(0))
}

func TestGetAllValues(t *testing.T) {
	RegisterTestingT(t)

	db := NewBoltDBCache(TestDB, []byte("bucketTestGetAllValues"))

	err := db.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	err = db.Set([]byte("foo2"), []byte("bar"))
	Expect(err).To(BeNil())

	vals, err := db.GetAllValues()
	Expect(err).To(BeNil())

	for i := 0; i < 2; i++ {
		Expect(string(vals[i])).To(Equal("bar"))
	}
}

func TestGetAllEntries(t *testing.T) {
	RegisterTestingT(t)

	db := NewBoltDBCache(TestDB, []byte("bucketTestGetAllEntries"))

	err := db.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	err = db.Set([]byte("foo2"), []byte("bar"))
	Expect(err).To(BeNil())

	vals, err := db.GetAllEntries()
	Expect(err).To(BeNil())

	for _, v := range vals {
		Expect(string(v)).To(Equal("bar"))
	}
}

func TestGetAllKeys(t *testing.T) {
	RegisterTestingT(t)

	db := NewBoltDBCache(TestDB, []byte("bucketTestGetAllKeys"))

	err := db.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	err = db.Set([]byte("foo2"), []byte("bar"))
	Expect(err).To(BeNil())

	keys, err := db.GetAllKeys()
	Expect(err).To(BeNil())

	Expect(keys).To(HaveKeyWithValue("foo", true))
	Expect(keys).To(HaveKeyWithValue("foo2", true))
	Expect(keys).ToNot(HaveKey("foo10"))
}

func TestDeleteRecords(t *testing.T) {
	RegisterTestingT(t)

	db := NewBoltDBCache(TestDB, []byte("bucketTestDeleteRecords"))

	err := db.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	err = db.Set([]byte("foo2"), []byte("bar"))
	Expect(err).To(BeNil())

	err = db.Set([]byte("foo3"), []byte("bar"))
	Expect(err).To(BeNil())

	ct, err := db.RecordsCount()
	Expect(err).To(BeNil())
	Expect(ct).To(Equal(3))

	err = db.DeleteData()
	Expect(err).To(BeNil())

	ctNew, err := db.RecordsCount()
	Expect(err).To(BeNil())
	Expect(ctNew).To(Equal(0))

}

func TestDeleteRecord(t *testing.T) {
	RegisterTestingT(t)

	db := NewBoltDBCache(TestDB, []byte("bucketTestDeleteRecord"))

	err := db.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	err = db.Delete([]byte("foo"))
	Expect(err).To(BeNil())

	_, err = db.Get([]byte("foo"))
	Expect(err).ToNot(BeNil())
}

func TestDeleteNotExisting(t *testing.T) {
	RegisterTestingT(t)

	db := NewBoltDBCache(TestDB, []byte("bucketTestDeleteNotExisting"))

	err := db.Delete([]byte("foo"))
	Expect(err).To(BeNil())
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
