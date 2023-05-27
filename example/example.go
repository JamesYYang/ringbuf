package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/JamesYYang/ringbuf"
)

type Message struct {
	Name  string
	Value string
}

func SendMessage(input chan<- Message) {
	for i := 0; i < 50; i++ {
		msg := Message{
			Name:  fmt.Sprintf("Message %d", i),
			Value: fmt.Sprintf("Value %d", i),
		}
		input <- msg
		log.Printf("sending message: %s\n", msg.Name)
		time.Sleep(1 * time.Second)
	}
	close(input)

}

func ReceiveMessage(output <-chan Message) {
	for msg := range output {
		log.Printf("receiving message: %s\n", msg.Name)
		time.Sleep(5 * time.Second)
	}
}

func main() {
	log.SetPrefix(fmt.Sprintf("[%d]: ", os.Getpid()))
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	rb := ringbuf.New[Message](10)

	go SendMessage(rb.In())
	go ReceiveMessage(rb.Out())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	log.Fatal("program interrupted")
}
