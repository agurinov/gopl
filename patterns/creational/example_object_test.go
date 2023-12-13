package creational_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	MyStruct struct {
		w io.ReadWriter
		s string
		i int
	}
	MyStructOption            c.Option[MyStruct]
	MyStructOptionWithContext = c.OptionWithContext[MyStruct]
)

var (
	MustNew        = c.MustNew[MyStruct, MyStructOption]
	NewWithContext = c.NewWithContext[MyStruct, MyStructOptionWithContext]
	NewExtended    = c.NewExtended[MyStruct, MyStructOptionWithContext]
)

// ExampleMustNew describes basic flow - just construct object.
// Validation and initialization are on responsibility of an engineer.
func ExampleMustNew() {
	var (
		ctx = context.TODO()
		w   bytes.Buffer
	)

	opts := []MyStructOption{
		WithW(&w),
		WithS("foobar"),
		WithI(100),
	}
	obj := MustNew(opts...)

	fmt.Printf("Is valid: %t\n", obj.Validate() == nil)
	fmt.Printf("Is inited: %t\n", obj.Init(ctx) == nil)
	fmt.Printf("%s\n", obj)
	// Output:
	// Is valid: true
	// Is inited: true
	// foobar100
}

// ExampleNewWithContext describes flow where attributes requires ctx to be created.
// For example - http.Request as attribute
func ExampleNewWithContext() {
	var (
		ctx = context.TODO()
		w   bytes.Buffer
	)

	opts := []MyStructOptionWithContext{
		WithWCtx(&w),
		WithSCtx("foobar"),
		WithICtx(100),
	}

	obj, err := NewWithContext(ctx, opts...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Is valid: %t\n", obj.Validate() == nil)
	fmt.Printf("Is inited: %t\n", obj.Init(ctx) == nil)
	fmt.Printf("%s\n", obj)
	// Output:
	// Is valid: true
	// Is inited: true
	// foobar100
}

// ExampleNewExtended describes flow guarantees that .Validate() and .Init() method succeed
// Methods must be present on generic type
func ExampleNewExtended() {
	var (
		ctx = context.TODO()
		w   bytes.Buffer
	)

	opts := []MyStructOptionWithContext{
		WithWCtx(&w),
		WithSCtx("foobar"),
		WithICtx(100),
	}

	obj, err := NewExtended(ctx, opts...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", obj)
	// Output:
	// foobar100
}

func (o MyStruct) String() string {
	var b strings.Builder

	if _, err := io.Copy(&b, o.w); err != nil {
		return err.Error()
	}

	return b.String()
}

func (o MyStruct) Validate() error {
	s := struct {
		W io.ReadWriter `validate:"required"`
		S string        `validate:"required"`
		I int           `validate:"required"`
	}{
		W: o.w,
		S: o.s,
		I: o.i,
	}

	if err := validator.New().Struct(s); err != nil {
		return err
	}

	return nil
}

func (o MyStruct) Init(_ context.Context) error {
	if o.i == 100500 {
		return io.EOF
	}

	if _, err := io.WriteString(o.w, o.s); err != nil {
		return err
	}

	if _, err := io.WriteString(o.w, strconv.Itoa(o.i)); err != nil {
		return err
	}

	return nil
}

func WithI(i int) MyStructOption {
	return func(o *MyStruct) error {
		o.i = i

		return nil
	}
}

func WithICtx(i int) MyStructOptionWithContext {
	return func(_ context.Context, o *MyStruct) error {
		o.i = i

		return nil
	}
}

func WithS(s string) MyStructOption {
	return func(o *MyStruct) error {
		o.s = s

		return nil
	}
}

func WithSCtx(s string) MyStructOptionWithContext {
	return func(_ context.Context, o *MyStruct) error {
		o.s = s

		return nil
	}
}

func WithW(w io.ReadWriter) MyStructOption {
	return func(o *MyStruct) error {
		o.w = w

		return nil
	}
}

func WithWCtx(w io.ReadWriter) MyStructOptionWithContext {
	return func(_ context.Context, o *MyStruct) error {
		o.w = w

		return nil
	}
}

func WithErr(err error) MyStructOption {
	return func(_ *MyStruct) error {
		return err
	}
}

func WithErrCtx(err error) MyStructOptionWithContext {
	return func(_ context.Context, o *MyStruct) error {
		return err
	}
}
