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

package schema

const (
	TypeString  = "string"
	TypeNumber  = "number"
	TypeInteger = "integer"
	TypeBoolean = "boolean"
)

// NextflowSchema ...
type NextflowSchema struct {
	Schema      string `json:"$schema"`
	ID          string `json:"$id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Definitions map[string]DefinitionSchema
}

// PropertySchema ...
type PropertySchema struct {
	Type        string      `json:"type"`
	Out         bool        `json:"out"`
	Format      string      `json:"format"`
	Description string      `json:"description"`
	MIMEType    string      `json:"mimetype"`
	Default     interface{} `json:"default"`
}

// DefinitionSchema ...
type DefinitionSchema struct {
	Title      string                    `json:"title"`
	Type       string                    `json:"type"`
	Required   []string                  `json:"required"`
	Properties map[string]PropertySchema `json:"properties"`
}
