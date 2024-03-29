package glog

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestInitGlobalLogger(t *testing.T) {
	err := InitGlobalLogger(&LoggerOptions{
		MinAllowLevel:    logrus.DebugLevel,
		HighPerformance:  false,
		OutputDir:        "./test-log",
		FilePrefix:       "cyto",
		SaveDay:          7,
		ExtLoggerWriter:  []io.Writer{os.Stdout},
		CustomTimeLayout: "20060102150405",
	})
	assert.Nil(t, err)
	assert.DirExists(t, "./test-log")
	assert.NotNil(t, sl)
	assert.NoErrorf(t, ShareLogger().Close(), "close logger failed")
	assert.NoErrorf(t, os.RemoveAll("./test-log"), "remove test log dir failed")
}

// Usage Example
func TestUsageExample(t *testing.T) {
	// Must init before use logger
	const kOutputDir = "./test-log"
	err := InitGlobalLogger(&LoggerOptions{
		MinAllowLevel:   logrus.DebugLevel,
		HighPerformance: true,
		OutputDir:       kOutputDir,
		FilePrefix:      "cyto",
		SaveDay:         7,
		ExtLoggerWriter: []io.Writer{os.Stdout},
	})
	assert.Nil(t, err)

	// global logger
	ShareLogger().Warn("error message.")
	param := 123
	ShareLogger().Infof("param: %d", param)

	// partner logger
	partnerLogger, err := NewLoggerHandle(&LoggerOptions{
		MinAllowLevel:    logrus.InfoLevel,
		HighPerformance:  true,
		OutputDir:        kOutputDir,
		FilePrefix:       "partner",
		SaveDay:          7,
		CustomTimeLayout: "2006/01/02 15:04:05",
	})
	assert.Nil(t, err)
	partnerLogger.Info("test partner log")

	// for middleware usage see test cases in specific modules.

	// You can commit next line to see the logging files.
	_ = os.RemoveAll(kOutputDir)
}
