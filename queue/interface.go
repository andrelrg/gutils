package queue

//Config Interface With Methods to be a database config
type Config interface {
	GetHost() string
	GetPort() int
	GetUser() string
	GetPassword() string
	GetDatabase() string
}