package controllers

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Processor interface {
	ProcessMessage(message string) ([]byte, error)
}

type Broker interface {
	GetStockQueueName() string
	Consume(key string) (<-chan amqp.Delivery, error)
}

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	processor  Processor
	rabbit     Broker
}

func NewWebsocketServer(p Processor, r Broker) *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		processor:  p,
		rabbit:     r,
	}
}

func (s *WsServer) Run() {
	go s.startMQConsumer()
	for {
		select {

		case client := <-s.register:
			s.registerClient(client)

		case client := <-s.unregister:
			s.unregisterClient(client)

		case message := <-s.broadcast:
			s.broadcastToClients(message)
		}

	}
}

func (s *WsServer) startMQConsumer() {
	msgs, err := s.rabbit.Consume(s.rabbit.GetStockQueueName())
	if err != nil {
		panic(fmt.Sprintf("starting rabbit consumer %v\n", err))
	}
	loop := make(chan bool)
	go func() {
		for msg := range msgs {
			for client := range s.clients {
				client.send <- msg.Body
			}
		}
	}()
	<-loop
}

func (s *WsServer) registerClient(client *Client) {
	s.clients[client] = true
}

func (s *WsServer) unregisterClient(client *Client) {
	if _, ok := s.clients[client]; ok {
		delete(s.clients, client)
	}
}

func (s *WsServer) broadcastToClients(message []byte) {
	m, err := s.processor.ProcessMessage(string(message))
	if err != nil {
		log.Default().Println(fmt.Sprintf("processing chat message error: %s", err.Error()))
	}
	if m != nil {
		for client := range s.clients {
			client.send <- m
		}
	}

}
