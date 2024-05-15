package logging

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/runtimebag"
	"github.com/rs/zerolog/diode"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type LoggerWriters []*LoggerWriter

func (l LoggerWriters) Graylog() io.Writer {
	for _, w := range l {
		if w.loggerType == constants.LogGraylog {
			return w.Writer
		}
	}
	return nil
}

func (l LoggerWriters) File() io.Writer {
	for _, w := range l {
		if w.loggerType == constants.LogFile {
			return w
		}
	}
	return nil
}

func (l LoggerWriters) Stdout() io.Writer {
	for _, w := range l {
		if w.loggerType == constants.LogStdOut {
			return w
		}
	}
	return nil
}

func (l LoggerWriters) Writers() []io.Writer {
	writers := make([]io.Writer, len(l))
	for i, w := range l {
		writers[i] = w.Writer
	}
	return writers
}

func (l LoggerWriters) DiodeWriters() []io.Writer {
	writers := make([]io.Writer, len(l))
	for i, w := range l {
		writers[i] = diode.NewWriter(w.Writer, 1000, 10*time.Millisecond, func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})
	}
	return writers
}

func (l LoggerWriters) Printf(format string, v ...interface{}) {
	for _, w := range l {
		w.Printf(format, v...)
	}
}

type LoggerWriter struct {
	io.Writer
	loggerType string
}

func (l LoggerWriter) Printf(format string, v ...interface{}) {
	l.Write([]byte(fmt.Sprintf(format, v...)))
}

func GreyLogLoggerWriter(serviceName string) (*LoggerWriter, error) {
	writer, err := NewWriter(
		runtimebag.GetEnvString(constants.GraylogHostname, "live-graylog"),
		runtimebag.GetEnvString(constants.GraylogPort, "12201"),
		serviceName)
	if err != nil {
		return nil, err
	}

	return &LoggerWriter{
		Writer:     writer,
		loggerType: constants.LogGraylog,
	}, nil
}

func GetLogWriters(filename string, maxBackups int, maxSize int, maxAge int, serviceName string) LoggerWriters {
	var writers LoggerWriters
	for _, driver := range strings.Split(runtimebag.GetEnvString(constants.LogDrivers, "file,stdout"), ",") {
		switch driver {
		case constants.LogGraylog:
			graylogWriter, err := GreyLogLoggerWriter(serviceName)
			if err != nil {
				log.Panic(err.Error())
			}
			writers = append(writers, graylogWriter)
			break

		case constants.LogStdOut:
			writers = append(writers, &LoggerWriter{
				Writer:     os.Stdout,
				loggerType: constants.LogStdOut,
			})
			break

		case constants.LogFile:
			writers = append(writers, &LoggerWriter{
				Writer: &lumberjack.Logger{
					Filename:   filename,
					MaxBackups: maxBackups, // files
					MaxSize:    maxSize,    // megabytes
					MaxAge:     maxAge,     // days
				},
				loggerType: constants.LogFile,
			})
			break
		}
	}

	return writers
}
