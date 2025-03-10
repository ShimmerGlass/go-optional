# go-optional [![.github/workflows/check.yml](https://github.com/shimmerglass/go-optional/actions/workflows/check.yml/badge.svg)](https://github.com/shimmerglass/go-optional/actions/workflows/check.yml) [![GoDoc](https://godoc.org/github.com/shimmerglass/go-optional?status.svg)](https://godoc.org/github.com/shimmerglass/go-optional)

A library that provides [Go Generics](https://go.dev/blog/generics-proposal) friendly "optional" features.

## Fork of [moznion/go-optional](https://github.com/moznion/go-optional)

This package is a fork of [moznion/go-optional](https://github.com/moznion/go-optional) with some changes.

### Option are now structs, not slices

Original implementation uses a `[]T` to represent `Option`s. This induces a 2x `int` + 1x `uintptr` (usually 24 bytes) memory overhead, plus a pointer dereference. This new implementation uses a struct with a boolean to store whether the `Option` is set or not.

As a result:

- Memory overhead reduced from 24 bytes to 1 byte
- No heap allocations
- `Option`s are comparable if the type they wrap is comparable and can be used as `map[]` keys

However:

- `nil` is no longer a valid `Option` value. `opt.None()` must be used instead. An `Option` default value is still `None()`
- The JSON tag `omitempty` no longer works on `Option`s

### Map*, FlatMap*, Zip*, Unzip* functions removed

For simplicity, and because they made less sense now that `Option`s are no longer slices.

## Synopsis

```go
some := opt.Some[int](123)
fmt.Printf("%v\n", some.IsSome()) // => true
fmt.Printf("%v\n", some.IsNone()) // => false

v, err := some.Take()
fmt.Printf("err is nil: %v\n", err == nil) // => err is nil: true
fmt.Printf("%d\n", v) // => 123
```

and more detailed examples are here: [./examples_test.go](./examples_test.go).

## Docs

[![GoDoc](https://godoc.org/github.com/shimmerglass/go-optional?status.svg)](https://godoc.org/github.com/shimmerglass/go-optional)

### Supported Operations

#### Value Factory Methods

- [Some[T]\() Option[T]](https://pkg.go.dev/github.com/shimmerglass/go-optional#Some)
- [None[T]\() Option[T]](https://pkg.go.dev/github.com/shimmerglass/go-optional#None)
- [FromNillable[T]\() Option[T]](https://pkg.go.dev/github.com/shimmerglass/go-optional#FromNillable)
- [PtrFromNillable[T]\() Option[T]](https://pkg.go.dev/github.com/shimmerglass/go-optional#PtrFromNillable)

#### Option value handler methods

- [Option[T]#IsNone() bool](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.IsNone)
- [Option[T]#IsSome() bool](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.IsSome)
- [Option[T]#Unwrap() T](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.Unwrap)
- [Option[T]#UnwrapAsPtr() \*T](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.UnwrapAsPtr)
- [Option[T]#Take() (T, error)](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.Take)
- [Option[T]#TakeOr(fallbackValue T) T](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.TakeOr)
- [Option[T]#TakeOrElse(fallbackFunc func() T) T](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.TakeOrElse)
- [Option[T]#Or(fallbackOptionValue Option[T]) Option[T]](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.Or)
- [Option[T]#OrElse(fallbackOptionFunc func() Option[T]) Option[T]](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.OrElse)
- [Option[T]#Filter(predicate func(v T) bool) Option[T]](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.Filter)
- [Option[T]#IfSome(f func(v T))](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.IfSome)
- [Option[T]#IfSomeWithError(f func(v T) error) error](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.IfSomeWithError)
- [Option[T]#IfNone(f func())](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.IfNone)
- [Option[T]#IfNoneWithError(f func() error) error](https://pkg.go.dev/github.com/shimmerglass/go-optional#Option.IfNoneWithError)

### JSON marshal/unmarshal support

This `Option[T]` type supports JSON marshal and unmarshal.

If the value wanted to marshal is `Some[T]` then it marshals that value into the JSON bytes simply, and in unmarshaling, if the given JSON string/bytes has the actual value on corresponded property, it unmarshals that value into `Some[T]` value.

example:

```go
type JSONStruct struct {
	Val opt.Option[int] `json:"val"`
}

some := opt.Some[int](123)
jsonStruct := &JSONStruct{Val: some}

marshal, err := json.Marshal(jsonStruct)
if err != nil {
	return err
}
fmt.Printf("%s\n", marshal) // => {"val":123}

var unmarshalJSONStruct JSONStruct
err = json.Unmarshal(marshal, &unmarshalJSONStruct)
if err != nil {
	return err
}
// unmarshalJSONStruct.Val == Some[int](123)
```

Elsewise, when the value is `None[T]`, the marshaller serializes that value as `null`. And if the unmarshaller gets the JSON `null` value on a property corresponding to the `Optional[T]` value, or the value of a property is missing, that deserializes that value as `None[T]`.

example:

```go
type JSONStruct struct {
	Val opt.Option[int] `json:"val"`
}

none := opt.None[int]()
jsonStruct := &JSONStruct{Val: none}

marshal, err := json.Marshal(jsonStruct)
if err != nil {
	return err
}
fmt.Printf("%s\n", marshal) // => {"val":null}

var unmarshalJSONStruct JSONStruct
err = json.Unmarshal(marshal, &unmarshalJSONStruct)
if err != nil {
	return err
}
// unmarshalJSONStruct.Val == None[int]()
```

### SQL Driver Support

`Option[T]` satisfies [sql/driver.Valuer](https://pkg.go.dev/database/sql/driver#Valuer) and [sql.Scanner](https://pkg.go.dev/database/sql#Scanner), so this type can be used by SQL interface on Golang.

example of the primitive usage:

```go
sqlStmt := "CREATE TABLE tbl (id INTEGER NOT NULL PRIMARY KEY, name VARCHAR(32));"
db.Exec(sqlStmt)

tx, _ := db.Begin()
func() {
    stmt, _ := tx.Prepare("INSERT INTO tbl(id, name) values(?, ?)")
    defer stmt.Close()
    stmt.Exec(1, "foo")
}()
func() {
    stmt, _ := tx.Prepare("INSERT INTO tbl(id) values(?)")
    defer stmt.Close()
    stmt.Exec(2) // name is NULL
}()
tx.Commit()

var maybeName opt.Option[string]

row := db.QueryRow("SELECT name FROM tbl WHERE id = 1")
row.Scan(&maybeName)
fmt.Println(maybeName) // Some[foo]

row := db.QueryRow("SELECT name FROM tbl WHERE id = 2")
row.Scan(&maybeName)
fmt.Println(maybeName) // None[]
```
