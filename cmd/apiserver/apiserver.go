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

package main

import (
	"math/rand"
	"os"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	bioos_server "github.com/Bio-OS/bioos/internal/apiserver"
)

//	@title			BioOS Apiserver
//	@version		1.0
//	@description	This is bioos apiserver using Hertz.

//	@contact.name	hertz-contrib
//	@contact.url	https://github.com/hertz-contrib

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/
//	@schemes					http https
//	@securityDefinitions.basic	basicAuth

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	ctx := ctrl.SetupSignalHandler()
	command := bioos_server.NewBioosServerCommand(ctx)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
