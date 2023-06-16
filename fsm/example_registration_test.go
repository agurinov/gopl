//nolint:testableexamples
package fsm_test

import (
	"context"

	"github.com/google/uuid"

	"github.com/agurinov/gopl/fsm"
)

type (
	RegistrationContext struct {
		UUID uuid.UUID
	}
	RegistrationStateMachine       = fsm.StateMachine[RegistrationContext]
	RegistrationStateMachineOption = fsm.StateMachineOption[RegistrationContext]
	RegistrationEvent              = fsm.Event[RegistrationContext]
)

type ChooseCountryContract struct{ country string }

func (c ChooseCountryContract) transition(ctx context.Context) (fsm.State, error) {
	// Push contract to db
	// db.SaveClientCountry(c.country)
	switch c.country {
	case "ru":
		return uploadPassportState, nil
	default:
		return brokenState, nil
	}
}

type PassportPhotoContract struct{ blob []byte }

func (c PassportPhotoContract) transition(ctx context.Context) (fsm.State, error) {
	// Push contract to reviewer
	// reviewer.UploadPassportPhoto(c.blog)
	return uploadSelfieState, nil
}

type SelfiePhotoContract struct{ blob []byte }

func (c SelfiePhotoContract) transition(ctx context.Context) (fsm.State, error) {
	// Push contract to reviewer
	// reviewer.UploadSelfiePhoto(c.blog)
	return uploadDriverLicenseState, nil
}

type DriverLicensePhotoContract struct{ blob []byte }

func (c DriverLicensePhotoContract) transition(ctx context.Context) (fsm.State, error) {
	// Push contract to reviewer
	// reviewer.UploadDriverLicensePhoto(c.blog)
	return reviewState, nil
}

type ReviewResponseContract struct {
	passportPhotoValid bool
	selfiePhotoValid   bool
	driverLicenseValid bool
}

func (c ReviewResponseContract) transition(ctx context.Context) (fsm.State, error) {
	switch {
	case !c.passportPhotoValid:
		return uploadPassportState, nil
	case !c.selfiePhotoValid:
		return uploadSelfieState, nil
	case !c.driverLicenseValid:
		return uploadDriverLicenseState, nil
	default:
		return approvedState, nil
	}
}

var (
	brokenState = fsm.State{
		Name:   "broken",
		Broken: true,
		OnTransition: func(_ context.Context) error {
			println("broken state faced: log and metric for investigation")

			return nil
		},
	}
	chooseCountryState = fsm.State{
		Name:    "choose_country",
		Initial: true,
		PossibleStates: fsm.MustNewStateMap(
			uploadPassportState,
		),
	}
	uploadPassportState = fsm.State{
		Name: "upload_passport",
		PossibleStates: fsm.MustNewStateMap(
			uploadSelfieState,
		),
	}
	uploadSelfieState = fsm.State{
		Name: "upload_selfie",
		PossibleStates: fsm.MustNewStateMap(
			uploadDriverLicenseState,
		),
	}
	uploadDriverLicenseState = fsm.State{
		Name: "upload_driver_license",
		PossibleStates: fsm.MustNewStateMap(
			reviewState,
		),
	}
	reviewState = fsm.State{
		Name: "review_requested",
		PossibleStates: fsm.MustNewStateMap(
			deniedState,
			approvedState,
		),
	}
	deniedState = fsm.State{
		Name: "denied",
	}
	approvedState = fsm.State{
		Name:  "approved",
		Final: true,
		OnTransition: func(_ context.Context) error {
			println("approved state achieved: log and metric for analytic")

			return nil
		},
	}
)

func Example() {
	opts := []RegistrationStateMachineOption{
		fsm.WithName[RegistrationContext]("registration_machine"),
		fsm.WithVersion[RegistrationContext]("v1"),
		fsm.WithStateStorage[RegistrationContext](nil),
		fsm.WithStateMap[RegistrationContext](
			brokenState,
			chooseCountryState,
			uploadPassportState,
			uploadSelfieState,
			uploadDriverLicenseState,
			reviewState,
			deniedState,
			approvedState,
		),
	}

	sm, err := fsm.New(opts...)
	if err != nil {
		panic(err)
	}

	var (
		ctx                 = context.Background()
		registrationContext = RegistrationContext{
			UUID: uuid.MustParse("10000000-0000-0000-0000-111111111111"),
		}
		event1 = ChooseCountryContract{country: "ru"}
		event2 = PassportPhotoContract{blob: []byte("passport photo")}
		event3 = SelfiePhotoContract{blob: []byte("selfie photo")}
		event4 = DriverLicensePhotoContract{blob: []byte("driver license")}
		event5 = ReviewResponseContract{
			passportPhotoValid: true,
			selfiePhotoValid:   true,
			driverLicenseValid: true,
		}
	)

	// Imitates event bus stream
	events := []RegistrationEvent{
		{Context: registrationContext, TransitionFunc: event1.transition},
		{Context: registrationContext, TransitionFunc: event2.transition},
		{Context: registrationContext, TransitionFunc: event3.transition},
		{Context: registrationContext, TransitionFunc: event4.transition},
		{Context: registrationContext, TransitionFunc: event5.transition},
	}

	var state fsm.State

	for i := range events {
		if state, err = sm.Transition(ctx, events[i]); err != nil {
			panic(err)
		}
	}

	if equal := state.Equal(approvedState); !equal {
		panic("unexpected state occurs")
	}
}
