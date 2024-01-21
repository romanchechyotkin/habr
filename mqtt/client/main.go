package main

import (
	"log"
	"net/http"
	"os"

	"github.com/romanchechyotkin/habr/mqtt/broker/subscriber"
)

const Address = "127.0.0.1:1883"

func main() {
	topic := os.Args[1]
	opts := &subscriber.Options{
		Topic:   topic,
		Address: Address,
	}

	sub, err := subscriber.New(opts)
	if err != nil {
		log.Fatal(err)
	}

	logging(sub)

	log.Fatal(http.ListenAndServe(":5000", nil))
}

func logging(sub *subscriber.Subscriber) {
	msgChan := sub.GetMsgChannel()

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("closed chan")
				return
			}
			log.Println("got msg in logging func", string(msg))
		}
	}
}
