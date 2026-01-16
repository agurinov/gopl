//go:generate -command stringer go run ./vendor/golang.org/x/tools/cmd/stringer
//go:generate -command mockery  go run ./vendor/github.com/vektra/mockery/v3

//go:generate -command protoc_gen_types    protoc      --go_out=paths=import:.
//go:generate -command protoc_gen_grpc     protoc --go-grpc_out=paths=import,require_unimplemented_servers=false:.
//go:generate -command protoc_gen_protoset protoc --include_imports --include_source_info --descriptor_set_out

//go:generate -command validate-gen go run ./cmd/validate-gen
package main

//go:generate mockery
