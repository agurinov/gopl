package telegram

import initdata "github.com/telegram-mini-apps/init-data-golang"

type User = initdata.User

func Dummy() User {
	return User{
		ID:        100500,
		Username:  "johndoe",
		FirstName: "John",
		LastName:  "Doe",
	}
}
