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

The package provide common logger middleware for gingonic/gin(https://github.com/gin-gonic/gin)
*/

package gingonic

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-melodic/glog"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Options The options of the common middleware
type Options struct {
	// BodyMaxSize Limit max characters of request body, default is 500
	BodyMaxSize uint
	// IgnoreExtensions Ignore some specific resources by extension in http request,
	// like ".html", ".jpg", etc.
	// default is DefaultIgnoreExtensions
	IgnoreExtensions []string
	// CustomRequest You can add customize log output from request, like some specific contents in header.
	CustomRequest func(r *http.Request) string
	// CustomResponseWriter Add customize log output from response, like some specific contents in header.
	CustomResponseWriter func(w http.ResponseWriter) string
}

// DefaultIgnoreExtensions Default ignore some specific resources by extension in http request,
// you can change it before call InjectLogger
var DefaultIgnoreExtensions = []string{".js", ".css", ".html", ".png", ".jpg",
	".jpeg", ".heic", ".gif", ".ico", ".mp3", ".mp4", ".mov", ".woff", ".ttf", ".webp", ".apng"}

// DefaultOptions For convenience usage
var DefaultOptions = &Options{
	BodyMaxSize:          500,
	IgnoreExtensions:     []string{},
	CustomRequest:        nil,
	CustomResponseWriter: nil,
}

// bodyWriter Use for hook http response
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (bw bodyWriter) Write(b []byte) (int, error) {
	bw.body.Write(b)
	return bw.ResponseWriter.Write(b)
}

func (bw bodyWriter) WriteString(s string) (int, error) {
	bw.body.WriteString(s)
	return bw.ResponseWriter.WriteString(s)
}

// InjectLogger common logger middleware for gingonic/gin
func InjectLogger(options *Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		// performance recording
		startReq := time.Now()
		path := c.Request.URL.Path
		// get request body in request with POST method
		body, err := parseRequestBody(c, options.BodyMaxSize)
		if err != nil {
			fmt.Println(err)
		}
		// log query param
		query := parseQueryParam(c.Request.URL.Query())
		requestExtInfo := ""
		if options.CustomRequest != nil {
			requestExtInfo = options.CustomRequest(c.Request) + " |"
		}
		// output request
		glog.ShareLogger().Infof("REQ -> | %15s | %s %s | %s | %s %s", c.ClientIP(), c.Request.Method,
			path, query, requestExtInfo, body)

		// parse response
		bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bw
		// deal request
		c.Next()
		excuteDurtion := time.Now().Sub(startReq)
		respBody := ""
		// ignore ext
		ignoreExt := append(DefaultIgnoreExtensions, options.IgnoreExtensions...)
		if !contains(ignoreExt, filepath.Ext(path)) {
			if c.Writer.Size() > int(options.BodyMaxSize) {
				respBody = bw.body.String()[:options.BodyMaxSize] + "..."
			} else {
				respBody = bw.body.String()
			}
		}
		responseExtInfo := ""
		if options.CustomResponseWriter != nil {
			responseExtInfo = options.CustomResponseWriter(c.Writer) + " |"
		}
		// output response
		logFunc := glog.ShareLogger().Infof
		if c.Writer.Status() > http.StatusBadRequest {
			logFunc = glog.ShareLogger().Errorf
		}
		logFunc("<- RESP | %15s | %3d | %13v | %s %s | %s %s", c.ClientIP(), c.Writer.Status(),
			excuteDurtion, c.Request.Method, path, responseExtInfo, respBody)
	}
}

func parseRequestBody(c *gin.Context, limit uint) (string, error) {
	if strings.ToUpper(c.Request.Method) != "POST" || c.Request.Body == nil {
		return "", nil
	}
	var body bytes.Buffer
	// ioutil.ReadAll for request body may clear it!!!
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return "", errors.WithMessagef(err, "[GINLOG]Read body in request %s error. %s",
			c.Request.URL.Path, err)
	}
	// resume request body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	// deal unicode char in body
	br := []rune(string(b[:]))
	body.WriteString(limitBeautyBody(br, limit))
	return body.String(), nil
}

func limitBeautyBody(body []rune, limit uint) string {
	l := len(body)
	if l <= 0 {
		return ""
	}
	// slice body content
	ellipsis := ""
	if l > int(limit) {
		l = int(limit)
		ellipsis = "..."
	}
	// make content more clear
	re := regexp.MustCompile(`\\(")|\n|\t|([{,\['])\s+`)
	return re.ReplaceAllString(string(body[:l]), "$1") + ellipsis
}

func parseQueryParam(hq url.Values) string {
	var query bytes.Buffer
	for k, v := range hq {
		if query.Len() > 0 {
			query.WriteByte('&')
		}
		query.WriteString(k)
		query.WriteByte('=')
		query.WriteString(strings.Join(v, ","))
	}
	if query.Len() <= 0 {
		query.WriteString("[EMPTY QUERY]")
	}
	return query.String()
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
