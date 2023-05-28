package domain

import (
	"reflect"
	"time"

	"github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/notebook"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type Factory struct {
	officialImages map[string]struct{}
	sizeOptions    []notebook.ResourceOption
}

func NewFactory(officialImages []string, sizeOptions []notebook.ResourceOption) *Factory {
	mapImage := make(map[string]struct{})
	for _, s := range officialImages {
		mapImage[s] = struct{}{}
	}
	return &Factory{
		officialImages: mapImage,
		sizeOptions:    sizeOptions,
	}
}

type CreateParam struct {
	WorkspaceID  string
	Image        string
	ResourceSize notebook.ResourceSize
	Volumes      []Volume
}

func (f *Factory) New(param *CreateParam) (*NotebookServer, error) {
	var nodeSelector map[string]string
	foundSize := false
	for i := range f.sizeOptions {
		if reflect.DeepEqual(param.ResourceSize, f.sizeOptions[i].ResourceSize) {
			nodeSelector = f.sizeOptions[i].NodeSelector
			foundSize = true
		}
	}
	if !foundSize {
		return nil, errors.NewInvalidError("notebookserver", "resource size", param.ResourceSize.String())
	}
	// TODO validate image exist and available
	// if _, ok := f.officialImages[image]; !ok {
	now := time.Now()
	return &NotebookServer{
		ID:          utils.GenNotebookServerID(),
		WorkspaceID: param.WorkspaceID,
		Settings: Settings{
			DockerImage:  param.Image,
			ResourceSize: param.ResourceSize,
			NodeSelector: nodeSelector,
		},
		Volumes:    param.Volumes,
		CreateTime: now,
		UpdateTime: now,
	}, nil
}

func (f *Factory) GetResourceSizes() []notebook.ResourceSize {
	sizes := make([]notebook.ResourceSize, len(f.sizeOptions))
	for i, opt := range f.sizeOptions {
		sizes[i] = opt.ResourceSize
	}
	return sizes
}
