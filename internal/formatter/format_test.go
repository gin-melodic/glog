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

package formatter

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func TestFormatter_Format(t *testing.T) {
	entry := &logrus.Entry{
		Logger:  nil,
		Level:   logrus.DebugLevel,
		Caller:  nil,
		Message: "logger content",
	}
	var f Formatter
	b, err := f.Format(entry)
	assert.Nil(t, err)
	println(string(b))
	bStr := string(b)

	assert.Contains(t, bStr, "DEBUG")	// check level
	assert.Contains(t, bStr, "PID")	// check pid
	assert.Contains(t, bStr, "logger content")	// check content

	// entry with caller
	time.Sleep(2 * time.Millisecond)
	entry.Caller = &runtime.Frame{File: "test.go", Line: 11211}
	entry.Logger = &logrus.Logger{
		ReportCaller: true,
	}
	b, err = f.Format(entry)
	assert.Nil(t, err)
	println(string(b))
	assert.Contains(t, string(b), "[test.go:11211]")	// check caller
}