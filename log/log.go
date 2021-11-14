package log

import (
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc/grpclog"
)

const (
	LogLevelFlag         = "log-level"
	LogLevelDefaultValue = "warning"
	LogLevelFlagUsage    = "Sets the logging level (debug, info, warning, error, fatal, panic)"
)

func InitLogs(c *cli.Context, output io.Writer) {
	formatter := log.TextFormatter{
		FullTimestamp:          true,
		DisableTimestamp:       false,
		DisableSorting:         true,
		DisableLevelTruncation: true,
		QuoteEmptyFields:       true,
	}

	log.SetFormatter(&formatter)

	log.SetReportCaller(true)

	log.SetOutput(output)

	// use cmd level if provided
	logLevel := c.String(LogLevelFlag)

	// Only logs with this severity or above will be issued.
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Errorf("invalid log level, setting to be warning: %v", err)
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(level)
	}
	setGrpcLogs(level)
}

// Set log level to Trace to see grpc logs
// We do this so that they don't mess with our stdio
func setGrpcLogs(level log.Level) {
	if level < log.TraceLevel {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, ioutil.Discard))
	} else {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stdout, os.Stdout, os.Stderr, 2))
	}
}
