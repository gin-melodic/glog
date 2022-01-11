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

package setup

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-melodic/glog/internal/formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

// Options logger setup options, BaseDir is requirement, RewriteDuration default 7 days
type Options struct {
	BaseDir string

	Level          logrus.Level
	ReportCaller   bool
	LogFilePrefix  string
	RotateDuration time.Duration
	// ExtLoggerWriter write to other output, like os.Stdout in dev. Default write to logFile
	ExtLoggerWriter  []io.Writer
	CustomTimeLayout string
}

func New(opt *Options) (*logrus.Logger, error) {
	mutePipe, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	writers := io.MultiWriter(append(opt.ExtLoggerWriter, mutePipe)...)
	lc := logrus.New()
	lc.SetLevel(opt.Level)
	lc.SetReportCaller(opt.ReportCaller)
	lc.SetFormatter(&formatter.Formatter{TimeStampLayout: opt.CustomTimeLayout})
	lc.Out = writers
	// check log base dir
	if opt.BaseDir == "" {
		return nil, errors.New("Must give a log file dir path.")
	}
	if _, err = os.Stat(opt.BaseDir); err != nil {
		if os.IsNotExist(err) {
			// try to create dir
			if err := os.Mkdir(opt.BaseDir, 0755); err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			return nil, errors.WithStack(err)
		}
	}
	prefix := opt.LogFilePrefix + "-"
	if opt.LogFilePrefix == "" {
		prefix = ""
	}
	// rotate max age
	maxAge := 7 * 24 * time.Hour
	if opt.RotateDuration > 0 {
		maxAge = opt.RotateDuration
	}
	// combine log
	cbFmt := opt.BaseDir + "/" + prefix + "combine-%Y%m%d.log"
	cbWriter, err := rotatelogs.New(cbFmt,
		rotatelogs.WithLinkName(opt.BaseDir+"/latest-combine-"+prefix+"log"),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		return nil, errors.WithMessage(err, "rotate combine log error")
	}
	// error log
	errorFmt := opt.BaseDir + "/" + prefix + "error-%Y%m%d.log"
	errorWriter, err := rotatelogs.New(errorFmt,
		rotatelogs.WithLinkName(opt.BaseDir+"/latest-error-"+prefix+"log"),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		return nil, errors.WithMessage(err, "rotate error log error")
	}
	lfsMap := lfshook.WriterMap{
		logrus.PanicLevel: io.MultiWriter(cbWriter, errorWriter),
		logrus.FatalLevel: io.MultiWriter(cbWriter, errorWriter),
		logrus.ErrorLevel: io.MultiWriter(cbWriter, errorWriter),
		logrus.WarnLevel:  cbWriter,
		logrus.InfoLevel:  cbWriter,
		logrus.DebugLevel: cbWriter,
		logrus.TraceLevel: cbWriter,
	}
	hook := lfshook.NewHook(lfsMap, &formatter.Formatter{TimeStampLayout: opt.CustomTimeLayout})
	lc.AddHook(hook)
	return lc, nil
}
