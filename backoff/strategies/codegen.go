//go:generate -command gen-validate go run ../../cmd/gen-validate
package strategies

//go:generate gen-validate --types=exponential --package=strategies --output=exponential_validate.gen.go
