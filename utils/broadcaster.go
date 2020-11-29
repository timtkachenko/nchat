package utils

import (
	"fmt"
	"sync"
)
type KVMessage struct {
	SenderAddr string
	Type       string
	Message    []byte
}

type Broadcaster struct {
	mutex     *sync.RWMutex
	input     <-chan KVMessage
	receivers map[chan KVMessage]bool
}

func NewBroadcaster(input <-chan KVMessage) *Broadcaster {
	return &Broadcaster{
		mutex:     &sync.RWMutex{},
		input:     input,
		receivers: make(map[chan KVMessage]bool),
	}
}
func (b *Broadcaster) Run() {
	for {
		select {
		case msg := <-b.input:
			receivers := make(map[chan KVMessage]bool)
			b.mutex.RLock()
			for k, v := range b.receivers {
				receivers[k] = v
			}
			b.mutex.RUnlock()
			for msgChan := range receivers {
				msgChan <- msg
			}
		}
	}
}
func (b *Broadcaster) Add(msgChan chan KVMessage) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.receivers[msgChan] = true
}

func (b *Broadcaster) Delete(msgChan chan KVMessage) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	fmt.Println(len(b.receivers))
	delete(b.receivers, msgChan)
	fmt.Println(len(b.receivers))
}
