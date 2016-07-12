gomol-console
============

[![GoDoc](https://godoc.org/github.com/aphistic/gomol-console?status.svg)](https://godoc.org/github.com/aphistic/gomol-console)
[![Build Status](https://img.shields.io/travis/aphistic/gomol-console.svg)](https://travis-ci.org/aphistic/gomol-console)
[![Code Coverage](https://img.shields.io/codecov/c/github/aphistic/gomol-console.svg)](http://codecov.io/github/aphistic/gomol-console?branch=master)

gomol-console is a logger for [gomol](https://github.com/aphistic/gomol) to support logging to the console.

Installation
============

The recommended way to install is via http://gopkg.in

    go get gopkg.in/aphistic/gomol-console.v0
    ...
    import "gopkg.in/aphistic/gomol-console.v0"

gomol-console can also be installed the standard way as well

    go get github.com/aphistic/gomol-console
    ...
    import "github.com/aphistic/gomol-console"

Examples
========

For brevity a lot of error checking has been omitted, be sure you do your checks!

This is a super basic example of adding a console logger to gomol and then logging a few messages:

```go
package main

import (
	"github.com/aphistic/gomol"
	gc "github.com/aphistic/gomol-console"
)

func main() {
	// Add a console logger
	consoleCfg := gc.NewConsoleLoggerConfig()
	consoleLogger, _ := gc.NewConsoleLogger(consoleCfg)
	gomol.AddLogger(consoleLogger)

	// Set some global attrs that will be added to all
	// messages automatically
	gomol.SetAttr("facility", "gomol.example")
	gomol.SetAttr("another_attr", 1234)

	// Initialize the loggers
	gomol.InitLoggers()
	defer gomol.ShutdownLoggers()

	// Log some debug messages with message-level attrs
	// that will be sent only with that message
	for idx := 1; idx <= 10; idx++ {
		gomol.Dbgm(map[string]interface{}{
			"msg_attr1": 4321,
		}, "Test message %v", idx)
	}
}

```
