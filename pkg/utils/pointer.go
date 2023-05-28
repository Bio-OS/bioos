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

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PointString convert string to point of string.
func PointString(s string) *string {
	return &s
}

// PointInt32 convert int32 to point of int32.
func PointInt32(i int32) *int32 {
	return &i
}

// PointInt64 convert int64 to point of int64.
func PointInt64(i int64) *int64 {
	return &i
}

// PointBool convert bool to point of bool.
func PointBool(b bool) *bool {
	return &b
}

// PointTime convert time.Time to point of time.Time.
func PointTime(b time.Time) *time.Time {
	return &b
}

// PointMetav1Time convert metav1.Time to point of metav1.Time.
func PointMetav1Time(b metav1.Time) *metav1.Time {
	return &b
}

// PointFloat64 convert int64 to point of float64.
func PointFloat64(f float64) *float64 {
	return &f
}
