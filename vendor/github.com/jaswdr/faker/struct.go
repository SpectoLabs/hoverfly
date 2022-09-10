package faker

import (
	"reflect"
	"strconv"
)

// Struct is a faker struct for Struct
type Struct struct {
	Faker *Faker
}

// Fill elements of a struct with random data
func (s Struct) Fill(v interface{}) {
	s.r(reflect.TypeOf(v), reflect.ValueOf(v), "", 0)
}

func (s Struct) r(t reflect.Type, v reflect.Value, function string, size int) {
	switch t.Kind() {
	case reflect.Ptr:
		s.rPointer(t, v, function)
	case reflect.Struct:
		s.rStruct(t, v)
	case reflect.String:
		s.rString(t, v, function)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s.rUint(t, v, function)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.rInt(t, v, function)
	case reflect.Float32, reflect.Float64:
		s.rFloat(t, v, function)
	case reflect.Bool:
		s.rBool(t, v)
	case reflect.Array, reflect.Slice:
		s.rSlice(t, v, function, size)
	}
}

func (s Struct) rStruct(t reflect.Type, v reflect.Value) {
	n := t.NumField()
	for i := 0; i < n; i++ {
		elementT := t.Field(i)
		elementV := v.Field(i)
		t, ok := elementT.Tag.Lookup("fake")
		if ok && t == "skip" {
			// Do nothing, skip it
		} else if elementV.CanSet() {
			// Check if fakesize is set
			size := -1 // Set to -1 to indicate fakesize was not set
			fs, ok := elementT.Tag.Lookup("fakesize")
			if ok {
				var err error
				size, err = strconv.Atoi(fs)
				if err != nil {
					size = s.Faker.IntBetween(1, 10)
				}
			}
			s.r(elementT.Type, elementV, t, size)
		}
	}
}

func (s Struct) rPointer(t reflect.Type, v reflect.Value, function string) {
	elemT := t.Elem()
	if v.IsNil() {
		nv := reflect.New(elemT)
		s.r(elemT, nv.Elem(), function, 0)
		v.Set(nv)
	} else {
		s.r(elemT, v.Elem(), function, 0)
	}
}

func (s Struct) rSlice(t reflect.Type, v reflect.Value, function string, size int) {
	// If you cant even set it dont even try
	if !v.CanSet() {
		return
	}

	// Grab original size to use if needed for sub arrays
	ogSize := size

	// If the value has a cap and is less than the size
	// use that instead of the requested size
	elemCap := v.Cap()
	if elemCap == 0 && size == -1 {
		size = s.Faker.IntBetween(1, 10)
	} else if elemCap != 0 && (size == -1 || elemCap < size) {
		size = elemCap
	}

	// Get the element type
	elemT := t.Elem()

	// If values are already set fill them up, otherwise append
	if v.Len() != 0 {
		// Loop through the elements length and set based upon the index
		for i := 0; i < size; i++ {
			nv := reflect.New(elemT)
			s.r(elemT, nv.Elem(), function, ogSize)
			v.Index(i).Set(reflect.Indirect(nv))
		}
	} else {
		// Loop through the size and append and set
		for i := 0; i < size; i++ {
			nv := reflect.New(elemT)
			s.r(elemT, nv.Elem(), function, ogSize)
			v.Set(reflect.Append(reflect.Indirect(v), reflect.Indirect(nv)))
		}
	}
}

func (s Struct) rString(t reflect.Type, v reflect.Value, function string) {
	if function != "" {
		v.SetString(s.Faker.Bothify(function))
	} else {
		v.SetString(s.Faker.UUID().V4())
	}
}

func (s Struct) rInt(t reflect.Type, v reflect.Value, function string) {
	if function != "" {
		i, err := strconv.ParseInt(s.Faker.Numerify(function), 10, 64)
		if err == nil {
			v.SetInt(i)
			return
		}
	}

	// If no function or error converting to int, set with random value
	switch t.Kind() {
	case reflect.Int:
		v.SetInt(s.Faker.Int64())
	case reflect.Int8:
		v.SetInt(int64(s.Faker.Int8()))
	case reflect.Int16:
		v.SetInt(int64(s.Faker.Int16()))
	case reflect.Int32:
		v.SetInt(int64(s.Faker.Int32()))
	case reflect.Int64:
		v.SetInt(s.Faker.Int64())
	}
}

func (s Struct) rUint(t reflect.Type, v reflect.Value, function string) {
	if function != "" {
		u, err := strconv.ParseUint(s.Faker.Numerify(function), 10, 64)
		if err == nil {
			v.SetUint(u)
			return
		}
	}

	// If no function or error converting to uint, set with random value
	switch t.Kind() {
	case reflect.Uint:
		v.SetUint(s.Faker.UInt64())
	case reflect.Uint8:
		v.SetUint(uint64(s.Faker.UInt8()))
	case reflect.Uint16:
		v.SetUint(uint64(s.Faker.UInt16()))
	case reflect.Uint32:
		v.SetUint(uint64(s.Faker.UInt32()))
	case reflect.Uint64:
		v.SetUint(s.Faker.UInt64())
	}
}

func (s Struct) rFloat(t reflect.Type, v reflect.Value, function string) {
	if function != "" {
		f, err := strconv.ParseFloat(s.Faker.Numerify(function), 64)
		if err == nil {
			v.SetFloat(f)
			return
		}
	}

	// If no function or error converting to float, set with random value
	switch t.Kind() {
	case reflect.Float64:
		v.SetFloat(s.Faker.Float64(2, 0, 100))
	case reflect.Float32:
		v.SetFloat(s.Faker.Float64(2, 0, 100))
	}
}

func (s Struct) rBool(t reflect.Type, v reflect.Value) {
	v.SetBool(s.Faker.Bool())
}
