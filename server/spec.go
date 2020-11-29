package main

import (
	"net"
	"sync"
)

type MapIntNetClient struct {
	mx    sync.RWMutex
	Store map[int]*netClient
}

func NewMapIntNetClient() *MapIntNetClient {
	return &MapIntNetClient{Store: map[int]*netClient{}}
}
func (qs *MapIntNetClient) Set(key int, val *netClient) {
	qs.mx.Lock()
	defer qs.mx.Unlock()
	qs.Store[key] = val
}

func (qs *MapIntNetClient) Get(key int) *netClient {
	qs.mx.RLock()
	defer qs.mx.RUnlock()
	return qs.Store[key]
}

func (qs *MapIntNetClient) Fetch() map[int]*netClient {
	qs.mx.RLock()
	defer qs.mx.RUnlock()
	ms := map[int]*netClient{}
	for k, v := range qs.Store {
		ms[k] = v
	}
	return ms
}

func (qs *MapIntNetClient) Drain(key int) *netClient {
	qs.mx.Lock()
	defer qs.mx.Unlock()
	defer delete(qs.Store, key)
	return qs.Store[key]
}



//
type MapAddrNetClient struct {
	mx    sync.RWMutex
	Store map[net.Addr]*netClient
}

func NewMapAddrNetClient() *MapAddrNetClient {
	return &MapAddrNetClient{Store: map[net.Addr]*netClient{}}
}

func (qs *MapAddrNetClient) Set(key net.Addr, val *netClient) {
	qs.mx.Lock()
	defer qs.mx.Unlock()
	qs.Store[key] = val
}

func (qs *MapAddrNetClient) Get(key net.Addr) *netClient {
	qs.mx.RLock()
	defer qs.mx.RUnlock()
	return qs.Store[key]
}

func (qs *MapAddrNetClient) Fetch() map[net.Addr]*netClient {
	qs.mx.RLock()
	defer qs.mx.RUnlock()
	ms := map[net.Addr]*netClient{}
	for k, v := range qs.Store {
		ms[k] = v
	}
	return ms
}

func (qs *MapAddrNetClient) Drain(key net.Addr) *netClient {
	qs.mx.Lock()
	defer qs.mx.Unlock()
	defer delete(qs.Store, key)
	return qs.Store[key]
}

