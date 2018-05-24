package cache_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/cache"
	. "github.com/onsi/gomega"
)

var (
	expectedKey1   = []byte{'1'}
	expectedKey2   = []byte{'2'}
	expectedValue1 = []byte{'A'}
	expectedValue2 = []byte{'B'}
)

func Test_InMemoryCache_Get_NonExistingKey(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewBoltDBCache(TestDB, []byte("somebucket"))

	// getting key
	_, err := unit.Get([]byte("should not be here"))
	Expect(err).ToNot(BeNil())
}

func Test_InMemoryCache_CacheGetIsEmptyByDefault(t *testing.T) {
	RegisterTestingT(t)

	unit := &cache.InMemoryCache{}

	expectedKey := []byte{'k'}

	actualValue, err := unit.Get(expectedKey)
	Expect(err).ToNot(BeNil())
	Expect(actualValue).To(HaveLen(0))
}

func Test_InMemoryCache_CacheGetAllEntriesIsEmptyByDefault(t *testing.T) {
	RegisterTestingT(t)

	unit := &cache.InMemoryCache{}

	actualValue, err := unit.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(actualValue).To(HaveLen(0))
}

func Test_InMemoryCache_CacheGetAllKeysIsEmptyByDefault(t *testing.T) {
	RegisterTestingT(t)

	unit := &cache.InMemoryCache{}

	actualValue, err := unit.GetAllKeys()
	Expect(err).To(BeNil())
	Expect(actualValue).To(HaveLen(0))
}

func Test_InMemoryCache_CacheGetAllValuesIsEmptyByDefault(t *testing.T) {
	RegisterTestingT(t)

	unit := &cache.InMemoryCache{}

	actualValue, err := unit.GetAllValues()
	Expect(err).To(BeNil())
	Expect(actualValue).To(HaveLen(0))
}

func Test_InMemoryCache_SetAndGet(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewInMemoryCache()

	err := unit.Set(expectedKey1, expectedValue1)
	Expect(err).To(BeNil())

	err = unit.Set(expectedKey2, expectedValue2)
	Expect(err).To(BeNil())

	actualValue, err := unit.Get(expectedKey1)
	Expect(err).To(BeNil())
	Expect(actualValue).To(Equal(expectedValue1))

	actualValue, err = unit.Get(expectedKey2)
	Expect(err).To(BeNil())
	Expect(actualValue).To(Equal(expectedValue2))
}

func Test_InMemoryCache_Get_DoesntLockIfFailed(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewInMemoryCache()

	_, err := unit.Get(expectedKey1)
	Expect(err).ToNot(BeNil())

	err = unit.Set(expectedKey1, expectedValue1)
	Expect(err).To(BeNil())

	actualValue, err := unit.Get(expectedKey1)
	Expect(err).To(BeNil())
	Expect(actualValue).To(Equal(expectedValue1))
}

func Test_InMemoryCache_GetAllKeysMem(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewInMemoryCache()

	unit.Set(expectedKey1, expectedValue1)
	unit.Set(expectedKey2, expectedValue2)

	keys, err := unit.GetAllKeys()
	Expect(err).To(BeNil())

	Expect(keys).To(HaveLen(2))
	Expect(keys).To(HaveKey(string(expectedKey1)))
	Expect(keys).To(HaveKey(string(expectedKey2)))
}

func Test_InMemoryCache_GetAllValuesMem(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewInMemoryCache()

	unit.Set(expectedKey1, expectedValue1)
	unit.Set(expectedKey2, expectedValue2)

	values, err := unit.GetAllValues()
	Expect(err).To(BeNil())

	Expect(values).To(HaveLen(2))
	Expect(values).To(ContainElement(expectedValue1))
	Expect(values).To(ContainElement(expectedValue2))
}

func Test_InMemoryCache_GetAllEntriesMem(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewInMemoryCache()

	unit.Set(expectedKey1, expectedValue1)
	unit.Set(expectedKey2, expectedValue2)

	entries, err := unit.GetAllEntries()
	Expect(err).To(BeNil())

	Expect(entries).To(HaveLen(2))
	Expect(entries).To(HaveKeyWithValue(string(expectedKey1), expectedValue1))
	Expect(entries).To(HaveKeyWithValue(string(expectedKey2), expectedValue2))
}

func Test_InMemoryCache_GetRecordCount(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewInMemoryCache()

	unit.Set(expectedKey1, expectedValue1)
	unit.Set(expectedKey2, expectedValue2)

	recordCount, err := unit.RecordsCount()
	Expect(err).To(BeNil())

	Expect(recordCount).To(Equal(2))
}

func Test_InMemoryCache_DeleteData(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewInMemoryCache()

	unit.Set(expectedKey1, expectedValue1)
	unit.Set(expectedKey2, expectedValue2)

	unit.DeleteData()

	recordCount, err := unit.RecordsCount()
	Expect(err).To(BeNil())

	Expect(recordCount).To(Equal(0))
}

func Test_InMemoryCache_DeleteKey(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewInMemoryCache()

	unit.Set(expectedKey1, expectedValue1)
	unit.Set(expectedKey2, expectedValue2)

	unit.Delete([]byte(expectedKey1))

	value, err := unit.Get(expectedKey1)
	Expect(err).ToNot(BeNil())

	Expect(value).To(HaveLen(0))

	count, err := unit.RecordsCount()
	Expect(err).To(BeNil())

	Expect(count).To(Equal(1))
}

func Test_InMemoryCache_ThrowErrorIfGettingANotExistingKey(t *testing.T) {
	RegisterTestingT(t)

	unit := cache.NewInMemoryCache()
	_, err := unit.Get([]byte("key"))
	Expect(err).ToNot(BeNil())
}
