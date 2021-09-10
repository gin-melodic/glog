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

package gingonic

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin7758258/glog"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"testing"
)

func testLogHandle(c *gin.Context) {}

func TestInjectLogger(t *testing.T) {
	const tDir = "./gingonic-log"
	_ = os.RemoveAll(tDir)
	err := glog.InitGlobalLogger(&glog.LoggerOptions{
		MinAllowLevel:    logrus.DebugLevel,
		OutputDir:        tDir,
		FilePrefix:       "gingonic-test",
		SaveDay:          1,
		ExtLoggerWriter:  []io.Writer{os.Stdout},
	})
	assert.Nil(t, err)
	// prepare
	router := gin.New()
	router.Use(InjectLogger(&Options{
		BodyMaxSize: 500,
		CustomRequest: func(r *http.Request) string {
			return "SPEC_HEADER=" + r.Header.Get("SPEC_HEADER")
		},
	}))
	srv := &http.Server{
		Addr:    ":8180",
		Handler: router,
	}

	go func() {
		c := &http.Client{}
		req, _ := http.NewRequest("GET", "http://127.0.0.1:8180/logTest", nil)

		_, _ = c.Do(req)
		// check log
		assert.FileExists(t, tDir + "/latest-combine-gingonic-test-log")
		_ = os.RemoveAll(tDir)

		// print request header
		req.Header.Set("SPEC_HEADER", "123")

		_ = srv.Shutdown(context.Background())

	}()
	router.GET("logTest", testLogHandle)

	_ = srv.ListenAndServe()
}
