package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"nchat/utils"
	"net"
	"sync"
)

type netClient struct {
	mx            sync.RWMutex
	conn          *net.TCPConn
	netClientList *MapAddrNetClient
	users         *MapIntNetClient
	id            int
	Friends       *MapIntNetClient
	dataCh        chan utils.KVMessage
	ctx           context.Context
}

func NewTcpClient(conn *net.TCPConn, netClientList *MapAddrNetClient, users *MapIntNetClient) *netClient {
	return &netClient{conn: conn, netClientList: netClientList, users: users, dataCh: make(chan utils.KVMessage, 1)}
}

func (nc *netClient) NoConn() bool {
	nc.mx.RLock()
	defer nc.mx.RUnlock()
	return nc.conn == nil
}
func (nc *netClient) CloseConn() {
	nc.mx.Lock()
	defer nc.mx.Unlock()
	if nc.conn == nil {
		return
	}
	nc.conn.Close()
	nc.conn = nil
}
func (nc *netClient) notify() {
	connected := utils.KVMessage{Type: "connected", Message: []byte(nc.conn.RemoteAddr().String())}
	msg, _ := json.Marshal(connected)
	node.send(msg)
}

func (nc *netClient) handleConnection() (userID int) {
	defer nc.conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nc.ctx = ctx

	go nc.waitData()
	//fmt.Printf("Serving %s\n", nc.conn.RemoteAddr().String())
	//defer fmt.Printf("Closed %s\n", nc.conn.RemoteAddr().String())
	for {
		data, err := bufio.NewReader(nc.conn).ReadString('}')
		if err != nil {
			if err == io.EOF {
				return
			}
			//fmt.Println("Error: ", err)
			return
		}
		var doc utils.UserRelations
		err = json.Unmarshal([]byte(data), &doc)
		if err != nil {
			fmt.Println("Unmarshal Error: ", err)
		}
		// protocol fallback
		if doc.UserID == 0 {
			nc.accept([]byte(data), "")
		} else {
			userID = doc.UserID
			if nc.users.Get(userID) != nil {
				// unique
				userID = -1
				return
			}
			nc.setup(doc)
		}
	}
}
func (nc *netClient) drop(userID int) {
	fmt.Println("user went offline ", userID)
	if userID == 0 {
		return
	}
	user := nc.users.Drain(userID)
	if user == nil {
		return
	}
	friendList := user.Friends.Fetch()
	for _, friend := range friendList {
		if friend.NoConn() {
			continue
		}
		en := json.NewEncoder(friend.conn)
		err := en.Encode(utils.State{FriendID: userID, Online: false})
		if err != nil {
			continue
		}
		if f := nc.users.Get(friend.id); f != nil {
			f.Friends.Drain(userID)
		}
	}
}

func (nc *netClient) setup(msg utils.UserRelations) {
	nc.id = msg.UserID
	nc.Friends = NewMapIntNetClient()
	state := utils.State{FriendID: msg.UserID, Online: true}
	for _, fid := range msg.Friends {
		if friend := nc.users.Get(fid); friend != nil {
			if friend.NoConn() {
				continue
			}
			en := json.NewEncoder(friend.conn)
			err := en.Encode(state)
			if err != nil {
				friend.CloseConn()
				if fid != 0 {
					go nc.drop(fid)
				}
				continue
			}
			friend.Friends.Set(nc.id, nc)
			nc.Friends.Set(friend.id, friend)
		} else {
			state.UserID = fid
			data, _ := json.Marshal(state)
			forwardToNode := utils.KVMessage{Type: "forward", Message: data}
			msg, _ := json.Marshal(forwardToNode)
			node.send(msg)
		}
	}
	nc.users.Set(msg.UserID, nc)
}

func (nc *netClient) accept(data []byte, msgType string) {
	var doc utils.State
	if err := json.Unmarshal(data, &doc); err != nil {
		fmt.Println("forward Unmarshal Error: ", err)
	}
	// assuming state
	if doc.UserID != 0 {
		if user := nc.users.Get(doc.UserID); user != nil {
			en := json.NewEncoder(user.conn)
			if err := en.Encode(utils.State{FriendID: doc.UserID, Online: true}); err != nil {
				fmt.Println("forward Error: ", err)
			}
			return
		}
	}
	// assuming data message
	if friend := nc.users.Get(doc.FriendID); friend != nil {
		en := json.NewEncoder(friend.conn)

		if err := en.Encode(utils.State{FriendID: nc.id, Online: true, Message: doc.Message}); err != nil {
			fmt.Println("forward Error: ", err)
		}
		return
	} else if msgType == "" {
		forwardToNode := utils.KVMessage{Type: "forward", Message: data}
		msg, _ := json.Marshal(forwardToNode)
		node.send(msg)
	}
}

func (nc *netClient) waitData() {
	for {
		select {
		case m := <-nc.dataCh:
			nc.accept(m.Message, m.Type)
		case <-nc.ctx.Done():
			return
		}
	}
}
