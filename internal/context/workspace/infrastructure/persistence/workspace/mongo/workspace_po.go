package mongo

import "time"

type workspacePO struct {
	ID          string           `json:"id" bson:"id"`
	Name        string           `json:"name" bson:"name"`
	Description string           `json:"description" bson:"description"`
	Storage     workspaceStorage `json:"storage" bson:"storage"`
	CreateTime  time.Time        `json:"createTime" bson:"createTime"`
	UpdateTime  time.Time        `json:"updateTime" bson:"updateTime"`
}

// WorkspaceStorage ...
type workspaceStorage struct {
	NFS *workspaceStorageNFS `json:"nfs" bson:"nfs,omitempty"`
}

// NFSWorkspaceStorage ...
type workspaceStorageNFS struct {
	MountPath string `json:"mountPath" bson:"mountPath"`
}
