package logger

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
)

type DefaultFormatter struct{}

func (l DefaultFormatter) ColoredLevel(entry *logrus.Entry) string {
	levelName := strings.ToUpper(entry.Level.String())

	switch entry.Level { //nolint:exhaustive
	case logrus.DebugLevel:
		return aurora.Yellow(levelName).String()
	case logrus.InfoLevel:
		return aurora.Green(levelName).String()
	case logrus.WarnLevel:
		return aurora.Magenta(levelName).String()
	case logrus.ErrorLevel:
		return aurora.Red(levelName).String()
	case logrus.FatalLevel:
		return aurora.BrightRed(levelName).String()
	default:
		return entry.Level.String()
	}
}

func (l DefaultFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(fmt.Sprintf("[%s] %s: %s",
		entry.Time.Format(time.RFC1123Z),
		l.ColoredLevel(entry),
		entry.Message))

	for key, value := range entry.Data {
		b.WriteString(fmt.Sprintf(" %s=%v", aurora.Blue(key), value))
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func New(level, format string) (*logrus.Logger, error) {
	logger := logrus.New()

	if strings.ToUpper(format) == "JSON" {
		logger.Formatter = &logrus.JSONFormatter{}
	} else {
		logger.Formatter = &DefaultFormatter{}
	}

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	logger.Level = logLevel

	return logger, nil
}
