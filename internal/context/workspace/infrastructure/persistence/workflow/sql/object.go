package sql

import (
	"encoding/json"
	"time"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	domain "github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

type workflow struct {
	ID            string `gorm:"primaryKey"`
	Name          string `gorm:"type:varchar(200);uniqueIndex:idx_name_ws"`
	WorkspaceID   string `gorm:"type:varchar(32);uniqueIndex:idx_name_ws"`
	Description   string
	LatestVersion string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func workflowDOToPO(wf *domain.Workflow) *workflow {
	if wf == nil {
		return nil
	}
	return &workflow{
		ID:            wf.ID,
		Name:          wf.Name,
		WorkspaceID:   wf.WorkspaceID,
		Description:   wf.Description,
		LatestVersion: wf.LatestVersion,
		CreatedAt:     wf.CreatedAt,
		UpdatedAt:     wf.UpdatedAt,
	}
}

func workflowVersionDOToPO(workflowID string, version *domain.WorkflowVersion) (*workflowVersion, error) {
	if version == nil {
		return nil, nil
	}
	metadata, err := json.Marshal(version.Metadata)
	if err != nil {
		return nil, err
	}
	inputs, err := json.Marshal(version.Inputs)
	if err != nil {
		return nil, err
	}
	outputs, err := json.Marshal(version.Outputs)
	if err != nil {
		return nil, err
	}

	return &workflowVersion{
		ID:               version.ID,
		WorkflowID:       workflowID,
		Status:           version.Status,
		Message:          version.Message,
		Language:         version.Language,
		LanguageVersion:  version.LanguageVersion,
		MainWorkflowPath: version.MainWorkflowPath,
		Inputs:           string(inputs),
		Outputs:          string(outputs),
		Graph:            version.Graph,
		Metadata:         string(metadata),
		Source:           version.Source,
		CreatedAt:        version.CreatedAt,
		UpdatedAt:        version.UpdatedAt,
	}, nil
}

func (w *workflow) toDO() *domain.Workflow {
	return &domain.Workflow{
		ID:            w.ID,
		Name:          w.Name,
		WorkspaceID:   w.WorkspaceID,
		Description:   w.Description,
		LatestVersion: w.LatestVersion,
		Versions:      map[string]*domain.WorkflowVersion{},
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}
}

func (w *workflow) toDTO() *query.Workflow {
	return &query.Workflow{
		ID:          w.ID,
		Name:        w.Name,
		WorkspaceID: w.WorkspaceID,
		Description: w.Description,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}

func (w *workflow) TableName() string {
	return "workflow"
}

type workflowVersion struct {
	ID               string `gorm:"primaryKey"`
	WorkflowID       string
	Status           string `gorm:"type:varchar(100)"`
	Message          string
	Language         string `gorm:"type:varchar(20)"`
	LanguageVersion  string `gorm:"type:varchar(20)"`
	MainWorkflowPath string `gorm:"type:varchar(255)"`
	Inputs           string
	Outputs          string
	Graph            string
	Metadata         string
	Source           string `gorm:"type:varchar(20)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (v *workflowVersion) TableName() string {
	return "workflow_version"
}
func (v *workflowVersion) toDO() (*domain.WorkflowVersion, error) {
	metadata := make(map[string]string)
	if len(v.Metadata) > 0 {
		if err := json.Unmarshal([]byte(v.Metadata), &metadata); err != nil {
			return nil, err
		}
	}
	inputs, err := workflowParamPOToDO(v.Inputs)
	if err != nil {
		return nil, err
	}
	outputs, err := workflowParamPOToDO(v.Outputs)
	if err != nil {
		return nil, err
	}

	return &domain.WorkflowVersion{
		ID:               v.ID,
		Status:           v.Status,
		Message:          v.Message,
		Language:         v.Language,
		LanguageVersion:  v.LanguageVersion,
		MainWorkflowPath: v.MainWorkflowPath,
		Inputs:           inputs,
		Outputs:          outputs,
		Graph:            v.Graph,
		Source:           v.Source,
		Files:            map[string]*domain.WorkflowFile{},
		Metadata:         metadata,
		CreatedAt:        v.CreatedAt,
		UpdatedAt:        v.UpdatedAt,
	}, nil
}
func (v *workflowVersion) toDTO() (*query.WorkflowVersion, error) {
	metadata := make(map[string]string)
	if len(v.Metadata) > 0 {
		if err := json.Unmarshal([]byte(v.Metadata), &metadata); err != nil {
			return nil, err
		}
	}

	inputs, err := workflowParamPOToDTO(v.Inputs)
	if err != nil {
		applog.Errorw("fail to unmarshal inputs", "inputs", v.Inputs)
		return nil, err
	}
	outputs, err := workflowParamPOToDTO(v.Outputs)
	if err != nil {
		applog.Errorw("fail to unmarshal outputs", "outputs", v.Outputs)
		return nil, err
	}

	return &query.WorkflowVersion{
		ID:               v.ID,
		Status:           v.Status,
		Message:          v.Message,
		Language:         v.Language,
		LanguageVersion:  v.LanguageVersion,
		MainWorkflowPath: v.MainWorkflowPath,
		Inputs:           inputs,
		Outputs:          outputs,
		Graph:            v.Graph,
		Source:           v.Source,
		Metadata:         metadata,
		CreatedAt:        v.CreatedAt,
		UpdatedAt:        v.UpdatedAt,
	}, nil
}

type workflowFile struct {
	ID                string `gorm:"primaryKey"`
	WorkflowVersionID string
	Path              string
	Content           string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (f *workflowFile) TableName() string {
	return "workflow_file"
}

func (f *workflowFile) toDTO() *query.WorkflowFile {
	return &query.WorkflowFile{
		ID:                f.ID,
		WorkflowVersionID: f.WorkflowVersionID,
		Path:              f.Path,
		Content:           f.Content,
		CreatedAt:         f.CreatedAt,
		UpdatedAt:         f.UpdatedAt,
	}
}

func (f *workflowFile) toDO() *domain.WorkflowFile {
	return &domain.WorkflowFile{
		ID:        f.ID,
		Path:      f.Path,
		Content:   f.Content,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func workflowFileDOToPO(versionID string, file *domain.WorkflowFile) *workflowFile {
	if file == nil {
		return nil
	}
	return &workflowFile{
		ID:                file.ID,
		WorkflowVersionID: versionID,
		Path:              file.Path,
		Content:           file.Content,
		CreatedAt:         file.CreatedAt,
		UpdatedAt:         file.UpdatedAt,
	}
}

func workflowParamPOToDTO(paramStr string) ([]query.WorkflowParam, error) {
	var params []query.WorkflowParam
	if len(paramStr) > 0 {
		var paramDOs []domain.WorkflowParam
		if err := json.Unmarshal([]byte(paramStr), &paramDOs); err != nil {
			return nil, err
		}
		for _, paramDO := range paramDOs {
			params = append(params, query.WorkflowParam{
				Name:     paramDO.Name,
				Type:     paramDO.Type,
				Optional: paramDO.Optional,
				Default:  paramDO.Default,
			})
		}
	}
	return params, nil
}

func workflowParamPOToDO(paramStr string) ([]domain.WorkflowParam, error) {
	var params []domain.WorkflowParam
	if len(paramStr) > 0 {
		if err := json.Unmarshal([]byte(paramStr), &params); err != nil {
			return nil, err
		}
	}
	return params, nil
}
