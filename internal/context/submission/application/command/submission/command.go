package submission

import (
	"github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	submissionquery "github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
)

type CreateSubmissionCommand struct {
	WorkspaceID    string  `validate:"required"`
	Name           string  `validate:"required,submissionName"`
	WorkflowID     string  `validate:"required"`
	Description    *string `validate:"omitempty,submissionDesc"`
	Type           string  `validate:"required,oneof=dataModel filePath"`
	Language       string  `validate:"required,oneof=WDL Nextflow"`
	Entity         *Entity
	ExposedOptions ExposedOptions
	InOutMaterial  *InOutMaterial
}

type Entity struct {
	DataModelID     string
	DataModelRowIDs []string `validate:"unique"`
	InputsTemplate  string
	OutputsTemplate string
}

type InOutMaterial struct {
	InputsMaterial  string
	OutputsMaterial string
}

type ExposedOptions struct {
	ReadFromCache bool
}

type DeleteSubmissionCommand struct {
	WorkspaceID string `validate:"required"`
	ID          string `validate:"required"`
}

type CancelSubmissionCommand struct {
	WorkspaceID string `validate:"required"`
	ID          string `validate:"required"`
}

type Commands struct {
	CreateSubmission CreateSubmissionHandler
	DeleteSubmission DeleteSubmissionHandler
	CancelSubmission CancelSubmissionHandler
}

func NewCommands(grpcFactory grpc.Factory, submissionRepo submission.Repository, submissionFactory *submission.Factory, eventBus eventbus.EventBus, submissionReadModel submissionquery.ReadModel, runReadModel run.ReadModel) *Commands {
	service := submission.NewService(grpcFactory, submissionRepo, eventBus, submissionReadModel, runReadModel)
	return &Commands{
		CreateSubmission: NewCreateSubmissionHandler(service, submissionFactory, eventBus),
		DeleteSubmission: NewDeleteSubmissionHandler(service, eventBus),
		CancelSubmission: NewCancelSubmissionHandler(service, eventBus),
	}
}
