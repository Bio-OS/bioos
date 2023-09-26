package convert

import workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"

type ListNotebooksRequest struct {
	WorkspaceID string `path:"workspace-id"`
}

func (req *ListNotebooksRequest) ToGRPC() *workspaceproto.ListNotebooksRequest {
	r := &workspaceproto.ListNotebooksRequest{
		WorkspaceID: req.WorkspaceID,
	}
	return r
}

type ListNotebooksResponse struct {
	Items []NotebookItem `json:"items"`
}

type NotebookItem struct {
	Name          string `json:"name"`
	ContentLength int64  `json:"contentLength"`
	UpdateTime    int64  `json:"updateTime"`
}

func (resp *ListNotebooksResponse) FromGRPC(protoResp *workspaceproto.ListNotebooksResponse) {
	resp.Items = make([]NotebookItem, len(protoResp.GetItems()))
	for i, item := range protoResp.Items {
		resp.Items[i] = NotebookItem{
			Name:          item.Name,
			ContentLength: item.Length,
			UpdateTime:    item.UpdatedAt.GetSeconds(),
		}
	}
	return
}

type GetNotebookRequest struct {
	WorkspaceID string `path:"workspace-id"`
	Name        string `path:"name"`
}

func (req *GetNotebookRequest) ToGRPC() *workspaceproto.GetNotebookRequest {
	r := &workspaceproto.GetNotebookRequest{
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
	}
	return r
}

type GetNotebookResponse struct {
	Content []byte
}

func (resp *GetNotebookResponse) FromGRPC(protoResp *workspaceproto.GetNotebookResponse) {
	resp.Content = protoResp.Content
	return
}
