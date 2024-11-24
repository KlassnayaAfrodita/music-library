package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.JSONFormatter{}) // Логирование в JSON-формате
	Log.SetLevel(logrus.DebugLevel)           // Уровень логирования
}
