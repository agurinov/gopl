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
	Worker struct {
		w io.ReadWriter
		s string
		i int
	}
	WorkerOption            c.Option[Worker]
	WorkerOptionWithContext = c.OptionWithContext[Worker]
)

var (
	New                    = c.New[Worker, WorkerOption]
	NewWithValidate        = c.NewWithValidate[Worker, WorkerOption]
	NewWithContext         = c.NewWithContext[Worker, WorkerOptionWithContext]
	NewWithContextValidate = c.NewWithContextValidate[Worker, WorkerOptionWithContext]
)

var (
	validW io.ReadWriter = new(bytes.Buffer)
	validS               = "foobar"
	validI               = 100500
)

var (
	invalidW io.ReadWriter = nil
	invalidS               = "lolkek"
	invalidI               = 999
)

func (obj Worker) String() string {
	var b strings.Builder

	b.WriteString("s=" + obj.s + ";")
	b.WriteString("i=" + strconv.Itoa(obj.i) + ";")
	b.WriteString("w_is_nil=" + strconv.FormatBool(obj.w == nil) + ";")
	b.WriteString("is_valid=" + strconv.FormatBool(obj.Validate() == nil) + ";")

	return b.String()
}

func (obj Worker) Validate() error {
	s := struct {
		W io.ReadWriter `validate:"required"`
		S string        `validate:"eq=foobar"`
		I int           `validate:"eq=100500"`
	}{
		W: obj.w,
		S: obj.s,
		I: obj.i,
	}

	if err := validator.New().Struct(s); err != nil {
		return err
	}

	return nil
}

// ExampleNew just configures object - no validate applied
func ExampleNew() {
	opts := []WorkerOption{
		WithW(invalidW), nil,
		WithS(invalidS), nil,
		WithI(invalidI), nil,
	}

	obj, err := New(opts...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", obj)
	// Output:
	// s=lolkek;i=999;w_is_nil=true;is_valid=false;
}

// ExampleNewWithValidate configures and validates object
func ExampleNewWithValidate() {
	opts := []WorkerOption{
		WithW(validW), nil,
		WithS(validS), nil,
		WithI(validI), nil,
	}

	obj, err := NewWithValidate(opts...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", obj)
	// Output:
	// s=foobar;i=100500;w_is_nil=false;is_valid=true;
}

// ExampleNewWithContext just configures object - no validate applied. Also ctx is required
// For example - http.Request as struct attribute
func ExampleNewWithContext() {
	ctx := context.TODO()

	opts := []WorkerOptionWithContext{
		WithWCtx(invalidW), nil,
		WithSCtx(invalidS), nil,
		WithICtx(invalidI), nil,
	}

	obj, err := NewWithContext(ctx, opts...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", obj)
	// Output:
	// s=lolkek;i=999;w_is_nil=true;is_valid=false;
}

// ExampleNewWithContextValidate configures and validates object. Also ctx is required
// For example - http.Request as struct attribute
func ExampleNewWithContextValidate() {
	ctx := context.TODO()

	opts := []WorkerOptionWithContext{
		WithWCtx(validW), nil,
		WithSCtx(validS), nil,
		WithICtx(validI), nil,
	}

	obj, err := NewWithContextValidate(ctx, opts...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", obj)
	// Output:
	// s=foobar;i=100500;w_is_nil=false;is_valid=true;
}

func WithI(i int) WorkerOption {
	return func(o *Worker) error {
		o.i = i

		return nil
	}
}

func WithS(s string) WorkerOption {
	return func(o *Worker) error {
		o.s = s

		return nil
	}
}

func WithW(w io.ReadWriter) WorkerOption {
	return func(o *Worker) error {
		o.w = w

		return nil
	}
}

func WithErr(err error) WorkerOption {
	return func(_ *Worker) error {
		return err
	}
}

func WithICtx(i int) WorkerOptionWithContext {
	return func(_ context.Context, o *Worker) error {
		o.i = i

		return nil
	}
}

func WithSCtx(s string) WorkerOptionWithContext {
	return func(_ context.Context, o *Worker) error {
		o.s = s

		return nil
	}
}

func WithWCtx(w io.ReadWriter) WorkerOptionWithContext {
	return func(_ context.Context, o *Worker) error {
		o.w = w

		return nil
	}
}

func WithErrCtx(err error) WorkerOptionWithContext {
	return func(_ context.Context, o *Worker) error {
		return err
	}
}
