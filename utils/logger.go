package utils

import "github.com/sirupsen/logrus"

var logger *logrus.Logger

func InitLogger() *logrus.Logger {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

func LogError(err error) {
	if err != nil {
		logger.Error(err)
	}
}
