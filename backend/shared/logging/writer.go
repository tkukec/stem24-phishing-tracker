package logging

import (
	"encoding/json"
	"fmt"
	"github.com/vovailchenko/go-zerolog-gelf/gelf"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	levelMap = map[string]int32{
		"debug":    gelf.LOG_DEBUG,
		"info":     gelf.LOG_INFO,
		"notice":   gelf.LOG_NOTICE,
		"warning":  gelf.LOG_WARNING,
		"error":    gelf.LOG_ERR,
		"critical": gelf.LOG_CRIT,
	}
)

type Writer struct {
	graylogWriter *gelf.Writer
	serviceName   string
	hostname      string
}

func NewWriter(graylogHost, graylogPort, serviceName string) (*Writer, error) {
	graylogAddr := fmt.Sprintf("%s:%s", graylogHost, graylogPort)
	gelfWriter, err := gelf.NewWriter(graylogAddr)
	if err != nil {
		return nil, err
	}
	hostname, _ := os.Hostname()
	return &Writer{
		graylogWriter: gelfWriter,
		serviceName:   serviceName,
		hostname:      hostname,
	}, nil
}

type IncomingMessage struct {
	XCorrelationID  string `json:"X-Correlation-ID"`
	XTenantID       string `json:"X-TENANT-ID"`
	UserId          string `json:"user_id"`
	UserDisplayName string `json:"display_name"`
	Message         string `json:"message"`
	Method          string `json:"method"`
	LevelCode       string `json:"level"`
	ErrorMessage    string `json:"error"`
	File            string `json:"file"`
	Line            int    `json:"line"`
}

func NewIncomingMessage(data []byte) (*IncomingMessage, error) {
	var incomingMessage *IncomingMessage
	err := json.Unmarshal(data, &incomingMessage)
	if err != nil {
		return nil, err
	}

	if incomingMessage.XCorrelationID == "" {
		var msg map[string]map[string]interface{}
		json.Unmarshal([]byte(incomingMessage.Message), &msg)
		if xCorrelationID, ok := msg["headers"]["X-Correlation-ID"]; ok {
			incomingMessage.XCorrelationID = xCorrelationID.(string)
		}
	}

	// convert message to new line format
	var msg map[string]interface{}
	err = json.Unmarshal([]byte(incomingMessage.Message), &msg)
	if err == nil {
		var multiLinerMsgString string
		for key, value := range msg {
			jsonValue, _ := json.Marshal(value)
			multiLinerMsgString += fmt.Sprintf("%s: %s\n", key, string(jsonValue))
		}
		incomingMessage.Message = multiLinerMsgString
	}

	return incomingMessage, nil
}

func NewIncomingMessageFromRecovery(data []byte) (*IncomingMessage, error) {
	dataString := string(data)
	if !strings.Contains(dataString, "Recovery") {
		return nil, fmt.Errorf("not a recovery message")
	}
	errorMessage := getErrorMessage(dataString)
	stackTrace := getStackTrace(dataString)
	return &IncomingMessage{
		Message:      stackTrace,
		LevelCode:    "error",
		ErrorMessage: errorMessage,
	}, nil
}

func (w *Writer) Printf(format string, v ...interface{}) {
	w.Write([]byte(fmt.Sprintf(format, v...)))
}

func (w *Writer) Write(p []byte) (n int, err error) {
	// UnixMillis produces int64 and graylog expects float
	// graylog thinks they are seconds, which is way into the future
	// and does not log the message
	floatTime := float64(time.Now().Unix())
	stringMillis := fmt.Sprintf("%d", time.Now().UnixMilli())
	if stringMillis != "" {
		s := stringMillis[:len(stringMillis)-3] + "." + stringMillis[len(stringMillis)-3:]
		if s != "" {
			floatTime, _ = strconv.ParseFloat(s, 64)
		}
	}

	m := gelf.Message{
		Version:  "1.1",
		Host:     w.serviceName,
		Full:     string(p),
		TimeUnix: floatTime,
		Extra: map[string]interface{}{
			"_hostname": w.hostname,
		},
	}

	// parse the message, get correlation id, etc.
	var recoveryError error
	incomingMessage, err := NewIncomingMessage(p)
	if err != nil {
		incomingMessage, recoveryError = NewIncomingMessageFromRecovery(p)
		if err != nil {
			m.Extra["_error"] = fmt.Sprintf("%s | %s", err.Error(), recoveryError.Error())
			m.Level = gelf.LOG_DEBUG
		}
	}

	if incomingMessage != nil {
		file := incomingMessage.File
		line := incomingMessage.Line
		if file == "" && line <= 0 {
			file, line = getCallerIgnoringLogMulti(1)
		}

		// resolve level
		level, ok := levelMap[incomingMessage.LevelCode]
		if !ok {
			level = gelf.LOG_DEBUG
		}

		m.Short = incomingMessage.Message
		m.Level = level
		m.Extra = map[string]interface{}{
			"_file":              file,
			"_line":              line,
			"_x_correlation_id":  incomingMessage.XCorrelationID,
			"_x_tenant_id":       incomingMessage.XTenantID,
			"_user_id":           incomingMessage.UserId,
			"_user_display_name": incomingMessage.UserDisplayName,
			"_method":            incomingMessage.Method,
			"_error":             incomingMessage.ErrorMessage,
			"_hostname":          w.hostname,
		}
	}

	if err = w.graylogWriter.WriteMessage(&m); err != nil {
		return 0, err
	}

	return len(p), nil
}

func getStackTrace(recoveryMessage string) string {
	beginIndex := strings.LastIndex(recoveryMessage, "\n\n")
	if beginIndex > -1 {
		s := recoveryMessage[beginIndex+2:]
		endIndex := strings.Index(s, "\n")
		stackTrace := s[endIndex:]
		stackTrace = strings.ReplaceAll(stackTrace, "\033[0m", "")
		return stackTrace
	}
	return ""
}

func getErrorMessage(recoveryMessage string) string {
	beginIndex := strings.LastIndex(recoveryMessage, "\n\n")
	if beginIndex > -1 {
		s := recoveryMessage[beginIndex+2:]
		endIndex := strings.Index(s, "\n")
		errMsg := s[:endIndex]
		return errMsg
	}
	return ""
}

func getCallerIgnoringLogMulti(callDepth int) (string, int) {
	// the +1 is to ignore this (getCallerIgnoringLogMulti) frame
	return getCaller(callDepth+1, "/src/log/log.go", "/src/io/multi.go", "/writer.go", "/event.go")
}

func getCaller(callDepth int, suffixesToIgnore ...string) (file string, line int) {
	// bump by 1 to ignore the getCaller (this) stackframe
	callDepth++
outer:
	for {
		var ok bool
		_, file, line, ok = runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
			break
		}

		for _, s := range suffixesToIgnore {
			if strings.HasSuffix(file, s) {
				callDepth++
				continue outer
			}
		}
		break
	}
	return
}
