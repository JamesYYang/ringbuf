package ringbuf

import (
	"container/list"
)

type RingBuf[T any] struct {
	input   chan T
	output  chan T
	buffer  *list.List
	bufSize int
}

func New[T any](size int) *RingBuf[T] {
	ch := &RingBuf[T]{
		input:   make(chan T),
		output:  make(chan T),
		buffer:  list.New(),
		bufSize: size,
	}
	go ch.ringBuffer()
	return ch
}

func (ch *RingBuf[T]) In() chan<- T {
	return ch.input
}

func (ch *RingBuf[T]) Out() <-chan T {
	return ch.output
}

func (ch *RingBuf[T]) push(ele T) {
	ch.buffer.PushBack(ele)
	if ch.buffer.Len() > ch.bufSize {
		ch.buffer.Remove(ch.buffer.Front())
	}
}

func (ch *RingBuf[T]) pop() *list.Element {
	ele := ch.buffer.Front()
	if ele != nil {
		ch.buffer.Remove(ele)
	}
	return ele
}

func (ch *RingBuf[T]) ringBuffer() {
	var input, output chan T
	var next T
	input = ch.input
	isSend := true

	for input != nil || output != nil {
		select {
		case output <- next:
			isSend = true
		default:
			select {
			case elem, open := <-input:
				if !open {
					input = nil
					break
				}
				ch.push(elem)
			case output <- next:
				isSend = true
			}
		}

		if ch.buffer.Len() == 0 {
			output = nil
			continue
		}

		output = ch.output
		if isSend {
			isSend = false
			ele := ch.pop()
			if ele != nil {
				next = ele.Value.(T)
			}
		}
	}

	close(ch.output)
}
