package telegram

type User struct {
	Username     string `validate:"required"`
	FirstName    string
	LastName     string
	AuthorityBot string `validate:"required"`
	ID           int64  `validate:"required"`
	IsBot        bool
}
