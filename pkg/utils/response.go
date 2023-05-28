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
	"errors"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

// WriteHertzErrorResponse for error response.
func WriteHertzErrorResponse(c *app.RequestContext, err error) {
	appError := new(apperrors.AppError)
	if !errors.As(err, &appError) {
		applog.Errorf("not apperror: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	switch appError.Code {
	case apperrors.InvalidCode, apperrors.BindErrorCode:
		c.JSON(http.StatusBadRequest, appError.Message)
	case apperrors.UnauthorizedCode:
		c.JSON(http.StatusUnauthorized, appError.Message)
	case apperrors.ForbiddenCode:
		c.JSON(http.StatusForbidden, appError.Message)
	case apperrors.NotFoundCode, apperrors.RouteNotFoundCode:
		c.JSON(http.StatusNotFound, appError.Message)
	default:
		applog.Errorf("internal error: %s", appError.Inner)
		c.JSON(http.StatusInternalServerError, appError.Message)
	}
}

// WriteHertzOKResponse for all success request.
func WriteHertzOKResponse(c *app.RequestContext, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// WriteHertzCreatedResponse for create resource.
func WriteHertzCreatedResponse(c *app.RequestContext, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// WriteHertzAcceptedResponse for soft delete resource.
func WriteHertzAcceptedResponse(c *app.RequestContext) {
	c.JSON(http.StatusAccepted, nil)
}

// ToGRPCError return grpc status error
func ToGRPCError(err error) error {
	appError := new(apperrors.AppError)
	if !errors.As(err, &appError) {
		applog.Errorf("not apperror: %s", err)
		return status.Convert(err).Err()
	}
	var res error
	switch appError.Code {
	case apperrors.InvalidCode, apperrors.BindErrorCode:
		res = status.Error(codes.InvalidArgument, appError.Message)
	case apperrors.UnauthorizedCode:
		res = status.Error(codes.PermissionDenied, appError.Message)
	case apperrors.ForbiddenCode:
		res = status.Error(codes.FailedPrecondition, appError.Message)
	case apperrors.NotFoundCode, apperrors.RouteNotFoundCode:
		res = status.Error(codes.NotFound, appError.Message)
	default:
		applog.Errorf("internal error: %s", appError.Inner)
		res = status.Error(codes.Internal, appError.Message)
	}
	return res
}
