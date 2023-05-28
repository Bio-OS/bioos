package sql

import (
	"time"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/query"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/pkg/notebook"
)

type notebookServer struct {
	ID           string `gorm:"primaryKey"`
	WorkspaceID  string
	DockerImage  string
	ResourceSize resourceSize `gorm:"serializer:json"`
	Volumes      []volume     `gorm:"serializer:json"`
	CreateTime   time.Time
	UpdateTime   time.Time
}

type gpuInfo struct {
	Model  string
	Card   float64
	Memory int64 // byte
}

type resourceSize struct {
	CPU    float64  `json:"cpu"`
	Memory int64    `json:"memory"`
	Disk   int64    `json:"disk"`
	GPU    *gpuInfo `json:"gpu,omitempty"`
}

type volume struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	Source            string `json:"source"`
	MountRelativePath string `json:"mountPath"`
}

func newNotebookServer(do *domain.NotebookServer) *notebookServer {
	var gpu *gpuInfo
	if do.Settings.ResourceSize.GPU != nil {
		gpu = &gpuInfo{
			Model:  do.Settings.ResourceSize.GPU.Model,
			Card:   do.Settings.ResourceSize.GPU.Card,
			Memory: do.Settings.ResourceSize.GPU.Memory,
		}
	}
	volumes := make([]volume, len(do.Volumes))
	for i := range do.Volumes {
		volumes[i].Name = do.Volumes[i].Name
		volumes[i].Type = do.Volumes[i].Type
		volumes[i].Source = do.Volumes[i].Source
		volumes[i].MountRelativePath = do.Volumes[i].MountRelativePath
	}
	return &notebookServer{
		ID:          do.ID,
		WorkspaceID: do.WorkspaceID,
		DockerImage: do.Settings.DockerImage,
		ResourceSize: resourceSize{
			CPU:    do.Settings.ResourceSize.CPU,
			Memory: do.Settings.ResourceSize.Memory,
			Disk:   do.Settings.ResourceSize.Disk,
			GPU:    gpu,
		},
		Volumes:    volumes,
		CreateTime: do.CreateTime,
		UpdateTime: do.UpdateTime,
	}
}

func (s *notebookServer) TableName() string {
	return "notebookserver"
}

func (s *notebookServer) toDO() *domain.NotebookServer {
	var gpu *notebook.GPU
	if s.ResourceSize.GPU != nil {
		gpu = &notebook.GPU{
			Model:  s.ResourceSize.GPU.Model,
			Card:   s.ResourceSize.GPU.Card,
			Memory: s.ResourceSize.GPU.Memory,
		}
	}
	volumes := make([]domain.Volume, len(s.Volumes))
	for i := range s.Volumes {
		volumes[i].Name = s.Volumes[i].Name
		volumes[i].Type = s.Volumes[i].Type
		volumes[i].Source = s.Volumes[i].Source
		volumes[i].MountRelativePath = s.Volumes[i].MountRelativePath
	}
	return &domain.NotebookServer{
		ID:          s.ID,
		WorkspaceID: s.WorkspaceID,
		Settings: domain.Settings{
			DockerImage: s.DockerImage,
			ResourceSize: notebook.ResourceSize{
				CPU:    s.ResourceSize.CPU,
				Memory: s.ResourceSize.Memory,
				Disk:   s.ResourceSize.Disk,
				GPU:    gpu,
			},
		},
		Volumes:    volumes,
		CreateTime: s.CreateTime,
		UpdateTime: s.UpdateTime,
	}
}

func (s *notebookServer) toDTO() *query.NotebookSettings {
	var gpu *notebook.GPU
	if s.ResourceSize.GPU != nil {
		gpu = &notebook.GPU{
			Model:  s.ResourceSize.GPU.Model,
			Card:   s.ResourceSize.GPU.Card,
			Memory: s.ResourceSize.GPU.Memory,
		}
	}
	return &query.NotebookSettings{
		ID:          s.ID,
		WorkspaceID: s.WorkspaceID,
		Image:       s.DockerImage,
		ResourceSize: notebook.ResourceSize{
			CPU:    s.ResourceSize.CPU,
			Memory: s.ResourceSize.Memory,
			Disk:   s.ResourceSize.Disk,
			GPU:    gpu,
		},
		CreateTime: s.CreateTime,
		UpdateTime: s.UpdateTime,
	}
}
