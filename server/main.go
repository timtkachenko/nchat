package main

import (
	"fmt"
	"github.com/libp2p/go-reuseport"
	"log"
	"nchat/utils"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

const PORT = 22111

var port int
var sCounter utils.CounterSync
var node *netNode
var dataCh chan utils.KVMessage

func main() {
	port = PORT
	arguments := os.Args
	if len(arguments) > 1 {
		fromArgs, _ := strconv.Atoi(arguments[1])
		port = fromArgs
	}

	var joinAddr string
	if len(arguments) > 2 {
		joinAddr = arguments[2]
	}
	log.Println("Starting...", port)
	maxListeners := 1
	//maxListeners := runtime.NumCPU() / 2
	runtime.GOMAXPROCS(maxListeners)
	fmt.Printf("maxListeners: %d\n", maxListeners)
	sCounter = utils.CounterSync{}
	go stats()

	dataCh = make(chan utils.KVMessage)
	node = NewNetNode(fmt.Sprintf("tcp://0.0.0.0:%d", port+1), dataCh)
	go node.start()
	go node.join(fmt.Sprintf("tcp://%s", joinAddr))

	for i := 0; i < maxListeners; i++ {
		go beginListen()
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-quit
	time.Sleep(time.Second)
	log.Println("Shutdown ...")
}

func beginListen() {
	addr := net.TCPAddr{
		Port: port,
		IP:   net.IP{0, 0, 0, 0},
	}

	listener, err := reuseport.Listen("tcp", addr.String())
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	netClientList := NewMapAddrNetClient()
	users := NewMapIntNetClient()
	manager := NewManager(netClientList, users, dataCh)

	maxFileDescriptors := 1000
	maxChan := make(chan bool, maxFileDescriptors)
	for {
		maxChan <- true
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go manager.Add(NewTcpClient(connection.(*net.TCPConn), netClientList, users), maxChan)
	}
}

func stats() {
	var history int
	passed := 1
	for {
		select {
		case <-time.After(time.Second):
			fmt.Printf("online: %d\n", sCounter.Get())
			fmt.Printf("--NumGoroutine: %d\n", runtime.NumGoroutine())
			passed++
			if passed > 5 && history == sCounter.Get() {
				//sCounter.Reset()
				passed = 0
			} else if passed > 5 {
				history = sCounter.Get()
				passed = 0
			}
			break
		}

	}
}
