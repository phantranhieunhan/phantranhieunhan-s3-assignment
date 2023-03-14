package rabbitmq

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/streadway/amqp"

	"github.com/phantranhieunhan/s3-assignment/common/logger"
)

// WaitTimeReconnect constants
const (
	WaitTimeReconnect      = 5
	DefaultPrefetchCount   = 10
	DefaultDeclarationFile = "./config/declaration.json"
)

type KVMessage struct {
	Key   string
	Value string
}

// Topic : Topic
type Topic struct {
	Name     string
	Exchange string
}

// Queue : Queue
type Queue struct {
	Name   string
	Topics []Topic
}

// Declaration : Declaration
type Declaration struct {
	Exchanges []string
	Queues    []Queue
}

type mqImpl struct {
	url              string
	connection       *amqp.Connection
	channel          *amqp.Channel
	declaration      *Declaration
	connectionClosed bool
	channelClosed    bool
	chanErr          chan *amqp.Error
}

/**========================================================================
 *                           INTERFACE IMPLEMENTATION
 *========================================================================**/

func (mq *mqImpl) PushKVMessage(exchange, routing string, data KVMessage) (err error) {
	bytes, _ := json.Marshal(data)
	return mq.pushMessage(exchange, routing, bytes)
}

func (mq *mqImpl) PushRawMessage(exchange, routing string, data []byte) (err error) {
	return mq.pushMessage(exchange, routing, data)
}

func (mq *mqImpl) Consume() (chan amqp.Delivery, chan error) {
	chanMsg := make(chan amqp.Delivery)
	chanErr := make(chan error)
	mq.consuming(chanMsg, chanErr)
	go func() {
		for {
			closedErr := <-mq.chanErr
			if closedErr != nil {
				logger.Errorf("[RabbitMQ] connection is closed, err: %v. Reconnecting...", closedErr)
				err := mq.reconnect()
				if err != nil {
					logger.Errorf("[RabbitMQ] failed to reconnect, err: %v", err)
					continue
				}
				mq.consuming(chanMsg, chanErr)
			}
		}
	}()
	return chanMsg, chanErr
}

/**========================================================================
 *                           PRIVATE METHODS
 *========================================================================**/

func (mq *mqImpl) newConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial(mq.url)
	for err != nil {
		logger.Errorf(
			"[RabbitMQ] failed to create new connection to AMQP, err: %v. Sleep %d seconds to reconnect",
			err,
			WaitTimeReconnect,
		)
		time.Sleep(WaitTimeReconnect * time.Second)
		conn, err = amqp.Dial(mq.url)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	mq.channel = ch
	mq.connection = conn
	mq.chanErr = make(chan *amqp.Error)
	mq.connection.NotifyClose(mq.chanErr)

	return conn, nil
}

func (mq *mqImpl) newChannel() (*amqp.Channel, error) {
	mq.ensureConnection()
	if mq.connection == nil || mq.connection.IsClosed() {
		return nil, errors.New("connection is not open")
	}

	channel, err := mq.connection.Channel()
	if err != nil {
		logger.Fatalf("[RabbitMQ] failed to new channel, err: %v", err)
		return nil, err
	}
	mq.channel = channel
	mq.channelClosed = false

	err = mq.channel.Qos(DefaultPrefetchCount, 0, false)
	if err != nil {
		logger.Fatalf("[RabbitMQ] failed to prefetch consumer channel, err: %v", err)
		return nil, err
	}

	logger.Info("[RabbitMQ] new channel successfully")
	return channel, nil
}

func (mq *mqImpl) ensureConnection() (err error) {
	if mq.connection == nil || mq.connection.IsClosed() {
		_, err = mq.newConnection()
		if err != nil {
			return err
		}
	}
	return nil
}

func (mq *mqImpl) closeChannel() error {
	if mq.connectionClosed || mq.channelClosed {
		return nil
	}
	logger.Info("[RabbitMQ] close channel.")
	if mq.channel != nil {
		_ = mq.channel.Close()
		mq.channel = nil
		mq.channelClosed = true
	}

	return nil
}

func (mq *mqImpl) closeConnection() (err error) {
	if !mq.connection.IsClosed() {
		err = mq.connection.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (mq *mqImpl) initFromConfigFile(declarationFile string) error {
	if declarationFile == "" {
		declarationFile = DefaultDeclarationFile
	}
	var dec Declaration
	decs, err := ioutil.ReadFile(declarationFile)
	if err != nil {
		logger.Fatalf("[RabbitMQ] failed to read file %s, err: %v", declarationFile, err)
	}

	if err = json.Unmarshal(decs, &dec); err != nil {
		logger.Fatalf("[RabbitMQ] failed to unmarshal declaration", err)
	}

	// Exchange
	for _, exName := range dec.Exchanges {
		// ExchangeDeclare: name, type, durable, autoDelete, internal, noWait, args
		err = mq.channel.ExchangeDeclare(exName, "topic", true, false, false, false, nil)
		if err != nil {
			logger.Fatalf("[RabbitMQ] failed to declare exchange, err: %v", err)
		}
	}

	// Queue
	for _, q := range dec.Queues {
		// QueueDeclare: name, durable, autoDelete, exclusive, noWait, args
		args := amqp.Table{
			"x-queue-mode": "lazy",
		}
		_, err := mq.channel.QueueDeclare(q.Name, true, false, false, false, args)
		if err != nil {
			logger.Fatalf("[RabbitMQ] failed to declare queue, err: %v", err)
		}

		// Binding
		for _, t := range q.Topics {
			// QueueBind: queue name, routing key, exchange, noWait, args
			err = mq.channel.QueueBind(q.Name, t.Name, t.Exchange, false, nil)
			if err != nil {
				logger.Fatalf("[RabbitMQ] failed to bind queue, err: %v", err)
			}
		}
	}

	mq.declaration = &dec
	return nil
}

func (mq *mqImpl) reconnect() (err error) {
	err = mq.closeAll()
	if err != nil {
		logger.Infof("[RabbitMQ] failed to close connection, err: %v", err)
	}

	var conn *amqp.Connection
	for {
		conn, err = amqp.Dial(mq.url)
		if err == nil {
			break
		}

		logger.Infof(
			"[RabbitMQ] failed to create new connection to AMQP: %s. Sleep %d seconds to reconnect.",
			err,
			WaitTimeReconnect,
		)
		time.Sleep(WaitTimeReconnect * time.Second)
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	mq.channel = ch
	mq.connection = conn
	mq.chanErr = make(chan *amqp.Error)
	conn.NotifyClose(mq.chanErr)

	logger.Info("[RabbitMQ] reconnect rabbitMQ successfully!!!")
	return nil
}

func (mq *mqImpl) chanelClosed() bool {
	if mq.channel == nil || mq.channelClosed {
		return true
	}
	return false
}

func (mq *mqImpl) pushMessage(exchange, routing string, data []byte) (err error) {
	mq.ensureConnection()
	channel, _ := mq.connection.Channel()
	defer channel.Close()

	err = mq.channel.Publish(
		exchange, // exchange
		routing,  // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		})
	return err
}

func (mq *mqImpl) consuming(chanMsg chan amqp.Delivery, chanErr chan error) {
	for _, queue := range mq.declaration.Queues {
		go func(qName string) {
			// Consume: queue, consumer, autoAck, exclusive, noLocal, noWait, args
			msgs, err := mq.channel.Consume(qName, "", false, false, false, false, nil)
			if err != nil {
				chanErr <- err
			}

			forever := make(chan bool)
			go func() {
				for d := range msgs {
					chanMsg <- d
				}
			}()
			<-forever
		}(queue.Name)
	}
}

// closeAll : Close connection and channel
func (mq *mqImpl) closeAll() (err error) {
	if !mq.connection.IsClosed() {
		err = mq.connection.Close()
		if err != nil {
			return err
		}
	}
	return nil
}