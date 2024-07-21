package telegram

type User struct {
	Username     string
	FirstName    string
	LastName     string
	AuthorityBot string `validate:"required"`
	ID           int64  `validate:"required"`
	IsBot        bool
}
