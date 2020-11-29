package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"nchat/utils"
	"net"
	"os"
	"strconv"
)

var AMOUNT = 100

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide arg.")
		return
	}
	CONNECT := arguments[1]
	if CONNECT == "" {
		CONNECT = "0.0.0.0:22111"
	}
	conn, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	uid, _ := strconv.Atoi(arguments[2])
	fid, _ := strconv.Atoi(arguments[3])
	AMOUNT = uid
	user := utils.UserRelations{
		UserID:  uid,
		Friends: []int{fid},
	}
	i := AMOUNT - 1
	for ; i > 0; i-- {
		user.Friends = append(user.Friends, i)
	}
	store, _ := json.Marshal(user)
	defer conn.Close()
	defer fmt.Printf("Closed %s\n", conn.RemoteAddr().String())
	conn.Write(store)
	fmt.Printf("Starting %v\n", string(store))
	go func(conn *net.TCPConn) {
		d := json.NewDecoder(conn)
		for {
			var doc utils.State
			err := d.Decode(&doc)
			if err != nil {
				if err == io.EOF {
					panic("no connection")
				}
				if err, ok := err.(net.Error); ok && err.Timeout() {
					continue
				}
				fmt.Println("Error: ", err)
			}
			fmt.Printf("ping %v\n", doc)
		}
	}(conn.(*net.TCPConn))

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		data, _ := reader.ReadString('\n')

		fid := user.Friends[0]
		//for _, fid := range user.Friends {
		state := utils.State{FriendID: fid, Online: true, Message: data[:len(data)-1]}
		msg, _ := json.Marshal(state)
		if _, err := conn.Write(msg); err != nil {
			fmt.Println(err)
		}
		//}
	}
}
