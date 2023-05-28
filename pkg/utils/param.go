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

package utils

import "encoding/json"

// UnmarshalParamValue ...
func UnmarshalParamValue(value string) interface{} {
	var res interface{}
	if err := json.Unmarshal([]byte(value), &res); err != nil {
		return value // err != nil means String/File type
	}
	return res
}

// MarshalParamValue ...
func MarshalParamValue(value interface{}) (string, error) {
	if valueStr, ok := value.(string); ok {
		return valueStr, nil
	}
	valueJSON, err := json.Marshal(value)
	return string(valueJSON), err
}
