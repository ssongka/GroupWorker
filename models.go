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
	Packet
}

type GroupInfo struct {
	GroupNo      int64
	GroupMembers map[int64]*User
}

type User struct {
	UserNo   int64
	UserName string
}

type Packet struct {
	Process WorkRequestProcess
}

type WorkRequestProcess func(req *Parameters, groupInfo *GroupInfo, work *WorkRequest) PacketResult

type PacketResult struct {
	ErrorCode int16       `json:"error_code"`
	Detail    interface{} `json:"detail"`
}

type result struct {
}

var _result *result = new(result)

//성공 처리
func (r *result) Success(res interface{}) PacketResult {
	create := r.create(OK)
	create.Detail = res

	return create
}

//실패 처리
func (r *result) Fail(errorCode int16) PacketResult {
	return r.create(errorCode)
}

func (r *result) create(code int16) PacketResult {
	var packetResult = PacketResult{}
	packetResult.ErrorCode = code
	return packetResult
}

//응답 코드 정의..
const (
	OK                   int16 = 0
	UNSUPPORTED_HANDLER  int16 = -1
	UNKNOWN_SYSTEM_ERROR int16 = -2
)
