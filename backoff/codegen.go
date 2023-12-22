//go:generate -command gen-validate go run ../cmd/gen-validate
package backoff

//go:generate gen-validate --types=Backoff --package=backoff --output=backoff_validate.gen.go
