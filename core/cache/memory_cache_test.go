package cache

import (
	"github.com/SpectoLabs/hoverfly/testutil"
	"reflect"
	"testing"
)

func TestEmptyByDefault(t *testing.T) {
	cache := &InMemoryCache{}

	expectedKey := []byte{'k'}

	if actualValue, err := cache.Get(expectedKey); err != nil {
		t.Fatalf("err: %v", err)
	} else if len(actualValue) != 0 {
		t.Fatal("Cache should be empty by default")
	}

	if actualValue, err := cache.GetAllEntries(); err != nil {
		t.Fatalf("err: %v", err)
	} else if len(actualValue) != 0 {
		t.Fatal("Cache should be empty by default")
	}

	if actualValue, err := cache.GetAllKeys(); err != nil {
		t.Fatalf("err: %v", err)
	} else if len(actualValue) != 0 {
		t.Fatal("Cache should be empty by default")
	}

	if actualValue, err := cache.GetAllValues(); err != nil {
		t.Fatalf("err: %v", err)
	} else if len(actualValue) != 0 {
		t.Fatal("Cache should be empty by default")
	}

	if actualValue, err := cache.GetAllEntries(); err != nil {
		t.Fatalf("err: %v", err)
	} else if len(actualValue) != 0 {
		t.Fatal("Cache should be empty by default")
	}
}

func TestSetAndGet(t *testing.T) {
	cache := NewInMemoryCache()

	expectedKey1 := []byte{'1'}
	expectedKey2 := []byte{'2'}
	expectedValue1 := []byte{'A'}
	expectedValue2 := []byte{'B'}

	if err := cache.Set(expectedKey1, expectedValue1); err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := cache.Set(expectedKey2, expectedValue2); err != nil {
		t.Fatalf("err: %v", err)
	}

	if actualValue, err := cache.Get(expectedKey1); err != nil {
		t.Fatalf("err: %v", err)
	} else if !reflect.DeepEqual(expectedValue1, actualValue) {
		t.Fatalf("Expected value %s does not match actual value %s", expectedValue1, actualValue)
	}

	if actualValue, err := cache.Get(expectedKey2); err != nil {
		t.Fatalf("err: %v", err)
	} else if !reflect.DeepEqual(expectedValue2, actualValue) {
		t.Fatalf("Expected value %s does not match actual value %s", expectedValue1, actualValue)
	}
}

func TestGetAllKeysMem(t *testing.T) {

	cache := NewInMemoryCache()

	expectedKey1 := []byte{'1'}
	expectedKey2 := []byte{'2'}
	expectedValue1 := []byte{'A'}
	expectedValue2 := []byte{'B'}

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	if actualKey, err := cache.GetAllKeys(); err != nil {
		t.Fatalf("err: %v", err)
	} else if len(actualKey) != 2 {
		t.Fatal("Cache should have two keys")
	} else if !actualKey[string(expectedKey1)] {
		t.Fatalf("Cache does not contain key %s", expectedKey1)
	} else if !actualKey[string(expectedKey2)] {
		t.Fatalf("Cache does not contain key %s", expectedKey2)
	}
}

func TestGetAllValuesMem(t *testing.T) {

	cache := NewInMemoryCache()

	expectedKey1 := []byte{'1'}
	expectedKey2 := []byte{'2'}
	expectedValue1 := []byte{'A'}
	expectedValue2 := []byte{'B'}

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	if values, err := cache.GetAllValues(); err != nil {
		t.Fatalf("err: %v", err)
	} else if len(values) != 2 {
		t.Fatalf("Expected %v values but got %s", 2, len(values))
	} else if !testutil.Contains(values, expectedValue1) {
		t.Fatalf("Value %s was not present in %s", expectedValue1, values)
	} else if !testutil.Contains(values, expectedValue2) {
		t.Fatalf("Value %s was not present in %s", expectedValue2, values)
	}
}

func TestGetAllEntriesMem(t *testing.T) {
	cache := NewInMemoryCache()

	expectedKey1 := []byte{'1'}
	expectedKey2 := []byte{'2'}
	expectedValue1 := []byte{'A'}
	expectedValue2 := []byte{'B'}

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	if entries, err := cache.GetAllEntries(); err != nil {
		t.Fatalf("err: %v", err)
	} else if len(entries) != 2 {
		t.Fatalf("Expected %v values but got %v", 2, len(entries))
	} else if reflect.DeepEqual(entries[string(expectedKey1)], string(expectedValue1)) {
		t.Fatalf("Entry %s %s was not present in %s", expectedKey1, expectedValue1, entries)
	} else if reflect.DeepEqual(entries[string(expectedKey1)], string(expectedValue2)) {
		t.Fatalf("Entry %s %s was not present in %s", expectedKey2, expectedValue2, entries)
	}
}

func TestGetRecordCount(t *testing.T) {
	cache := NewInMemoryCache()

	expectedKey1 := []byte{'1'}
	expectedKey2 := []byte{'2'}
	expectedValue1 := []byte{'A'}
	expectedValue2 := []byte{'B'}

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	if recordCount, _ := cache.RecordsCount(); recordCount != 2 {
		t.Fatalf("Expected %v records but got %v", 2, recordCount)
	}
}

func TestDeleteData(t *testing.T) {
	cache := NewInMemoryCache()

	expectedKey1 := []byte{'1'}
	expectedKey2 := []byte{'2'}
	expectedValue1 := []byte{'A'}
	expectedValue2 := []byte{'B'}

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	cache.DeleteData()

	if recordCount, _ := cache.RecordsCount(); recordCount != 0 {
		t.Fatalf("Expected %v records but got %v", 0, recordCount)
	}
}

func TestDeleteKey(t *testing.T) {
	cache := NewInMemoryCache()

	expectedKey1 := []byte{'1'}
	expectedKey2 := []byte{'2'}
	expectedValue1 := []byte{'A'}
	expectedValue2 := []byte{'B'}

	cache.Set(expectedKey1, expectedValue1)
	cache.Set(expectedKey2, expectedValue2)

	cache.Delete([]byte(expectedKey1))

	if value, err := cache.Get(expectedKey1); err != nil {
		t.Fatalf("Error %v", err)
	} else if len(value) != 0 {
		t.Fatalf("Expected nil value but got %v", value)
	} else if value, _ := cache.RecordsCount(); value != 1 {
		t.Fatalf("Expected %v records but got %v", 1, cache.RecordsCount)
	}
}
