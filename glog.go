// Package log implements structured logging for Google App Engine, Cloud Run
// and Cloud Functions. The API is compatible with the standard library "log" module.
//
// All the severities conform to the Google Cloud Logging API v2 as described in
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logseverity.
// These severity levels are: DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL, ALERT, EMERGENCY.
//
// The ERROR, CRITICAL, ALERT, EMERGENCY logs are written to the standard error stream, while
// the remaining logs are written to the standard output.
package log

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var std Logger

// ProjectID should be set to the Google Cloud project ID.
var ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

// Debug logs detailed information that could mainly be used to catch unforeseen problems.
// Arguments are handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	std.Debug(v...)
}

// Debugln logs detailed information that could mainly be used to catch unforeseen problems.
// Arguments are handled in the manner of fmt.Println.
func Debugln(v ...interface{}) {
	std.Debugln(v...)
}

// Debugf logs detailed information that could mainly be used to catch unforeseen problems.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

// Debugj logs detailed information that could mainly be used to catch unforeseen problems.
// Argument v becomes jsonPayload field in the log entry.
func Debugj(msg string, v interface{}) {
	std.Debugj(msg, v)
}

// Info logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	std.Info(v...)
}

// Infoln logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Println.
func Infoln(v ...interface{}) {
	std.Infoln(v...)
}

// Infof logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

// Infoj logs routine information, such as ongoing status or performance.
// Arguments become jsonPayload in the log entry.
func Infoj(msg string, v interface{}) {
	std.Infoj(msg, v)
}

// Notice logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Print.
func Notice(v ...interface{}) {
	std.Notice(v...)
}

// Noticeln logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Println.
func Noticeln(v ...interface{}) {
	std.Noticeln(v...)
}

// Noticef logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Printf.
func Noticef(format string, v ...interface{}) {
	std.Noticef(format, v...)
}

// Noticej logs normal but significant events, such as start up, shut down, or configuration.
// Arguments become jsonPayload in the log entry.
func Noticej(msg string, v interface{}) {
	std.Noticej(msg, v)
}

// Warning logs events that might cause problems.
// Arguments are handled in the manner of fmt.Print.
func Warning(v ...interface{}) {
	std.Warning(v...)
}

// Warningln logs events that might cause problems.
// Arguments are handled in the manner of fmt.Println.
func Warningln(v ...interface{}) {
	std.Warningln(v...)
}

// Warningf logs events that might cause problems.
// Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, v ...interface{}) {
	std.Warningf(format, v...)
}

// Warningj logs events that might cause problems.
// Arguments become jsonPayload in the log entry.
func Warningj(msg string, v interface{}) {
	std.Warningj(msg, v)
}

// Error logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	std.Error(v...)
}

// Errorln logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	std.Errorln(v...)
}

// Errorf logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

// Errorj logs events likely to cause problems.
// Arguments become jsonPayload in the log entry.
func Errorj(msg string, v interface{}) {
	std.Errorj(msg, v)
}

// Critical logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Print.
func Critical(v ...interface{}) {
	std.Critical(v...)
}

// Criticalln logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Println.
func Criticalln(v ...interface{}) {
	std.Criticalln(v...)
}

// Criticalf logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Printf.
func Criticalf(format string, v ...interface{}) {
	std.Criticalf(format, v...)
}

// Criticalj logs events that cause more severe problems or outages.
// Arguments become jsonPayload in the log entry.
func Criticalj(msg string, v interface{}) {
	std.Criticalj(msg, v)
}

// Print logs routine information, such as ongoing status or performance, same as Info().
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	std.Print(v...)
}

// Println logs routine information, such as ongoing status or performance, same as Infoln().
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	std.Println(v...)
}

// Printf logs routine information, such as ongoing status or performance, same as Infof().
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	std.Printf(format, v...)
}

// Printj logs routine information, such as ongoing status or performance, same as Infoj().
// Arguments become jsonPayload in the log entry.
func Printj(msg string, v interface{}) {
	std.Printj(msg, v)
}

// Fatal is equivalent to a call to Critical() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Fatal(v...)
}

// Fatalln is equivalent to a call to Criticalln() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	std.Fatalln(v...)
}

// Fatalf is equivalent to a call to Criticalf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

// Fatalj is equivalent to a call to Criticalj() followed by a call to os.Exit(1).
func Fatalj(msg string, v interface{}) {
	std.Fatalj(msg, v)
}

// Panic is equivalent to a call to Critical() followed by a call to panic().
func Panic(v ...interface{}) {
	std.Panic(v...)
}

// Panicln is equivalent to a call to Criticalln() followed by a call to panic().
func Panicln(v ...interface{}) {
	std.Panicln(v...)
}

// Panicf is equivalent to a call to Criticalf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	std.Panicf(format, v...)
}

// Panicj is equivalent to a call to Criticalj() followed by a call to panic().
func Panicj(msg string, v interface{}) {
	std.Panicj(msg, v)
}

// Alert logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Print.
func Alert(v ...interface{}) {
	std.Alert(v...)
}

// Alertln logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Println.
func Alertln(v ...interface{}) {
	std.Alertln(v...)
}

// Alertf logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Printf.
func Alertf(format string, v ...interface{}) {
	std.Alertf(format, v...)
}

// Alertj logs when a person must take an action immediately.
// Arguments become jsonPayload in the log entry.
func Alertj(msg string, v interface{}) {
	std.Alertj(msg, v)
}

// Emergency logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Print.
func Emergency(v ...interface{}) {
	std.Emergency(v...)
}

// Emergencyln logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Println.
func Emergencyln(v ...interface{}) {
	std.Emergencyln(v...)
}

// Emergencyf logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Printf.
func Emergencyf(format string, v ...interface{}) {
	std.Emergencyf(format, v...)
}

// Emergencyj logs when one or more systems are unusable.
// Arguments become jsonPayload in the log entry.
func Emergencyj(msg string, v interface{}) {
	std.Emergencyj(msg, v)
}

// Debug logs detailed information that could mainly be used to catch unforeseen problems.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Debug(v ...interface{}) {
	log(debugsev, l, v...)
}

// Debugln logs detailed information that could mainly be used to catch unforeseen problems.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Debugln(v ...interface{}) {
	logln(debugsev, l, v...)
}

// Debugf logs detailed information that could mainly be used to catch unforeseen problems.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	logf(debugsev, l, format, v...)
}

// Debugj logs detailed information that could mainly be used to catch unforeseen problems.
// Argument v becomes jsonPayload field in the log entry.
func (l *Logger) Debugj(msg string, v interface{}) {
	logj(debugsev, l, msg, v)
}

// Info logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	log(infosev, l, v...)
}

// Infoln logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Infoln(v ...interface{}) {
	logln(infosev, l, v...)
}

// Infof logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	logf(infosev, l, format, v...)
}

// Infoj logs routine information, such as ongoing status or performance.
// Arguments become jsonPayload in the log entry.
func (l *Logger) Infoj(msg string, v interface{}) {
	logj(infosev, l, msg, v)
}

// Notice logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Notice(v ...interface{}) {
	log(noticesev, l, v...)
}

// Noticeln logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Noticeln(v ...interface{}) {
	logln(noticesev, l, v...)
}

// Noticef logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Noticef(format string, v ...interface{}) {
	logf(noticesev, l, format, v...)
}

// Noticej logs normal but significant events, such as start up, shut down, or configuration.
// Arguments become jsonPayload in the log entry.
func (l *Logger) Noticej(msg string, v interface{}) {
	logj(noticesev, l, msg, v)
}

// Warning logs events that might cause problems.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Warning(v ...interface{}) {
	log(warningsev, l, v...)
}

// Warningln logs events that might cause problems.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Warningln(v ...interface{}) {
	logln(warningsev, l, v...)
}

// Warningf logs events that might cause problems.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warningf(format string, v ...interface{}) {
	logf(warningsev, l, format, v...)
}

// Warningj logs events that might cause problems.
// Arguments become jsonPayload in the log entry.
func (l *Logger) Warningj(msg string, v interface{}) {
	logj(warningsev, l, msg, v)
}

// Error logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	log(errorsev, l, v...)
}

// Errorln logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Errorln(v ...interface{}) {
	logln(errorsev, l, v...)
}

// Errorf logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	logf(errorsev, l, format, v...)
}

// Errorj logs events likely to cause problems.
// Arguments become jsonPayload in the log entry.
func (l *Logger) Errorj(msg string, v interface{}) {
	logj(errorsev, l, msg, v)
}

// Critical logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Critical(v ...interface{}) {
	log(criticalsev, l, v...)
}

// Criticalln logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Criticalln(v ...interface{}) {
	logln(criticalsev, l, v...)
}

// Criticalf logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Criticalf(format string, v ...interface{}) {
	logf(criticalsev, l, format, v...)
}

// Criticalj logs events that cause more severe problems or outages.
// Arguments become jsonPayload in the log entry.
func (l *Logger) Criticalj(msg string, v interface{}) {
	logj(criticalsev, l, msg, v)
}

// Alert logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Alert(v ...interface{}) {
	log(alertsev, l, v...)
}

// Alertln logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Alertln(v ...interface{}) {
	logln(alertsev, l, v...)
}

// Alertf logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Alertf(format string, v ...interface{}) {
	logf(alertsev, l, format, v...)
}

// Alertj logs when a person must take an action immediately.
// Arguments become jsonPayload in the log entry.
func (l *Logger) Alertj(msg string, v interface{}) {
	logj(alertsev, l, msg, v)
}

// Emergency logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Emergency(v ...interface{}) {
	log(emergencysev, l, v...)
}

// Emergencyln logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Emergencyln(v ...interface{}) {
	logln(emergencysev, l, v...)
}

// Emergencyf logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Emergencyf(format string, v ...interface{}) {
	logf(emergencysev, l, format, v...)
}

// Emergencyj logs when one or more systems are unusable.
// Arguments become jsonPayload in the log entry.
func (l *Logger) Emergencyj(msg string, v interface{}) {
	logj(emergencysev, l, msg, v)
}

// Print logs routine information, such as ongoing status or performance, same as l.Info().
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	log(infosev, l, v...)
}

// Println logs routine information, such as ongoing status or performance, same as l.Infoln().
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...interface{}) {
	logln(infosev, l, v...)
}

// Printf logs routine information, such as ongoing status or performance, same as l.Infof().
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	logf(infosev, l, format, v...)
}

// Printj logs routine information, such as ongoing status or performance, same as l.Infoj().
// Arguments become jsonPayload in the log entry.
func (l *Logger) Printj(msg string, v interface{}) {
	logj(infosev, l, msg, v)
}

// Fatal is equivalent to a call to l.Critical() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	log(criticalsev, l, v...)
	os.Exit(1)
}

// Fatalln is equivalent to a call to l.Criticalln() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	logln(criticalsev, l, v...)
	os.Exit(1)
}

// Fatalf is equivalent to a call to l.Criticalf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	logf(criticalsev, l, format, v...)
	os.Exit(1)
}

// Fatalj is equivalent to a call to l.Criticalj() followed by a call to os.Exit(1).
func (l *Logger) Fatalj(msg string, v interface{}) {
	logj(criticalsev, l, msg, v)
	os.Exit(1)
}

// Panic is equivalent to a call to l.Critical() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	panic(log(criticalsev, l, v...))
}

// Panicln is equivalent to a call to l.Criticalln() followed by a call to panic().
func (l *Logger) Panicln(v ...interface{}) {
	panic(logln(criticalsev, l, v...))
}

// Panicf is equivalent to a call to l.Criticalf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	panic(logf(criticalsev, l, format, v...))
}

// Panicj is equivalent to a call to l.Criticalj() followed by a call to panic().
func (l *Logger) Panicj(msg string, v interface{}) {
	logj(criticalsev, l, msg, v)
	panic(v)
}

type Logger struct {
	out   io.Writer
	err   io.Writer
	mu    sync.Mutex
	trace string
}

// ForRequest creates a new Logger. All the messages logged
// through it will trace back to the HTTP request r.
func ForRequest(r *http.Request) *Logger {
	l := &Logger{}

	if ProjectID != "" {
		h := r.Header.Get("X-Cloud-Trace-Context")
		if i := strings.IndexByte(h, '/'); i > 0 {
			if t := h[:i]; strings.Count(t, "0") != len(t) {
				l.trace = fmt.Sprintf("projects/%s/traces/%s", ProjectID, t)
			}
		}
	}

	return l
}

// New is for interface-level compatibility with standard library's
// "log" package. All the logged messages are written to w, one write per message.
// Remaining arguments are ignored.
// The ForRequest() constructor is more useful.
func New(w io.Writer, dummy2 string, dummy3 int) *Logger {
	return &Logger{
		out: w,
		err: w,
	}
}

func (l *Logger) writer(s severity) io.Writer {
	if s.IsErrorish() {
		if l.err != nil {
			return l.err
		}

		return os.Stderr
	}

	if l.out != nil {
		return l.out
	}

	return os.Stdout
}

type severity int32

const (
	debugsev severity = iota * 100
	infosev
	noticesev
	warningsev
	errorsev
	criticalsev
	alertsev
	emergencysev
)

func (s severity) MarshalJSON() ([]byte, error) {
	switch s {
	default:
		return []byte(`"UNKNOWN"`), fmt.Errorf("unknown severity: %d", s)
	case debugsev:
		return []byte(`"DEBUG"`), nil
	case infosev:
		return []byte(`"INFO"`), nil
	case noticesev:
		return []byte(`"NOTICE"`), nil
	case warningsev:
		return []byte(`"WARNING"`), nil
	case errorsev:
		return []byte(`"ERROR"`), nil
	case criticalsev:
		return []byte(`"CRITICAL"`), nil
	case alertsev:
		return []byte(`"ALERT"`), nil
	case emergencysev:
		return []byte(`"EMERGENCY"`), nil
	}
}

// IsErrorish returns true for severity ERROR and above it.
func (s severity) IsErrorish() bool {
	return s >= errorsev
}

func log(s severity, l *Logger, v ...interface{}) string {
	return logs(s, l, fmt.Sprint(v...))
}

func logln(s severity, l *Logger, v ...interface{}) string {
	return logs(s, l, fmt.Sprintln(v...))
}

func logf(s severity, l *Logger, format string, v ...interface{}) string {
	return logs(s, l, fmt.Sprintf(format, v...))
}

type entry struct {
	Message  string   `json:"message"`
	Severity severity `json:"severity,omitempty"`
	Trace    string   `json:"logging.googleapis.com/trace,omitempty"`
}

func logs(s severity, l *Logger, msg string) string {
	entry := entry{msg, s, l.trace}

	l.mu.Lock()
	defer l.mu.Unlock()
	_ = json.NewEncoder(l.writer(s)).Encode(entry)

	return msg
}

func logj(s severity, l *Logger, msg string, j interface{}) {
	if j == nil {
		log(s, l, msg)
		return
	}
	// TODO maybe use reflect to catch things like map[string]string, but it'd be
	// best effort only.
	// This allows to check for duplicate fields, e.g. "message", if a user throws
	// at us a map which they don't bother to sanitize.
	//
	// This is a dead end, it doesn't catch map[string]string:
	// switch v := j.(type) {
	// case map[string]interface{}:
	// 	if _, ok := v["message"]; ok {
	// 	}
	// }
	buf, err := json.Marshal(j)
	if err != nil {
		// Ignore. This is just logging. Do not panic the entire app
		// just because the caller had botched their x.Marshal(), etc.
		return
	}

	logRawJSON(s, l, msg, buf)
}

// logRawJSON writes the buf to the l logger. The buf should be
// an encoded JSON and its first byte must be '{'.
// The s and msg are brutally inserted as "severity" and "message" top-level JSON fields.
// The buf should not contain "severity", "message", or "logging.googleapis.com/trace"
// top-level JSON fields.
// No attempt is made to check whether the resulting string does not have these fields
// duplicated and whether it is a valid JSON. Spoiler alert: GCP Logging API seems to be
// quite gracefully handling malformed JSON entries with such duplicate fields.
func logRawJSON(s severity, l *Logger, msg string, buf []byte) {
	var msgj, sevj, tracej []byte
	var err error

	if msg != "" {
		msgj, err = json.Marshal(msg)
		if err != nil {
			return
		}
	}

	if l.trace != "" {
		tracej, err = json.Marshal(l.trace)
		if err != nil {
			return
		}
	}

	w := l.writer(s)
	jsonStruct := len(buf) > 0 && buf[0] == '{'

	if jsonStruct {
		buf = buf[1:]
	}

	// Critical Section
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, err := w.Write([]byte("{")); err != nil {
		return
	}

	comma := []byte{}

	if msg != "" {
		if _, err := w.Write([]byte("\"message\":")); err != nil {
			return
		}
		if _, err := w.Write(msgj); err != nil {
			return
		}

		comma = []byte(",")
	}

	sevj, err = s.MarshalJSON()
	if err == nil {
		if _, err := w.Write(comma); err != nil {
			return
		}
		if _, err := w.Write([]byte("\"severity\":")); err != nil {
			return
		}
		if _, err := w.Write(sevj); err != nil {
			return
		}

		comma = []byte(",")
	}

	if l.trace != "" {
		if _, err := w.Write(comma); err != nil {
			return
		}
		if _, err := w.Write([]byte("\"logging.googleapis.com/trace\":")); err != nil {
			return
		}
		if _, err := w.Write(tracej); err != nil {
			return
		}

		comma = []byte(",")
	}

	if !jsonStruct {
		if _, err := w.Write(comma); err != nil {
			return
		}
		if _, err := w.Write([]byte("\"value\":")); err != nil {
			return
		}
		if _, err := w.Write(buf); err != nil {
			return
		}
		_, _ = w.Write([]byte("}\n"))

		return
	}

	if len(buf) > 1 {
		if _, err := w.Write(comma); err != nil {
			return
		}
	}
	if _, err := w.Write(buf); err != nil {
		return
	}

	if buf[len(buf)-1] != '\n' {
		_, _ = w.Write([]byte("\n"))
	}
}
