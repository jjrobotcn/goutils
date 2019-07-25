package logrusenv

/*
ENV:
    - level
    - LOG_LEVEL

import (
    _ "github.com/jjrobotcn/goutils/logrusenv"
)
*/

import (
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	parseAndSet("level")
	parseAndSet("LOG_LEVEL")
}

func parseAndSet(envKey string) {
	logLevel, err := logrus.ParseLevel(os.Getenv(envKey))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)
}
