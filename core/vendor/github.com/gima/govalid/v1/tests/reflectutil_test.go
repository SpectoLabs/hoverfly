package govalid_test

import (
	"github.com/gima/govalid/v1/internal"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

var (
	s  = "asd"
	np *string
)

func TestReflectOrIndirect(t *testing.T) {
	_, err := internal.ReflectOrIndirect(nil)
	require.Error(t, err, "@nil")

	v, err := internal.ReflectOrIndirect(s)
	require.NoError(t, err, "@value")
	require.Equal(t, reflect.String, v.Kind(), "value kind")
	require.Equal(t, s, v.Interface(), "value")

	v, err = internal.ReflectOrIndirect(&s)
	require.NoError(t, err, "@ptr")
	require.Equal(t, reflect.String, v.Kind(), "ptr kind")
	require.Equal(t, s, v.Interface(), "ptr")

	_, err = internal.ReflectOrIndirect(np)
	require.Error(t, err, "@nil ptr")
}

func TestReflectPtrOrFabricate(t *testing.T) {
	_, err := internal.ReflectPtrOrFabricate(nil)
	require.Error(t, err, "@nil")

	v, err := internal.ReflectPtrOrFabricate(s)
	require.NoError(t, err, "@value")
	require.Equal(t, reflect.Ptr, v.Kind(), "value kind: "+v.Kind().String())
	require.Equal(t, s, v.Elem().Interface(), "value")

	v, err = internal.ReflectPtrOrFabricate(np)
	require.Error(t, err, "@nil ptr")
}
