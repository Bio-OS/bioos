package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	query "github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/utils"
)

func SubmissionPOToSubmissionDTO(ctx context.Context, submission *Submission) (*query.SubmissionItem, error) {
	item := &query.SubmissionItem{
		ID:                submission.ID,
		Name:              submission.Name,
		Description:       submission.Description,
		Type:              submission.Type,
		Status:            submission.Status,
		Language:          submission.Language,
		StartTime:         submission.StartTime.Unix(),
		WorkflowID:        submission.WorkflowID,
		WorkflowVersionID: submission.WorkflowVersionID,
		WorkspaceID:       submission.WorkspaceID,
		ExposedOptions: query.ExposedOptions{
			ReadFromCache: submission.ExposedOptions.ReadFromCache,
		},
	}
	if submission.FinishTime != nil {
		item.FinishTime = utils.PointInt64(submission.FinishTime.Unix())
		item.Duration = submission.FinishTime.Unix() - submission.StartTime.Unix()
	} else {
		item.Duration = time.Now().Unix() - submission.StartTime.Unix()
	}

	switch submission.Type {
	case consts.DataModelTypeSubmission:
		if submission.DataModelRowIDs == nil || submission.DataModelID == nil {
			return nil, fmt.Errorf("data model id & row ids should not empty")
		}
		var rowIDs []string
		var inputs, outputs []byte
		var err error
		if err := json.Unmarshal([]byte(*submission.DataModelRowIDs), &rowIDs); err != nil {
			return nil, err
		}
		if inputs, err = json.Marshal(submission.Inputs); err != nil {
			return nil, err
		}
		if outputs, err = json.Marshal(submission.Outputs); err != nil {
			return nil, err
		}
		item.Entity = &query.Entity{
			DataModelID:     *submission.DataModelID,
			DataModelRowIDs: rowIDs,
			InputsTemplate:  string(inputs),
			OutputsTemplate: string(outputs),
		}
	case consts.FilePathTypeSubmission:
		inputs, err := json.Marshal(submission.Inputs)
		if err != nil {
			return nil, err
		}
		outputsInString, err := utils.MarshalParamValue(submission.Outputs)
		if err != nil {
			return nil, err
		}
		item.InOutMaterial = &query.InOutMaterial{
			InputsMaterial:  string(inputs),
			OutputsMaterial: outputsInString,
		}
	}
	return item, nil
}

func SubmissionPOToSubmissionDO(ctx context.Context, sb *Submission) (*submission.Submission, error) {
	var rowIDs []string
	if sb.DataModelRowIDs != nil {
		if err := json.Unmarshal([]byte(*sb.DataModelRowIDs), &rowIDs); err != nil {
			return nil, err
		}
	}
	return &submission.Submission{
		ID:                sb.ID,
		Name:              sb.Name,
		Description:       sb.Description,
		WorkflowID:        sb.WorkflowID,
		WorkflowVersionID: sb.WorkflowVersionID,
		DataModelID:       sb.DataModelID,
		DataModelRowIDs:   rowIDs,
		WorkspaceID:       sb.WorkspaceID,
		Type:              sb.Type,
		Inputs:            sb.Inputs,
		Outputs:           sb.Outputs,
		ExposedOptions: submission.ExposedOptions{
			ReadFromCache: sb.ExposedOptions.ReadFromCache,
		},
		Status:     sb.Status,
		StartTime:  sb.StartTime,
		FinishTime: sb.FinishTime,
	}, nil
}

func SubmissionDOToSubmissionPO(ctx context.Context, sb *submission.Submission) (*Submission, error) {
	rowIDsInBytes, err := json.Marshal(sb.DataModelRowIDs)
	if err != nil {
		return nil, err
	}
	rowIDs := utils.PointString(string(rowIDsInBytes))
	if len(sb.DataModelRowIDs) == 0 {
		rowIDs = nil
	}
	return &Submission{
		ID:                sb.ID,
		Name:              sb.Name,
		Description:       sb.Description,
		WorkflowID:        sb.WorkflowID,
		WorkflowVersionID: sb.WorkflowVersionID,
		DataModelID:       sb.DataModelID,
		DataModelRowIDs:   rowIDs,
		Type:              sb.Type,
		Inputs:            sb.Inputs,
		Outputs:           sb.Outputs,
		WorkspaceID:       sb.WorkspaceID,
		ExposedOptions: ExposedOptions{
			ReadFromCache: sb.ExposedOptions.ReadFromCache,
		},
		Status:     sb.Status,
		Language:   sb.Language,
		StartTime:  sb.StartTime,
		FinishTime: sb.FinishTime,
	}, nil
}
