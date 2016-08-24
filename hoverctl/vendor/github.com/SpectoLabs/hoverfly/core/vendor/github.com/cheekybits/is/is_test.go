package is

import (
	"errors"
	"strings"
	"testing"
)

type customErr struct{}

func (e *customErr) Error() string {
	return "Oops"
}

type mockT struct {
	failed bool
}

func (m *mockT) FailNow() {
	m.failed = true
}
func (m *mockT) Failed() bool {
	return m.failed
}

func TestIs(t *testing.T) {

	for _, test := range []struct {
		N     string
		F     func(is I)
		Fails []string
	}{
		{
			N: "Fail('msg')",
			F: func(is I) {
				is.Fail("something")
			},
			Fails: []string{"failed: something"},
		},
		{
			N: "Failf('%d is wrong',123)",
			F: func(is I) {
				is.Failf("%d is wrong", 123)
			},
			Fails: []string{"failed: 123 is wrong"},
		},
		// is.Nil
		{
			N: "Nil(nil)",
			F: func(is I) {
				is.Nil(nil)
			},
		},
		{
			N: "Nil(\"nope\")",
			F: func(is I) {
				is.Nil("nope")
			},
			Fails: []string{"expected nil: \"nope\""},
		},
		// is.NotNil
		{
			N: "NotNil(\"nope\")",
			F: func(is I) {
				is.NotNil("nope")
			},
		},
		{
			N: "NotNil(nil)",
			F: func(is I) {
				is.NotNil(nil)
			},
			Fails: []string{"unexpected nil"},
		},
		// is.OK
		{
			N: "OK(false)",
			F: func(is I) {
				is.OK(false)
			},
			Fails: []string{"unexpected false"},
		}, {
			N: "OK(true)",
			F: func(is I) {
				is.OK(true)
			},
		}, {
			N: "OK(nil)",
			F: func(is I) {
				is.OK(nil)
			},
			Fails: []string{"unexpected nil"},
		}, {
			N: "OK(1,2,3)",
			F: func(is I) {
				is.OK(1, 2, 3)
			},
		}, {
			N: "OK(0)",
			F: func(is I) {
				is.OK(0)
			},
			Fails: []string{"unexpected zero"},
		}, {
			N: "OK(1)",
			F: func(is I) {
				is.OK(1)
			},
		}, {
			N: "OK(\"\")",
			F: func(is I) {
				is.OK("")
			},
			Fails: []string{"unexpected \"\""},
		},
		// NoErr
		{
			N: "NoErr(errors.New(\"an error\"))",
			F: func(is I) {
				is.NoErr(errors.New("an error"))
			},
			Fails: []string{"unexpected error: an error"},
		}, {
			N: "NoErr(&customErr{})",
			F: func(is I) {
				is.NoErr(&customErr{})
			},
			Fails: []string{"unexpected error: Oops"},
		}, {
			N: "NoErr(error(nil))",
			F: func(is I) {
				var err error
				is.NoErr(err)
			},
		},
		{
			N: "NoErr(err1, err2, err3)",
			F: func(is I) {
				is.NoErr(&customErr{}, &customErr{}, &customErr{})
			},
			Fails: []string{"unexpected error: Oops"},
		},
		{
			N: "NoErr(err1, err2, err3)",
			F: func(is I) {
				var err1 error
				var err2 error
				var err3 error
				is.NoErr(err1, err2, err3)
			},
		},
		// Err
		{
			N: "Err(errors.New(\"an error\"))",
			F: func(is I) {
				is.Err(errors.New("an error"))
			},
		}, {
			N: "Err(&customErr{})",
			F: func(is I) {
				is.Err(&customErr{})
			},
		}, {
			N: "Err(error(nil))",
			F: func(is I) {
				var err error
				is.Err(err)
			},
			Fails: []string{"error expected"},
		},
		{
			N: "Err(customErr1, customErr2, customErr3)",
			F: func(is I) {
				is.Err(&customErr{}, &customErr{}, &customErr{})
			},
		},
		{
			N: "Err(err1, err2, err3)",
			F: func(is I) {
				var err1 error
				var err2 error
				var err3 error
				is.Err(err1, err2, err3)
			},
			Fails: []string{"error expected"},
		},
		// OK
		{
			N: "OK(customErr(nil))",
			F: func(is I) {
				var err *customErr
				is.NoErr(err)
			},
		}, {
			N: "OK(func) panic",
			F: func(is I) {
				is.OK(func() {
					panic("panic message")
				})
			},
			Fails: []string{"unexpected panic: panic message"},
		}, {
			N: "OK(func) no panic",
			F: func(is I) {
				is.OK(func() {})
			},
		},
		// is.Panic
		{
			N: "PanicWith(\"panic message\", func(){ panic() })",
			F: func(is I) {
				is.PanicWith("panic message", func() {
					panic("panic message")
				})
			},
		},
		{
			N: "PanicWith(\"panic message\", func(){ /* no panic */ })",
			F: func(is I) {
				is.PanicWith("panic message", func() {
				})
			},
			Fails: []string{"expected panic: \"panic message\""},
		},
		{
			N: "Panic(func(){ panic() })",
			F: func(is I) {
				is.Panic(func() {
					panic("panic message")
				})
			},
		},
		{
			N: "Panic(func(){ /* no panic */ })",
			F: func(is I) {
				is.Panic(func() {
				})
			},
			Fails: []string{"expected panic"},
		},
		// is.Equal
		{
			N: "Equal(msi,msi) nil maps",
			F: func(is I) {
				var m1 map[string]interface{}
				var m2 map[string]interface{}
				is.Equal(m1, m2)
			},
		},
		{
			N: "Equal(1,1)",
			F: func(is I) {
				is.Equal(1, 1)
			},
		}, {
			N: "Equal(1,2)",
			F: func(is I) {
				is.Equal(1, 2)
			},
			Fails: []string{"1 != 2"},
		}, {
			N: "Equal(1,nil)",
			F: func(is I) {
				is.Equal(1, nil)
			},
			Fails: []string{"1 != <nil>"},
		}, {
			N: "Equal(nil,1)",
			F: func(is I) {
				is.Equal(nil, 1)
			},
			Fails: []string{"<nil> != 1"},
		},
		{
			N: "Equal(uint64(1),int64(1))",
			F: func(is I) {
				is.Equal(uint64(1), int64(1))
			},
		},
		{
			N: "Equal(false,false)",
			F: func(is I) {
				is.Equal(false, false)
			},
		}, {
			N: "Equal(map1,map2)",
			F: func(is I) {
				is.Equal(
					map[string]interface{}{"package": "is"},
					map[string]interface{}{"package": "is"},
				)
			},
		},

		// is.True
		{
			N: "True(true,true)",
			F: func(is I) {
				is.True(true, true)
			},
		}, {
			N: "True(false,false)",
			F: func(is I) {
				is.True(false, false)
			},
			Fails: []string{"true!=false"},
		},
		// is.False
		{
			N: "False(false,false)",
			F: func(is I) {
				is.False(false, false)
			},
		}, {
			N: "False(true,true)",
			F: func(is I) {
				is.False(true, true)
			},
			Fails: []string{"false!=true"},
		},

		// is.NotEqual
		{
			N: "NotEqual(1,2)",
			F: func(is I) {
				is.NotEqual(1, 2)
			},
		}, {
			N: "NotEqual(1,1)",
			F: func(is I) {
				is.NotEqual(1, 1)
			},
			Fails: []string{"1 == 1"},
		}, {
			N: "NotEqual(1,nil)",
			F: func(is I) {
				is.NotEqual(1, nil)
			},
		}, {
			N: "NotEqual(nil,1)",
			F: func(is I) {
				is.NotEqual(nil, 1)
			},
		}, {
			N: "NotEqual(false,false)",
			F: func(is I) {
				is.NotEqual(false, false)
			},
			Fails: []string{"false == false"},
		}, {
			N: "NotEqual(map1,map2)",
			F: func(is I) {
				is.NotEqual(
					map[string]interface{}{"package": "is"},
					map[string]interface{}{"package": "isn't"},
				)
			},
		}} {

		tt := new(mockT)
		is := New(tt)

		func() {
			defer func() {
				recover()
			}()
			test.F(is)
		}()

		if len(test.Fails) > 0 {
			for n, fail := range test.Fails {
				if !tt.Failed() {
					t.Errorf("%s should fail", test.N)
				}
				if test.Fails[n] != fail {
					t.Errorf("expected fail \"%s\" but was \"%s\".", test.Fails[n], fail)
				}
			}
		} else {
			if tt.Failed() {
				t.Errorf("%s shouldn't fail but: %s", test.N, strings.Join(test.Fails, ", "))
			}
		}

	}

}

func TestNewStrict(t *testing.T) {
	tt := new(mockT)
	is := Relaxed(tt)

	is.OK(nil)
	is.Equal(1, 2)
	is.NoErr(errors.New("nope"))

	if tt.Failed() {
		t.Error("Relaxed should not call FailNow")
	}

}
