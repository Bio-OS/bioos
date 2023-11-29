package wes

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func NewNextflowClient(options *Options) Client {
	client := resty.NewWithClient(http.DefaultClient).SetTimeout(time.Duration(options.Timeout) * time.Second).SetHeaders(commonClientHeaders).SetRetryCount(options.Retry)

	return &NextflowImpl{
		endpoint:   options.Endpoint,
		basePath:   options.BasePath,
		httpClient: client,
	}
}

type NextflowImpl struct {
	endpoint   string
	basePath   string
	httpClient *resty.Client
}

// ListRuns ...
func (i *NextflowImpl) ListRuns(ctx context.Context, req *ListRunsRequest) (*ListRunsResponse, error) {
	res := &listRunsResponseWithError{}
	newReq := i.httpClient.R().SetContext(ctx).SetResult(res).SetError(&res.ErrorResp)
	if req.PageSize != nil {
		newReq.SetQueryParam(pageSize, strconv.FormatInt(*req.PageSize, 10))
	}
	if req.PageToken != nil {
		newReq.SetQueryParam(pageToken, *req.PageToken)
	}
	if req.TagFilter != nil {
		newReq.SetQueryParam(tagFilter, *req.TagFilter)
	}
	resp, err := newReq.Get(fmt.Sprintf("%s%s%s", i.endpoint, i.basePath, listRunsPath))
	if err != nil {
		return nil, err
	}
	if resp.IsSuccess() {
		if res.StatusCode >= 400 {
			return nil, res.ErrorResp
		}
		return &res.ListRunsResponse, nil
	} else if resp.IsError() {
		return nil, res.ErrorResp
	}
	return nil, fmt.Errorf("unknown http status: %s", resp.Status())
}

// RunWorkflow ...
func (i *NextflowImpl) RunWorkflow(ctx context.Context, req *RunWorkflowRequest) (*RunWorkflowResponse, error) {
	res := &runWorkflowResponseWithError{}
	newReq := i.httpClient.R().SetContext(ctx).SetResult(res).SetError(&res.ErrorResp)
	filesPath := []string{}
	for key := range req.WorkflowAttachment {
		filesPath = append(filesPath, key)
	}
	if len(filesPath) == 0 {
		return nil, newBadRequestError("workflowAttachment is empty")
	}
	prefix := longestCommonPrefix(filesPath)
	if len(filesPath) == 1 {
		prefix = fmt.Sprintf("%s/", path.Dir(filesPath[0]))
	}
	// fix bug that func longestCommonPrefix() only return string level longestCommonPrefix
	// However we need path level longestCommonPrefix.
	// |   Main File  |Attachment File|String Level Prefix|Path Level Prefix|
	// |/app/tasks.wdl| /app/test.wdl |      /app/t       |      /app/      |
	if !strings.HasSuffix(prefix, "/") {
		prefix = fmt.Sprintf("%s/", path.Dir(prefix))
	}
	for _, filePath := range filesPath {
		decodeContent, err := base64.StdEncoding.DecodeString(req.WorkflowAttachment[filePath])
		if err != nil {
			return nil, fmt.Errorf("wrong wdl file: %w", err)
		}
		newReq = newReq.SetFileReader(workflowAttachment, filePath[len(prefix):], bytes.NewReader(decodeContent))
	}
	formData, err := runRequest2FormData(&req.RunRequest)
	if err != nil {
		return nil, newBadRequestError(err.Error())
	}
	newReq = newReq.SetFormData(formData)
	resp, err := newReq.Post(fmt.Sprintf("%s%s%s", i.endpoint, i.basePath, listRunsPath))
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

// GetRunLog ...
func (i *NextflowImpl) GetRunLog(ctx context.Context, req *GetRunLogRequest) (*GetRunLogResponse, error) {
	res := &getRunLogResponseWithError{}
	resp, err := i.httpClient.R().SetContext(ctx).SetResult(res).SetError(&res.ErrorResp).
		Get(fmt.Sprintf("%s%s%s", i.endpoint, i.basePath, fmt.Sprintf(getRunLogPath, req.RunID)))
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

// CancelRun ...
func (i *NextflowImpl) CancelRun(ctx context.Context, req *CancelRunRequest) (*CancelRunResponse, error) {
	res := &cancelRunResponseWithError{}
	resp, err := i.httpClient.R().SetContext(ctx).SetResult(res).SetError(&res.ErrorResp).
		Post(fmt.Sprintf("%s%s%s", i.endpoint, i.basePath, fmt.Sprintf(cancelRunPath, req.RunID)))
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
