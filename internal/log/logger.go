package log

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/fatih/color"
)

var logger = log.New(os.Stdout, "METEO-AIRPORT", log.LstdFlags)

var (
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	reset  = color.New(color.Reset).SprintFunc()
)

func logMessage(message, level string) string {
	_, file, line, ok := runtime.Caller(2)
	if ok && level == red(bold("[Error] ")) {
		return fmt.Sprintf("\n\t%s%s:%d %s%s", level, file, line, message, reset())
	}

	return fmt.Sprintf("\n\t%s%s%s", level, message, reset())
}

func Error(format string, v ...interface{}) {
	logger.Print(logMessage(red(fmt.Sprintf(format, v...)), red(bold("[Error] "))))
}

func Warn(format string, v ...interface{}) {
	logger.Print(logMessage(yellow(fmt.Sprintf(format, v...)), yellow(bold("[Warn] "))))
}

func Info(format string, v ...interface{}) {
	logger.Print(logMessage(green(fmt.Sprintf(format, v...)), green(bold("[Info] "))))
}
