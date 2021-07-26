package formatter

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	FieldTime           = "time"
	FieldLevel          = "level"
	FieldCalledFile     = "called_file"
	FieldCalledFunction = "called_function"
	FieldMessage        = "message"
)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	/**
	FieldsOrder
	Default:
	- time
	- level
	- called_file (if report caller was enabled)
	- called_function (if report caller was enabled)
	- other fields (sorted alphabetically)
	- message
	*/
	FieldsOrder []string

	// TimestampFormat - default: "2006-01-02 15:04:05.000"
	TimestampFormat string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - disable colors
	NoColors bool

	// NoFieldsColors - apply colors only to the level, default is level + fields
	NoFieldsColors bool

	// NoFieldsSpace - no space between fields
	NoFieldsSpace bool

	// ShowFullLevel - show a full level [WARNING] instead of [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel - no upper case for level value
	NoUppercaseLevel bool

	// TrimMessages - trim white spaces on messages
	TrimMessages bool
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.FieldsOrder == nil {
		f.FieldsOrder = []string{
			FieldTime,
			FieldLevel,
		}

		if entry.HasCaller() {
			f.FieldsOrder = append(f.FieldsOrder, FieldCalledFile, FieldCalledFunction)
		}

		f.FieldsOrder = append(f.FieldsOrder, FieldMessage)
	}

	// output buffer
	b := &bytes.Buffer{}

	// write fields
	f.writeOrderedFields(b, entry)

	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	defaultFields := f.getDefaultFields(entry)

	length := len(defaultFields) + len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry.Data, field, false, true)
		} else if _, ok := defaultFields[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, defaultFields, field, false, false)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if !foundFieldsMap[field] {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for i, field := range notFoundFields {
			f.writeField(b, entry.Data, field, i == 0, true)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, fields logrus.Fields, field string, firstNotFoundField, brackets bool) {
	prefix := ""
	if brackets {
		prefix += "["
	}

	suffix := ""
	if brackets {
		suffix += "]"
	}

	if !f.NoFieldsSpace && firstNotFoundField {
		b.WriteString(" ")
	}

	if f.HideKeys {
		fmt.Fprintf(b, "%s%v%s", prefix, fields[field], suffix)
	} else {
		fmt.Fprintf(b, "%s%s:%v%s", prefix, field, fields[field], suffix)
	}

	l := len(f.FieldsOrder)
	if !f.NoFieldsSpace && field != f.FieldsOrder[l-1] {
		b.WriteString(" ")
	}
}

func (f *Formatter) getDefaultFields(entry *logrus.Entry) logrus.Fields {
	defaultFields := logrus.Fields{}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = "2006-01-02 15:04:05.000"
	}
	defaultFields[FieldTime] = entry.Time.Format(timestampFormat)

	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	if !f.ShowFullLevel {
		level = level[:4]
	}

	var levelString string
	if f.NoColors {
		levelString = fmt.Sprintf("[%s]", level)
	} else {
		levelColor := getColorByLevel(entry.Level)
		levelString = fmt.Sprintf("\x1b[%dm[%s]", levelColor, level)
	}

	if !f.NoColors && f.NoFieldsColors {
		levelString += "\x1b[0m"
	}
	defaultFields[FieldLevel] = levelString

	if entry.HasCaller() {
		caller := entry.Caller
		s := strings.Split(caller.Function, ".")
		defaultFields[FieldCalledFile] = fmt.Sprintf("%s:%d", caller.File, caller.Line)
		defaultFields[FieldCalledFunction] = s[len(s)-1]
	}

	var message string
	if f.TrimMessages {
		message = strings.TrimSpace(entry.Message)
	} else {
		message = entry.Message
	}
	defaultFields[FieldMessage] = message

	return defaultFields
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
