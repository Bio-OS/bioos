package run

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
)

// save submission -> save runs
type EventHandlerCreateRuns struct {
	runRepo         Repository
	runFactory      Factory
	dataModelClient grpc.DataModelClient
	eventBus        eventbus.EventBus
}

func NewEventHandlerCreateRuns(repo Repository, client grpc.DataModelClient, eventBus eventbus.EventBus, runFactory Factory) *EventHandlerCreateRuns {
	return &EventHandlerCreateRuns{
		runRepo:         repo,
		runFactory:      runFactory,
		dataModelClient: client,
		eventBus:        eventBus,
	}
}

func (e *EventHandlerCreateRuns) Handle(ctx context.Context, event *submission.EventCreateRuns) error {

	// gen dataList
	dataList, err := e.genDataList(ctx, event)
	if err != nil {
		return err
	}

	// gen run
	runList, err := e.genRunList(ctx, dataList, event)
	if err != nil {
		return err
	}

	for _, run := range runList {
		// public submit run
		if err = e.Publish(ctx, run, event); err != nil {
			return err
		}
	}

	return nil
}

func (e *EventHandlerCreateRuns) genRunList(ctx context.Context, dataList *dataList, event *submission.EventCreateRuns) ([]*Run, error) {
	res := make([]*Run, 0)
	// only contains ws data model implement submissionType is filePath
	if len(event.DataModelRowIDs) == 0 {
		if event.SubmissionType != consts.FilePathTypeSubmission {
			return nil, apperrors.NewInternalError(fmt.Errorf("only submissionType is filepath when with no datamodel"))
		}
		for name, vInputs := range event.InputsTemplate {
			vInputsMap, ok := vInputs.(map[string]interface{})
			if !ok {
				return nil, apperrors.NewInternalError(fmt.Errorf("wrong input of type filePath: %v", vInputs))
			}
			inputStr, err := genInput(dataList, vInputsMap, "")
			if err != nil {
				return nil, err
			}
			tempRun, err := e.runFactory.CreateWithRunParam(CreateRunParam{
				SubmissionID: event.SubmissionID,
				Name:         name,
				Inputs:       inputStr,
				Status:       consts.RunPending,
			})
			if err != nil {
				return nil, err
			}
			res = append(res, tempRun)
		}
		return res, nil
	}
	// with entity data model
	for _, rowID := range event.DataModelRowIDs {
		tempRowID := rowID
		inputStr, err := genInput(dataList, event.InputsTemplate, tempRowID)
		if err != nil {
			return nil, err
		}
		tempRun, err := e.runFactory.CreateWithRunParam(CreateRunParam{
			SubmissionID: event.SubmissionID,
			Name:         rowID,
			Inputs:       inputStr,
			Status:       consts.RunPending,
		})
		if err != nil {
			return nil, err
		}
		res = append(res, tempRun)
	}
	return res, nil
}

func (e *EventHandlerCreateRuns) genDataList(ctx context.Context, event *submission.EventCreateRuns) (*dataList, error) {
	wsModelData, err := e.genWsModelData(ctx, event.WorkspaceID, event.InputsTemplate, event.SubmissionType)
	if err != nil {
		return nil, err
	}
	if event.DataModelID == nil {
		return &dataList{
			workspaceModel: wsModelData,
		}, nil
	}
	entityModel, setModelList, err := e.genEntityModelData(ctx, event.WorkspaceID, *event.DataModelID, event.DataModelRowIDs)
	if err != nil {
		return nil, err
	}
	return &dataList{
		workspaceModel: wsModelData,
		entityModel:    entityModel,
		setModelChain:  setModelList,
	}, nil
}

func (e *EventHandlerCreateRuns) genEntityModelData(ctx context.Context, workspaceID, dataModelID string, rowIDs []string) (*dataModel, []setModel, error) {
	originDataModelResp, err := e.dataModelClient.GetDataModel(ctx, &workspaceproto.GetDataModelRequest{
		WorkspaceID: workspaceID,
		Id:          dataModelID,
	})
	if err != nil {
		return nil, nil, apperrors.NewInternalError(err)
	}
	if originDataModelResp == nil || originDataModelResp.DataModel == nil {
		return nil, nil, apperrors.NewInternalError(err)
	}
	originDataModel := originDataModelResp.DataModel
	originRowsResp, err := e.dataModelClient.ListDataModelRows(ctx, &workspaceproto.ListDataModelRowsRequest{
		Id:          originDataModel.Id,
		WorkspaceID: workspaceID,
		RowIDs:      rowIDs,
	})
	if err != nil {
		return nil, nil, apperrors.NewInternalError(err)
	}
	switch originDataModel.Type {
	case consts.DataModelTypeEntity:
		return NewDataModel(originDataModel, originRowsResp.Headers, originRowsResp.Rows), nil, nil
	case consts.DataModelTypeEntitySet:
		// get all set model name and final entity name
		setModelNameList := []string{}
		modelName := originDataModel.Name
		for strings.HasSuffix(modelName, consts.DataModelEntitySetNameSuffix) {
			modelName = strings.TrimSuffix(modelName, consts.DataModelEntitySetNameSuffix)
			setModelNameList = append(setModelNameList, modelName)
		}
		// get branch of items of all datamodels
		m, err := e.dataModelClient.ListDataModels(ctx, &workspaceproto.ListDataModelsRequest{
			WorkspaceID: workspaceID,
			SearchWord:  modelName,
		})
		if err != nil {
			return nil, nil, apperrors.NewInternalError(err)
		}
		setModelNameMap := make(map[string]*workspaceproto.DataModel, 0)
		for _, dm := range m.Items {
			tempDm := dm
			if _, ok := setModelNameMap[tempDm.Name]; !ok {
				setModelNameMap[tempDm.Name] = tempDm
			}
		}
		var finalDataModels *dataModel
		setModelChain := []setModel{
			{
				NewDataModel(originDataModel, originRowsResp.Headers, nil),
				parseSetIDList(originDataModel.Name, originRowsResp.Headers, originRowsResp.Rows),
			},
		}
		for _, name := range setModelNameList {
			dm, ok := setModelNameMap[name]
			if !ok {
				return nil, nil, apperrors.NewInternalError(fmt.Errorf("cannot get datamodel with name %s in workspace %s", name, workspaceID))
			}
			searchRowIDs := collectIDList(setModelChain[len(setModelChain)-1].sets)
			resp, err := e.dataModelClient.ListDataModelRows(ctx, &workspaceproto.ListDataModelRowsRequest{
				WorkspaceID: workspaceID,
				Id:          dm.Id,
				RowIDs:      searchRowIDs,
				Page:        1,
				Size:        int32(len(searchRowIDs)),
			})
			if err != nil {
				return nil, nil, apperrors.NewInternalError(err)
			}
			if dm.Type != consts.DataModelTypeEntitySet {
				finalDataModels = NewDataModel(dm, resp.Headers, resp.Rows)
				break
			}
			setModelChain = append(setModelChain, setModel{NewDataModel(dm, resp.Headers, nil), parseSetIDList(dm.Name, resp.Headers, resp.Rows)})
		}
		return finalDataModels, setModelChain, nil
	default:
		applog.Errorf("this type %s should not be handle here", originDataModel.Type)
		return nil, nil, apperrors.NewInternalError(err)
	}
}

func (e *EventHandlerCreateRuns) genWsModelData(ctx context.Context, workspaceID string, inputs map[string]interface{}, submissionType string) (*dataModel, error) {
	// get ws datamodel
	wsDataModelResp, err := e.dataModelClient.ListDataModels(ctx, &workspaceproto.ListDataModelsRequest{
		WorkspaceID: workspaceID,
		Types:       []string{consts.DataModelTypeWorkspace},
	})
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	// each workspace can only have single datamodal with type workspace
	if wsDataModelResp == nil || len(wsDataModelResp.Items) > 1 {
		return nil, apperrors.NewInternalError(err)
	}

	if len(wsDataModelResp.Items) == 0 {
		applog.Infof("the workspace with ID %s has no workspace data model", workspaceID)
		return nil, nil
	}

	var wsDataIDs []string
	switch submissionType {
	case consts.FilePathTypeSubmission:
		for _, vInputs := range inputs {
			vInputsMap, ok := vInputs.(map[string]interface{})
			if !ok {
				return nil, apperrors.NewInternalError(fmt.Errorf("wrong input of type filePath: %v", vInputs))
			}
			for _, v := range vInputsMap {
				value, err := utils.MarshalParamValue(v)
				if err != nil {
					return nil, apperrors.NewInternalError(err)
				}
				if strings.HasPrefix(value, consts.WorkspaceTypeDataModelRefPrefix) {
					wsDataIDs = append(wsDataIDs, strings.Split(value, ".")[1])
				}
			}
		}

	case consts.DataModelTypeSubmission:
		for _, v := range inputs {
			value, err := utils.MarshalParamValue(v)
			if err != nil {
				return nil, apperrors.NewInternalError(err)
			}
			if strings.HasPrefix(value, consts.WorkspaceTypeDataModelRefPrefix) {
				wsDataIDs = append(wsDataIDs, strings.Split(value, ".")[1])
			}
		}
	default:
		return nil, apperrors.NewInternalError(fmt.Errorf("not support submission type %s", submissionType))
	}

	if len(wsDataIDs) == 0 {
		applog.Infof("the submission does not need workspace data model")
		return nil, nil
	}
	// todo remove repetition
	wsDataModelRows, err := e.dataModelClient.ListDataModelRows(ctx, &workspaceproto.ListDataModelRowsRequest{
		WorkspaceID: workspaceID,
		Id:          wsDataModelResp.Items[0].Id,
		RowIDs:      wsDataIDs,
		Page:        1,
		Size:        int32(len(wsDataIDs)),
	})
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return NewDataModel(wsDataModelResp.Items[0], wsDataModelRows.Headers, wsDataModelRows.Rows), nil
}

func (e *EventHandlerCreateRuns) Publish(ctx context.Context, run *Run, event *submission.EventCreateRuns) error {
	if err := e.runRepo.Save(ctx, run); err != nil {
		return err
	}
	eventSubmitRun := submission.NewEventSubmitRun(run.ID, event.RunConfig)
	if err := e.eventBus.Publish(ctx, eventSubmitRun); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func genInput(dataList *dataList, inputMap map[string]interface{}, rowID string) (map[string]interface{}, error) {
	renderedMap := make(map[string]interface{}, len(inputMap))
	for key, value := range inputMap {
		if valueStr, ok := value.(string); ok {
			rendered, err := dataList.pickInputValue(rowID, valueStr)
			if err != nil {
				return nil, err
			}
			renderedMap[key] = rendered
		} else {
			renderedMap[key] = value
		}
	}
	return renderedMap, nil
}

func parseSetIDList(modelName string, headers []string, rows []*workspaceproto.Row) map[string][]string {
	if len(rows) < 1 || headers[1]+consts.DataModelEntitySetNameSuffix != modelName {
		return nil
	}
	res := make(map[string][]string)
	for _, r := range rows {
		var list []string
		if err := json.Unmarshal([]byte(r.Grids[1]), &list); err != nil {
			// data will check before writing so should not go into this
			applog.Warnf("can not decode data model %s id list of %s: %s", modelName, r.Grids[0], err)
		} else {
			res[r.Grids[0]] = list
		}
	}
	return res
}

func collectIDList(m map[string][]string) (res []string) {
	if len(m) > 0 {
		check := make(map[string]bool)
		res = make([]string, 0)
		for _, v := range m {
			for _, s := range v {
				if _, ok := check[s]; !ok {
					check[s] = true
					res = append(res, s)
				}
			}
		}
	}
	return
}

// todo move it to a seperate place

type dataModel struct {
	ID           string
	Name         string
	Type         string
	RowHeaderMap map[string]int
	Rows         []*row
}
type row struct {
	Grids []string
}

type setModel struct {
	model *dataModel
	// set id -> ref id list
	sets map[string][]string
}

func (d *dataModel) rowLength() int {
	return len(d.Rows)
}

func (d *dataModel) pickData(columnID, rowID string) (interface{}, bool) {
	if d.rowLength() == 0 {
		return nil, false
	}
	columnIndex, ok := d.RowHeaderMap[columnID]
	if !ok {
		return nil, false
	}

	for _, row := range d.Rows {
		tempRow := row
		if len(tempRow.Grids) == 0 {
			continue
		}
		if rowID == string(tempRow.Grids[0]) {
			if len(tempRow.Grids) < columnIndex+1 {
				return nil, false
			}
			dataString := tempRow.Grids[columnIndex]
			return utils.UnmarshalParamValue(dataString), true
		}
	}

	return nil, false
}

func NewDataModel(srcDataModel *workspaceproto.DataModel, srcHeader []string, srcRows []*workspaceproto.Row) *dataModel {
	if srcDataModel == nil {
		return nil
	}
	dataModel := &dataModel{
		ID:   srcDataModel.Id,
		Name: srcDataModel.Name,
		Type: srcDataModel.Type,
	}

	dataModel.RowHeaderMap = convertHeader(srcHeader)

	rows := []*row{}
	for _, srcRow := range srcRows {
		tempRow := srcRow
		if tempRow != nil {
			rows = append(rows, convertRow(tempRow))
		}
	}

	dataModel.Rows = rows
	return dataModel
}

func convertHeader(srcHeader []string) map[string]int {
	res := make(map[string]int)
	if srcHeader == nil {
		return res
	}
	for i, grid := range srcHeader {
		res[grid] = i
	}
	return res
}

func convertRow(src *workspaceproto.Row) *row {
	if src == nil {
		return &row{}
	}
	var grids []string
	grids = append(grids, src.Grids...)
	return &row{Grids: grids}
}

// system
type dataList struct {
	entityModel    *dataModel
	workspaceModel *dataModel
	setModelChain  []setModel
}

func (c *dataList) pickInputValue(rowID, valueStr string) (interface{}, error) {
	// type workspace
	if strings.HasPrefix(valueStr, consts.WorkspaceTypeDataModelRefPrefix) {
		res, ok := c.workspaceModel.pickData(consts.WsDataModelValueHeader, strings.TrimPrefix(valueStr, consts.WorkspaceTypeDataModelRefPrefix))
		if !ok {
			return nil, apperrors.NewInternalError(fmt.Errorf("ref expression '%s' invalid, not found its value in DataModel 'workspace_data'", valueStr))
		}
		return res, nil
	}
	// type entity or set
	if strings.HasPrefix(valueStr, consts.DataModelRefPrefix) {
		refStage := strings.Split(strings.TrimPrefix(valueStr, consts.DataModelRefPrefix), ".")
		if len(refStage) < 1 {
			return nil, apperrors.NewInternalError(fmt.Errorf("ref expression '%s' incomplete", valueStr))
		}
		if len(refStage) > (len(c.setModelChain) + 1) {
			return nil, apperrors.NewInternalError(fmt.Errorf("ref expression '%s' does not match datamodel set(count %d)", valueStr, len(c.setModelChain)))
		}
		return pick(rowID, refStage, c.setModelChain, c.entityModel)
	}
	// others
	return valueStr, nil
}

func pick(rowID string, refStage []string, chain []setModel, dm *dataModel) (interface{}, error) {
	if len(refStage) == 0 || dm == nil {
		return nil, apperrors.NewInternalError(fmt.Errorf("ref expression is empty"))
	}

	curHeader := refStage[0]
	if len(refStage) == 1 {
		// Distinguish between set_id and entity
		if len(chain) > 0 {
			// set_id
			if curHeader == fmt.Sprintf("%s_id", chain[0].model.Name) {
				return rowID, nil
			}
			// entity_list
			if _, ok := chain[0].model.RowHeaderMap[curHeader]; ok {
				idList, ok := chain[0].sets[rowID]
				if !ok {
					return nil, apperrors.NewInternalError(fmt.Errorf("row id %s not found", rowID))
				}
				return idList, nil
			}
			return nil, apperrors.NewInternalError(fmt.Errorf("header %s not exist", curHeader))
		}

		// entity
		res, ok := dm.pickData(curHeader, rowID)
		if !ok {
			return nil, apperrors.NewInternalError(fmt.Errorf("can not find data with row %s and column %s ", rowID, curHeader))
		}
		return res, nil
	}

	if len(chain) < 1 {
		return nil, apperrors.NewInternalError(fmt.Errorf("none datamodel set to reference"))
	}

	idList, ok := chain[0].sets[rowID]
	if !ok {
		return nil, apperrors.NewInternalError(fmt.Errorf("row id '%s' not found", rowID))
	}
	if len(idList) < 1 {
		return nil, nil
	}
	res := make([]interface{}, len(idList))
	for i := range idList {
		var err error
		res[i], err = pick(idList[i], refStage[1:], chain[1:], dm)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
	}
	return res, nil

}
