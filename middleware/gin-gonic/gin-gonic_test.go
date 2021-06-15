package gin_gonic

import (
	"context"
	"github.com/gin-gonic/gin"
	gin_logger "github.com/gin7758258/gin-logger.git"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"testing"
)

func testLogHandle(c *gin.Context) {}

func TestInjectLogger(t *testing.T) {
	const tDir = "./gin-gonic-log"
	_ = os.RemoveAll(tDir)
	err := gin_logger.InitGlobalLogger(&gin_logger.LoggerOptions{
		MinAllowLevel:    logrus.DebugLevel,
		OutputDir:        tDir,
		FilePrefix:       "gin-gonic-test",
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
		assert.FileExists(t, tDir + "/latest-combine-gin-gonic-test-log")
		_ = os.RemoveAll(tDir)

		// print request header
		req.Header.Set("SPEC_HEADER", "123")

		_ = srv.Shutdown(context.Background())

	}()
	router.GET("logTest", testLogHandle)

	_ = srv.ListenAndServe()
}
