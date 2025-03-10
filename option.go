package opt

import (
	"errors"
	"fmt"
)

var (
	// ErrNoneValueTaken represents the error that is raised when None value is taken.
	ErrNoneValueTaken = errors.New("none value taken")
)

// Option is a data type that must be Some (i.e. having a value) or None (i.e. doesn't have a value).
// This type implements database/sql/driver.Valuer and database/sql.Scanner.
type Option[T any] struct {
	value  T
	isSome bool
}

// Some is a function to make an Option type value with the actual value.
func Some[T any](v T) Option[T] {
	return Option[T]{
		value:  v,
		isSome: true,
	}
}

// None is a function to make an Option type value that doesn't have a value.
func None[T any]() Option[T] {
	return Option[T]{
		isSome: false,
	}
}

// FromNillable is a function to make an Option type value with the nillable value with value de-referencing.
// If the given value is not nil, this returns Some[T] value. On the other hand, if the value is nil, this returns None[T].
// This function does "dereference" for the value on packing that into Option value. If this value is not preferable, please consider using PtrFromNillable() instead.
func FromNillable[T any](v *T) Option[T] {
	if v == nil {
		return None[T]()
	}
	return Some[T](*v)
}

// PtrFromNillable is a function to make an Option type value with the nillable value without value de-referencing.
// If the given value is not nil, this returns Some[*T] value. On the other hand, if the value is nil, this returns None[*T].
// This function doesn't "dereference" the value on packing that into the Option value; in other words, this puts the as-is pointer value into the Option envelope.
// This behavior contrasts with the FromNillable() function's one.
func PtrFromNillable[T any](v *T) Option[*T] {
	if v == nil {
		return None[*T]()
	}
	return Some[*T](v)
}

// IsNone returns whether the Option *doesn't* have a value or not.
func (o Option[T]) IsNone() bool {
	return !o.isSome
}

// IsSome returns whether the Option has a value or not.
func (o Option[T]) IsSome() bool {
	return o.isSome
}

// Unwrap returns the value regardless of Some/None status.
// If the Option value is Some, this method returns the actual value.
// On the other hand, if the Option value is None, this method returns the *default* value according to the type.
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		var defaultValue T
		return defaultValue
	}
	return o.value
}

// UnwrapAsPtr returns the contained value in receiver Option as a pointer.
// This is similar to `Unwrap()` method but the difference is this method returns a pointer value instead of the actual value.
// If the receiver Option value is None, this method returns nil.
func (o Option[T]) UnwrapAsPtr() *T {
	if o.IsNone() {
		return nil
	}
	return &o.value
}

// Take takes the contained value in Option.
// If Option value is Some, this returns the value that is contained in Option.
// On the other hand, this returns an ErrNoneValueTaken as the second return value.
func (o Option[T]) Take() (T, error) {
	if o.IsNone() {
		var defaultValue T
		return defaultValue, ErrNoneValueTaken
	}
	return o.value, nil
}

// TakeOr returns the actual value if the Option has a value.
// On the other hand, this returns fallbackValue.
func (o Option[T]) TakeOr(fallbackValue T) T {
	if o.IsNone() {
		return fallbackValue
	}
	return o.value
}

// TakeOrElse returns the actual value if the Option has a value.
// On the other hand, this executes fallbackFunc and returns the result value of that function.
func (o Option[T]) TakeOrElse(fallbackFunc func() T) T {
	if o.IsNone() {
		return fallbackFunc()
	}
	return o.value
}

// Or returns the Option value according to the actual value existence.
// If the receiver's Option value is Some, this function pass-through that to return. Otherwise, this value returns the `fallbackOptionValue`.
func (o Option[T]) Or(fallbackOptionValue Option[T]) Option[T] {
	if o.IsNone() {
		return fallbackOptionValue
	}
	return o
}

// OrElse returns the Option value according to the actual value existence.
// If the receiver's Option value is Some, this function pass-through that to return. Otherwise, this executes `fallbackOptionFunc` and returns the result value of that function.
func (o Option[T]) OrElse(fallbackOptionFunc func() Option[T]) Option[T] {
	if o.IsNone() {
		return fallbackOptionFunc()
	}
	return o
}

// IfSome calls given function with the value of Option if the receiver value is Some.
func (o Option[T]) IfSome(f func(v T)) {
	if o.IsNone() {
		return
	}
	f(o.value)
}

// IfSomeWithError calls given function with the value of Option if the receiver value is Some.
// This method propagates the error of given function, and if the receiver value is None, this returns nil error.
func (o Option[T]) IfSomeWithError(f func(v T) error) error {
	if o.IsNone() {
		return nil
	}
	return f(o.value)
}

// IfNone calls given function if the receiver value is None.
func (o Option[T]) IfNone(f func()) {
	if o.IsSome() {
		return
	}
	f()
}

// IfNoneWithError calls given function if the receiver value is None.
// This method propagates the error of given function, and if the receiver value is Some, this returns nil error.
func (o Option[T]) IfNoneWithError(f func() error) error {
	if o.IsSome() {
		return nil
	}
	return f()
}

func (o Option[T]) String() string {
	if o.IsNone() {
		return "None[]"
	}

	v := o.Unwrap()
	if stringer, ok := interface{}(v).(fmt.Stringer); ok {
		return fmt.Sprintf("Some[%s]", stringer)
	}
	return fmt.Sprintf("Some[%v]", v)
}
