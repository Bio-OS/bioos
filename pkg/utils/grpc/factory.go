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

package grpc

import (
	"github.com/Bio-OS/bioos/pkg/client"
)

type Factory interface {
	WorkspaceClient() (WorkspaceClient, error)
	WorkflowClient() (WorkflowClient, error)
	DataModelClient() (DataModelClient, error)
	VersionClient() (VersionClient, error)
}

func NewFactory(opts *client.Options) Factory {
	return factoryImpl{
		opts: opts,
	}
}

var _ Factory = factoryImpl{}

type factoryImpl struct {
	opts *client.Options
}

func (f factoryImpl) WorkspaceClient() (WorkspaceClient, error) {
	return NewWorkspaceClient(f.opts)
}

func (f factoryImpl) WorkflowClient() (WorkflowClient, error) {
	return NewWorkflowClient(f.opts)
}

func (f factoryImpl) DataModelClient() (DataModelClient, error) {
	return NewDataModelClient(f.opts)
}

func (f factoryImpl) VersionClient() (VersionClient, error) {
	return NewVersionClient(f.opts)
}
