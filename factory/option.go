package pl_factory

type (
	Option[O object]    func(*O) error
	OptionSet[O object] []func(*O) error
)
