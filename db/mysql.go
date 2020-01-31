package db

import (
	"fmt"
	_ "github.com/astrolink/gorm/dialects/mysql" //This is to get the mysql driver
	"log"
)

// NewMySQL makes a new instance of Database and connect to a MySQL database.
func NewMySQL(config Config) *Database {

	connectionLine := "%s:%s@tcp(%s:%s)/%s"
	connectionLine = fmt.Sprintf(connectionLine,
		config.GetUser(), config.GetPassword(), config.GetHost(), config.GetPort(), config.GetDatabase())
	mysql := Database{
		ConnectionLine: connectionLine,
		Driver: "mysql",
	}
	var err error
	err = mysql.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	return &mysql
}
