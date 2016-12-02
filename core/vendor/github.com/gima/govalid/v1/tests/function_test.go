package govalid_test

import (
	"fmt"
	v "github.com/gima/govalid/v1"
	"reflect"
	"strings"
	"testing"
)

func TestFunction(t *testing.T) {
	counter, countingValidator := createCountingValidator()

	test(t, "combination1", false, v.Function(dogeValidator), "cat")
	test(t, "combination2", true, v.Function(dogeValidator), "cate")
	test(t, "combination3", true, v.Function(countingValidator), "doge")
	test(t, "combination4", true, v.Function(countingValidator), "doge")

	if *counter != 2 {
		t.Fatalf("counting validator count should be 2, is %d", *counter)
	}
}

func dogeValidator(data interface{}) (path string, err error) {
	s, ok := data.(string)
	if !ok {
		return "doge-validator", fmt.Errorf("expected string, got %s", reflect.TypeOf(data))
	}

	if !strings.HasSuffix(strings.ToLower(s), "e") {
		return "doge-validator", fmt.Errorf("expected string to end in 'e'")
	}

	return "", nil
}

func createCountingValidator() (counter *int, _ v.ValidatorFunc) {
	counter = new(int)
	return counter, func(data interface{}) (path string, err error) {
		*counter++
		return "", nil
	}
}
