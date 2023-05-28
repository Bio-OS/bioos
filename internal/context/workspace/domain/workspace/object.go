package workspace

import (
	"fmt"
	"time"
)

// Workspace aggregate root DO.
type Workspace struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Storage     Storage
}

type Storage struct {
	NFS *NFSStorage
}

type NFSStorage struct {
	MountPath string
}

func (w *Workspace) GetID() string {
	return w.ID
}

func (w *Workspace) GetName() string {
	return w.Name
}

func (w *Workspace) UpdateName(name string) {
	if w.Name != name {
		w.UpdatedAt = time.Now()
		w.Name = name
	}
}

func (w *Workspace) GetCreatedAt() time.Time {
	return w.CreatedAt
}

func (w *Workspace) GetUpdatedAt() time.Time {
	return w.UpdatedAt
}

func (w *Workspace) GetDescription() string {
	return w.Description
}

func (w *Workspace) UpdateDescription(description string) {
	if w.Description != description {
		w.UpdatedAt = time.Now()
		w.Description = description
	}
}

func (w *Workspace) GetStorage() Storage {
	return w.Storage
}

func (w *Workspace) String() string {
	return fmt.Sprintf("workspace id:%s name:%s description:%s", w.ID, w.Name, w.Description)
}
