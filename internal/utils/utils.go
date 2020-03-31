package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(false)

	switch fmt := os.Getenv("LOG_FORMAT"); fmt {
	case "text":
		setTextFormat()
	case "json":
		setJSONLogFormat()
	default:
		log.Printf("unknown LOG_FORMAT value: '%s'", fmt)
		setJSONLogFormat()
		// setTextFormat()
	}

	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = log.InfoLevel
		// log.WithError(err).Printf("defaulted to %s", level.String())
	}
	log.SetLevel(level)
	// log.Printf("log level set to %s", log.GetLevel().String())
}

func setTextFormat() {
	log.SetFormatter(&log.TextFormatter{
		// DisableColors: false,
		ForceColors:   true,
		FullTimestamp: true,
	})

}
func setJSONLogFormat() {
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint:      false,
		DisableTimestamp: false,
	})
	// log.Info("set JSON log format")
}

// GetLogger needs to be called once to ensure logrus is configured correctly.
func GetLogger() *log.Logger {
	return log.StandardLogger()
}
