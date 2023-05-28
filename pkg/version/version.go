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

package version

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gosuri/uitable"
)

var (
	// Version is the code version.
	Version string
	// GitBranch code branch.
	GitBranch string
	// GitCommit is the git commit.
	GitCommit string
	// GitTreeState clean/dirty.
	GitTreeState string
	// BuildTime is the build time.
	BuildTime string
)

type Info struct {
	Version      string `json:"version"`
	GitBranch    string `json:"gitBranch"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildTime    string `json:"buildTime"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// String return the string of Info.
func (info Info) String() string {
	table := uitable.New()
	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "
	table.AddRow("version:", info.Version)
	table.AddRow("gitCommit:", info.GitCommit)
	table.AddRow("gitBranch:", info.GitBranch)
	table.AddRow("gitTreeState:", info.GitTreeState)
	table.AddRow("buildTime:", info.BuildTime)
	table.AddRow("goVersion:", info.GoVersion)
	table.AddRow("compiler:", info.Compiler)
	table.AddRow("platform:", info.Platform)

	return table.String()
}

func Get() Info {
	return Info{
		Version:      Version,
		GitBranch:    GitBranch,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildTime:    BuildTime,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// PrintVersionOrContinue will print git commit and exit with os.Exit(0) if CLI v flag is present.
func PrintVersionOrContinue() {
	fmt.Printf("%s\n", Get())
	if versionFlag {
		os.Exit(0)
	}
}
