package db

import (
	"database/sql"
)

type ConnectionManager struct {
	db *Database
}

//New is responsible of creating a new instance of consultQuestionRepository for mysql database
func NewMySqlConnectionManager(config Config) *ConnectionManager {
	mysql := NewMySQL(config)
	return &ConnectionManager{db: mysql}
}

//New is responsible of creating a new instance of consultQuestionRepository for postgres database
func NewPgSqlConnectionManager(config Config) *ConnectionManager {
	pgsql := NewPgSQL(config)
	return &ConnectionManager{db: pgsql}
}

// Close closes database connection
func (c *ConnectionManager) StartTransaction() (*sql.Tx, error) {
	return c.db.StartTransaction()
}

// Close closes database connection
func (c *ConnectionManager) GetConnection() *sql.DB {
	return c.db.Conn
}

// Close closes database connection
func (c *ConnectionManager) Close() {
	c.db.Close()
}
