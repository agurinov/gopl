package fsm_test

import (
	"github.com/golang/mock/gomock"

	"github.com/agurinov/gopl/fsm/mock"
)

type mocks struct {
	storage *mock.StateStorage[RegistrationContext]
}

func NewMocks(ctrl *gomock.Controller) mocks {
	return mocks{
		storage: mock.NewStateStorage[RegistrationContext](ctrl),
	}
}
