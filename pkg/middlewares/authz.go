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

package middlewares

import (
	"fmt"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/spf13/viper"

	"github.com/Bio-OS/bioos/pkg/auth/authz"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

type Authorizer interface {
	Authorize(sub, obj, act string) (bool, error)
}

var (
	DefaultAuthorizer  Authorizer
	registerAuthorizer sync.Once
)

// RegisterAuthorizer register global authorizer.
func RegisterAuthorizer(opts *authz.Options) {
	registerAuthorizer.Do(func() {
		if opts.Casbin != nil && opts.Casbin.Enabled() {
			var err error
			DefaultAuthorizer, err = newCasbinAuthorizer(opts.Casbin)
			if err != nil {
				applog.Fatalf("failed to register global authorizer: %v", err)
			}
			applog.Infof("casbin authz enabled")
		} else {
			applog.Warnf("authorization disabled")
			DefaultAuthorizer = newNoneAuthorizer()
		}
	})
}

type casbinAuthZ struct {
	*casbin.SyncedEnforcer
}

func newCasbinAuthorizer(opts *authz.CasbinOption) (Authorizer, error) {
	// Initialize a Gorm adapter and use it in a Casbin enforcer
	var enforcer *casbin.SyncedEnforcer
	var err error
	if opts.PolicyFile != "" {
		enforcer, err = casbin.NewSyncedEnforcer(opts.ModelFile, opts.PolicyFile)
		if err != nil {
			return nil, err
		}
	} else if opts.MySQL != nil && opts.Driver != "" {
		applog.Debugw("try to new orm adapter", "driver", opts.Driver, "host", opts.MySQL.Host)
		ormAdapter, err := gormadapter.NewAdapter(opts.Driver, fmt.Sprintf("%s:%s@tcp(%s:%s)/",
			viper.GetString(opts.MySQL.Username),
			viper.GetString(opts.MySQL.Password),
			viper.GetString(opts.MySQL.Host),
			viper.GetString(opts.MySQL.Port),
		))
		if err != nil {
			return nil, err
		}

		m := model.NewModel()
		if err = m.LoadModel(opts.ModelFile); err != nil {
			return nil, fmt.Errorf("load casbin model file '%s' fail: %w", opts.ModelFile, err)
		}

		enforcer, err = casbin.NewSyncedEnforcer(m, ormAdapter)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("none casbin options to new enforcer")
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	enforcer.StartAutoLoadPolicy(5 * time.Second)

	return &casbinAuthZ{enforcer}, nil
}

// Authorize authorization impl.
func (a *casbinAuthZ) Authorize(sub, obj, act string) (bool, error) {
	return a.Enforce(sub, obj, act)
}

//__________________

type noneAuthZ struct {
}

func newNoneAuthorizer() Authorizer {
	return &noneAuthZ{}
}

func (a *noneAuthZ) Authorize(sub, obj, act string) (bool, error) {
	return true, nil
}
