//go:generate -command mockgen go run ../vendor/go.uber.org/mock/mockgen
package fsm

//go:generate mockgen -source=state_storage_iface.go -destination=gomock/state_storage_iface.go -package=gomock -mock_names=StateStorage=StateStorage
