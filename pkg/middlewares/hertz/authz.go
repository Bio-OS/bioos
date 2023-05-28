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
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/shaj13/go-guardian/v2/auth"

	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/middlewares"
)

type PermissionFunc func(ctx context.Context, c *app.RequestContext) string

// Authz authorization middleware for hertz.
func Authz(permissionFunc PermissionFunc) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		permission := permissionFunc(ctx, c)
		applog.Infow("try to auth", "permission", permission)
		if permission == "" {
			c.Next(ctx)
			return
		}
		// Look up current subject.
		user := auth.UserFromCtx(ctx)
		if user == nil {
			applog.Debugw("user not authenticated")
			c.AbortWithStatus(consts.StatusUnauthorized)
			return
		}

		permissions := strings.Split(permission, ":")
		if len(permissions) < 2 {
			c.AbortWithStatus(consts.StatusInternalServerError)
			return
		}

		allowed, err := middlewares.DefaultAuthorizer.Authorize(user.GetUserName(), permissions[0], permissions[1])
		if err != nil {
			applog.Errorf("authorize fail %s", err)
			c.AbortWithStatus(consts.StatusInternalServerError)
			return
		}
		if !allowed {
			applog.Debugw("user not allowed", "obj", permissions[0], "action", permissions[1], "user", user)
			c.AbortWithStatus(consts.StatusForbidden)
			return
		}

		applog.Debugw("user permitted to enter", "obj", permissions[0], "action", permissions[1], "user", user)
		c.Next(ctx)
		return
	}
}
