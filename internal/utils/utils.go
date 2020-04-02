package utils

import (
	"math/rand"
	"os"
	"strings"
	"text/scanner"

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
		// log.Printf("unknown LOG_FORMAT value: '%s'", fmt)
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

func ParseCommaSeparatedFiles(commaSeparatedFilenames string) []string {
	var stringSlice []string
	stringSlice = strings.Split(commaSeparatedFilenames, ",")
	return stringSlice
}

func ParseCommaSeparatedStrings(commaSeparatedStrings string) []string {
	var s scanner.Scanner
	s.Init(strings.NewReader(commaSeparatedStrings))
	s.Whitespace = 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' ' | 1<<','
	stringSlice := []string{}

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		stringSlice = append(stringSlice, s.TokenText())
	}
	return stringSlice
}

func SelectRandomString(stringValues []string) string {
	sLen := len(stringValues)
	if sLen <= 0 {
		return ""
	}
	return stringValues[rand.Intn(sLen)]
}
