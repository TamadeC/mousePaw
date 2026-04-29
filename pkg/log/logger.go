package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogEntry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type Logger struct {
	mu       sync.RWMutex
	entries  []LogEntry
	maxSize  int
	file     *os.File
	onAppend func(LogEntry)
}

func NewLogger() *Logger {
	l := &Logger{
		entries: make([]LogEntry, 0),
		maxSize: 100,
	}
	l.initFile()
	return l
}

func (l *Logger) initFile() {
	exePath, err := os.Executable()
	if err != nil {
		return
	}
	logPath := filepath.Join(filepath.Dir(exePath), "mousepaw.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	l.file = f
}

func (l *Logger) SetOnAppend(fn func(LogEntry)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.onAppend = fn
}

func (l *Logger) log(level, msg string) {
	entry := LogEntry{
		Time:    time.Now().Format("15:04:05"),
		Level:   level,
		Message: msg,
	}

	l.mu.Lock()
	l.entries = append(l.entries, entry)
	if len(l.entries) > l.maxSize {
		l.entries = l.entries[len(l.entries)-l.maxSize:]
	}
	onAppend := l.onAppend
	l.mu.Unlock()

	if l.file != nil {
		fmt.Fprintf(l.file, "[%s] [%s] %s\n", entry.Time, entry.Level, entry.Message)
	}

	if onAppend != nil {
		onAppend(entry)
	}
}

func (l *Logger) Info(msg string) {
	l.log("INFO", msg)
}

func (l *Logger) Error(msg string) {
	l.log("ERROR", msg)
}

func (l *Logger) GetEntries() []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	result := make([]LogEntry, len(l.entries))
	copy(result, l.entries)
	return result
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
