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

package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

const AuthUser = "oauth2"

func Clone(dir, url, token, tagOrBranch string) error {
	var auth transport.AuthMethod
	if token != "" {
		auth = &http.BasicAuth{
			Username: AuthUser,
			Password: token,
		}
	}

	// Create the remote with repository URL
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})
	refs, err := rem.List(&git.ListOptions{
		Auth: auth,
	})
	if err != nil {
		return err
	}

	var tag plumbing.ReferenceName
	var branch plumbing.ReferenceName
	for _, ref := range refs {
		if ref.Name().Short() == tagOrBranch {
			if ref.Name().IsTag() {
				tag = ref.Name()
			} else if ref.Name().IsBranch() {
				branch = ref.Name()
			}
		}
	}

	if tag == "" && branch == "" {
		return fmt.Errorf("can not found any tag or branch name: %s", tagOrBranch)
	}

	var ref = tag
	if ref == "" {
		ref = branch
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:               url,
		Auth:              auth,
		Progress:          os.Stdout,
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ReferenceName:     ref,
	})
	return err
}
