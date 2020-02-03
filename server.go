package src

import (
	"net/http"
	"runtime"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	InitializeServer(wg)

	_logger.Info("ServerInitialize well done..")
	wg.Wait()

	http.HandleFunc("/GroupWork")
}

//some component init here..
func InitializeServer(wg *sync.WaitGroup) {

	ServerLoggerInitialize("seelog.xml")

	ConfigInitialize("config.ini")

	HandlerInitialize()

	WorkerInitialize()

	runtime.GOMAXPROCS(runtime.NumCPU())

	//초기화 완료 처리
	wg.Done()
}
