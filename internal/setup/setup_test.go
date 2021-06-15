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
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	_ = os.RemoveAll("./test-logs")
	l, err := New(&Options{
		Level: logrus.DebugLevel,
		ReportCaller: true,
		BaseDir: "./test-logs",
		LogFilePrefix: "test",
		ExtLoggerWriter: []io.Writer{os.Stdout},
	})
	assert.Nil(t, err)
	l.Traceln("123")	// this stmt won't create log file
	_, err = os.Stat("./test-logs/latest-combine-test-log")
	assert.True(t, os.IsNotExist(err))

	// only combine file
	l.Debugln("123")
	_, err = os.Stat("./test-logs/latest-combine-test-log")
	assert.Nil(t, err)
	_, err = os.Stat("./test-logs/latest-error-test-log")
	assert.True(t, os.IsNotExist(err))

	// all log file
	l.Errorln("123")
	_, err = os.Stat("./test-logs/latest-combine-test-log")
	assert.Nil(t, err)
	_, err = os.Stat("./test-logs/latest-error-test-log")
	assert.Nil(t, err)

	_ = os.RemoveAll("./test-logs")
}