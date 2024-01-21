package main

import (
	"log"

	"github.com/romanchechyotkin/habr/mqtt/broker/publisher"
)

const Address = "127.0.0.1:1883"

func main() {
	opts := &publisher.Options{
		Topic:   "test",
		Address: Address,
	}

	pub, err := publisher.New(opts)
	if err != nil {
		log.Println("failed to create new instance")
		return
	}

	err = pub.Write([]byte("test topiasdasdasc asd \n wrasdite"))
	if err != nil {
		log.Println("failed to write")
		return
	}
}
