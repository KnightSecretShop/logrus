package logrus

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

type Entry struct {
	Logger *Logger
	Data   Fields
}

var baseTimestamp time.Time

func init() {
	baseTimestamp = time.Now()
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is three fields, give a little extra room
		Data: make(Fields, 5),
	}
}

func (entry *Entry) Reader() (*bytes.Buffer, error) {
	serialized, err := entry.Logger.Formatter.Format(entry)
	return bytes.NewBuffer(serialized), err
}

func (entry *Entry) String() (string, error) {
	reader, err := entry.Reader()
	if err != nil {
		return "", err
	}

	return reader.String(), err
}

func (entry *Entry) WithField(key string, value interface{}) *Entry {
	entry.Data[key] = value
	return entry
}

func (entry *Entry) WithFields(fields Fields) *Entry {
	for key, value := range fields {
		entry.WithField(key, value)
	}
	return entry
}

func (entry *Entry) log(level string, levelInt Level, msg string) string {
	entry.Data["time"] = time.Now().String()
	entry.Data["level"] = level
	entry.Data["msg"] = msg

	if err := entry.Logger.Hooks.Fire(levelInt, entry); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fire hook", err)
	}

	reader, err := entry.Reader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v", err)
	}

	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()

	_, err = io.Copy(entry.Logger.Out, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v", err)
	}

	return reader.String()
}

func (entry *Entry) Debug(args ...interface{}) {
	if entry.Logger.Level >= Debug {
		entry.log("debug", Debug, fmt.Sprint(args...))
		entry.Logger.Hooks.Fire(Debug, entry)
	}
}

func (entry *Entry) Print(args ...interface{}) {
	entry.Info(args...)
}

func (entry *Entry) Info(args ...interface{}) {
	if entry.Logger.Level >= Info {
		entry.log("info", Info, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warn(args ...interface{}) {
	if entry.Logger.Level >= Warn {
		entry.log("warning", Warn, fmt.Sprint(args...))
	}
}

func (entry *Entry) Error(args ...interface{}) {
	if entry.Logger.Level >= Error {
		entry.log("error", Error, fmt.Sprint(args...))
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Logger.Level >= Fatal {
		entry.log("fatal", Fatal, fmt.Sprint(args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	if entry.Logger.Level >= Panic {
		msg := entry.log("panic", Panic, fmt.Sprint(args...))
		panic(msg)
	}
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Logger.Level >= Debug {
		entry.Debug(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Logger.Level >= Info {
		entry.Info(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Logger.Level >= Warn {
		entry.Warn(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Logger.Level >= Error {
		entry.Error(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.Level >= Fatal {
		entry.Fatal(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Logger.Level >= Panic {
		entry.Panic(fmt.Sprintf(format, args...))
	}
}

// Entry Println family functions

func (entry *Entry) Debugln(args ...interface{}) {
	entry.Debug(args...)
}

func (entry *Entry) Infoln(args ...interface{}) {
	entry.Info(args...)
}

func (entry *Entry) Println(args ...interface{}) {
	entry.Info(args...)
}

func (entry *Entry) Warnln(args ...interface{}) {
	entry.Warn(args...)
}

func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warn(args...)
}

func (entry *Entry) Errorln(args ...interface{}) {
	entry.Error(args...)
}

func (entry *Entry) Fatalln(args ...interface{}) {
	entry.Fatal(args...)
}

func (entry *Entry) Panicln(args ...interface{}) {
	entry.Panic(args...)
}
