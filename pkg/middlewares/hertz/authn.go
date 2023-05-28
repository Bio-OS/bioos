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
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/shaj13/go-guardian/v2/auth"

	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/middlewares"
)

func Authn() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if middlewares.DefaultAuthenticator == nil {
			c.Next(ctx)
			return
		}
		req, err := adaptor.GetCompatRequest(&c.Request)
		if err != nil {
			applog.Errorw("GetCompatRequest failed", "err", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		user, err := middlewares.DefaultAuthenticator.Authenticate(ctx, req)
		if err != nil {
			applog.Errorw("AuthenticateRequest failed", "err", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userCtx := auth.CtxWithUser(ctx, user)

		c.Next(userCtx)
	}
}
