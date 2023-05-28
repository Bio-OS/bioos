package submission

import (
	"context"
	"reflect"

	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
)

const WESTag = "wes"

type CreateHandler struct {
	workflowClient grpc.WorkflowClient
	repository     Repository
	eventbus       eventbus.EventBus
}

func NewCreateHandler(repository Repository, eventbus eventbus.EventBus, workflowClient grpc.WorkflowClient) *CreateHandler {
	return &CreateHandler{
		repository:     repository,
		workflowClient: workflowClient,
		eventbus:       eventbus,
	}
}

func (h *CreateHandler) Handle(ctx context.Context, event *CreateEvent) (err error) {
	if event == nil {
		return nil
	}
	sub, err := h.repository.Get(ctx, event.SubmissionID)
	if err != nil {
		return err
	}
	// todo we should store the data model & workflow used in submission
	createRunEvent, err := h.genCreateRunEvent(ctx, event, sub)
	if err != nil {
		return err
	}

	return h.eventbus.Publish(ctx, createRunEvent)
}

func (h *CreateHandler) genCreateRunEvent(ctx context.Context, event *CreateEvent, sub *Submission) (*EventCreateRuns, error) {
	getWorkflowResp, err := h.workflowClient.GetWorkflow(ctx, &workspaceproto.GetWorkflowRequest{
		Id:          event.SourceWorkflowID,
		WorkspaceID: event.WorkspaceID,
	})
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	getWorkflowVersionResp, err := h.workflowClient.GetWorkflowVersion(ctx, &workspaceproto.GetWorkflowVersionRequest{
		Id:          getWorkflowResp.Workflow.LatestVersion.Id,
		WorkflowID:  event.SourceWorkflowID,
		WorkspaceID: event.WorkspaceID,
	})
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	workflowEngineParameters := exposedOptions2Map(&sub.ExposedOptions)
	files, err := h.genWorkflowFiles(ctx, getWorkflowVersionResp.Version, event)
	if err != nil {
		return nil, err
	}
	return NewEventCreateRuns(sub.WorkspaceID, sub.ID, sub.Type, sub.Inputs, sub.Outputs, sub.DataModelID, sub.DataModelRowIDs, &RunConfig{
		Language:                 getWorkflowVersionResp.Version.Language,
		MainWorkflowFilePath:     getWorkflowVersionResp.Version.MainWorkflowPath,
		WorkflowContents:         files,
		WorkflowEngineParameters: workflowEngineParameters,
		Version:                  getWorkflowVersionResp.Version.LanguageVersion,
	}), nil

}

func (h *CreateHandler) genWorkflowFiles(ctx context.Context, workflowVersion *workspaceproto.WorkflowVersion, event *CreateEvent) (workflowFiles map[string]string, err error) {
	ids := make([]string, 0)
	for _, fileInfo := range workflowVersion.Files {
		ids = append(ids, fileInfo.Id)
	}
	ListWorkflowFilesResponse, err := h.workflowClient.ListWorkflowFiles(ctx, &workspaceproto.ListWorkflowFilesRequest{
		Page:              1,
		Size:              int32(len(ids)),
		Ids:               ids,
		WorkspaceID:       event.WorkspaceID,
		WorkflowID:        event.SourceWorkflowID,
		WorkflowVersionID: utils.PointString(workflowVersion.Id),
	})
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	workflowFiles = make(map[string]string)
	for _, file := range ListWorkflowFilesResponse.Files {
		workflowFiles[file.Path] = file.Content
	}
	return workflowFiles, nil
}

// exposedOptions2Map ...
func exposedOptions2Map(exposedOptions *ExposedOptions) map[string]interface{} {
	res := make(map[string]interface{})
	t := reflect.TypeOf(exposedOptions).Elem()
	v := reflect.ValueOf(exposedOptions).Elem()
	for i := 0; i < v.NumField(); i++ {
		res[t.Field(i).Tag.Get(WESTag)] = v.Field(i).Interface()
	}
	return res
}
