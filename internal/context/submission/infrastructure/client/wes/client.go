package wes

import (
	"context"
	"fmt"

	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
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

var _ Client = &adaptorImpl{}

type adaptorImpl struct {
	wdlImpl      Client
	nextflowImpl Client
}

func (a *adaptorImpl) ListRuns(ctx context.Context, req *ListRunsRequest) (*ListRunsResponse, error) {
	//return a.wdlImpl.ListRuns(ctx, req)
	return nil, nil
}

func (a *adaptorImpl) RunWorkflow(ctx context.Context, req *RunWorkflowRequest) (*RunWorkflowResponse, error) {
	applog.Debugw(
		"run workflow",
		"workflowParams", req.WorkflowParams,
		"workflowType", req.WorkflowType,
		"workflowTypeVersion", req.WorkflowTypeVersion,
		"tags", req.Tags,
		"workflowEngineParameters", req.WorkflowEngineParameters,
	)
	switch req.WorkflowType {
	case consts.WorkflowTypeWDL:
		return a.wdlImpl.RunWorkflow(ctx, req)
	case consts.WorkflowTypeNextflow:
		return a.nextflowImpl.RunWorkflow(ctx, req)
	default:
		return nil, apperrors.NewInternalError(fmt.Errorf("%s not supported", req.WorkflowType))
	}
}

func (a *adaptorImpl) GetRunLog(ctx context.Context, req *GetRunLogRequest) (*GetRunLogResponse, error) {
	applog.Debugw(
		"get run log",
		"runID", req.RunID,
		"workflowType", req.WorkflowType,
	)
	switch req.WorkflowType {
	case consts.WorkflowTypeWDL:
		return a.wdlImpl.GetRunLog(ctx, req)
	case consts.WorkflowTypeNextflow:
		return a.nextflowImpl.GetRunLog(ctx, req)
	default:
		return nil, apperrors.NewInternalError(fmt.Errorf("%s not supported", req.WorkflowType))
	}
}

func (a *adaptorImpl) CancelRun(ctx context.Context, req *CancelRunRequest) (*CancelRunResponse, error) {
	applog.Debugw(
		"cancel run",
		"runID", req.RunID,
		"workflowType", req.WorkflowType,
	)
	switch req.WorkflowType {
	case consts.WorkflowTypeWDL:
		return a.wdlImpl.CancelRun(ctx, req)
	case consts.WorkflowTypeNextflow:
		return a.nextflowImpl.CancelRun(ctx, req)
	default:
		return nil, apperrors.NewInternalError(fmt.Errorf("%s not supported", req.WorkflowType))
	}
}

// NewClient ...
func NewClient(options *Options) Client {
	return &adaptorImpl{
		wdlImpl:      newWDLClient(options),
		nextflowImpl: newNextflowClient(options),
	}
}
