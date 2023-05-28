package query

import (
	"time"

	"github.com/Bio-OS/bioos/pkg/notebook"
)

type NotebookServer struct {
	NotebookSettings
	Status    string
	AccessURL string
}

type NotebookSettings struct {
	ID           string
	WorkspaceID  string
	Image        string
	ResourceSize notebook.ResourceSize
	CreateTime   time.Time
	UpdateTime   time.Time
}
