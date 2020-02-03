package src

import (
	"bytes"
	"fmt"
	seelog "github.com/cihub/seelog"
	"net/http"
	"strconv"
	"time"
)

type util struct {
}

var _util *util = new(util)

func (u util) ToInt(s string) int {
	if len(s) == 0 {
		return 0
	}
	ret, err := strconv.Atoi(s)
	if nil != err {
		_ = seelog.Error(err.Error())
		return 0
	}

	return ret
}

func (u util) ToLong(s string) int64 {
	if len(s) == 0 {
		return 0
	}

	ret, err := strconv.ParseInt(s, 10, 64)
	if nil != err {
		_ = seelog.Error(err.Error())
		return 0
	}

	return ret
}

func (u util) LongToStr(n int64) string {
	return strconv.FormatInt(n, 10)
}

func (u util) IntToStr(n int) string {
	return strconv.Itoa(n)
}

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

func newParameters(req *http.Request) Parameters {
	_ = req.ParseForm()
	p := Parameters{}
	p.data = make(map[string]string)

	for k := range req.Form {
		v := req.FormValue(k)
		p.data[k] = v
	}

	return p
}

func (p *Parameters) PrintDebug() {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\n\n>>>>>> REQUEST[%s][gNo:%s][uID:%s] ", p.data["cmd"], p.data["groupNo"], p.data["userNo"]))
	for k, v := range p.data {
		buf.WriteString(fmt.Sprintf("%s=%s,", k, v))
	}

	seelog.Debug(buf.String())
}

func (p *Parameters) GetParamDebugString() string {
	var buf bytes.Buffer
	for k, v := range p.data {
		buf.WriteString(fmt.Sprintf("%s=%s", k, v))
	}

	return buf.String()
}

func (p *Parameters) Get(key string) string {
	return p.data[key]
}

func (p *Parameters) GetInt(key string) int {
	return _util.ToInt(p.data[key])
}

func (p *Parameters) GetLong(key string) int64 {
	return _util.ToLong(p.data[key])
}
