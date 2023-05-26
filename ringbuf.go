package ringbuf

type RingBuf struct {
	input  chan interface{}
	output chan interface{}
	buffer interface{}
}

func New() *RingBuf {
	ch := &RingBuf{
		input:  make(chan interface{}),
		output: make(chan interface{}),
	}
	go ch.ringBuffer()
	return ch
}

func (ch *RingBuf) in() chan<- interface{} {
	return ch.input
}

func (ch *RingBuf) out() <-chan interface{} {
	return ch.output
}

func (ch *RingBuf) ringBuffer() {
	var input, output chan interface{}
	var next interface{}
	input = ch.input

	for input != nil || output != nil {
		select {
		case output <- next:
			ch.buffer = nil
		default:
			select {
			case elem, open := <-input:
				if !open {
					input = nil
					break
				}
				ch.buffer = &elem
			case output <- next:
				ch.buffer = nil
			}
		}

		if ch.buffer == nil {
			output = nil
			continue
		}

		output = ch.output
		next = ch.buffer
	}

	close(ch.output)
}
