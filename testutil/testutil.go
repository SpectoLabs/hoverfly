package testutil

import (
	"reflect"
	"testing"
)

func Expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Contains(slice [][]byte, element []byte) bool {
	for _, v := range slice {
		if reflect.DeepEqual(element, v) {
			return true
		}
	}
	return false
}
