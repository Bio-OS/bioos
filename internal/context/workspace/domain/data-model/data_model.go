package datamodel

import "time"

type DataModel struct {
	WorkspaceID string
	ID          string
	Name        string
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Headers     []string   // all headers in this data model
	RowIDs      []string   // all rowIDs in this data model, only used in delete data model
	Rows        [][]string // the row should be saved(create/update)
}

type Row struct {
	Grids []*Grid `json:"grids"`
}

type Grid struct {
	Value []byte `json:"value"`
	Type  string `json:"type"`
}
