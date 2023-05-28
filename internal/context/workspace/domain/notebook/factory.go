package notebook

import (
	"time"

	"github.com/Bio-OS/bioos/pkg/errors"
)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

type CreateParam struct {
	Name        string
	WorkspaceID string
	Content     []byte
}

func (f *Factory) New(param *CreateParam) (*Notebook, error) {
	if err := validateIPythonNotebook(param.Content); err != nil {
		return nil, errors.NewInvalidError("notebook", "content", err.Error())
	}
	return &Notebook{
		Name:       param.Name,
		Namespace:  param.WorkspaceID,
		Content:    param.Content,
		UpdateTime: time.Now(),
	}, nil
}
