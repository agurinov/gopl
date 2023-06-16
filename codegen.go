//go:build neverbuild

//go:generate -command stringer go run ./vendor/golang.org/x/tools/cmd/stringer
//go:generate -command mockgen go run ./vendor/github.com/golang/mock/mockgen

//go:generate -command protoc_gen_types    protoc      --go_out=paths=import:.
//go:generate -command protoc_gen_grpc     protoc --go-grpc_out=paths=import,require_unimplemented_servers=false:.
//go:generate -command protoc_gen_protoset protoc --include_imports --include_source_info --descriptor_set_out
package main

//go:generate mockgen -source=fsm/state_storage_iface.go -destination=fsm/mock/state_storage_iface.go -package=mock -mock_names=StateStorage=StateStorage
