package filesystem

import (
	"time"

	dto "github.com/Bio-OS/bioos/internal/context/workspace/application/query/notebook"
	do "github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
)

type notebookPO struct {
	path       string
	content    []byte
	size       int64
	updateTime time.Time
}

func newPO(do *do.Notebook) *notebookPO {
	return &notebookPO{
		path:       do.Path(),
		content:    do.Content,
		updateTime: do.UpdateTime,
	}
}

func (po *notebookPO) toDO() *do.Notebook {
	return &do.Notebook{
		Name:       do.NameOfPath(po.path),
		Namespace:  do.NamespaceOfPath(po.path),
		Content:    po.content,
		UpdateTime: po.updateTime,
	}
}

func (po *notebookPO) toDTO() *dto.Notebook {
	return &dto.Notebook{
		Name:        do.NameOfPath(po.path),
		WorkspaceID: do.NamespaceOfPath(po.path),
		Content:     po.content,
		Size:        po.size,
		UpdateTime:  po.updateTime,
	}
}
