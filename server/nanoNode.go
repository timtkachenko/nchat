package main

import (
	"encoding/json"
	"fmt"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/bus"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
	"log"
	"nchat/utils"
)

type netNode struct {
	addr   string
	socket mangos.Socket
	dataCh chan utils.KVMessage
}

func NewNetNode(addr string, dataCh chan utils.KVMessage) *netNode {
	log.Println("Starting node...", addr)
	sock, err := bus.NewSocket()
	if err != nil {
		fmt.Printf("bus.NewSocket: %s", err)
	}
	if err = sock.Listen(addr); err != nil {
		fmt.Printf("nano Listen: %s", err.Error())
	}
	return &netNode{addr, sock, dataCh}
}

func (n *netNode) send(msg []byte) {
	fmt.Printf("nano %s", string(msg))
	err := n.socket.Send(msg)
	if err != nil {
		fmt.Printf("nano Send: %s", err.Error())
	}
}
func (n *netNode) start() {
	for {
		data, err := n.socket.Recv()
		if err != nil {
			fmt.Printf("nano Recv: %s", err.Error())
		}
		var msg utils.KVMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			fmt.Println("nano Unmarshal Error: ", err)
		}
		if msg.Type == "forward" {
			n.dataCh <- msg
		}
		//fmt.Printf("nano %s\n", string(data))
	}
}

func (n *netNode) join(addr string) {
	err := n.socket.Dial(addr)
	if err != nil {
		fmt.Printf("nano Dial: %s", err.Error())
		return
	}
}
