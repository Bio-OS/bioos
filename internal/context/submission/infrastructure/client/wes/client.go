package wes

import (
	"context"
	"fmt"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
)

// Request Params of Wes ...
const (
	pageSize                 = "page_size"
	pageToken                = "page_token"
	tagFilter                = "tag_filter"
	workflowParams           = "workflow_params"
	workflowType             = "workflow_type"
	workflowTypeVersion      = "workflow_type_version"
	tags                     = "tags"
	workflowEngineParameters = "workflow_engine_parameters"
	workflowURL              = "workflow_url"
	workflowAttachment       = "workflow_attachment"
)

// Path of Wes ...
const (
	listRunsPath     = "/runs"
	runWorkflowPath  = "/runs"
	getRunLogPath    = "/runs/%s"
	cancelRunPath    = "/runs/%s/cancel"
	getRunStatusPath = "/runs/%s/status"
)

// CommonClientHeaders ...
var commonClientHeaders = map[string]string{
	"Accept":          "*/*",
	"Accept-Encoding": "gzip, deflate, br",
	"Connection":      "keep-alive",
}

// Client ...
type Client interface {
	ListRuns(ctx context.Context, req *ListRunsRequest) (*ListRunsResponse, error)
	RunWorkflow(ctx context.Context, req *RunWorkflowRequest) (*RunWorkflowResponse, error)
	GetRunLog(ctx context.Context, req *GetRunLogRequest) (*GetRunLogResponse, error)
	CancelRun(ctx context.Context, req *CancelRunRequest) (*CancelRunResponse, error)
}

func NewClient(options *Options) Client {
	return &impl{
		WDLImpl:      NewWDLClient(options),
		NextflowImpl: NewNextflowClient(options),
	}
}

type impl struct {
	WDLImpl      Client
	NextflowImpl Client
}

func (i *impl) ListRuns(ctx context.Context, req *ListRunsRequest) (*ListRunsResponse, error) {
	switch req.WorkflowType {
	case workflow.LanguageWDL:
		return i.WDLImpl.ListRuns(ctx, req)
	case workflow.LanguageNextflow:
		return i.NextflowImpl.ListRuns(ctx, req)
	default:
		return nil, apperrors.NewInternalError(fmt.Errorf("%s not supported", req.WorkflowType))
	}
}

func (i *impl) RunWorkflow(ctx context.Context, req *RunWorkflowRequest) (*RunWorkflowResponse, error) {
	switch req.RunRequest.WorkflowType {
	case workflow.LanguageWDL:
		return i.WDLImpl.RunWorkflow(ctx, req)
	case workflow.LanguageNextflow:
		return i.NextflowImpl.RunWorkflow(ctx, req)
	default:
		return nil, apperrors.NewInternalError(fmt.Errorf("%s not supported", req.WorkflowType))
	}
}

func (i *impl) GetRunLog(ctx context.Context, req *GetRunLogRequest) (*GetRunLogResponse, error) {
	switch req.WorkflowType {
	case workflow.LanguageWDL:
		return i.WDLImpl.GetRunLog(ctx, req)
	case workflow.LanguageNextflow:
		return i.NextflowImpl.GetRunLog(ctx, req)
	default:
		return nil, apperrors.NewInternalError(fmt.Errorf("%s not supported", req.WorkflowType))
	}
}

func (i *impl) CancelRun(ctx context.Context, req *CancelRunRequest) (*CancelRunResponse, error) {
	switch req.WorkflowType {
	case workflow.LanguageWDL:
		return i.WDLImpl.CancelRun(ctx, req)
	case workflow.LanguageNextflow:
		return i.NextflowImpl.CancelRun(ctx, req)
	default:
		return nil, apperrors.NewInternalError(fmt.Errorf("%s not supported", req.WorkflowType))
	}
}
