package opt

import (
	"errors"
	"fmt"
)

func ExampleOption_IsNone() {
	some := Some[int](1)
	fmt.Printf("%v\n", some.IsNone())
	none := None[int]()
	fmt.Printf("%v\n", none.IsNone())

	num := 123
	some = FromNillable[int](&num)
	fmt.Printf("%v\n", some.IsNone())
	none = FromNillable[int](nil)
	fmt.Printf("%v\n", none.IsNone())

	ptrSome := PtrFromNillable[int](&num)
	fmt.Printf("%v\n", ptrSome.IsNone())
	ptrNone := PtrFromNillable[int](nil)
	fmt.Printf("%v\n", ptrNone.IsNone())

	var nilValue Option[int]
	fmt.Printf("%v\n", nilValue.IsNone())

	// Output:
	// false
	// true
	// false
	// true
	// false
	// true
	// true
}

func ExampleOption_IsSome() {
	some := Some[int](1)
	fmt.Printf("%v\n", some.IsSome())
	none := None[int]()
	fmt.Printf("%v\n", none.IsSome())

	num := 123
	some = FromNillable[int](&num)
	fmt.Printf("%v\n", some.IsSome())
	none = FromNillable[int](nil)
	fmt.Printf("%v\n", none.IsSome())

	ptrSome := PtrFromNillable[int](&num)
	fmt.Printf("%v\n", ptrSome.IsSome())
	ptrNone := PtrFromNillable[int](nil)
	fmt.Printf("%v\n", ptrNone.IsSome())

	var nilValue Option[int]
	fmt.Printf("%v\n", nilValue.IsSome())

	// Output:
	// true
	// false
	// true
	// false
	// true
	// false
	// false
}

func ExampleOption_Unwrap() {
	fmt.Printf("%v\n", Some[int](12345).Unwrap())
	fmt.Printf("%v\n", None[int]().Unwrap())
	fmt.Printf("%v\n", None[*int]().Unwrap())

	num := 123
	fmt.Printf("%v\n", FromNillable[int](&num).Unwrap())
	fmt.Printf("%v\n", FromNillable[int](nil).Unwrap())
	fmt.Printf("%v\n", *PtrFromNillable[int](&num).Unwrap()) // NOTE: this dereferences tha unwrapped value
	fmt.Printf("%v\n", PtrFromNillable[int](nil).Unwrap())
	// Output:
	// 12345
	// 0
	// <nil>
	// 123
	// 0
	// 123
	// <nil>
}

func ExampleOption_UnwrapAsPtr() {
	fmt.Printf("%v\n", *Some[int](12345).UnwrapAsPtr())
	fmt.Printf("%v\n", None[int]().UnwrapAsPtr())
	fmt.Printf("%v\n", None[*int]().UnwrapAsPtr())

	num := 123
	fmt.Printf("%v\n", *FromNillable[int](&num).UnwrapAsPtr())
	fmt.Printf("%v\n", FromNillable[int](nil).UnwrapAsPtr())
	fmt.Printf("%v\n", **PtrFromNillable[int](&num).UnwrapAsPtr()) // NOTE: this dereferences tha unwrapped value
	fmt.Printf("%v\n", PtrFromNillable[int](nil).UnwrapAsPtr())
	// Output:
	// 12345
	// <nil>
	// <nil>
	// 123
	// <nil>
	// 123
	// <nil>
}

func ExampleOption_Take() {
	some := Some[int](1)
	v, err := some.Take()
	fmt.Printf("%d\n", v)
	fmt.Printf("%v\n", err == nil)

	none := None[int]()
	_, err = none.Take()
	fmt.Printf("%v\n", err == nil)

	// Output:
	// 1
	// true
	// false
}

func ExampleOption_TakeOr() {
	some := Some[int](1)
	v := some.TakeOr(666)
	fmt.Printf("%d\n", v)

	none := None[int]()
	v = none.TakeOr(666)
	fmt.Printf("%d\n", v)

	// Output:
	// 1
	// 666
}

func ExampleOption_TakeOrElse() {
	some := Some[int](1)
	v := some.TakeOrElse(func() int {
		return 666
	})
	fmt.Printf("%d\n", v)

	none := None[int]()
	v = none.TakeOrElse(func() int {
		return 666
	})
	fmt.Printf("%d\n", v)

	// Output:
	// 1
	// 666
}

func ExampleOption_IfSome() {
	Some("foo").IfSome(func(val string) {
		fmt.Println(val)
	})

	None[string]().IfSome(func(val string) {
		fmt.Println("do not show this message")
	})

	// Output:
	// foo
}

func ExampleOption_IfSomeWithError() {
	err := Some("foo").IfSomeWithError(func(val string) error {
		fmt.Println(val)
		return nil
	})
	if err != nil {
		fmt.Println(err) // no error
	}

	err = Some("bar").IfSomeWithError(func(val string) error {
		fmt.Println(val)
		return errors.New("^^^ error occurred")
	})
	if err != nil {
		fmt.Println(err)
	}

	err = None[string]().IfSomeWithError(func(val string) error {
		return errors.New("do not show this error")
	})
	if err != nil {
		fmt.Println(err) // must not show this error
	}

	// Output:
	// foo
	// bar
	// ^^^ error occurred
}

func ExampleOption_IfNone() {
	None[string]().IfNone(func() {
		fmt.Println("value is none")
	})

	Some("foo").IfNone(func() {
		fmt.Println("do not show this message")
	})

	// Output:
	// value is none
}

func ExampleOption_IfNoneWithError() {
	err := None[string]().IfNoneWithError(func() error {
		fmt.Println("value is none")
		return nil
	})
	if err != nil {
		fmt.Println(err) // no error
	}

	err = None[string]().IfNoneWithError(func() error {
		fmt.Println("value is none!!")
		return errors.New("^^^ error occurred")
	})
	if err != nil {
		fmt.Println(err)
	}

	err = Some("foo").IfNoneWithError(func() error {
		return errors.New("do not show this error")
	})
	if err != nil {
		fmt.Println(err) // must not show this error
	}

	// Output:
	// value is none
	// value is none!!
	// ^^^ error occurred
}

func ExampleOption_Or() {
	fallback := Some[string]("fallback")

	some := Some[string]("actual")
	fmt.Printf("%s\n", some.Or(fallback))

	none := None[string]()
	fmt.Printf("%s\n", none.Or(fallback))

	// Output:
	// Some[actual]
	// Some[fallback]
}

func ExampleOption_OrElse() {
	fallbackFunc := func() Option[string] { return Some[string]("fallback") }

	some := Some[string]("actual")
	fmt.Printf("%s\n", some.OrElse(fallbackFunc))

	none := None[string]()
	fmt.Printf("%s\n", none.OrElse(fallbackFunc))

	// Output:
	// Some[actual]
	// Some[fallback]
}
