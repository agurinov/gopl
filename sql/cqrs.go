package sql

import c "github.com/agurinov/gopl/patterns/creational"

type (
	CQRS interface {
		RW() DB
		RO() DB
	}
	cqrs struct {
		rw DB
		ro DB
	}
	ClientOption c.Option[cqrs]
)

var NewCQRS = c.NewWithValidate[cqrs, ClientOption]

func (cq cqrs) RW() DB { return cq.rw }
func (cq cqrs) RO() DB { return cq.ro }
