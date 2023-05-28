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

package validator

import (
	"encoding/json"
	"fmt"
	"path"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/utils"

	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

var gValidator = validator.New()

func init() {
	RegisterValidators()
}

// RegisterValidators registers k8s validators.
func RegisterValidators() error {
	validates := []struct {
		tag string
		fn  validator.Func
	}{
		{"resName", validateResName},
		{"workspaceDesc", validateWorkspaceDesc},
		{"nfsMountPath", validateNFSMountPath},
		{"dataModelName", validateDataModelName},
		{"dataModelHeaders", validateDataModelHeaders},
		{"deleteDataModelHeaders", validateDeleteDataModelHeaders},
		{"dataModelRows", validateDataModelRows},
		{"submissionName", validateSubmissionName},
		{"submissionDesc", validateSubmissionDesc},
	}

	for _, v := range validates {
		if err := gValidator.RegisterValidation(v.tag, v.fn); err != nil {
			return fmt.Errorf("register %s validator: %w", v.tag, err)
		}
	}
	return nil
}

// Validate can validate struct field with validate tag.
func Validate(s interface{}) error {
	err := gValidator.Struct(s)
	if err != nil {
		applog.Errorw("validation error", "err", err)
		if validationErrors := make(validator.ValidationErrors, 0); errors.As(err, &validationErrors) {
			fields := make([]string, 0, len(validationErrors))
			for _, validationErr := range validationErrors {
				fields = append(fields, validationErr.Field())
			}
			return apperrors.NewInvalidError(fields...)
		}
		return apperrors.NewInvalidError()
	}
	return nil
}

func validateResName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	return ValidateResNameInString(name)
}

// ValidateResNameInString ...
func ValidateResNameInString(name string) bool {
	length := utf8.RuneCountInString(name)
	if length < MinResNameLength || length > MaxResNameLength {
		return false
	}
	return ResNameRegex.Match([]byte(name))
}

func validateWorkspaceDesc(fl validator.FieldLevel) bool {
	desc := fl.Field().String()
	return ValidateWorkspaceDescInString(desc)
}

// ValidateWorkspaceDescInString ...
func ValidateWorkspaceDescInString(desc string) bool {
	length := utf8.RuneCountInString(desc)
	return length >= MinWorkspaceDescLength && length <= MaxWorkspaceDescLength
}

func validateNFSMountPath(fl validator.FieldLevel) bool {
	p := fl.Field().String()
	return ValidateNFSMountPathInString(p)
}

// ValidateNFSMountPathInString ...
func ValidateNFSMountPathInString(p string) bool {
	return path.IsAbs(p)
}

func validateDataModelName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	return ValidateDataModelNameInString(name)
}

// ValidateDataModelNameInString ...
func ValidateDataModelNameInString(name string) bool {
	length := utf8.RuneCountInString(name)
	if length < MinDataModelNameLength {
		return false
	}
	if utils.GetDataModelType(name) == consts.DataModelTypeEntitySet {
		if length > MaxEntitySetDataModelNameLength {
			return false
		}
	} else if length > MaxDataModelNameLength {
		return false
	}
	return DataModelNameReg.Match([]byte(name))
}

func validateDataModelHeaders(fl validator.FieldLevel) bool {
	count := fl.Field().Len()
	if count == 0 {
		return false
	}
	haveID := false
	isWorkspaceType := false
	dataModelName := fl.Parent().FieldByName("Name").String()
	switch utils.GetDataModelType(dataModelName) {
	case consts.DataModelTypeWorkspace:
		isWorkspaceType = validateWorkspaceTypeDataModelHeaders(count, fl)
	case consts.DataModelTypeEntitySet:
		if count < 2 || (fl.Field().Index(1).String()+consts.DataModelEntitySetNameSuffix) != dataModelName {
			return false
		}
		fallthrough
	default:
		haveID = fl.Field().Index(0).String() == utils.GenDataModelHeaderOfID(dataModelName)
	}
	if !haveID && !isWorkspaceType {
		return false
	}
	for i := 0; i < count; i++ {
		if !validateDataModelHeader(fl.Field().Index(i).String()) {
			return false
		}
	}
	return true
}

func validateDataModelHeader(header string) bool {
	matched := DataModelHeaderReg.Match([]byte(header))
	if !matched {
		return false
	}
	return len(header) <= MaxDataModelHeaderLength
}

func validateWorkspaceTypeDataModelHeaders(count int, fl validator.FieldLevel) bool {
	return count == consts.WorkspaceTypeDataModelMaxHeaderNum &&
		fl.Field().Index(0).String() == consts.WorkspaceTypeDataModelHeaderKey &&
		fl.Field().Index(1).String() == consts.WorkspaceTypeDataModelHeaderValue
}

func validateDataModelRows(fl validator.FieldLevel) bool {
	count := fl.Field().Len()
	if count == 0 {
		return false
	}
	dataModelName := fl.Parent().FieldByName("Name").String()
	typ := utils.GetDataModelType(dataModelName)
	for i := 0; i < count; i++ {
		row := fl.Field().Index(i)
		if row.Len() == 0 {
			return false
		}
		switch typ {
		case consts.DataModelTypeEntitySet:
			if row.Len() < 2 || row.Index(1).Len() == 0 {
				return false
			}
			var array []string // must be array of row id
			if err := json.Unmarshal([]byte(row.Index(1).String()), &array); err != nil {
				return false
			}
		case consts.DataModelTypeWorkspace:
			if row.Len() < 2 || row.Index(1).Len() == 0 {
				return false
			}
		}
		rowIDLength := utf8.RuneCountInString(row.Index(0).String())
		if rowIDLength == 0 || rowIDLength > MaxDataModelRowIDLength {
			return false
		}
	}
	return true
}

func validateDeleteDataModelHeaders(fl validator.FieldLevel) bool {
	for i := 0; i < fl.Field().Len(); i++ {
		if len(fl.Field().Index(i).String()) == 0 {
			return false
		}
	}
	return true
}

func validateSubmissionName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	length := utf8.RuneCountInString(name)
	if length < MinSubmissionNameLength || length > MaxSubmissionNameLength {
		return false
	}
	return ResNameRegex.Match([]byte(name))
}

func validateSubmissionDesc(fl validator.FieldLevel) bool {
	desc := fl.Field().String()
	length := utf8.RuneCountInString(desc)
	return length <= MaxSubmissionDescLength
}
