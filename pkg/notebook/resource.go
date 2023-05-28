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

package notebook

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	GPUVendorUnknown = "Unknown"
	GPUVendorNvidia  = "Nvidia"
	GPUVendorAMD     = "AMD"   // not support
	GPUVendorIntel   = "Intel" // not support
)

func GetGPUVendor(model string) string {
	if strings.HasPrefix(model, "Nvidia") {
		return GPUVendorNvidia
	}
	return GPUVendorUnknown
}

type GPU struct {
	Model  string  `json:"model" mapstructure:"model"`
	Card   float64 `json:"card" mapstructure:"card"` // float is for mgpu
	Memory int64   `json:"memory" mapstructure:"memory"`
}

func (g *GPU) Vendor() string {
	return GetGPUVendor(g.Model)
}

type ResourceSize struct {
	CPU    float64 `json:"cpu" mapstructure:"cpu"`
	Memory int64   `json:"memory" mapstructure:"memory"`
	Disk   int64   `json:"disk" mapstructure:"disk"`
	GPU    *GPU    `json:"gpu,omitempty" mapstructure:"gpu"`
}

type configFormat struct {
	CPU    float64           `json:"cpu" mapstructure:"cpu"`
	Memory resource.Quantity `json:"memory" mapstructure:"memory"`
	Disk   resource.Quantity `json:"disk" mapstructure:"disk"`
	GPU    *struct {
		Model  string            `json:"model" mapstructure:"model"`
		Card   float64           `json:"card" mapstructure:"card"`
		Memory resource.Quantity `json:"memory" mapstructure:"memory"`
	} `json:"gpu,omitempty" mapstructure:"gpu"`
}

func (v *configFormat) toResourceSize(r *ResourceSize) {
	r.CPU = v.CPU
	r.Memory = v.Memory.Value()
	r.Disk = v.Disk.Value()
	if v.GPU != nil {
		r.GPU = &GPU{
			Model:  v.GPU.Model,
			Card:   v.GPU.Card,
			Memory: v.GPU.Memory.Value(),
		}
	}
}

func (r *ResourceSize) String() string {
	return fmt.Sprintf("%+v", *r)
}

func (r *ResourceSize) UnmarshalJSON(b []byte) error {
	var v configFormat
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	v.toResourceSize(r)
	return nil
}

func (r *ResourceSize) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v configFormat
	if err := unmarshal(&v); err != nil {
		return err
	}
	v.toResourceSize(r)
	return nil
}

// ResourceQuantityStringToInt64HookFunc converts resource.Quantity format string to int64.
func ResourceQuantityStringToInt64HookFunc(f, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	if t.Kind() != reflect.Int64 {
		return data, nil
	}

	// Convert it by parsing
	q, err := resource.ParseQuantity(data.(string))
	if err != nil {
		return nil, err
	}
	return q.Value(), nil
}
