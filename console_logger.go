// Package gomolconsole implements a console logger for gomol.
package gomolconsole

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aphistic/gomol"
	"github.com/mgutz/ansi"
)

type ConsoleLoggerConfig struct {
	// Colorize specifies whether the output will include ANSI colors or not. Defaults to true
	Colorize bool

	// Writer is the io.Writer to be used if not overridden in other options. Defaults to os.Stdout
	Writer io.Writer

	// DebugWriter is the io.Writer debug messages will be written to. Defaults to Writer
	DebugWriter io.Writer
	// InfoWriter is the io.Writer info messages will be written to. Defaults to Writer
	InfoWriter io.Writer
	// WarningWriter is the io.Writer warning messages will be written to. Defaults to Writer
	WarningWriter io.Writer
	// ErrorWriter is the io.Writer error messages will be written to. Defaults to Writer
	ErrorWriter io.Writer
	// FatalWriter is the io.Writer fatal messages will be written to. Defaults to Writer
	FatalWriter io.Writer
}

type ConsoleLogger struct {
	base          *gomol.Base
	writers       map[gomol.LogLevel]io.Writer
	tpl           *gomol.Template
	isInitialized bool
	config        *ConsoleLoggerConfig
}
type consoleWriter interface {
	Print(msg string)
}

func NewConsoleLoggerConfig() *ConsoleLoggerConfig {
	return &ConsoleLoggerConfig{
		Colorize: true,
	}
}

func NewConsoleLogger(config *ConsoleLoggerConfig) (*ConsoleLogger, error) {
	if config == nil {
		config = NewConsoleLoggerConfig()
	}

	l := &ConsoleLogger{
		writers: make(map[gomol.LogLevel]io.Writer),
		tpl:     NewTemplateDefault(),
		config:  config,
	}

	l.populateWriters(config)

	return l, nil
}

var printclean = func(msg string) string {
	return msg
}
var printdbg = ansi.ColorFunc("cyan")
var printinfo = ansi.ColorFunc("green")
var printwarn = ansi.ColorFunc("yellow")
var printerr = ansi.ColorFunc("red")
var printfatal = ansi.ColorFunc("red+b")

func (l *ConsoleLogger) populateWriters(cfg *ConsoleLoggerConfig) {
	var defaultWriter io.Writer = os.Stdout
	if cfg.Writer != nil {
		defaultWriter = cfg.Writer
	}

	l.writers[gomol.LevelDebug] = defaultWriter
	l.writers[gomol.LevelInfo] = defaultWriter
	l.writers[gomol.LevelWarning] = defaultWriter
	l.writers[gomol.LevelError] = defaultWriter
	l.writers[gomol.LevelFatal] = defaultWriter

	if cfg.DebugWriter != nil {
		l.writers[gomol.LevelDebug] = cfg.DebugWriter
	}
	if cfg.InfoWriter != nil {
		l.writers[gomol.LevelInfo] = cfg.InfoWriter
	}
	if cfg.WarningWriter != nil {
		l.writers[gomol.LevelWarning] = cfg.WarningWriter
	}
	if cfg.ErrorWriter != nil {
		l.writers[gomol.LevelError] = cfg.ErrorWriter
	}
	if cfg.FatalWriter != nil {
		l.writers[gomol.LevelFatal] = cfg.FatalWriter
	}
}

func (l *ConsoleLogger) SetBase(base *gomol.Base) {
	l.base = base
}

func (l *ConsoleLogger) SetTemplate(tpl *gomol.Template) error {
	if tpl == nil {
		return errors.New("A template must be provided")
	}
	l.tpl = tpl

	return nil
}

func (l *ConsoleLogger) InitLogger() error {
	l.isInitialized = true
	return nil
}
func (l *ConsoleLogger) IsInitialized() bool {
	return l.isInitialized
}

func (l *ConsoleLogger) ShutdownLogger() error {
	l.isInitialized = false
	return nil
}

func (l *ConsoleLogger) Logm(timestamp time.Time, level gomol.LogLevel, attrs map[string]interface{}, msg string) error {
	mergedAttrs := make(map[string]interface{})

	if l.base != nil && l.base.BaseAttrs != nil {
		for key, val := range l.base.BaseAttrs.Attrs() {
			mergedAttrs[key] = val
		}
	}

	for key, val := range attrs {
		mergedAttrs[key] = val
	}

	nMsg := gomol.NewTemplateMsg(timestamp, level, mergedAttrs, msg)
	out, err := l.tpl.Execute(nMsg, l.config.Colorize)
	if err != nil {
		return err
	}

	w, ok := l.writers[level]
	if !ok {
		return fmt.Errorf("unsupported log level")
	}
	w.Write([]byte(out + "\n"))

	return nil
}
