//go:build neverbuild

//go:generate -command mockgen go run ../vendor/github.com/golang/mock/mockgen
package fsm

//go:generate mockgen -source=state_storage_iface.go -destination=mock/state_storage_iface.go -package=mock -mock_names=StateStorage=StateStorage
