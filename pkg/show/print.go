package show

import (
	"fmt"
)

// Define ANSI escape code colors
const (
	ResetColor = "\033[0m"
	DebugColor = "\033[1;34m" // Blue
	InfoColor  = "\033[1;32m" // Green
	WarnColor  = "\033[1;33m" // Yellow
	ErrorColor = "\033[1;31m" // Red
)

// Define icons for log levels
const (
	DebugIcon = "🐞"
	InfoIcon  = "ℹ️"
	WarnIcon  = "⚠️"
	ErrorIcon = "❌ "
)

// Output log with color and icon
func logWithColor(level, color, icon, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Printf("%s%s %s%s %s\n", color, icon, level, ResetColor, message)
}

// Debug outputs debug information
func Debug(format string, a ...interface{}) {
	logWithColor("DEBUG", DebugColor, DebugIcon, format, a...)
}

// Info outputs informational message
func Info(format string, a ...interface{}) {
	logWithColor("INFO", InfoColor, InfoIcon, format, a...)
}

// Warn outputs warning message
func Warn(format string, a ...interface{}) {
	logWithColor("WARN", WarnColor, WarnIcon, format, a...)
}

// Error outputs error message
func Error(format string, a ...interface{}) {
	logWithColor("ERROR", ErrorColor, ErrorIcon, format, a...)
}

func EmptyLine() {
	fmt.Println()
}
