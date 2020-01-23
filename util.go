package src

import (
	seelog "github.com/cihub/seelog"
	"time"
)

type util struct {
}

var _util *util = new(util)

type Parameters struct {
	data map[string]string
}

type SlowTimer struct {
	Lengthy int64
	Before  int64
	IsSlow  bool
	UserNo  int64
	GroupNo int64
}

func (u util) NewSlowChecker(lengthy, userNo, groupNo int64) SlowTimer {
	if 0 == lengthy {
		lengthy = 1000
	}

	return SlowTimer{
		Lengthy: lengthy,
		Before:  _util.Now(),
		IsSlow:  false,
		UserNo:  userNo,
		GroupNo: groupNo,
	}
}

func (s *SlowTimer) Check(log string) {
	now := _util.Now()

	before := s.Before
	delay := now - before

	if delay > s.Lengthy {
		seelog.Info("[SLOW] >>>>>> [%s] delay=%d, groupNo=%d, userNo=%d // start=%d, end=%d",
			log, delay, s.GroupNo, s.UserNo, before, now)
		s.IsSlow = true
	}

	s.Before = now
}

func (u util) Now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
