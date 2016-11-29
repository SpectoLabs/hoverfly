package cache

import (
	"testing"

	. "github.com/onsi/gomega"
)

var (
	expectedKey1   = []byte{'1'}
	expectedKey2   = []byte{'2'}
	expectedValue1 = []byte{'A'}
	expectedValue2 = []byte{'B'}
)

func TestCacheGetIsEmptyByDefault(t *testing.T) {
	RegisterTestingT(t)

	cache := &InMemoryCache{}

	expectedKey := []byte{'k'}

	actualValue, err := cache.Get(expectedKey)
	Expect(err).ToNot(BeNil())
	Expect(actualValue).To(HaveLen(0))
}

func TestCacheGetAllEntriesIsEmptyByDefault(t *testing.T) {
	RegisterTestingT(t)

	cache := &InMemoryCache{}

	actualValue, err := cache.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(actualValue).To(HaveLen(0))
}

func TestCacheGetAllKeysIsEmptyByDefault(t *testing.T) {
	RegisterTestingT(t)

	cache := &InMemoryCache{}

	actualValue, err := cache.GetAllKeys()
	Expect(err).To(BeNil())
	Expect(actualValue).To(HaveLen(0))
}

func TestCacheGetAllValuesIsEmptyByDefault(t *testing.T) {
	RegisterTestingT(t)

	cache := &InMemoryCache{}

	actualValue, err := cache.GetAllValues()
	Expect(err).To(BeNil())
	Expect(actualValue).To(HaveLen(0))
}

func TestSetAndGet(t *testing.T) {
	RegisterTestingT(t)

	cache := NewInMemoryCache()

	err := cache.Set(expectedKey1, expectedValue1)
	Expect(err).To(BeNil())

	err = cache.Set(expectedKey2, expectedValue2)
	Expect(err).To(BeNil())

	actualValue, err := cache.Get(expectedKey1)
	Expect(err).To(BeNil())
	Expect(actualValue).To(Equal(expectedValue1))

	actualValue, err = cache.Get(expectedKey2)
	Expect(err).To(BeNil())
	Expect(actualValue).To(Equal(expectedValue2))
}

func TestInMemoryCache_Get_DoesntLockIfFailed(t *testing.T) {
	RegisterTestingT(t)

	cache := NewInMemoryCache()

	_, err := cache.Get(expectedKey1)
	Expect(err).ToNot(BeNil())

	err = cache.Set(expectedKey1, expectedValue1)
	Expect(err).To(BeNil())

	actualValue, err := cache.Get(expectedKey1)
	Expect(err).To(BeNil())
	Expect(actualValue).To(Equal(expectedValue1))
}

func TestGetAllKeysMem(t *testing.T) {
	RegisterTestingT(t)

	cache := NewInMemoryCache()

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	keys, err := cache.GetAllKeys()
	Expect(err).To(BeNil())

	Expect(keys).To(HaveLen(2))
	Expect(keys).To(HaveKey(string(expectedKey1)))
	Expect(keys).To(HaveKey(string(expectedKey2)))
}

func TestGetAllValuesMem(t *testing.T) {
	RegisterTestingT(t)

	cache := NewInMemoryCache()

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	values, err := cache.GetAllValues()
	Expect(err).To(BeNil())

	Expect(values).To(HaveLen(2))
	Expect(values).To(ContainElement(expectedValue1))
	Expect(values).To(ContainElement(expectedValue2))
}

func TestGetAllEntriesMem(t *testing.T) {
	RegisterTestingT(t)

	cache := NewInMemoryCache()

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	entries, err := cache.GetAllEntries()
	Expect(err).To(BeNil())

	Expect(entries).To(HaveLen(2))
	Expect(entries).To(HaveKeyWithValue(string(expectedKey1), expectedValue1))
	Expect(entries).To(HaveKeyWithValue(string(expectedKey2), expectedValue2))
}

func TestGetRecordCount(t *testing.T) {
	RegisterTestingT(t)

	cache := NewInMemoryCache()

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	recordCount, err := cache.RecordsCount()
	Expect(err).To(BeNil())

	Expect(recordCount).To(Equal(2))
}

func TestDeleteData(t *testing.T) {
	RegisterTestingT(t)

	cache := NewInMemoryCache()

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	cache.DeleteData()

	recordCount, err := cache.RecordsCount()
	Expect(err).To(BeNil())

	Expect(recordCount).To(Equal(0))
}

func TestDeleteKey(t *testing.T) {
	RegisterTestingT(t)

	cache := NewInMemoryCache()

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	cache.Delete([]byte(expectedKey1))

	value, err := cache.Get(expectedKey1)
	Expect(err).ToNot(BeNil())

	Expect(value).To(HaveLen(0))

	count, err := cache.RecordsCount()
	Expect(err).To(BeNil())

	Expect(count).To(Equal(1))
}

func TestThrowErrorIfGettingANotExistingKey(t *testing.T) {
	RegisterTestingT(t)

	cache := NewInMemoryCache()
	_, err := cache.Get([]byte("key"))
	Expect(err).ToNot(BeNil())
}
