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

func DeleteStrSliceElms(sl []string, elms ...string) []string {
	if len(sl) == 0 || len(elms) == 0 {
		return sl
	}
	m := make(map[string]struct{})
	for _, v := range elms {
		m[v] = struct{}{}
	}
	res := make([]string, 0, len(sl))
	for _, v := range sl {
		if _, ok := m[v]; !ok {
			res = append(res, v)
		}
	}
	return res
}

// In ...
func In(elm string, elms []string) bool {
	for _, item := range elms {
		if elm == item {
			return true
		}
	}
	return false
}
