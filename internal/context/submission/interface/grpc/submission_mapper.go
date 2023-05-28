package grpc

import (
	command "github.com/Bio-OS/bioos/internal/context/submission/application/command/submission"
	runquery "github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	query "github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	pb "github.com/Bio-OS/bioos/internal/context/submission/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/utils"
)

func submissionsItemDTOToVO(item *query.SubmissionItem) *pb.SubmissionItem {
	ret := &pb.SubmissionItem{
		Id:              item.ID,
		Name:            item.Name,
		Type:            item.Type,
		Status:          item.Status,
		StartTime:       item.StartTime,
		Duration:        item.Duration,
		WorkflowVersion: queryWorkflowVersionDTOToVO(item.WorkflowID, item.WorkflowVersionID),
		RunStatus:       submissionQueryStatusDTOToVO(item.RunStatus),
		Entity:          queryEntityDTOToVO(item.Entity),
		ExposedOptions:  queryExposedOptionsDTOToVO(item.ExposedOptions),
		InOutMaterial:   queryInOutMaterialDTOToVO(item.InOutMaterial),
	}

	if item.Description != nil {
		ret.Description = *item.Description
	}
	if item.FinishTime != nil {
		ret.FinishTime = *item.FinishTime
	}

	return ret
}

func queryWorkflowVersionDTOToVO(id, version string) *pb.WorkflowVersionInfo {
	return &pb.WorkflowVersionInfo{
		Id:        id,
		VersionID: version,
	}
}

func submissionQueryStatusDTOToVO(status query.Status) *pb.Status {
	return &pb.Status{
		Count:      status.Count,
		Pending:    status.Pending,
		Succeeded:  status.Succeeded,
		Failed:     status.Failed,
		Running:    status.Running,
		Cancelling: status.Cancelling,
		Cancelled:  status.Cancelled,
	}
}

func queryEntityDTOToVO(entity *query.Entity) *pb.Entity {
	if entity == nil {
		return nil
	}
	return &pb.Entity{
		DataModelID:     entity.DataModelID,
		DataModelRowIDs: entity.DataModelRowIDs,
		InputsTemplate:  entity.InputsTemplate,
		OutputsTemplate: entity.OutputsTemplate,
	}
}

func queryExposedOptionsDTOToVO(options query.ExposedOptions) *pb.ExposedOptions {
	return &pb.ExposedOptions{
		ReadFromCache: options.ReadFromCache,
	}
}

func queryInOutMaterialDTOToVO(material *query.InOutMaterial) *pb.InOutMaterial {
	if material == nil {
		return nil
	}
	return &pb.InOutMaterial{
		InputsMaterial:  material.InputsMaterial,
		OutputsMaterial: material.OutputsMaterial,
	}
}

func createSubmissionVOToDTO(req *pb.CreateSubmissionRequest) *command.CreateSubmissionCommand {
	return &command.CreateSubmissionCommand{
		WorkspaceID:    req.WorkspaceID,
		Name:           req.Name,
		WorkflowID:     req.WorkflowID,
		Description:    &req.Description,
		Type:           req.Type,
		Entity:         commandEntityVoToDto(req.Entity),
		ExposedOptions: commandExposedOptionsVoToDto(req.ExposedOptions),
		InOutMaterial:  commandInOutMaterialVoToDto(req.InOutMaterial),
	}
}

func commandEntityVoToDto(entity *pb.Entity) *command.Entity {
	if entity == nil {
		return nil
	}
	return &command.Entity{
		DataModelID:     entity.DataModelID,
		DataModelRowIDs: entity.DataModelRowIDs,
		InputsTemplate:  entity.InputsTemplate,
		OutputsTemplate: entity.OutputsTemplate,
	}
}

func commandExposedOptionsVoToDto(options *pb.ExposedOptions) command.ExposedOptions {
	if options == nil {
		return command.ExposedOptions{}
	}
	return command.ExposedOptions{
		ReadFromCache: options.ReadFromCache,
	}
}

func commandInOutMaterialVoToDto(material *pb.InOutMaterial) *command.InOutMaterial {
	if material == nil {
		return nil
	}
	return &command.InOutMaterial{
		InputsMaterial:  material.InputsMaterial,
		OutputsMaterial: material.OutputsMaterial,
	}
}

func runItemDTOToVO(item *runquery.RunItem) *pb.RunItem {
	ret := &pb.RunItem{
		Id:          item.ID,
		Name:        item.Name,
		Status:      item.Status,
		StartTime:   item.StartTime,
		Duration:    item.Duration,
		EngineRunID: item.EngineRunID,
		Inputs:      item.Inputs,
		Outputs:     item.Outputs,
		TaskStatus:  runQueryStatusDtoToVo(item.TaskStatus),
	}
	if item.FinishTime != nil {
		ret.FinishTime = *item.FinishTime
	}
	if item.Log != nil {
		ret.Log = *item.Log
	}
	if item.Message != nil {
		ret.Message = *item.Message
	}

	return ret
}

func runQueryStatusDtoToVo(status runquery.Status) *pb.Status {
	return &pb.Status{
		Count:        status.Count,
		Succeeded:    status.Succeeded,
		Failed:       status.Failed,
		Running:      status.Running,
		Cancelling:   status.Cancelling,
		Cancelled:    status.Cancelled,
		Queued:       status.Queued,
		Initializing: status.Initializing,
	}
}

func taskItemDTOToVO(item *runquery.TaskItem) *pb.TaskItem {
	ret := &pb.TaskItem{
		Name:      item.Name,
		RunID:     item.RunID,
		Status:    item.Stdout,
		StartTime: item.StartTime,
		Duration:  item.Duration,
		Stdout:    item.Stdout,
		Stderr:    item.Stderr,
	}
	if item.FinishTime != nil {
		ret.FinishTime = *item.FinishTime
	}
	return ret
}

func listSubmissionsVOToDTO(req *pb.ListSubmissionsRequest) (*query.ListQuery, error) {
	pg := utils.NewPagination(int(req.GetSize()), int(req.GetPage()))
	if err := pg.SetOrderBy(req.GetOrderBy()); err != nil {
		return nil, err
	}
	filter := &query.ListSubmissionsFilter{}
	if req.GetSearchWord() != "" {
		filter.SearchWord = req.GetSearchWord()
		filter.Exact = req.GetExact()
	}
	if req.GetWorkflowID() != "" {
		filter.WorkflowID = req.GetWorkflowID()
	}
	if len(req.GetIds()) != 0 {
		filter.IDs = req.GetIds()
	}
	if len(req.GetStatus()) != 0 {
		filter.Status = req.GetStatus()
	}
	return &query.ListQuery{
		WorkspaceID: req.GetWorkspaceID(),
		Pg:          pg,
		Filter:      filter,
	}, nil
}

func listRunsVOToDTO(req *pb.ListRunsRequest) (*runquery.ListRunsQuery, error) {
	pg := utils.NewPagination(int(req.GetSize()), int(req.GetPage()))
	if err := pg.SetOrderBy(req.GetOrderBy()); err != nil {
		return nil, err
	}
	filter := &runquery.ListRunsFilter{}
	if req.GetSearchWord() != "" {
		filter.SearchWord = req.GetSearchWord()
	}
	if len(req.GetIds()) != 0 {
		filter.IDs = req.GetIds()
	}
	if len(req.GetStatus()) != 0 {
		filter.Status = req.GetStatus()
	}
	return &runquery.ListRunsQuery{
		WorkspaceID:  req.GetWorkspaceID(),
		SubmissionID: req.GetSubmissionID(),
		Pg:           pg,
		Filter:       filter,
	}, nil
}

func listTasksVOToDTO(req *pb.ListTasksRequest) (*runquery.ListTasksQuery, error) {
	pg := utils.NewPagination(int(req.GetSize()), int(req.GetPage()))
	if err := pg.SetOrderBy(req.GetOrderBy()); err != nil {
		return nil, err
	}
	return &runquery.ListTasksQuery{
		WorkspaceID:  req.GetWorkspaceID(),
		SubmissionID: req.GetSubmissionID(),
		RunID:        req.GetRunID(),
		Pg:           pg,
	}, nil
}
