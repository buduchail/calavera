package logger

import (
	"io"
	"github.com/buduchail/catrina"
	"github.com/sirupsen/logrus"
)

type (
	// Thin wrapper around logrus that simplifies
	// logging structured data (logger context)
	StructuredLogrus struct {
		logger         *logrus.Logger
		defaultContext catrina.LoggerContext
	}
)

func NewLogrus(context catrina.LoggerContext) (logger *StructuredLogrus) {
	logger = &StructuredLogrus{}
	logger.logger = logrus.New()
	logger.defaultContext = context
	return
}

func (l *StructuredLogrus) SetFormatter(formatter logrus.Formatter) {
	l.logger.Formatter = formatter
}

func (l *StructuredLogrus) SetOutput(out io.Writer) {
	l.logger.Out = out
}

func (l *StructuredLogrus) SetLevel(level logrus.Level) {
	l.logger.Level = level
}

func (l *StructuredLogrus) getFields(context *catrina.LoggerContext) logrus.Fields {
	var fields map[string]interface{}
	if context != nil {
		fields = *context
	} else {
		fields = logrus.Fields{}
	}
	for field, value := range l.defaultContext {
		_, exists := fields[field]
		if !exists {
			fields[field] = value
		}
	}
	return logrus.Fields(fields)
}

func (l *StructuredLogrus) Debug(message string, context *catrina.LoggerContext) {
	l.logger.WithFields(l.getFields(context)).
		Debug(message)
}

func (l *StructuredLogrus) Info(message string, context *catrina.LoggerContext) {
	l.logger.WithFields(l.getFields(context)).
		Info(message)
}
func (l *StructuredLogrus) Print(message string, context *catrina.LoggerContext) {
	l.logger.WithFields(l.getFields(context)).
		Print(message)
}

func (l *StructuredLogrus) Warn(message string, context *catrina.LoggerContext) {
	l.logger.WithFields(l.getFields(context)).
		Warn(message)
}

func (l *StructuredLogrus) Warning(message string, context *catrina.LoggerContext) {
	l.logger.WithFields(l.getFields(context)).
		Warning(message)
}

func (l *StructuredLogrus) Error(message string, context *catrina.LoggerContext) {
	l.logger.WithFields(l.getFields(context)).
		Error(message)
}

func (l *StructuredLogrus) Fatal(message string, context *catrina.LoggerContext) {
	l.logger.WithFields(l.getFields(context)).
		Fatal(message)
}

func (l *StructuredLogrus) Panic(message string, context *catrina.LoggerContext) {
	l.logger.WithFields(l.getFields(context)).
		Panic(message)
}
