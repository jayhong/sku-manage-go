package server

import (
	"bytes"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

// LoggerEntry is the structure
// passed to the template.
type LoggerEntry struct {
	StartTime string
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
}

// LoggerDefaultFormat is the format
// logged used by the default Logger instance.
var LoggerDefaultFormat = "{{.StartTime}} | {{.Status}} | \t {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}} "

// LoggerDefaultDateFormat is the
// format used for date by the
// default Logger instance.
var LoggerDefaultDateFormat = time.RFC3339

// ALogger interface
type ALogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

type PearLogger struct {
}

func NewPearLogger() *PearLogger {
	return &PearLogger{}
}

func (p *PearLogger) Println(v ...interface{}) {
	logrus.Infoln(v...)
}

func (p *PearLogger) Printf(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// ALogger implements just enough log.Logger interface to be compatible with other implementations
	ALogger
	dateFormat string
	template   *template.Template
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	logger := &Logger{ALogger: NewPearLogger(), dateFormat: LoggerDefaultDateFormat}
	logger.SetFormat(LoggerDefaultFormat)
	return logger
}

func (l *Logger) SetFormat(format string) {
	l.template = template.Must(template.New("negroni_parser").Parse(format))
}

func (l *Logger) SetDateFormat(format string) {
	l.dateFormat = format
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	if strings.Contains(r.URL.Path, "/health") {
		return
	}

	res := rw.(negroni.ResponseWriter)
	log := LoggerEntry{
		StartTime: start.Format(l.dateFormat),
		Status:    res.Status(),
		Duration:  time.Since(start),
		Hostname:  r.Host,
		Method:    r.Method,
		Path:      r.URL.Path,
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)
	l.Printf(buff.String())
}
