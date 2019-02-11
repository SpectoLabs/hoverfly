package cache_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/cache"
	. "github.com/onsi/gomega"
)

var (
	testKey1   = '1'
	testKey2   = '2'
	testValue1 = 'A'
	testValue2 = 'B'
)

func Test_LRUFastCache_Get_NonExistingKey(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewDefaultLRUCache()

	value, found := unit.Get("should not be here")
	Expect(value).To(BeNil())
	Expect(found).To(BeFalse())
}


func Test_LRUFastCache_CacheGetAllEntriesIsEmptyByDefault(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewDefaultLRUCache()

	actualValue, err := unit.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(actualValue).To(HaveLen(0))
}


func Test_LRUFastCache_SetAndGet(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewDefaultLRUCache()

	err := unit.Set(testKey1, testValue1)
	Expect(err).To(BeNil())

	err = unit.Set(testKey2, testValue2)
	Expect(err).To(BeNil())

	actualValue, found := unit.Get(testKey1)
	Expect(found).To(BeTrue())
	Expect(actualValue).To(Equal(testValue1))

	actualValue, found = unit.Get(testKey2)
	Expect(found).To(BeTrue())
	Expect(actualValue).To(Equal(testValue2))
}

func Test_LRUFastCache_GetAllEntries(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewDefaultLRUCache()

	unit.Set(testKey1, testValue1)
	unit.Set(testKey2, testValue2)

	entries, err := unit.GetAllEntries()
	Expect(err).To(BeNil())

	Expect(entries).To(HaveLen(2))
	Expect(entries).To(HaveKeyWithValue(testKey1, testValue1))
	Expect(entries).To(HaveKeyWithValue(testKey2, testValue2))
}

func Test_LRUFastCache_GetRecordCount(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewDefaultLRUCache()

	unit.Set(testKey1, testValue1)
	unit.Set(testKey2, testValue2)

	recordCount, err := unit.RecordsCount()
	Expect(err).To(BeNil())

	Expect(recordCount).To(Equal(2))
}

func Test_LRUFastCache_DeleteData(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewDefaultLRUCache()

	unit.Set(testKey1, testValue1)
	unit.Set(testKey2, testValue2)

	unit.DeleteData()

	recordCount, err := unit.RecordsCount()
	Expect(err).To(BeNil())

	Expect(recordCount).To(Equal(0))
}
