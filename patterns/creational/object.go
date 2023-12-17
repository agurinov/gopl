package creational

type (
	Object          interface{ any }
	ObjectValidable interface{ Validate() error }
)
