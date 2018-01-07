package gomolconsole

import (
	"testing"
	"time"

	"github.com/aphistic/gomol"
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.AddSuite(&GomolSuite{})
	})
}

type GomolSuite struct{}

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

func (s *GomolSuite) TestTestConsoleWriter(t sweet.T) {
	w := newTestConsoleWriter()
	Expect(w.Output).ToNot(BeNil())
	Expect(w.Output).To(HaveLen(0))

	w.Print("print1")
	Expect(w.Output).To(HaveLen(1))

	w.Print("print2")
	Expect(w.Output).To(HaveLen(2))
}

// Issue-specific tests

func (s *GomolSuite) TestIssue5StringFormatting(t sweet.T) {
	b := gomol.NewBase()
	b.InitLoggers()

	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	l, err := NewConsoleLogger(cfg)
	Expect(err).To(BeNil())
	w := newTestConsoleWriter()
	l.setWriter(w)
	b.AddLogger(l)

	b.Dbgf("msg %v%%", 100)

	b.ShutdownLoggers()

	Expect(w.Output).To(HaveLen(1))
	Expect(w.Output[0]).To(Equal("[DEBUG] msg 100%\n"))
}

func (s *GomolSuite) TestAttrsMergedFromBase(t sweet.T) {
	b := gomol.NewBase()
	b.SetAttr("base_attr", "foo")
	b.InitLoggers()

	cfg := NewConsoleLoggerConfig()
	cfg.Colorize = false
	l, err := NewConsoleLogger(cfg)

	testTpl, err := gomol.NewTemplate(
		"[{{color}}{{ucase .LevelName}}{{reset}}] {{.Message}}" +
			"{{if .Attrs}}{{range $key, $val := .Attrs}}\n   {{$key}}: {{$val}}{{end}}{{end}}",
	)

	l.SetTemplate(testTpl)
	Expect(err).To(BeNil())
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

	Expect(w.Output).To(HaveLen(1))
	Expect(w.Output[0]).To(Equal("[DEBUG] msg 100%\n   adapter_attr: bar\n   base_attr: foo\n   log_attr: baz\n"))
}

// General tests

func (s *GomolSuite) TestConsoleSetTemplate(t sweet.T) {
	cl, err := NewConsoleLogger(nil)
	Expect(cl.tpl).ToNot(BeNil())

	err = cl.SetTemplate(nil)
	Expect(err).ToNot(BeNil())

	tpl, err := gomol.NewTemplate("")
	Expect(err).To(BeNil())
	err = cl.SetTemplate(tpl)
	Expect(err).To(BeNil())
}

func (s *GomolSuite) TestConsoleInitLogger(t sweet.T) {
	cl, err := NewConsoleLogger(nil)
	Expect(err).To(BeNil())
	Expect(cl.IsInitialized()).To(BeFalse())
	cl.InitLogger()
	Expect(cl.IsInitialized()).To(BeTrue())
}

func (s *GomolSuite) TestConsoleShutdownLogger(t sweet.T) {
	cl, _ := NewConsoleLogger(nil)
	cl.InitLogger()
	Expect(cl.IsInitialized()).To(BeTrue())
	cl.ShutdownLogger()
	Expect(cl.IsInitialized()).To(BeFalse())
}

func (s *GomolSuite) TestConsoleColorLogm(t sweet.T) {
	cfg := NewConsoleLoggerConfig()
	cl, _ := NewConsoleLogger(cfg)
	w := newTestConsoleWriter()
	cl.setWriter(w)
	cl.Logm(time.Now(), gomol.LevelFatal, nil, "test")
	Expect(w.Output).To(HaveLen(1))
	Expect(w.Output[0]).To(Equal("[\x1b[1;31mFATAL\x1b[0m] test\n"))
}

func (s *GomolSuite) TestConsoleLogm(t sweet.T) {
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
	Expect(w.Output).To(HaveLen(1))
	Expect(w.Output[0]).To(Equal("[FATAL] test 1234\n"))
}

func (s *GomolSuite) TestConsoleBaseAttrs(t sweet.T) {
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
	Expect(w.Output).To(HaveLen(1))
	Expect(w.Output[0]).To(Equal("[DEBUG] test 1234\n"))
}
