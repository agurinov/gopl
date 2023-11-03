package fsm_test

import (
	"github.com/golang/mock/gomock"

	mock "github.com/agurinov/gopl/fsm/gomock"
)

type mocks struct {
	storage *mock.StateStorage[RegistrationContext]
}

func NewMocks(ctrl *gomock.Controller) mocks {
	return mocks{
		storage: mock.NewStateStorage[RegistrationContext](ctrl),
	}
}
