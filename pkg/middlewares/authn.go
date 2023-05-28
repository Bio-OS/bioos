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
	"context"
	"net/http"
	"sync"
	"time"

	guardauth "github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	guardjwt "github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/shaj13/libcache"
	_ "github.com/shaj13/libcache/fifo"

	"github.com/Bio-OS/bioos/pkg/auth/authn"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

type Authenticator interface {
	Authenticate(ctx context.Context, r *http.Request) (guardauth.Info, error)
}

var (
	DefaultAuthenticator  Authenticator
	validUsers            authn.BasicUsers
	registerAuthenticator sync.Once
)

// RegisterAuthenticator register global authenticator.
func RegisterAuthenticator(opts *authn.Options) {
	registerAuthenticator.Do(func() {
		DefaultAuthenticator = newAuthenticator(opts)
	})
}

func newAuthenticator(opts *authn.Options) Authenticator {
	strategies := make([]guardauth.Strategy, 0)
	cache := libcache.FIFO.New(0)
	cache.SetTTL(time.Minute * 5)

	if opts.JWT != nil && opts.JWT.Enabled() {
		keeper := guardjwt.StaticSecret{
			ID:        opts.JWT.ID,
			Secret:    []byte(opts.JWT.Secret),
			Algorithm: opts.JWT.Algorithm,
		}

		strategies = append(strategies, guardjwt.New(cache, keeper))
	}
	if opts.Basic != nil && opts.Basic.Enabled() {
		validUsers = opts.Basic.Users
		applog.Infow("basic users", "users", opts.Basic.Users)
		strategies = append(strategies, basic.NewCached(validateUser, cache))
	}

	if len(strategies) > 0 {
		return &guardianAuthN{union.New(strategies...)}
	}
	applog.Warnf("authentication disabled and use default user name 'nobody'")
	return &nobodyAuthN{}
}

func validateUser(ctx context.Context, r *http.Request, userName, password string) (guardauth.Info, error) {
	for _, user := range validUsers {
		if user.Name == userName && user.Password == password {
			return guardauth.NewUserInfo(user.Name, user.ID, user.Groups, user.Extensions), nil
		}
	}

	return nil, apperrors.NewUnauthorizedError("invalid credentials")
}

type guardianAuthN struct {
	strategy union.Union
}

func (a *guardianAuthN) Authenticate(ctx context.Context, r *http.Request) (guardauth.Info, error) {
	_, info, err := a.strategy.AuthenticateRequest(r)
	return info, err
}

//___________________

type nobodyAuthN struct {
}

func (a *nobodyAuthN) Authenticate(ctx context.Context, r *http.Request) (guardauth.Info, error) {
	const nobody string = "nobody"
	return guardauth.NewDefaultUser(nobody, nobody, []string{nobody}, nil), nil
}
