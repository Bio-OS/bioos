package filesystem

import (
	"path"

	"github.com/Bio-OS/bioos/pkg/notebook"
)

func filename(basedir string, po *notebookPO) string {
	return path.Join(basedir, po.path) + notebook.NotebookFileExt
}
