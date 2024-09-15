package telegram

import (
	"cmp"
	"strconv"
	"strings"
)

type (
	User struct {
		Username     string
		FirstName    string
		LastName     string
		AuthorityBot string `validate:"required"`
		ID           int64  `validate:"required"`
		IsBot        bool
		PrivateChat  PrivateChat
	}
	PrivateChat struct {
		ID      int64
		Enabled bool
	}
)

func (u User) String() string {
	return cmp.Or(
		Username(u.Username),
		strconv.FormatInt(u.ID, 10),
	)
}

func Username(username string) string {
	switch {
	case username == "":
		return ""
	case strings.HasPrefix(username, "@"):
		return username
	default:
		return "@" + username
	}
}
