package main

import (
	"nchat/utils"
)

type manager struct {
	nClients    *MapAddrNetClient
	users       *MapIntNetClient
	broadcaster *utils.Broadcaster
}

func NewManager(netClientList *MapAddrNetClient, users *MapIntNetClient, inputDataCh chan utils.KVMessage) *manager {
	br := utils.NewBroadcaster(inputDataCh)
	go br.Run()
	return &manager{netClientList, users, br}
}

func (m *manager) Add(nc *netClient, maxChan chan bool) {
	defer func(maxChan chan bool) { <-maxChan }(maxChan)
	sCounter.Incr()
	defer sCounter.Decr()

	m.nClients.Set(nc.conn.RemoteAddr(), nc)
	m.broadcaster.Add(nc.dataCh)
	nc.drop(nc.handleConnection())
	m.broadcaster.Delete(nc.dataCh)
}
