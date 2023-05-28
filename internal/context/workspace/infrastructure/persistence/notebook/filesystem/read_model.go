package filesystem

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/notebook"
	domain "github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/notebook"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type readModel struct {
	basedir string
}

func NewReadModel(basedir string) (query.ReadModel, error) {
	if err := utils.ValidateFSDirectory(basedir); err != nil {
		return nil, err
	}
	return &readModel{basedir}, nil
}

func (r *readModel) ListByWorkspace(ctx context.Context, workspaceID string) ([]*query.Notebook, error) {
	dirFS := os.DirFS(r.basedir)
	list := []*notebookPO{}
	if err := fs.WalkDir(dirFS, workspaceID, func(root string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return fmt.Errorf("foreach file fail: %w", err)
		}
		if !d.IsDir() && path.Ext(root) == notebook.NotebookFileExt && path.Base(path.Dir(root)) == workspaceID && !strings.HasPrefix(path.Base(root), ".") {
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("get file '%s' info fail: %w", d.Name(), err)
			}
			po := notebookPO{
				path:       r.getRelativePath(root),
				size:       info.Size(),
				updateTime: info.ModTime(),
			}
			list = append(list, &po)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("walk dir fail: %w", err)
	}
	res := make([]*query.Notebook, len(list))
	for i := range list {
		res[i] = list[i].toDTO()
	}
	return res, nil
}

func (r *readModel) Get(ctx context.Context, workspaceID, name string) (*query.Notebook, error) {
	po := &notebookPO{path: domain.Path(workspaceID, name)}
	fname := r.filename(po)
	var err error
	if po.content, err = os.ReadFile(fname); err != nil {
		if os.IsNotExist(err) {
			applog.Errorf("notebook file '%s' no exist", fname)
			return nil, nil
		}
		return nil, fmt.Errorf("read file '%s' fail: %w", name, err)
	}
	po.size = int64(len(po.content))
	// TODO get modify time
	return po.toDTO(), nil
}

func (r *readModel) filename(po *notebookPO) string {
	return filename(r.basedir, po)
}

func (r *readModel) getRelativePath(p string) string {
	return strings.TrimSuffix(strings.TrimPrefix(p, r.basedir), notebook.NotebookFileExt)
}
