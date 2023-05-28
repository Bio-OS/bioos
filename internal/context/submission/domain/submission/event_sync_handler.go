package submission

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
)

type SyncHandler struct {
	runReadModel    run.ReadModel
	repository      Repository
	dataModelClient grpc.DataModelClient
}

func NewSyncHandler(repository Repository, runReadModel run.ReadModel, dataModelClient grpc.DataModelClient) *SyncHandler {
	return &SyncHandler{
		repository:      repository,
		runReadModel:    runReadModel,
		dataModelClient: dataModelClient,
	}
}

func (h *SyncHandler) Handle(ctx context.Context, event *EventSubmission) (err error) {
	if event == nil {
		return nil
	}
	sub, err := h.repository.Get(ctx, event.SubmissionID)
	if err != nil {
		return err
	}
	statusCount, err := h.runReadModel.CountRunsResult(ctx, sub.ID)
	if err != nil {
		return err
	}
	existPending := false
	existRunning := false
	existCancelling := false
	existCancelled := false

	for _, v := range statusCount {
		switch v.Status {
		case consts.RunPending:
			existPending = true
		case consts.RunRunning:
			existRunning = true
		case consts.RunCancelling:
			existCancelling = true
		case consts.RunCancelled:
			existCancelled = true
		}
	}
	if existPending {
		sub.Status = consts.SubmissionPending
	} else if existRunning {
		sub.Status = consts.SubmissionRunning
	} else if existCancelling {
		sub.Status = consts.SubmissionCancelling
	} else if existCancelled {
		sub.Status = consts.SubmissionCancelled
	} else {
		sub.Status = consts.SubmissionFinished
	}
	if existPending || existRunning || existCancelling {
		return h.repository.Save(ctx, sub)
	}

	if sub.FinishTime == nil {
		sub.FinishTime = utils.PointTime(time.Now())
	}

	// update datamodel when all run finished
	if sub.Status == consts.SubmissionFinished && sub.Type == consts.DataModelTypeSubmission {
		if sub.DataModelID == nil {
			// should not be here, if this situation happened, it should be reported error before
			return apperrors.NewInternalError(fmt.Errorf("datamodelID must exist while submission type is datamodel"))
		}
		runList, err := h.runReadModel.ListRuns(ctx, event.SubmissionID, &utils.Pagination{}, nil)
		if err != nil {
			return err
		}
		outputsMap := make(map[string]map[string]interface{}, 0)
		for _, item := range runList {
			var tempOutput map[string]interface{}
			if item.Name != "" && item.Outputs != "" {
				if err := json.Unmarshal([]byte(item.Outputs), &tempOutput); err != nil {
					return apperrors.NewInvalidError(err.Error())
				}
				outputsMap[item.Name] = tempOutput
			}
		}
		if err := h.updateDataModelRows(ctx, outputsMap, sub.Outputs, sub.WorkspaceID, *sub.DataModelID); err != nil {
			return err
		}
	}

	return h.repository.Save(ctx, sub)
}

func (h *SyncHandler) updateDataModelRows(ctx context.Context, outputsMap map[string]map[string]interface{}, outputsCfg map[string]interface{}, workspaceID, dataModelID string) error {

	if len(outputsMap) == 0 {
		return nil
	}
	// gen headers
	headers := []string{""}
	headersMap := make(map[string]int, 0)
	for key, value := range outputsCfg {
		header, ok := checkHeaderOfOutput(value)
		if ok {
			headers = append(headers, header)
			headersMap[key] = len(headers) - 1
		}
	}
	if len(headers) == 1 {
		applog.Infof("no need to write back datamodel")
		return nil
	}
	// get datamodel name
	originDataModelResp, err := h.dataModelClient.GetDataModel(ctx, &workspaceproto.GetDataModelRequest{
		WorkspaceID: workspaceID,
		Id:          dataModelID,
	})
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	if originDataModelResp == nil || originDataModelResp.DataModel == nil {
		return apperrors.NewInternalError(fmt.Errorf("can not find datamodel with id %s in workspace %s", dataModelID, workspaceID))
	}
	if originDataModelResp.DataModel.Type != consts.DataModelTypeEntity {
		return apperrors.NewInternalError(fmt.Errorf("only support to write back datamodel with type entity"))
	}

	dmName := originDataModelResp.DataModel.Name
	headers[0] = fmt.Sprintf("%s_id", dmName)

	rows := []*workspaceproto.Row{}
	for rowID, outputs := range outputsMap {
		row := make([]string, len(headers))
		row[0] = rowID

		for name, value := range outputs {
			index, ok := headersMap[name]
			if !ok {
				continue
			}
			valueStr, err := utils.MarshalParamValue(value)
			if err != nil {
				return err
			}
			row[index] = valueStr
		}
		rows = append(rows, &workspaceproto.Row{
			Grids: row,
		})

	}

	req := &workspaceproto.PatchDataModelRequest{
		WorkspaceID: workspaceID,
		Name:        dmName,
		Headers:     headers,
		Rows:        rows,
	}

	if _, err := h.dataModelClient.PatchDataModel(ctx, req); err != nil {
		return err
	}
	return nil

}

func checkHeaderOfOutput(param interface{}) (string, bool) {
	paramStr, ok := param.(string)
	if !ok {
		return "", false
	}
	if !strings.HasPrefix(paramStr, consts.DataModelRefPrefix) {
		return "", false
	}
	header := strings.TrimPrefix(paramStr, consts.DataModelRefPrefix)
	return header, true
}
