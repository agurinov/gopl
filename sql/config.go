package sql

import (
	"fmt"
	"net"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Driver   string `validate:"oneof=noop mysql pgx"`
	Host     string `validate:"required"`
	Database string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Port     int64  `validate:"gt=1000,lt=65536"`
}

func (c Config) DSN() string {
	switch c.Driver {
	case "mysql":
		driverConfig := mysql.Config{
			User:   c.User,
			Passwd: c.Password,
			Net:    "tcp",
			Addr: net.JoinHostPort(
				c.Host,
				strconv.FormatInt(c.Port, 10),
			),
			DBName: c.Database,
		}

		return driverConfig.FormatDSN()
	case "pgx":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Host,
			c.Port,
			c.User,
			c.Password,
			c.Database,
		)
	default:
		return ""
	}
}
