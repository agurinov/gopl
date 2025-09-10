package sql

import (
	"cmp"
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
	Enabled  bool
}

func (c Config) MergeWith(other Config) Config {
	return Config{
		Driver:   cmp.Or(other.Driver, c.Driver),
		Host:     cmp.Or(other.Host, c.Host),
		Database: cmp.Or(other.Database, c.Database),
		User:     cmp.Or(other.User, c.User),
		Password: cmp.Or(other.Password, c.Password),
		Port:     cmp.Or(other.Port, c.Port),
		Enabled:  cmp.Or(other.Enabled, c.Enabled),
	}
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
