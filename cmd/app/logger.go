package app

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func initLogger(_ *Config) *logrus.Logger {
	logWriter := io.Writer(os.Stderr)
	// if conf.Logging.Debug {
	// 	// setup logfile
	// 	f, err := setupLogFile(afs)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer f.Close()
	// 	log.Println("log file:", f.Name())
	// 	logWriter = io.MultiWriter(f, os.Stderr)
	// }
	logLevel := logrus.Level(logrus.DebugLevel)
	logger := logrus.StandardLogger()
	logger.SetOutput(logWriter)
	logger.Level = logLevel
	return logger
}
