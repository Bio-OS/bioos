//
// Copyright 2023 Beijing Volcano Engine Technology Ltd.
// Copyright 2023 Guangzhou Laboratory
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hertz

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/requestid"

	"github.com/Bio-OS/bioos/pkg/log"
)

// Logger log middleware for hertz.
func Logger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		startTime := time.Now()
		c.Next(ctx) // call the next middleware(handler)
		requestID := requestid.Get(c)
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := string(c.Request.Method())
		reqUri := string(c.Request.RequestURI())
		statusCode := c.Response.StatusCode()
		clientIP := c.ClientIP()
		log.Debugw("request log", "status_code", statusCode,
			"latency_time", latencyTime,
			"client_ip", clientIP,
			"req_method", reqMethod,
			"req_uri", reqUri,
			"requestID", requestID,
		)
	}
}
