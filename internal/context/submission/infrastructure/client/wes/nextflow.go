package wes

import (
	"context"
	"fmt"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"time"
)

const (
	defaultWorkflowTypeVersion = "21.04.0"
)

var _ Client = &nextflowImpl{}

type nextflowImpl struct {
	endpoint   string
	basePath   string
	httpClient *resty.Client
}

// newNextflowClient ...
func newNextflowClient(options *Options) Client {
	client := resty.NewWithClient(http.DefaultClient).SetTimeout(time.Duration(options.Timeout) * time.Second).SetHeaders(commonClientHeaders).SetRetryCount(options.Retry)

	return &nextflowImpl{
		endpoint:   options.Nextflow.Endpoint,
		basePath:   options.Nextflow.BasePath,
		httpClient: client,
	}
}

func (n *nextflowImpl) ListRuns(ctx context.Context, req *ListRunsRequest) (*ListRunsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (n *nextflowImpl) RunWorkflow(ctx context.Context, req *RunWorkflowRequest) (*RunWorkflowResponse, error) {
	applog.Infow("run nextflow workflow")
	res := &runWorkflowResponseWithError{}
	newReq := n.httpClient.R().SetContext(ctx).SetResult(res).SetError(&res.ErrorResp)

	// nextflow engine needs upper workflow type
	req.WorkflowType = strings.ToUpper(req.WorkflowType)
	// hack here, because of the ga4gh-starter-kit-wes
	req.WorkflowTypeVersion = defaultWorkflowTypeVersion

	formData, err := runRequest2FormData(&req.RunRequest)
	if err != nil {
		return nil, newBadRequestError(err.Error())
	}
	newReq = newReq.SetFormData(formData)
	resp, err := newReq.Post(fmt.Sprintf("%s%s%s", n.endpoint, n.basePath, runWorkflowPath))
	if err != nil {
		return nil, err
	}
	if resp.IsSuccess() {
		if res.StatusCode >= 400 {
			return nil, res.ErrorResp
		}
		return &res.RunWorkflowResponse, nil
	} else if resp.IsError() {
		return nil, res.ErrorResp
	}
	return nil, fmt.Errorf("unknown http status: %s", resp.Status())
}

func (n *nextflowImpl) GetRunLog(ctx context.Context, req *GetRunLogRequest) (*GetRunLogResponse, error) {
	res := &getRunLogResponseWithError{}
	resp, err := n.httpClient.R().SetContext(ctx).SetResult(res).SetError(&res.ErrorResp).
		Get(fmt.Sprintf("%s%s%s", n.endpoint, n.basePath, fmt.Sprintf(getRunLogPath, req.RunID)))
	if err != nil {
		return nil, err
	}
	if resp.IsSuccess() {
		if res.StatusCode >= 400 {
			return nil, res.ErrorResp
		}
		return &res.GetRunLogResponse, nil
	} else if resp.IsError() {
		return nil, res.ErrorResp
	}
	return nil, fmt.Errorf("unknown http status: %s", resp.Status())
}

func (n *nextflowImpl) CancelRun(ctx context.Context, req *CancelRunRequest) (*CancelRunResponse, error) {
	res := &cancelRunResponseWithError{}
	resp, err := n.httpClient.R().SetContext(ctx).SetResult(res).SetError(&res.ErrorResp).
		Post(fmt.Sprintf("%s%s%s", n.endpoint, n.basePath, fmt.Sprintf(cancelRunPath, req.RunID)))
	if err != nil {
		return nil, err
	}
	if resp.IsSuccess() {
		if res.StatusCode >= 400 {
			return nil, res.ErrorResp
		}
		return &res.CancelRunResponse, nil
	} else if resp.IsError() {
		return nil, res.ErrorResp
	}
	return nil, fmt.Errorf("unknown http status: %s", resp.Status())
}
