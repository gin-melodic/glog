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

The package provide common logger middleware for go-gorm/gorm(https://github.com/go-gorm/gorm)

Usage (e.g. sqlite):

```go
import (
	gormLogger "github.com/gin7758258/glog/middleware/gorm"
)

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	Logger: gormLogger.New(gormLogger.Options{}),
})
```
*/

package gorm

import (
	"context"
	"errors"
	"github.com/gin7758258/glog"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type Options struct {
	SlowThreshold         time.Duration
	SourceField           string
	IgnoreRecordNotFoundError bool
}

type sqlLogger struct {
	Options
}

func New(opt Options) *sqlLogger {
	return &sqlLogger{
		opt,
	}
}

func (l *sqlLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *sqlLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	glog.ShareLogger().WithContext(ctx).Infof(msg, data...)
}

func (l *sqlLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	glog.ShareLogger().WithContext(ctx).Warnf(msg, data...)
}

func (l *sqlLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	glog.ShareLogger().WithContext(ctx).Errorf(msg, data...)
}

func (l *sqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	dt := time.Since(begin)
	sql, _ := fc()
	f := logrus.Fields{}
	// gorm source field config
	if l.SourceField != "" {
		f[l.SourceField] = utils.FileWithLineNum()
	}
	// throw error, ignore empty result error
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.IgnoreRecordNotFoundError) {
		f[logrus.ErrorKey] = err
		glog.ShareLogger().WithContext(ctx).WithFields(f).Errorf("[SQL Error][cost %s] %s", dt, sql)
		return
	}
	// check slow threshold
	if l.SlowThreshold > 0 && dt > l.SlowThreshold {
		glog.ShareLogger().WithContext(ctx).WithFields(f).Warnf("[Slow SQL][cost %s] %s", dt, sql)
		return
	}
	// debug mode
	glog.ShareLogger().WithContext(ctx).WithFields(f).Debugf("[SQL][cost %s] %s", dt, sql)
}


