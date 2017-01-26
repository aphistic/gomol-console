package gomolconsole

import (
	"testing"
	"time"

	"github.com/aphistic/gomol"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GomolSuite struct{}

var _ = Suite(&GomolSuite{})

type testConsoleWriter struct {
	Output []string
}

func newTestConsoleWriter() *testConsoleWriter {
	return &testConsoleWriter{
		Output: make([]string, 0),
	}
}

func (w *testConsoleWriter) Print(msg string) {
	w.Output = append(w.Output, msg)
}

func (s *GomolSuite) TestTestConsoleWriter(c *C) {
	w := newTestConsoleWriter()
	c.Check(w.Output, NotNil)
	c.Check(w.Output, HasLen, 0)

	w.Print("print1")
	c.Check(w.Output, HasLen, 1)

	w.Print("print2")
	c.Check(w.Output, HasLen, 2)
}

// Issue-specific tests

func (s *GomolSuite) TestIssue5StringFormatting(c *C) {
	b := gomol.NewBase()
	b.InitLoggers()

	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	l, err := NewConsoleLogger(cfg)
	c.Assert(err, IsNil)
	w := newTestConsoleWriter()
	l.setWriter(w)
	b.AddLogger(l)

	b.Dbgf("msg %v%%", 100)

	b.ShutdownLoggers()

	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[DEBUG] msg 100%\n")
}

func (s *GomolSuite) TestAttrsMergedFromBase(c *C) {
	b := gomol.NewBase()
	b.SetAttr("base_attr", "foo")
	b.InitLoggers()

	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	l, err := NewConsoleLogger(cfg)

	testTpl, _ := gomol.NewTemplate("[{{color}}{{ucase .LevelName}}{{reset}}] {{.Message}}" +
		"{{if .Attrs}}{{range $key, $val := .Attrs}}\n   {{$key}}: {{$val}}{{end}}{{end}}")

	l.SetTemplate(testTpl)
	c.Assert(err, IsNil)
	w := newTestConsoleWriter()
	l.setWriter(w)
	b.AddLogger(l)

	la := b.NewLogAdapter(gomol.NewAttrsFromMap(map[string]interface{}{
		"adapter_attr": "bar",
	}))

	la.Dbgm(gomol.NewAttrsFromMap(map[string]interface{}{
		"log_attr": "baz",
	}), "msg %v%%", 100)

	b.ShutdownLoggers()

	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[DEBUG] msg 100%\n   adapter_attr: bar\n   base_attr: foo\n   log_attr: baz\n")
}

// General tests

func (s *GomolSuite) TestConsoleSetTemplate(c *C) {
	cl, err := NewConsoleLogger(nil)
	c.Check(cl.tpl, NotNil)

	err = cl.SetTemplate(nil)
	c.Check(err, NotNil)

	tpl, err := gomol.NewTemplate("")
	c.Assert(err, IsNil)
	err = cl.SetTemplate(tpl)
	c.Check(err, IsNil)
}

func (s *GomolSuite) TestConsoleInitLogger(c *C) {
	cl, err := NewConsoleLogger(nil)
	c.Assert(err, IsNil)
	c.Check(cl.IsInitialized(), Equals, false)
	cl.InitLogger()
	c.Check(cl.IsInitialized(), Equals, true)
}

func (s *GomolSuite) TestConsoleShutdownLogger(c *C) {
	cl, _ := NewConsoleLogger(nil)
	cl.InitLogger()
	c.Check(cl.IsInitialized(), Equals, true)
	cl.ShutdownLogger()
	c.Check(cl.IsInitialized(), Equals, false)
}

func (s *GomolSuite) TestConsoleColorLogm(c *C) {
	cfg := NewConsoleLoggerConfig()
	cl, _ := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Logm(time.Now(), gomol.LevelFatal, nil, "test")
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[\x1b[1;31mFATAL\x1b[0m] test\n")
}

func (s *GomolSuite) TestConsoleLogm(c *C) {
	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl, _ := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Logm(
		time.Now(),
		gomol.LevelFatal,
		map[string]interface{}{
			"attr1": 4321,
		},
		"test 1234")
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[FATAL] test 1234\n")
}

func (s *GomolSuite) TestConsoleBaseAttrs(c *C) {
	b := gomol.NewBase()
	b.SetAttr("attr1", 7890)
	b.SetAttr("attr2", "val2")

	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	cl, _ := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	b.AddLogger(cl)
	cl.Logm(
		time.Now(),
		gomol.LevelDebug,
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test 1234")
	c.Assert(w.Output, HasLen, 1)
	c.Check(w.Output[0], Equals, "[DEBUG] test 1234\n")
}
