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

package errors

import "fmt"

type Error interface {
	error

	GetCode() int

	Error() string
}

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Inner   error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%d: %s: %v", e.Code, e.Message, e.Inner)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e *AppError) GetCode() int {
	return e.Code
}

func NewInvalidError(params ...string) *AppError {
	return &AppError{
		Code:    InvalidCode,
		Message: fmt.Sprintf("invalid param %v", params),
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    UnauthorizedCode,
		Message: message,
	}
}

func NewForbiddenError() *AppError {
	return &AppError{
		Code:    ForbiddenCode,
		Message: fmt.Sprintf("forbidden"),
	}
}

func NewNotFoundError(resourceType, content string) *AppError {
	return &AppError{
		Code:    NotFoundCode,
		Message: fmt.Sprintf("%s %s not found", resourceType, content),
	}
}

func NewAlreadyExistError(resourceType, content string) *AppError {
	return &AppError{
		Code:    AlreadyExistCode,
		Message: fmt.Sprintf("%s %s already exist", resourceType, content),
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:    InternalCode,
		Message: "internal system error",
		Inner:   err,
	}
}

func NewValidateFailedError(err error) *AppError {
	return &AppError{
		Code:    ValidateFailedCode,
		Message: "validate failed",
		Inner:   err,
	}
}

func NewHertzBindError(err error) *AppError {
	return &AppError{
		Code:    BindErrorCode,
		Message: "hertz bind error",
		Inner:   err,
	}
}

func NewHertzFormFileError(err error) *AppError {
	return &AppError{
		Code:    FormFileErrorCode,
		Message: "hertz get form file error",
		Inner:   err,
	}
}
