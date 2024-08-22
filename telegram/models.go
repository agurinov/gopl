package telegram

import (
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
	}
	PersonalChat struct {
		ID      int64 `validate:"required"`
		Enabled bool
	}
)

func (u User) String() string {
	if u.Username != "" {
		return u.Username
	}

	return strconv.FormatInt(u.ID, 10)
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
