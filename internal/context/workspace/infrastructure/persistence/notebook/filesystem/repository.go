package filesystem

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type repository struct {
	basedir string
}

func NewRepository(basedir string) (notebook.Repository, error) {
	if err := utils.ValidateFSDirectory(basedir); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err = os.MkdirAll(basedir, 0750); err != nil {
			return nil, fmt.Errorf("create dir '%s' fail: %w", basedir, err)
		}
	}
	return &repository{basedir}, nil
}

func (r *repository) Save(_ context.Context, nb *notebook.Notebook) error {
	po := newPO(nb)
	name := r.filename(po)
	dir := path.Dir(name)
	if err := utils.ValidateFSDirectory(dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err = os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("create dir '%s' fail: %w", dir, err)
		}
	}
	if err := os.WriteFile(name, po.content, 0660); err != nil {
		return fmt.Errorf("write file '%s' fail: %w", name, err)
	}
	return nil
}

func (r *repository) Get(_ context.Context, path string) (*notebook.Notebook, error) {
	po := &notebookPO{path: path}
	name := r.filename(po)
	var err error
	if po.content, err = os.ReadFile(name); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read file '%s' fail: %w", name, err)
	}
	return po.toDO(), nil
}

func (r *repository) Delete(_ context.Context, nb *notebook.Notebook) error {
	po := newPO(nb)
	name := r.filename(po)
	if err := os.Remove(name); err != nil {
		return fmt.Errorf("remove file '%s' fail: %w", name, err)
	}
	return nil
}

func (r *repository) filename(po *notebookPO) string {
	return filename(r.basedir, po)
}
