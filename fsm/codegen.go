//go:generate -command mockgen go run ../vendor/go.uber.org/mock/mockgen
//go:generate -command gen-validate go run ../cmd/gen-validate
package fsm

//go:generate mockgen -source=state_storage_iface.go -destination=gomock/state_storage_iface.go -package=gomock -mock_names=StateStorage=StateStorage
//go:generate gen-validate --types=StateMachine --package=fsm --output=machine_validate.gen.go
