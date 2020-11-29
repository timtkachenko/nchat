package utils

import "sync"

type State struct {
	UserID   int    `json:"user_id"`
	FriendID int    `json:"friend_id"`
	Online   bool   `json:"online"`
	Message  string `json:"message"`
}

type UserRelations struct {
	UserID  int   `json:"user_id"`
	Friends []int `json:"friends"`
}

type CounterSync struct {
	Mx    sync.RWMutex
	Store int
}

func (cs *CounterSync) Incr() {
	cs.Mx.Lock()
	defer cs.Mx.Unlock()
	cs.Store++
}
func (cs *CounterSync) Get() int {
	cs.Mx.RLock()
	defer cs.Mx.RUnlock()
	return cs.Store
}
func (cs *CounterSync) Reset() {
	cs.Mx.Lock()
	defer cs.Mx.Unlock()
	cs.Store = 0
}

func (cs *CounterSync) Decr() {
	cs.Mx.Lock()
	defer cs.Mx.Unlock()
	cs.Store--
}

type element struct {
	data interface{}
	next *element
}

type Stack struct {
	lock *sync.Mutex
	head *element
	Size int
}

func (stk *Stack) Push(data interface{}) {
	stk.lock.Lock()

	element := new(element)
	element.data = data
	temp := stk.head
	element.next = temp
	stk.head = element
	stk.Size++

	stk.lock.Unlock()
}

func (stk *Stack) Pop() interface{} {
	if stk.head == nil {
		return nil
	}
	stk.lock.Lock()
	r := stk.head.data
	stk.head = stk.head.next
	stk.Size--

	stk.lock.Unlock()

	return r
}

func NewStack() *Stack {
	stk := new(Stack)
	stk.lock = &sync.Mutex{}

	return stk
}
