package sql

import (
	"fmt"
	"net"
	"strconv"

	"github.com/go-sql-driver/mysql"

	"github.com/agurinov/gopl/x"
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
		Driver:   x.Coalesce(other.Driver, c.Driver),
		Host:     x.Coalesce(other.Host, c.Host),
		Database: x.Coalesce(other.Database, c.Database),
		User:     x.Coalesce(other.User, c.User),
		Password: x.Coalesce(other.Password, c.Password),
		Port:     x.Coalesce(other.Port, c.Port),
		Enabled:  x.Coalesce(other.Enabled, c.Enabled),
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
