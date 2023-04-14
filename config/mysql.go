package config

import "time"

// MySQL is the configuration of the MySQL DB
type MySQL struct {
	// TestMode indicate if is a dummy connection or not (accessing to a real database)
	// Environment variable: APP_MYSQL_TESTMODE
	TestMode bool `default:"false"`
	// Connection is the connection string passed to `sql.Open`.
	// Environment variable: APP_MYSQL_CONNECTION
	Connection string
	// ConnRetrySleep is the time to wait before retrying a connection
	// Environment variable: APP_MYSQL_CONNRETRYSLEEP
	ConnRetrySleep time.Duration `default:"3s"`
	// ConnMaxLifetime is the maximum connection lifetime
	// Environment variable: APP_MYSQL_CONNMAXLIFETIME
	ConnMaxLifetime time.Duration
	// MaxIdleConns is the maxium idle connections
	// Environment variable: APP_MYSQL_MAXIDLECONNS
	MaxIdleConns uint `validate:"min=0"`
	// MaxOpenConns is the maximum open connections
	// Environment variable: APP_MYSQL_MAXOPENCONNS
	MaxOpenConns uint `validate:"min=0"`
}
