package domain

import (
	"time"

	"github.com/Bio-OS/bioos/pkg/notebook"
)

const (
	ServerStatusTerminated  = "Terminated"
	ServerStatusTerminating = "Terminating"
	ServerStatusPending     = "Pending"
	ServerStatusRunning     = "Running"
	ServerStatusUnknown     = "Unknown"
)

const (
	VolumeTypeNFS = "NFS"
	VolumeTypeTOS = "TOS" // not support
	VolumeTypeS3  = "S3"  // not support
)

type Settings struct {
	DockerImage  string
	ResourceSize notebook.ResourceSize
	NodeSelector map[string]string // no need to persistence
}

type Status struct {
	Status    string
	AccessURL string
}

// Volume describe these storage needs except single user HOME persistence (runtime specified):
//  1. workspace related storage
//  2. ipynb file mount path
type Volume struct {
	Name              string
	Type              string
	Source            string
	MountRelativePath string
}

type NotebookServer struct {
	ID          string
	WorkspaceID string
	Settings    Settings
	Status      Status
	Volumes     []Volume
	CreateTime  time.Time
	UpdateTime  time.Time
}
