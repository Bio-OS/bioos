package mongo

import (
	"time"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/query"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/pkg/notebook"
)

type notebookServer struct {
	ID          string    `bson:"id"`
	WorkspaceID string    `bson:"workspaceID"`
	Settings    settings  `bson:"settings"`
	Volumes     []volume  `bson:"volumes"`
	CreateTime  time.Time `bson:"createTime"`
	UpdateTime  time.Time `bson:"updateTime"`
}

type gpuInfo struct {
	Model  string  `bson:"model"`
	Card   float64 `bson:"card"`
	Memory int64   `bson:"memory"`
}

type resourceSize struct {
	CPU    float64  `bson:"cpu"`
	Memory int64    `bson:"memory"`
	Disk   int64    `bson:"disk"`
	GPU    *gpuInfo `bson:"gpu,omitempty"`
}

type settings struct {
	DockerImage  string       `bson:"dockerImage"`
	ResourceSize resourceSize `bson:"resourceSize"`
}

type volume struct {
	Name              string `bson:"name"`
	Type              string `bson:"type"`
	Source            string `bson:"source"`
	MountRelativePath string `bson:"mountPath"`
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
		Settings: settings{
			DockerImage: do.Settings.DockerImage,
			ResourceSize: resourceSize{
				CPU:    do.Settings.ResourceSize.CPU,
				Memory: do.Settings.ResourceSize.Memory,
				Disk:   do.Settings.ResourceSize.Disk,
				GPU:    gpu,
			},
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
	if s.Settings.ResourceSize.GPU != nil {
		gpu = &notebook.GPU{
			Model:  s.Settings.ResourceSize.GPU.Model,
			Card:   s.Settings.ResourceSize.GPU.Card,
			Memory: s.Settings.ResourceSize.GPU.Memory,
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
			DockerImage: s.Settings.DockerImage,
			ResourceSize: notebook.ResourceSize{
				CPU:    s.Settings.ResourceSize.CPU,
				Memory: s.Settings.ResourceSize.Memory,
				Disk:   s.Settings.ResourceSize.Disk,
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
	if s.Settings.ResourceSize.GPU != nil {
		gpu = &notebook.GPU{
			Model:  s.Settings.ResourceSize.GPU.Model,
			Card:   s.Settings.ResourceSize.GPU.Card,
			Memory: s.Settings.ResourceSize.GPU.Memory,
		}
	}
	return &query.NotebookSettings{
		ID:          s.ID,
		WorkspaceID: s.WorkspaceID,
		Image:       s.Settings.DockerImage,
		ResourceSize: notebook.ResourceSize{
			CPU:    s.Settings.ResourceSize.CPU,
			Memory: s.Settings.ResourceSize.Memory,
			Disk:   s.Settings.ResourceSize.Disk,
			GPU:    gpu,
		},
		CreateTime: s.CreateTime,
		UpdateTime: s.UpdateTime,
	}
}
