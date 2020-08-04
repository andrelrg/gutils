package db

import (
	"database/sql"
)

type ConnectionManager struct {
	db *Database
}

//New is responsible of creating a new instance of consultQuestionRepository
func NewConnectionManager(config Config) *ConnectionManager {
	mysql := NewMySQL(config)
	return &ConnectionManager{db: mysql}
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
