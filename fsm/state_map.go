package fsm

import "strings"

type StateMap struct {
	m       map[string]State
	initial string
	final   string
	broken  string
}

var EmptyStates = StateMap{}

func (sm StateMap) get(key string) (State, bool) {
	if key == "" {
		return EmptyState, false
	}

	state, exists := sm.m[key]

	return state, exists
}

func (sm *StateMap) set(state State) error {
	if _, exists := sm.m[state.Name]; exists {
		return StateAlreadyPresentError{StateName: state.Name}
	}

	sm.m[state.Name] = state

	if state.Initial {
		if sm.initial != "" {
			// TODO: doubts
			return StateAlreadyPresentError{StateName: sm.initial}
		}

		sm.initial = state.Name
	}

	if state.Final {
		if sm.final != "" {
			// TODO: doubts
			return StateAlreadyPresentError{StateName: sm.final}
		}

		sm.final = state.Name
	}

	if state.Broken {
		if sm.broken != "" {
			// TODO: doubts
			return StateAlreadyPresentError{StateName: sm.broken}
		}

		sm.broken = state.Name
	}

	return nil
}

func (sm StateMap) Len() int               { return len(sm.m) }
func (sm StateMap) Initial() (State, bool) { return sm.get(sm.initial) }
func (sm StateMap) Final() (State, bool)   { return sm.get(sm.final) }
func (sm StateMap) Broken() (State, bool)  { return sm.get(sm.broken) }

func (sm StateMap) StateNames() string {
	states := make([]string, 0, len(sm.m))

	for k := range sm.m {
		states = append(states, sm.m[k].Name)
	}

	return strings.Join(states, ",")
}

func (sm StateMap) Validate() error {
	var (
		initialState, initialStateExists = sm.Initial()
		finalState, finalStateExists     = sm.Final()
		brokenState, brokenStateExists   = sm.Broken()
	)

	switch {
	case !initialStateExists:
		return StateNotPresentError{StateName: "initial"}
	case !finalStateExists:
		return StateNotPresentError{StateName: "final"}
	case !brokenStateExists:
		return StateNotPresentError{StateName: "broken"}
	}

	switch {
	case initialState.PossibleStates.Len() == 0:
		return UnexpectedPossibleStatesError{
			StateName:          initialState.Name,
			MustBePresent:      true,
			PossibleStateNames: initialState.PossibleStates.StateNames(),
		}
	case finalState.PossibleStates.Len() > 0:
		return UnexpectedPossibleStatesError{
			StateName:          finalState.Name,
			MustBePresent:      false,
			PossibleStateNames: finalState.PossibleStates.StateNames(),
		}
	case brokenState.PossibleStates.Len() > 0:
		return UnexpectedPossibleStatesError{
			StateName:          brokenState.Name,
			MustBePresent:      false,
			PossibleStateNames: brokenState.PossibleStates.StateNames(),
		}
	}

	// TODO(a.gurinov): check reachability from initial to final

	return nil
}

func NewStateMap(states ...State) (StateMap, error) {
	smap := StateMap{
		m: make(map[string]State, len(states)),
	}

	for i := range states {
		if err := smap.set(states[i]); err != nil {
			return EmptyStates, err
		}
	}

	return smap, nil
}

func MustNewStateMap(states ...State) StateMap {
	smap, err := NewStateMap(states...)
	if err != nil {
		panic(err)
	}

	return smap
}
