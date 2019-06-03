package mysql

import (
	"database/sql"
	"fmt"
	"net/url"
)

// Config holds everything you need to
// connect and interact with a MySQL DB.
type MysqlConfig struct {
	Pw              string `envconfig:"MYSQL_PW"`
	User            string `envconfig:"MYSQL_USER"`
	Port            int    `envconfig:"MYSQL_PORT"`
	DBName          string `envconfig:"MYSQL_DB_NAME"`
	Location        string `envconfig:"MYSQL_LOCATION"`
	Host            string `envconfig:"MYSQL_HOST_NAME"`
	ReadTimeout     string `envconfig:"MYSQL_READ_TIMEOUT"`
	WriteTimeout    string `envconfig:"MYSQL_WRITE_TIMEOUT"`
	AddtlDSNOptions string `envconfig:"MYSQL_ADDTL_DSN_OPTIONS"`
}

const (
	// DefaultLocation is the default location for MySQL connections.
	DefaultLocation = "localhost"
	// DefaultMySQLPort is the default port for MySQL connections.
	DefaultMySQLPort = 3306
)

var (
	// MaxOpenConns will be used to set a MySQL
	// drivers MaxOpenConns value.
	MaxOpenConns = 1
	// MaxIdleConns will be used to set a MySQL
	// drivers MaxIdleConns value.
	MaxIdleConns = 1
)

// DB will attempt to open a sql connection with
// the credentials and the current MySQLMaxOpenConns
// and MySQLMaxIdleConns values.
// Users must import a mysql driver in their
// main to use this.
func (m *MysqlConfig) DB() (*sql.DB, error) {
	db, err := sql.Open("mysql", m.String())
	if err != nil {
		return db, err
	}
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	return db, nil
}

// String will return the MySQL connection string.
func (m *MysqlConfig) String() string {
	if m.Port == 0 {
		m.Port = DefaultMySQLPort
	}

	if m.Location != "" {
		m.Location = url.QueryEscape(m.Location)
	} else {
		m.Location = url.QueryEscape(DefaultLocation)
	}

	args, _ := url.ParseQuery(m.AddtlDSNOptions)

	args.Set("parseTime", "true")

	if m.ReadTimeout != "" {
		args.Set("readTimeout", m.ReadTimeout)
	}
	if m.WriteTimeout != "" {
		args.Set("writeTimeout", m.WriteTimeout)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		m.User,
		m.Pw,
		m.Host,
		m.Port,
		m.DBName,
	)
}
