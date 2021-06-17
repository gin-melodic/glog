/**
Copyright 2021 Gin Van

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gin_logger

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gin7758258/glog/internal/setup"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

var globalOnce = sync.Once{}
var sl *logrus.Logger

// LoggerOptions Init options
type LoggerOptions struct {
	MinAllowLevel	logrus.Level
	// When HighPerformance is true, logger won't auto record report file & line
	HighPerformance	bool
	OutputDir		string
	FilePrefix		string
	// Logs older than the specified number of SaveDay
	// will be automatically cleared
	SaveDay			time.Duration
	// ExtLoggerWriter write to other output,
	// like os.Stdout in dev. Default write to logFile
	ExtLoggerWriter  []io.Writer
	CustomTimeLayout string
}

// InitGlobalLogger Module entry function
// MUST call it before ShareLogger
func InitGlobalLogger(opt *LoggerOptions) (initErr error) {
	globalOnce.Do(func() {
		if sl != nil { return }
		l, err := setup.New(&setup.Options{
			BaseDir:          opt.OutputDir,
			Level:            opt.MinAllowLevel,
			ReportCaller:     !opt.HighPerformance,
			LogFilePrefix:    opt.FilePrefix,
			RotateDuration:   opt.SaveDay * 24 * time.Hour,
			ExtLoggerWriter:  opt.ExtLoggerWriter,
			CustomTimeLayout: opt.CustomTimeLayout,
		})
		if err != nil {
			initErr = errors.WithMessage(err, "[GINLOG]Init error.")
			return
		}
		sl = l
	})
	return
}

// ShareLogger Get global logger handle
// MUST InitGlobalLogger before call it
func ShareLogger() *logrus.Logger {
	if sl == nil {
		fmt.Println("[GINLOG]Please call InitGlobalLogger first.")
		return nil
	}
	return sl
}

// NewLoggerHandle Sometimes, when you need a log instance to print some
// specific log to a file, this method can provide that functionality
func NewLoggerHandle(opt *LoggerOptions) (logger *logrus.Logger, err error) {
	logger, err = setup.New(&setup.Options{
		BaseDir:          opt.OutputDir,
		Level:            opt.MinAllowLevel,
		ReportCaller:     !opt.HighPerformance,
		LogFilePrefix:    opt.FilePrefix,
		RotateDuration:   opt.SaveDay * 24 * time.Hour,
		ExtLoggerWriter:  opt.ExtLoggerWriter,
		CustomTimeLayout: opt.CustomTimeLayout,
	})
	return
}