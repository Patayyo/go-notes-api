package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type Fields = logrus.Fields

func Init() {
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
