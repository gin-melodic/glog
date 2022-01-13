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
	"bytes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Formatter struct {
	// timestamp layout, default is RFC3339Nano
	TimeStampLayout string
}

// Format extend logrus.Formatter, format logger content
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	layout := time.RFC3339Nano
	if f.TimeStampLayout != "" {
		layout = f.TimeStampLayout
	}
	ts := time.Now().Local().Format(layout)
	pid, err := getPid()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// for good performance, see https://github.com/hatlonely/hellogolang/blob/master/internal/buildin/string_test.go
	var msg bytes.Buffer
	// timestamp
	msg.WriteString(ts)
	// padding
	msg.WriteByte(' ')
	// pid
	msg.WriteString("[PID:")
	msg.WriteString(strconv.Itoa(int(pid)))
	msg.WriteByte(']')
	if entry.HasCaller() {
		// log with caller info
		msg.WriteByte('[')
		msg.WriteString(filepath.Base(entry.Caller.File))
		msg.WriteByte(':')
		msg.WriteString(strconv.Itoa(entry.Caller.Line))
		msg.WriteByte(']')
	}
	// level
	msg.WriteByte('[')
	msg.WriteString(strings.ToUpper(entry.Level.String()))
	msg.WriteByte(']')
	// logger content
	msg.WriteString(entry.Message)
	msg.WriteByte('\n')
	return msg.Bytes(), nil
}

func getPid() (pid uint64, pidErr error) {
	pb := make([]byte, 64)
	pb = pb[:runtime.Stack(pb, false)]
	pb = bytes.TrimPrefix(pb, []byte("goroutine "))
	pb = pb[:bytes.IndexByte(pb, ' ')]
	if pid, pidErr = strconv.ParseUint(string(pb), 10, 64); pidErr != nil {
		return 0, errors.WithStack(pidErr)
	}
	return
}
