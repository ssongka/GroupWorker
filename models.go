package src

import (
	"net/http"
	"sync"
)

type workerContext struct {
	WorkerQueueMap map[int64]chan chan *WorkRequest
	WorkQueue      chan *WorkRequest
}

type Worker struct {
	GroupNo        int64
	GroupInfo      *GroupInfo
	Work           chan *WorkRequest
	WorkerQueue    chan chan *WorkRequest
	ForceStop      chan bool
	LastAccessTime int64
}

type WorkRequest struct {
	GroupNo         int64
	ForceStop       bool
	Req             *Parameters
	Res             http.ResponseWriter
	PacketResultBin []byte
	SlowCheck       SlowTimer //패킷 슬로우 체크용
	UserNo          int64     //유저번호
	WaitGroup       *sync.WaitGroup
	//Packet
}

type GroupInfo struct {
	GroupNo      int64
	GroupMembers map[int64]*User
}

type User struct {
	UserNo   int64
	UserName string
}
