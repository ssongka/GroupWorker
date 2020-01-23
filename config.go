package src

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type config struct {
	configMap map[string]string
}

var _config = new(config)

func ConfigInitialize(configFile string) bool {
	_logger.Info("====================================")
	_logger.Info("Configuration Initialize Start!!!!")
	_logger.Info("=====================================")

	_config.configMap = make(map[string]string)

	_logger.Infof("[Config] configFile load %s", configFile)

	file, err := os.Open(configFile)
	if nil != err {
		fmt.Println(err)
		return false
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		r := scanner.Text()

		if i := strings.Index(r, ";"); i == 0 {
			continue
		}
		if i := strings.Index(r, "#"); i == 0 {
			continue
		}

		if strings.Contains(r, "=") {
			idx := strings.Index(r, "=")
			key := r[0:idx]
			value := r[idx+1:]

			_logger.Infof("load config... %s=%s", key, value)
			_config.configMap[key] = value
		}
	}

	_logger.Info("====================================")
	_logger.Info("Configuration Initialize Finish!!!!")
	_logger.Info("=====================================")
	return true
}

func (c *config) Get(k string, defVal string) string {
	v, ok := c.configMap[k]
	if !ok {
		return defVal
	}

	return v
}

func (c *config) GetInt(k string, defVal int) int {
	v, _ := strconv.Atoi(c.Get(k, strconv.Itoa(defVal)))
	return v
}

func (c *config) GetLong(k string, defVal int) int64 {
	return int64(c.GetInt(k, defVal))
}

func (c *config) IsDebuggable() bool {
	serverLocation := c.Get("server_location", "dev_")
	if serverLocation != "live_" {
		return true
	} else {
		return false
	}
}
