package src

import (
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
	"sync"
)

type handler struct {
	mappingPacks map[string]Packet
}

var _handler *handler

func HandlerInitialize() {
	_logger.Info("====================================")
	_logger.Info("Handler Initialize Start!!!!")
	_logger.Info("====================================\n")

	_handler = new(handler)
	_handler.mappingPacks = make(map[string]Packet)

	_handler.mappingPacks["TestProcess"] = Packet{Process: TestProcess} //호출 테스트용

	_logger.Info("====================================")
	_logger.Info("Handler Initialize Finish!!!!")
	_logger.Info("====================================\n")
}

func RequestHandle(res http.ResponseWriter, req *http.Request) {

	//필수 파라미터인 cmd와 guildNo를 검사한다
	param := newParameters(req)
	param.PrintDebug()

	cmd := param.Get("cmd")
	groupNo := param.GetLong("groupNo")
	userNo := param.GetLong("userNo")

	//여기서 핸들러 호출 도중에 일어나는 모든 패닉을 잡고 응답 처리를 하도록 하자
	defer _handler.panicCatch(res)

	//매핑되어 있는 핸들러 호출인가?
	mappingPacket, ok := _handler.mappingPacks[cmd]
	if !ok {
		ret := _result.Fail(UNSUPPORTED_HANDLER)
		panic(ret)
	}

	//새로운 작업을 큐에 등록한다
	slowCheck := _util.NewSlowChecker(_config.GetLong("SlowCheck", 1000), userNo, groupNo)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	work := WorkRequest{GroupNo: groupNo, Req: &param, Res: res, Packet: mappingPacket, SlowCheck: slowCheck, WaitGroup: wg, UserNo: userNo}
	_worker.WorkQueue <- &work

	wg.Wait()
	if work.GroupNo == groupNo {
		_, _ = res.Write(work.PacketResultBin)
	} else {
		_logger.Errorf("[Write] not same response!!!!!")
		return
	}
}

func (h *handler) panicCatch(res http.ResponseWriter) {
	if r := recover(); r != nil {

		var response []byte
		switch r.(type) {
		case PacketResult:
			response, _ = ffjson.Marshal(&r)

		case int16:
			_logger.Errorf("[Handler]Exception:%s", r)

			ret := _result.Fail(UNKNOWN_SYSTEM_ERROR)
			response, _ = ffjson.Marshal(&ret)
		}

		_, _ = res.Write(response)
	}
}
