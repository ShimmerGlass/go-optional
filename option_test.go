package opt

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption_IsNone(t *testing.T) {
	assert.True(t, None[int]().IsNone())
	assert.False(t, Some[int](123).IsNone())

	var nilValue Option[int]
	assert.True(t, nilValue.IsNone())

	i := 0
	assert.False(t, FromNillable[int](&i).IsNone())
	assert.True(t, FromNillable[int](nil).IsNone())
}

func TestOption_IsSome(t *testing.T) {
	assert.False(t, None[int]().IsSome())
	assert.True(t, Some[int](123).IsSome())

	var nilValue Option[int]
	assert.False(t, nilValue.IsSome())

	i := 0
	assert.True(t, FromNillable[int](&i).IsSome())
	assert.False(t, FromNillable[int](nil).IsSome())
}

func TestOption_Unwrap(t *testing.T) {
	assert.Equal(t, "foo", Some[string]("foo").Unwrap())
	assert.Equal(t, "", None[string]().Unwrap())
	assert.Nil(t, None[*string]().Unwrap())

	i := 123
	assert.Equal(t, i, FromNillable[int](&i).Unwrap())
	assert.Equal(t, 0, FromNillable[int](nil).Unwrap())
	assert.Equal(t, i, *PtrFromNillable[int](&i).Unwrap())
	assert.Nil(t, PtrFromNillable[int](nil).Unwrap())
}

func TestOption_UnwrapAsPointer(t *testing.T) {
	str := "foo"
	refStr := &str
	assert.EqualValues(t, &str, Some[string](str).UnwrapAsPtr())
	assert.EqualValues(t, &refStr, Some[*string](refStr).UnwrapAsPtr())
	assert.Nil(t, None[string]().UnwrapAsPtr())
	assert.Nil(t, None[*string]().UnwrapAsPtr())

	i := 123
	assert.Equal(t, &i, FromNillable[int](&i).UnwrapAsPtr())
	assert.Nil(t, FromNillable[int](nil).UnwrapAsPtr())
	assert.Equal(t, &i, *PtrFromNillable[int](&i).UnwrapAsPtr())
	assert.Nil(t, PtrFromNillable[int](nil).UnwrapAsPtr())
}

func TestOption_Take(t *testing.T) {
	v, err := Some[int](123).Take()
	assert.NoError(t, err)
	assert.Equal(t, 123, v)

	v, err = None[int]().Take()
	assert.ErrorIs(t, err, ErrNoneValueTaken)
	assert.Equal(t, 0, v)
}

func TestOption_TakeOr(t *testing.T) {
	v := Some[int](123).TakeOr(666)
	assert.Equal(t, 123, v)

	v = None[int]().TakeOr(666)
	assert.Equal(t, 666, v)
}

func TestOption_TakeOrElse(t *testing.T) {
	v := Some[int](123).TakeOrElse(func() int {
		return 666
	})
	assert.Equal(t, 123, v)

	v = None[int]().TakeOrElse(func() int {
		return 666
	})
	assert.Equal(t, 666, v)
}

func TestOption_IfSome(t *testing.T) {
	callingValue := ""
	Some("foo").IfSome(func(s string) {
		callingValue = s
	})
	assert.Equal(t, "foo", callingValue)

	callingValue = ""
	None[string]().IfSome(func(s string) {
		callingValue = s
	})
	assert.Equal(t, "", callingValue)
}

func TestOption_IfSomeWithError(t *testing.T) {
	err := Some("foo").IfSomeWithError(func(s string) error {
		return nil
	})
	assert.NoError(t, err)

	err = Some("foo").IfSomeWithError(func(s string) error {
		return errors.New(s)
	})
	assert.EqualError(t, err, "foo")

	err = None[string]().IfSomeWithError(func(s string) error {
		return errors.New(s)
	})
	assert.NoError(t, err)
}

func TestOption_IfNone(t *testing.T) {
	called := false
	None[string]().IfNone(func() {
		called = true
	})
	assert.True(t, called)

	called = false
	Some("string").IfNone(func() {
		called = true
	})
	assert.False(t, called)
}

func TestOption_IfNoneWithError(t *testing.T) {
	err := None[string]().IfNoneWithError(func() error {
		return nil
	})
	assert.NoError(t, err)

	err = None[string]().IfNoneWithError(func() error {
		return errors.New("err")
	})
	assert.EqualError(t, err, "err")

	err = Some("foo").IfNoneWithError(func() error {
		return errors.New("err")
	})
	assert.NoError(t, err)
}

type MyStringer struct {
}

func (m *MyStringer) String() string {
	return "mystr"
}

func TestOption_String(t *testing.T) {
	assert.Equal(t, "Some[123]", Some[int](123).String())
	assert.Equal(t, "None[]", None[int]().String())

	assert.Equal(t, "Some[mystr]", Some[*MyStringer](&MyStringer{}).String())
	assert.Equal(t, "None[]", None[*MyStringer]().String())
}

func TestOption_Or(t *testing.T) {
	fallback := Some[string]("fallback")

	assert.EqualValues(t, Some[string]("actual").Or(fallback).Unwrap(), "actual")
	assert.EqualValues(t, None[string]().Or(fallback).Unwrap(), "fallback")
}

func TestOption_OrElse(t *testing.T) {
	fallbackFunc := func() Option[string] { return Some[string]("fallback") }

	assert.EqualValues(t, Some[string]("actual").OrElse(fallbackFunc).Unwrap(), "actual")
	assert.EqualValues(t, None[string]().OrElse(fallbackFunc).Unwrap(), "fallback")
}

func TestOption_Comparable(t *testing.T) {
	_ = map[Option[int]]struct{}{}
}
