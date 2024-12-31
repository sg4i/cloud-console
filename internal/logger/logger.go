package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
	Log.SetLevel(logrus.InfoLevel)
}

// Info logs a message at level Info
func Info(msg string, fields ...interface{}) {
	Log.WithFields(makeFields(fields...)).Info(msg)
}

// makeFields converts a slice of interfaces to logrus.Fields
func makeFields(fields ...interface{}) logrus.Fields {
	f := make(logrus.Fields)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			f[fields[i].(string)] = fields[i+1]
		}
	}
	return f
}
