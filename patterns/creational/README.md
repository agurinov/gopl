# Creational toolset

Functional options generic toolset to get rid of bloilerplate constructor functions.

```go
// component.go

import (
	"go.uber.org/zap"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	MyComponent struct {
		logger *zap.Logger
		s      string
		i      int
	}

	// Or it can alias as: MyComponentOption = c.Option[Closer]
	MyComponentOption c.Option[Closer]


	// Or it can alias as: MyCtxComponentOption = c.OptionWithContext[Closer]
	MyCtxComponentOption c.OptionWithContext[Closer]
)

var (
	// Creating new oject with values from provided options.
	New = c.New[MyComponent, MyComponentOption]

	// Creating and validating resulting object.
	// Object must have .Valdate() method.
	New = c.NewWithValidate[MyComponent, MyComponentOption]

	// Creating new oject with available context to create some socket.
	New = c.NewWithContext[MyComponent, MyCtxComponentOption]

	// Creating and validating new oject with available context to create some socket.
	// Object must have .Valdate() method.
	New = c.NewWithContextValidate[MyComponent, MyCtxComponentOption]
)

// Also there are pair of methods Construct* which can patch already existing object.
// Useful in case of known default state.
// Feel free to combine ctx and validate mechanics as you want.
func New(opts ...MyComponentOption) (*MyComponent, error) {
	c := MyComponent{
		s: "foobar",
		i: 100500,
	}

	obj, err := c.ConstructWithValidate(c, opts...)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
```

```go
// component_f_opt.go

import (
	"go.uber.org/zap"
)

func WithLogger(logger *zap.Logger) MyComponentOption {
	return func(c *MyComponent) error {
		if logger == nil {
			return nil
		}

		c.logger = logger.Named("my.component")

		return nil
	}
}

func WithContextLogger(logger *zap.Logger) MyComponentOption {
	return func(ctx context.Context, c *MyComponent) error {
		if logger == nil {
			return nil
		}

		c.logger = logger.Named("my.component")

		return nil
	}
}
```
