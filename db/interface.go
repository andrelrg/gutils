package db

import "database/sql"

//Database Interface With Methods to be a database handler
type DatabaseInterface interface {
	Connect() error
	Ping() error
	Execute(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	GetJSON(sqlString string) (map[string]interface{}, error)
	MapScan(sqlString string, args ...interface{}) (map[string]interface{}, error)
	GetJSONList(sqlString string) ([]map[string]interface{}, error)
	Close()
}

//Config Interface With Methods to be a database config
type Config interface {
	GetHost() string
	GetPort() int
	GetUser() string
	GetPassword() string
	GetDatabase() string
}
