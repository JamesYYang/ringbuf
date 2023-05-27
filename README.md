# ringbuf
Ring buffer implementation by GO

## Concept

Ringbuf implements ring buffer base on channel so never blocks the writer.

If a value is written to ring buffer when its buffer is full then the oldest value in the buffer is discarded.

## Install

```
go get github.com/JamesYYang/ringbuf
```

## Example

```go

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

}

func ReceiveMessage(output <-chan Message) {
	for msg := range output {
		log.Printf("receiving message: %s\n", msg.Name)
		time.Sleep(5 * time.Second)
	}
}

func main() {

	rb := ringbuf.New[Message](10)

	go SendMessage(rb.In())
	go ReceiveMessage(rb.Out())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	log.Fatal("program interrupted")
}
```