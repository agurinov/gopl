package telegram

type User struct {
	Username     string `validate:"required"`
	FirstName    string `validate:"required"`
	LastName     string `validate:"required"`
	AuthorityBot string `validate:"required"`
	ID           int64  `validate:"required"`
	IsBot        bool
}

func DummyUser() User {
	return User{
		ID:           100500,
		Username:     "johndoe",
		FirstName:    "John",
		LastName:     "Doe",
		IsBot:        false,
		AuthorityBot: "DummyBot",
	}
}
