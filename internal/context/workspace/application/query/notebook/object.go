package notebook

import "time"

// Notebook is DTO
type Notebook struct {
	Name        string
	WorkspaceID string
	Content     []byte
	Size        int64
	UpdateTime  time.Time
}
