package handlers

import (
	runcommand "github.com/Bio-OS/bioos/internal/context/submission/application/command/run"
	submissioncommand "github.com/Bio-OS/bioos/internal/context/submission/application/command/submission"
	runquery "github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	submissionquery "github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	"github.com/Bio-OS/bioos/pkg/utils"
)

func createSubmissionVoToDto(req CreateSubmissionRequest) *submissioncommand.CreateSubmissionCommand {
	return &submissioncommand.CreateSubmissionCommand{
		WorkspaceID:    req.WorkspaceID,
		Name:           req.Name,
		WorkflowID:     req.WorkflowID,
		Description:    req.Description,
		Type:           req.Type,
		Language:       req.Language,
		Entity:         commandEntityVoToDto(req.Entity),
		ExposedOptions: commandExposedOptionsVoToDto(req.ExposedOptions),
		InOutMaterial:  commandInOutMaterialVoToDto(req.InOutMaterial),
	}
}

func commandEntityVoToDto(entity *Entity) *submissioncommand.Entity {
	if entity == nil {
		return nil
	}
	return &submissioncommand.Entity{
		DataModelID:     entity.DataModelID,
		DataModelRowIDs: entity.DataModelRowIDs,
		InputsTemplate:  entity.InputsTemplate,
		OutputsTemplate: entity.OutputsTemplate,
	}
}

func commandExposedOptionsVoToDto(options ExposedOptions) submissioncommand.ExposedOptions {
	return submissioncommand.ExposedOptions{
		ReadFromCache: options.ReadFromCache,
	}
}

func commandInOutMaterialVoToDto(material *InOutMaterial) *submissioncommand.InOutMaterial {
	if material == nil {
		return nil
	}
	return &submissioncommand.InOutMaterial{
		InputsMaterial:  material.InputsMaterial,
		OutputsMaterial: material.OutputsMaterial,
	}
}

func cancelSubmissionVoToDto(req CancelSubmissionRequest) *submissioncommand.CancelSubmissionCommand {
	return &submissioncommand.CancelSubmissionCommand{
		WorkspaceID: req.WorkspaceID,
		ID:          req.ID,
	}
}

func deleteSubmissionVoToDto(req DeleteSubmissionRequest) *submissioncommand.DeleteSubmissionCommand {
	return &submissioncommand.DeleteSubmissionCommand{
		WorkspaceID: req.WorkspaceID,
		ID:          req.ID,
	}
}

func checkSubmissionVoToDto(req CheckSubmissionRequest) *submissionquery.CheckQuery {
	return &submissionquery.CheckQuery{
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
	}
}

func listSubmissionsVoToDto(req ListSubmissionsRequest) (*submissionquery.ListQuery, error) {
	pg := utils.NewPagination(req.Size, req.Page)
	if err := pg.SetOrderBy(req.OrderBy); err != nil {
		return nil, err
	}
	filter := &submissionquery.ListSubmissionsFilter{}
	if req.SearchWord != "" {
		filter.SearchWord = req.SearchWord
		filter.Exact = req.Exact
	}
	if req.WorkflowID != "" {
		filter.WorkflowID = req.WorkflowID
	}
	if len(req.IDs) != 0 {
		filter.IDs = req.IDs
	}
	if len(req.Status) != 0 {
		filter.Status = req.Status
	}
	if len(req.Language) != 0 {
		filter.Language = req.Language
	}
	return &submissionquery.ListQuery{
		WorkspaceID: req.WorkspaceID,
		Pg:          pg,
		Filter:      filter,
	}, nil
}

func submissionItemDtoToVo(item *submissionquery.SubmissionItem) SubmissionItem {
	return SubmissionItem{
		ID:              item.ID,
		Name:            item.Name,
		Description:     item.Description,
		Type:            item.Type,
		Language:        item.Language,
		Status:          item.Status,
		StartTime:       item.StartTime,
		FinishTime:      item.FinishTime,
		Duration:        item.Duration,
		WorkflowVersion: queryWorkflowVersionDtoToVo(item.WorkflowID, item.WorkflowVersionID),
		RunStatus:       submissionQueryStatusDtoToVo(item.RunStatus),
		Entity:          queryEntityDtoToVo(item.Entity),
		ExposedOptions:  queryExposedOptionsDtoToVo(item.ExposedOptions),
		InOutMaterial:   queryInOutMaterialDtoToVo(item.InOutMaterial),
	}
}

func queryWorkflowVersionDtoToVo(id, version string) WorkflowVersion {
	return WorkflowVersion{
		ID:        id,
		VersionID: version,
	}
}

func submissionQueryStatusDtoToVo(status submissionquery.Status) Status {
	return Status{
		Count:      status.Count,
		Pending:    status.Pending,
		Succeeded:  status.Succeeded,
		Failed:     status.Failed,
		Running:    status.Running,
		Cancelling: status.Cancelling,
		Cancelled:  status.Cancelled,
	}
}

func queryEntityDtoToVo(entity *submissionquery.Entity) *Entity {
	if entity == nil {
		return nil
	}
	return &Entity{
		DataModelID:     entity.DataModelID,
		DataModelRowIDs: entity.DataModelRowIDs,
		InputsTemplate:  entity.InputsTemplate,
		OutputsTemplate: entity.OutputsTemplate,
	}
}

func queryExposedOptionsDtoToVo(options submissionquery.ExposedOptions) ExposedOptions {
	return ExposedOptions{
		ReadFromCache: options.ReadFromCache,
	}
}

func queryInOutMaterialDtoToVo(material *submissionquery.InOutMaterial) *InOutMaterial {
	if material == nil {
		return nil
	}
	return &InOutMaterial{
		InputsMaterial:  material.InputsMaterial,
		OutputsMaterial: material.OutputsMaterial,
	}
}

func cancelRunVoToDto(req CancelRunRequest) *runcommand.CancelRunCommand {
	return &runcommand.CancelRunCommand{
		WorkspaceID:  req.WorkspaceID,
		SubmissionID: req.SubmissionID,
		ID:           req.ID,
	}
}

func listRunsVoToDto(req ListRunsRequest) (*runquery.ListRunsQuery, error) {
	pg := utils.NewPagination(req.Size, req.Page)
	if err := pg.SetOrderBy(req.OrderBy); err != nil {
		return nil, err
	}
	filter := &runquery.ListRunsFilter{}
	if req.SearchWord != "" {
		filter.SearchWord = req.SearchWord
	}
	if len(req.IDs) != 0 {
		filter.IDs = req.IDs
	}
	if len(req.Status) != 0 {
		filter.Status = req.Status
	}
	return &runquery.ListRunsQuery{
		WorkspaceID:  req.WorkspaceID,
		SubmissionID: req.SubmissionID,
		Pg:           pg,
		Filter:       filter,
	}, nil
}

func runItemDtoToVo(item *runquery.RunItem) RunItem {
	return RunItem{
		ID:          item.ID,
		Name:        item.Name,
		Status:      item.Status,
		StartTime:   item.StartTime,
		FinishTime:  item.FinishTime,
		Duration:    item.Duration,
		EngineRunID: item.EngineRunID,
		Inputs:      item.Inputs,
		Outputs:     item.Outputs,
		TaskStatus:  runQueryStatusDtoToVo(item.TaskStatus),
		Log:         item.Log,
		Message:     item.Message,
	}
}

func runQueryStatusDtoToVo(status runquery.Status) Status {
	return Status{
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

func listTasksVoToDto(req ListTasksRequest) (*runquery.ListTasksQuery, error) {
	pg := utils.NewPagination(req.Size, req.Page)
	if err := pg.SetOrderBy(req.OrderBy); err != nil {
		return nil, err
	}
	return &runquery.ListTasksQuery{
		WorkspaceID:  req.WorkspaceID,
		SubmissionID: req.SubmissionID,
		RunID:        req.RunID,
		Pg:           pg,
	}, nil
}

func taskItemDtoToVo(item *runquery.TaskItem) TaskItem {
	return TaskItem{
		Name:       item.Name,
		RunID:      item.RunID,
		Status:     item.Status,
		StartTime:  item.StartTime,
		FinishTime: item.FinishTime,
		Duration:   item.Duration,
		Stdout:     item.Stdout,
		Stderr:     item.Stderr,
	}
}

func runStatusDtoToSubmissionRunStatus(runStatus *runquery.Status) submissionquery.Status {
	return submissionquery.Status{
		Count:      runStatus.Count,
		Pending:    runStatus.Pending,
		Succeeded:  runStatus.Succeeded,
		Failed:     runStatus.Failed,
		Running:    runStatus.Running,
		Cancelling: runStatus.Cancelling,
		Cancelled:  runStatus.Cancelled,
	}
}
