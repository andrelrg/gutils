package queue


import (
	"encoding/json"
	"fmt"
	"github.com/astrolink/amqp"
	"log"
	"strconv"
)

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	queue amqp.Queue
	config Config
}

// NewRabbitMQ creates a nem instance on RabbitMQ
func NewRabbitMQ(config Config) (*RabbitMQ, error) {
	var err error
	r := RabbitMQ{config: config}

	err = r.connect()

	if err != nil {
		err = fmt.Errorf("error opening rabbit connection, %s", err.Error())
		log.Println(err)
		return &r, err
	}

	// creating channel
	r.channel, err = r.conn.Channel()

	if err != nil {
		err = fmt.Errorf("error opening channel, %s", err.Error())
		log.Println(err)
		return &r, err
	}

	// creating queue
	// in here we took the database attribute as queue name just for reuse purposes
	name := config.GetDatabase()

	r.queue, err = r.channel.QueueDeclare(name, true, false, false, false, nil)

	if err != nil {
		err = fmt.Errorf("error declaring exchange, %s", err.Error())
		log.Println(err)
		return &r, err
	}


	return &r, nil
}


// connect open a connection to rabbit server
func (r *RabbitMQ) connect() error {
	var url string

	if r.config.GetUser() != "" {
		url = r.config.GetUser()
	}
	if r.config.GetPassword() != "" {
		url += ":" + r.config.GetPassword()
	}
	if r.config.GetUser() != "" || r.config.GetPassword() != "" {
		url += "@"
	}
	url += r.config.GetHost() + ":" + strconv.Itoa(r.config.GetPort())

	conn, err := amqp.Dial("amqp://" + url)

	if err != nil {
		return err
	}

	r.conn = conn
	return nil
}

// Publish sends a message to Rabbit queue
func (r *RabbitMQ) Publish(data interface{}) error {
	var err error

	body, err := json.Marshal(data)

	if err != nil {
		err = fmt.Errorf("error converting data struct to json, %s", err.Error())
		log.Println(err)
		return err
	}

	message := amqp.Publishing{ContentType: "application/json", Body: body}

	exchange := ""
	mandatory := false
	immediate := false

	err = r.channel.Publish(exchange, r.queue.Name, mandatory, immediate, message)

	if err != nil {
		err = fmt.Errorf("error publishing message on queue, %s", err.Error())
		log.Println(err)
		return err
	}

	return nil
}

// TestRabbitMQConnection tries to connect to specified rabbitQM broker
func TestRabbitMQConnection(config config.Interface) error {
	r := RabbitMQ{config: config}
	var err error

	err = r.connect()
	defer r.Close()

	if err != nil {
		log.Println(err)
	}

	return err
}

func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.channel
}

// Close closes the rabbit connection and the channel
func (r *RabbitMQ) Close()  {
	r.channel.Close()
	r.conn.Close()
}