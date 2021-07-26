package formatter_test

import (
	"bytes"
	"testing"

	formatter "github.com/SafeStudio/logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestFormatter(t *testing.T) {
	suite.Run(t, new(FormatterTestSuite))
}

type FormatterTestSuite struct {
	suite.Suite
}

func (s FormatterTestSuite) TestFormatter_Format_with_report_caller() {
	output := bytes.NewBuffer([]byte{})

	l := logrus.New()
	l.SetOutput(output)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&formatter.Formatter{
		HideKeys:        true,
		NoColors:        true,
		TimestampFormat: "-",
	})
	l.SetReportCaller(true)

	l.Debug("test1")

	assert.Regexp(s.T(), "- \\[DEBU\\] .+\\.go:[0-9]+ .+ test1\n$", output.String())
}

func (s FormatterTestSuite) TestFormatter_Format_without_colors_and_without_fieldsColor() {
	output := bytes.NewBuffer([]byte{})

	l := logrus.New()
	l.SetOutput(output)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&formatter.Formatter{
		HideKeys:        true,
		NoColors:        false,
		NoFieldsColors:  false,
		FieldsOrder: []string{"time", "level", "server", "region", "message"},
		TimestampFormat: "-",
	})

	l.
		WithFields(logrus.Fields{
			"server": "server1",
			"region": "region1",
			"service": "auth",
		}).
		Debug("test1")

	assert.Regexp(s.T(), "- \\x1b\\[37m\\[DEBU\\] \\[server1\\] \\[region1\\] test1 \\[auth\\] \\x1b\\[0m\n$", output.String())
}

func (s FormatterTestSuite) TestFormatter_Format_with_custom_field_order() {
	output := bytes.NewBuffer([]byte{})

	l := logrus.New()
	l.SetOutput(output)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&formatter.Formatter{
		HideKeys:    true,
		NoColors:    true,
		FieldsOrder: []string{"time", "called_file", "called_function", "level", "message"},
	})
	l.SetReportCaller(true)

	l.Debug("test1")

	assert.Regexp(s.T(), "\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}.\\d{3} .+\\.go:[0-9]+ .+ \\[DEBU\\] test1\n$", output.String())
}
