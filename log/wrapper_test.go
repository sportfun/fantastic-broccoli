package log

import (
	. "github.com/onsi/gomega"
	"github.com/sportfun/gakisitor/config"
	"testing"
)

func TestNewArgumentBinder(t *testing.T) {
	RegisterTestingT(t)

	argBinder := NewArgumentBinder("argument '%s'").(*argumentBinderImpl)
	Expect(argBinder.format).Should(Equal("argument '%s'"))
	Expect(argBinder.binds).Should(BeEmpty())
	Expect(argBinder.more).Should(BeEmpty())
}

func TestArgumentBinderImpl(t *testing.T) {
	RegisterTestingT(t)

	testCases := []struct {
		Format   string
		Argument []interface{}
		More     map[string]interface{}
		Result   string
	}{
		{Format: "", Argument: []interface{}{}, More: map[string]interface{}{}, Result: ""},
		{Format: "message", Argument: []interface{}{}, More: map[string]interface{}{}, Result: "message"},

		{Format: "%s|%d", Argument: []interface{}{"message", 24}, More: map[string]interface{}{}, Result: "message|24"},
		{Format: "", Argument: []interface{}{"message"}, More: map[string]interface{}{}, Result: "%!(EXTRA string=message)"},

		{Format: "%s", Argument: []interface{}{"message"}, More: map[string]interface{}{"info": 0}, Result: "message"},
		{Format: "", Argument: []interface{}{"message"}, More: map[string]interface{}{"info": 0}, Result: "%!(EXTRA string=message)"},
	}

	for _, tc := range testCases {
		argBinder := NewArgumentBinder(tc.Format)

		argBinder.Bind(tc.Argument...)
		for index, data := range tc.More {
			argBinder.More(index, data)
		}

		Expect(argBinder.getMessage()).Should(Equal(tc.Result))
		Expect(argBinder.getMoreInfo()).Should(HaveLen(len(tc.More)))
		for index, data := range tc.More {
			Expect(argBinder.getMoreInfo()).Should(HaveKey(index))
			Expect(argBinder.getMoreInfo()[index]).Should(Equal(data))
			argBinder.More(index, data)
		}
	}
}

func TestNewProduction(t *testing.T) {
	RegisterTestingT(t)

	testCases := []struct {
		Conf          []config.LogDefinition
		ForProduction bool
	}{
		{Conf: []config.LogDefinition{}, ForProduction: false},
		{Conf: []config.LogDefinition{{File: "stdout", Encoding: "json", Level: "info"}}, ForProduction: true},
		{Conf: []config.LogDefinition{{File: "stdout", Encoding: "json", Level: "info"}, {File: "log", Encoding: "text", Level: "none"}}, ForProduction: true},
	}

	for _, tc := range testCases {
		logger := NewProduction(tc.Conf...).(*loggerImpl)
		Expect(logger.forProduction).Should(Equal(tc.ForProduction))
	}
}

func TestNewDevelopment(t *testing.T) {
	RegisterTestingT(t)
	logger := NewDevelopment().(*loggerImpl)

	Expect(logger.forProduction).Should(BeFalse())
}

func TestNewTest(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	logger := NewTest(&buffer)
	logger.Infof("log")
	Expect(buffer).Should(Equal("INFO	log"))
}

func TestLoggerImpl(t *testing.T) {
	RegisterTestingT(t)

	buffer := ""
	logger := NewTest(&buffer)

	testCases := []struct {
		AFnc     func(argumentBinder)
		FFnc     func(string, ...interface{})
		Format   string
		Argument []interface{}
		More     map[string]interface{}
		Result   string
	}{
		{AFnc: logger.Error, Format: "log error: '%s'", Argument: []interface{}{"message"}, More: map[string]interface{}{"type": "error"}, Result: `ERROR	log error: 'message'	{"type": "error"}`},
		{FFnc: logger.Errorf, Format: "log error: '%s'", Argument: []interface{}{"message"}, Result: `ERROR	log error: 'message'`},

		{AFnc: logger.Warn, Format: "log warn: '%s'", Argument: []interface{}{"message"}, More: map[string]interface{}{"type": "warn"}, Result: `WARN	log warn: 'message'	{"type": "warn"}`},
		{FFnc: logger.Warnf, Format: "log warn: '%s'", Argument: []interface{}{"message"}, Result: `WARN	log warn: 'message'`},

		{AFnc: logger.Info, Format: "log info: '%s'", Argument: []interface{}{"message"}, More: map[string]interface{}{"type": "info"}, Result: `INFO	log info: 'message'	{"type": "info"}`},
		{FFnc: logger.Infof, Format: "log info: '%s'", Argument: []interface{}{"message"}, Result: `INFO	log info: 'message'`},

		{AFnc: logger.Debug, Format: "log debug: '%s'", Argument: []interface{}{"message"}, More: map[string]interface{}{"type": "debug"}, Result: `DEBUG	log debug: 'message'	{"type": "debug"}`},
		{FFnc: logger.Debugf, Format: "log debug: '%s'", Argument: []interface{}{"message"}, Result: `DEBUG	log debug: 'message'`},

		{AFnc: logger.Debug, Format: "log with struct", Argument: []interface{}{}, More: map[string]interface{}{"struct": struct{ A string }{""}}, Result: `DEBUG	log with struct	{"struct": {"A":""}}`},
		{AFnc: logger.Debug, Format: "log with nil", Argument: []interface{}{}, More: map[string]interface{}{"nil": nil}, Result: `DEBUG	log with nil`},

		{AFnc: logger.Debug, Format: "log with type 'String'", Argument: []interface{}{}, More: map[string]interface{}{"String": ""}, Result: `DEBUG	log with type 'String'	{"String": ""}`},
		{AFnc: logger.Debug, Format: "log with type 'Bool'", Argument: []interface{}{}, More: map[string]interface{}{"Bool": false}, Result: `DEBUG	log with type 'Bool'	{"Bool": false}`},
		{AFnc: logger.Debug, Format: "log with type 'Int'", Argument: []interface{}{}, More: map[string]interface{}{"Int": int(0)}, Result: `DEBUG	log with type 'Int'	{"Int": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Int8'", Argument: []interface{}{}, More: map[string]interface{}{"Int8": int8(0)}, Result: `DEBUG	log with type 'Int8'	{"Int8": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Int16'", Argument: []interface{}{}, More: map[string]interface{}{"Int16": int16(0)}, Result: `DEBUG	log with type 'Int16'	{"Int16": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Int32'", Argument: []interface{}{}, More: map[string]interface{}{"Int32": int32(0)}, Result: `DEBUG	log with type 'Int32'	{"Int32": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Int64'", Argument: []interface{}{}, More: map[string]interface{}{"Int64": int64(0)}, Result: `DEBUG	log with type 'Int64'	{"Int64": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Uint'", Argument: []interface{}{}, More: map[string]interface{}{"Uint": uint(0)}, Result: `DEBUG	log with type 'Uint'	{"Uint": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Uint8'", Argument: []interface{}{}, More: map[string]interface{}{"Uint8": uint8(0)}, Result: `DEBUG	log with type 'Uint8'	{"Uint8": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Uint16'", Argument: []interface{}{}, More: map[string]interface{}{"Uint16": uint16(0)}, Result: `DEBUG	log with type 'Uint16'	{"Uint16": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Uint32'", Argument: []interface{}{}, More: map[string]interface{}{"Uint32": uint32(0)}, Result: `DEBUG	log with type 'Uint32'	{"Uint32": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Uint64'", Argument: []interface{}{}, More: map[string]interface{}{"Uint64": uint64(0)}, Result: `DEBUG	log with type 'Uint64'	{"Uint64": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Float32'", Argument: []interface{}{}, More: map[string]interface{}{"Float32": float32(0)}, Result: `DEBUG	log with type 'Float32'	{"Float32": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Float64'", Argument: []interface{}{}, More: map[string]interface{}{"Float64": float64(0)}, Result: `DEBUG	log with type 'Float64'	{"Float64": 0}`},
		{AFnc: logger.Debug, Format: "log with type 'Complex64'", Argument: []interface{}{}, More: map[string]interface{}{"Complex64": complex64(0)}, Result: `DEBUG	log with type 'Complex64'	{"Complex64": "0+0i"}`},
		{AFnc: logger.Debug, Format: "log with type 'Complex128'", Argument: []interface{}{}, More: map[string]interface{}{"Complex128": complex128(0)}, Result: `DEBUG	log with type 'Complex128'	{"Complex128": "0+0i"}`},
	}

	for _, tc := range testCases {
		buffer = ""

		switch {
		case tc.AFnc != nil:
			args := NewArgumentBinder(tc.Format).Bind(tc.Argument...)
			for k, v := range tc.More {
				args.More(k, v)
			}

			tc.AFnc(args)
		case tc.FFnc != nil:
			tc.FFnc(tc.Format, tc.Argument...)
		}

		Expect(buffer).Should(Equal(tc.Result))
	}
}
