package cache_test

import (
	"os"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/cache"
	. "github.com/onsi/gomega"

	log "github.com/sirupsen/logrus"
	"github.com/boltdb/bolt"
)

// TestDB - holds connection to database during tests
var TestDB *bolt.DB

var testingDatabaseName = "bolt_test.db"

func Test_BoltDBCache_Get_NonExistingKey(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewBoltDBCache(TestDB, []byte("somebucket"))

	// getting key
	_, err := unit.Get([]byte("should not be here"))
	Expect(err).ToNot(BeNil())
}

func Test_BoltDBCache_Get_NonExistingBucket(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewBoltDBCache(TestDB, []byte("somebucket"))

	_, err := unit.Get([]byte("whatever"))
	Expect(err).ToNot(BeNil())
	Expect(err).To(MatchError("Bucket \"somebucket\" not found!"))
}

func Test_BoltDBCache_Set_GetValue(t *testing.T) {
	RegisterTestingT(t)
	unit := cache.NewBoltDBCache(TestDB, []byte("bucket1"))

	err := unit.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	val, err := unit.Get([]byte("foo"))
	Expect(string(val)).To(Equal("bar"))
}

func Test_BoltDBCache_RecordsCountZero(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewBoltDBCache(TestDB, []byte("bucketRecordsCountZero"))

	ct, err := unit.RecordsCount()
	Expect(err).To(BeNil())
	Expect(ct).To(Equal(0))
}

func Test_BoltDBCache_GetAllValues(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewBoltDBCache(TestDB, []byte("bucketTestGetAllValues"))

	err := unit.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	err = unit.Set([]byte("foo2"), []byte("bar"))
	Expect(err).To(BeNil())

	vals, err := unit.GetAllValues()
	Expect(err).To(BeNil())

	for i := 0; i < 2; i++ {
		Expect(string(vals[i])).To(Equal("bar"))
	}
}

func Test_BoltDBCache_GetAllEntries(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewBoltDBCache(TestDB, []byte("bucketTestGetAllEntries"))

	err := unit.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	err = unit.Set([]byte("foo2"), []byte("bar"))
	Expect(err).To(BeNil())

	vals, err := unit.GetAllEntries()
	Expect(err).To(BeNil())

	for _, v := range vals {
		Expect(string(v)).To(Equal("bar"))
	}
}

func Test_BoltDBCache_GetAllKeys(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewBoltDBCache(TestDB, []byte("bucketTestGetAllKeys"))

	err := unit.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	err = unit.Set([]byte("foo2"), []byte("bar"))
	Expect(err).To(BeNil())

	keys, err := unit.GetAllKeys()
	Expect(err).To(BeNil())

	Expect(keys).To(HaveKeyWithValue("foo", true))
	Expect(keys).To(HaveKeyWithValue("foo2", true))
	Expect(keys).ToNot(HaveKey("foo10"))
}

func Test_BoltDBCache_DeleteData(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewBoltDBCache(TestDB, []byte("bucketTestDeleteRecords"))

	err := unit.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	err = unit.Set([]byte("foo2"), []byte("bar"))
	Expect(err).To(BeNil())

	err = unit.Set([]byte("foo3"), []byte("bar"))
	Expect(err).To(BeNil())

	ct, err := unit.RecordsCount()
	Expect(err).To(BeNil())
	Expect(ct).To(Equal(3))

	err = unit.DeleteData()
	Expect(err).To(BeNil())

	ctNew, err := unit.RecordsCount()
	Expect(err).To(BeNil())
	Expect(ctNew).To(Equal(0))

}

func Test_BoltDBCache_Delete(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewBoltDBCache(TestDB, []byte("bucketTestDeleteRecord"))

	err := unit.Set([]byte("foo"), []byte("bar"))
	Expect(err).To(BeNil())

	err = unit.Delete([]byte("foo"))
	Expect(err).To(BeNil())

	_, err = unit.Get([]byte("foo"))
	Expect(err).ToNot(BeNil())
}

func Test_BoltDBCache_Delete_NotExisting(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewBoltDBCache(TestDB, []byte("bucketTestDeleteNotExisting"))

	err := unit.Delete([]byte("foo"))
	Expect(err).To(BeNil())
}

func setup() {
	// we don't really want to see what's happening
	log.SetLevel(log.FatalLevel)
	unit := cache.GetDB(testingDatabaseName)
	TestDB = unit
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
