package logging

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strconv"
	"strings"
	"time"
)

func NewDbLogger(writer logger.Writer, config logger.Config) logger.Interface {
	return &DbLogger{
		Writer: writer,
		Config: config,
		str:    "{\"level\":\"%s\", \"message\":\"%s\", \"file\":\"%s\", \"line\":%d, \"error\":\"%s\"}",
	}
}

type DbLogger struct {
	logger.Writer
	logger.Config
	str string
}

func (d DbLogger) LogMode(level logger.LogLevel) logger.Interface {
	d.LogLevel = level
	return d
}

func (d DbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if d.LogLevel >= logger.Info {
		file, line := splitFileLine(utils.FileWithLineNum())
		d.Printf(d.str, "info", fmt.Sprintf("%s\\n%v", msg, data), file, line, "")
	}
}

func (d DbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if d.LogLevel >= logger.Warn {
		file, line := splitFileLine(utils.FileWithLineNum())
		d.Printf(d.str, "warning", fmt.Sprintf("%s\\n%v", msg, data), file, line, msg)
	}
}

func (d DbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if d.LogLevel >= logger.Error {
		file, line := splitFileLine(utils.FileWithLineNum())
		d.Printf(d.str, "error", fmt.Sprintf("%s\\n%v", msg, data), file, line, msg)
	}
}

func (d DbLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

	if d.LogLevel <= logger.Silent {
		return
	}

	file, line := splitFileLine(utils.FileWithLineNum())
	sql, rows := fc()
	if rows == -1 {
		rows = 0
	}

	elapsed := time.Since(begin)

	message := fmt.Sprintf("duration: %.3fms\\nrows: %v\\nsql:%s", float64(elapsed.Nanoseconds())/1e6, rows, escape(sql))
	switch {
	case err != nil && d.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !d.IgnoreRecordNotFoundError):
		d.Printf(d.str, "error", message, file, line, escape(err.Error()))
	case elapsed > d.SlowThreshold && d.SlowThreshold != 0 && d.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", d.SlowThreshold)
		d.Printf(d.str, "error", message, file, line, slowLog)
	case d.LogLevel == logger.Info:
		d.Printf(d.str, "info", message, file, line, "")
	}
}

func escape(s string) string {
	str := strconv.Quote(s)
	remFirst := str[1:]
	remLast := remFirst[:len(remFirst)-1]
	return remLast
}

func splitFileLine(fileWithLine string) (string, int) {
	index := strings.LastIndex(fileWithLine, ":")
	if index == -1 {
		return "", 0
	}
	file := fileWithLine[:index]
	line := fileWithLine[index+1:]
	lineR, err := strconv.Atoi(line)
	if err != nil {
		lineR = 0
	}
	return file, lineR
}
