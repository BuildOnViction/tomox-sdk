package utils

import (
	"io"
	"os"

	"github.com/op/go-logging"
)

const (
	LogPrefix = ""
)

var Logger = logging.MustGetLogger("main")
var StdoutLogger = logging.MustGetLogger("main")

func InitLogger(logLevel string) {
	level, err := logging.LogLevel(logLevel)
	if err != nil {
		level = logging.ERROR
	}

	var format = logging.MustStringFormatter(
		`%{level:.4s} %{time:15:04:05} at %{shortpkg}/%{shortfile} in %{shortfunc}():%{message}`,
	)

	writer := io.MultiWriter(os.Stdout)
	backend := logging.NewLogBackend(writer, LogPrefix, 0)

	formattedBackend := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(formattedBackend)

	leveledBackend.SetLevel(level, "")

	logging.SetBackend(leveledBackend)
}
