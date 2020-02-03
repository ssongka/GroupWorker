package src

var _worker *workerContext

func WorkerInitialize() {
	_worker = new(workerContext)
	_worker.WorkerQueueMap = make(map[int64]chan chan *WorkRequest, 30)
	_worker.WorkQueue = make(chan *WorkRequest, 1)
	_worker.StartWorkDispatch()
}

func (w *workerContext) NewWorker(GroupNo int64, workerQueue chan chan *WorkRequest) Worker {
	worker := Worker{
		GroupNo:     GroupNo,
		Work:        make(chan *WorkRequest),
		WorkerQueue: workerQueue,
		ForceStop:   make(chan bool),
	}

	return worker
}

func DestroyWorker(GroupNo int64) {

	_, ok := _worker.WorkerQueueMap[GroupNo]
	if ok {
		_logger.Infof("[Worker ]>>> Stop Job Queue GroupNo=%d <<<  it's will be Destroy", GroupNo)
		delete(_worker.WorkerQueueMap, GroupNo)
	}
}

//작업큐 시작
func (w *Worker) Start() {

	//그룹 정보 초기화
	slowChecker := _util.NewSlowChecker(_config.GetLong("SlowCheck", 1000), 0, w.GroupNo)
	groupInfo, _ := InitializeGroupInfo(w.GroupNo)

	slowChecker.Check("InitializeGroupInfo")

	//초기화에 성공하면 워커에 정보를 넣는다
	w.GroupInfo = groupInfo

	go func() {
		for {
			//자기 자신 큐에 워커를 추가해주고..
			w.WorkerQueue <- w.Work
			select {
			//작업 요청을 받는다..
			case work := <-w.Work:

				work.SlowCheck.Check("work := <-w.Work")
				_logger.Debugf("[Worker] receive work packet.. GroupNo=%d", work.GroupNo)

				//todo 핸들러 호출
				//_handler.executeProcess(work, w)

				//예외 상황이 발생하여 강제 설정값이 활성화 되었다면 채널을 닫는다
				if work.ForceStop {
					DestroyWorker(work.GroupNo)
					return
				}
			//작업 중지
			case <-w.ForceStop:
				_logger.Infof("[Worker] >>> Stop Job Queue GroupNo=%d <<< ", w.GroupNo)
				return
			}
		}
	}()
}

/**
그룹 큐의 정보를 초기화 후 반환한다
*/
func InitializeGroupInfo(groupNo int64) (*GroupInfo, interface{}) {
	return &GroupInfo{
		GroupNo:      groupNo,
		GroupMembers: make(map[int64]*User),
	}, nil
}

//작업큐 중지
func (w *Worker) Stop() {
	_logger.Infof("[Worker] Stop Job Queue it's will be destroy.. GroupNo=%d", w.GroupNo)
	w.ForceStop <- true

	DestroyWorker(w.GroupNo)
}

func (w *workerContext) StartWorkDispatch() {
	_logger.Info("====================================")
	_logger.Info("WorkDispatch Start!!!!")
	_logger.Info("====================================\n")

	go func() {
		for {
			select {
			case work := <-w.WorkQueue:
				_logger.Debugf("[WorkQueue] work requeust.. GroupNo=%d", work.GroupNo)
				workerQueue, ok := w.WorkerQueueMap[work.GroupNo]
				if !ok {
					_logger.Infof("[WorkQueue] Generate New WorkerQueue.. GroupNo=%d", work.GroupNo)

					workerQueue = make(chan chan *WorkRequest)
					w.WorkerQueueMap[work.GroupNo] = workerQueue

					worker := w.NewWorker(work.GroupNo, workerQueue)
					worker.Start()
				}

				go func() {
					worker := <-workerQueue
					_logger.Debugf("[WorkQueue] Dispatching work request.. GroupNo=%d", work.GroupNo)
					worker <- work
				}()
			}
		}
	}()
}
