package datamodel

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/consts"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ImportDataModelsHandler struct {
	repo     Repository
	eventbus eventbus.EventBus
	factory  *Factory
}

func NewImportDataModelsHandler(repo Repository, bus eventbus.EventBus, factory *Factory) *ImportDataModelsHandler {
	return &ImportDataModelsHandler{
		repo:     repo,
		eventbus: bus,
		factory:  factory,
	}
}

func (h *ImportDataModelsHandler) Handle(ctx context.Context, event *ImportDataModelsEvent) error {
	dataModelMap := make(map[string]struct{}, 0)
	for _, dataModel := range event.Schemas {
		dataModelPath := strings.ReplaceAll(dataModel.Path, " ", "")
		if _, ok := dataModelMap[dataModel.Name]; ok {
			return fmt.Errorf("data model name[%s] is not unique ", dataModel.Name)
		}
		dataModelMap[dataModel.Name] = struct{}{}
		if !validator.ValidateDataModelNameInString(dataModel.Name) {
			return fmt.Errorf("data model name[%s] not passed the validation ", dataModel.Name)
		}
		_, dataModelFileName := filepath.Split(dataModelPath)
		if dataModelFileName != fmt.Sprintf("%s.csv", dataModel.Name) {
			return fmt.Errorf("data model name[%s] not equal with data model csv file[%s] ", dataModel.Name, dataModelPath)
		}
		headers, rows, err := utils.ReadDataModelFromCSV(path.Join(event.ImportFileBaseDir, dataModelPath))
		if err != nil {
			return fmt.Errorf("read csv file failed: %w", err)
		}
		newDataModel := h.factory.New(&CreateParam{
			WorkspaceID: event.WorkspaceID,
			Name:        dataModel.Name,
			Type:        consts.DataModelTypeEntity,
			Headers:     headers,
			Rows:        rows,
		})
		err = h.repo.Save(ctx, newDataModel)
		if err != nil {
			return err
		}
	}

	// clean files
	for _, dataModel := range event.Schemas {
		err := os.Remove(path.Join(event.ImportFileBaseDir, dataModel.Path))
		if err != nil {
			//remove file error should not lead to import fail
			applog.Errorf("remove file failed: %w", err)
		}
	}
	return nil
}
