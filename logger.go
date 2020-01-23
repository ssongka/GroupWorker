package src

import "github.com/cihub/seelog"

var _logger *logger

type logger struct {
	seelog seelog.LoggerInterface
}

func ServerLoggerInitialize(configFile string) {
	_logger = new(logger)

	logger, err := seelog.LoggerFromConfigAsFile(configFile)
	if nil != err {
		panic(err)
	}

	_ = seelog.ReplaceLogger(logger)
	_logger.seelog = logger
}

func (log *logger) Infof(format string, params ...interface{}) {
	log.seelog.Infof(format, params...)
}

func (log *logger) Info(format string) {
	log.seelog.Info(format)
}

func (log *logger) Debugf(format string, params ...interface{}) {
	log.seelog.Debugf(format, params...)
}

func (log *logger) Errorf(format string, params ...interface{}) {
	_ = log.seelog.Errorf(format, params...)
}
