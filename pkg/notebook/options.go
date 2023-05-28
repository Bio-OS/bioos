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
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

type Image struct {
	Name        string    `json:"name" mapstructure:"name"`
	Version     string    `json:"version" mapstructure:"version"`
	Description string    `json:"description" mapstructure:"description"`
	Image       string    `json:"image" mapstructure:"image"`
	UpdateTime  time.Time `json:"updateTime" mapstructure:"updateTime"`
}

type jupyterhubInKubeConfig struct {
	KubeconfigPath string `json:"kubeconfigPath" mapstructure:"kubeconfigPath"`
	MasterURL      string `json:"masterURL" mapstructure:"masterURL"`
	Namespace      string `json:"namespace" mapstructure:"namespace"`
	StorageClass   string `json:"storageClass" mapstructure:"storageClass"`
}

type JupyterhubConfig struct {
	Endpoint   string                  `json:"endpoint" mapstructure:"endpoint"`
	AdminToken string                  `json:"adminToken" mapstructure:"adminToken"`
	Kubernetes *jupyterhubInKubeConfig `json:"kubernetes" mapstructure:"kubernetes"`
}

type ResourceOption struct {
	ResourceSize `json:",inline" mapstructure:",squash"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty" mapstructure:"nodeSelector"`
}

type Options struct {
	OfficialImages   []Image          `json:"officialImages" mapstructure:"officialImages"`
	ResourceSizes    []ResourceOption `json:"resourceOptions" mapstructure:"resourceOptions"`
	StaticJupyterhub JupyterhubConfig `json:"staticJupyterhub" mapstructure:"staticJupyterhub"`
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Validate() error {
	if o.StaticJupyterhub.Endpoint != "" {
		if o.StaticJupyterhub.AdminToken == "" {
			return fmt.Errorf("staticJupyterhub required adminToken")
		}
	}
	if len(o.ResourceSizes) == 0 {
		return fmt.Errorf("none notebook resource size options")
	}
	for _, s := range o.ResourceSizes {
		if s.GPU != nil && s.GPU.Vendor() == GPUVendorUnknown {
			return fmt.Errorf("GPU '%s' vendor is unknown", s.GPU.Model)
		}
	}
	return nil
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.StaticJupyterhub.Endpoint, "jupyterhub-endpoint", "", "static jupyterhub endpoint e.g. http://localhost/hub")
	fs.StringVar(&o.StaticJupyterhub.AdminToken, "jupyterhub-token", "", "static jupyterhub token")
}

func (o *Options) ListOfficialImages() []string {
	res := make([]string, len(o.OfficialImages))
	for i := range o.OfficialImages {
		res[i] = o.OfficialImages[i].Image
	}
	return res
}

func (o *Options) ListResourceSizes() []ResourceSize {
	res := make([]ResourceSize, len(o.ResourceSizes))
	for i := range o.ResourceSizes {
		res[i] = o.ResourceSizes[i].ResourceSize
	}
	return res
}
