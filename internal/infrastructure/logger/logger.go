package logger

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type LogType string

const (
	System  LogType = "system"
	Warning LogType = "warning"
	Success LogType = "success"
	Error   LogType = "error"
)

type Debug struct {
	LogType LogType
}

func New(logType ...LogType) *Debug {
	lt := System
	if len(logType) > 0 && logType[0] != "" {
		lt = logType[0]
	}
	return &Debug{LogType: lt}
}

func (d *Debug) Print(message any) {
	var emoji string
	var colorFunc func(format string, a ...any) string

	switch strings.ToLower(string(d.LogType)) {
	case "warning":
		emoji = "⚠️"
		colorFunc = color.YellowString
	case "success":
		emoji = "✅"
		colorFunc = color.GreenString
	case "error":
		emoji = "❌"
		colorFunc = color.RedString
	case "system":
		fallthrough
	default:
		d.LogType = System
		emoji = "🔧"
		colorFunc = color.BlueString
	}

	msgStr := fmt.Sprintf("%v", message)
	fmt.Printf("[%s | %s] ==> %s\n", emoji, strings.ToUpper(string(d.LogType)), colorFunc(msgStr))
}

func Printf(logType LogType, format string, a ...any) {
	d := New(logType)
	d.Print(fmt.Sprintf(format, a...))
}
