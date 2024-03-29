// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.24.2
// source: internal/context/workspace/interface/grpc/proto/notebook.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateNotebookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkspaceID string `protobuf:"bytes,1,opt,name=workspaceID,proto3" json:"workspaceID,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Content     []byte `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *CreateNotebookRequest) Reset() {
	*x = CreateNotebookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateNotebookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNotebookRequest) ProtoMessage() {}

func (x *CreateNotebookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNotebookRequest.ProtoReflect.Descriptor instead.
func (*CreateNotebookRequest) Descriptor() ([]byte, []int) {
	return file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescGZIP(), []int{0}
}

func (x *CreateNotebookRequest) GetWorkspaceID() string {
	if x != nil {
		return x.WorkspaceID
	}
	return ""
}

func (x *CreateNotebookRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateNotebookRequest) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type CreateNotebookResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateNotebookResponse) Reset() {
	*x = CreateNotebookResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateNotebookResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNotebookResponse) ProtoMessage() {}

func (x *CreateNotebookResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNotebookResponse.ProtoReflect.Descriptor instead.
func (*CreateNotebookResponse) Descriptor() ([]byte, []int) {
	return file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescGZIP(), []int{1}
}

type DeleteNotebookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkspaceID string `protobuf:"bytes,1,opt,name=workspaceID,proto3" json:"workspaceID,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DeleteNotebookRequest) Reset() {
	*x = DeleteNotebookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteNotebookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteNotebookRequest) ProtoMessage() {}

func (x *DeleteNotebookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteNotebookRequest.ProtoReflect.Descriptor instead.
func (*DeleteNotebookRequest) Descriptor() ([]byte, []int) {
	return file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteNotebookRequest) GetWorkspaceID() string {
	if x != nil {
		return x.WorkspaceID
	}
	return ""
}

func (x *DeleteNotebookRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type DeleteNotebookResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteNotebookResponse) Reset() {
	*x = DeleteNotebookResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteNotebookResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteNotebookResponse) ProtoMessage() {}

func (x *DeleteNotebookResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteNotebookResponse.ProtoReflect.Descriptor instead.
func (*DeleteNotebookResponse) Descriptor() ([]byte, []int) {
	return file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescGZIP(), []int{3}
}

type ListNotebooksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkspaceID string `protobuf:"bytes,1,opt,name=workspaceID,proto3" json:"workspaceID,omitempty"`
}

func (x *ListNotebooksRequest) Reset() {
	*x = ListNotebooksRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListNotebooksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListNotebooksRequest) ProtoMessage() {}

func (x *ListNotebooksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListNotebooksRequest.ProtoReflect.Descriptor instead.
func (*ListNotebooksRequest) Descriptor() ([]byte, []int) {
	return file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescGZIP(), []int{4}
}

func (x *ListNotebooksRequest) GetWorkspaceID() string {
	if x != nil {
		return x.WorkspaceID
	}
	return ""
}

type Notebook struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Length    int64                  `protobuf:"varint,2,opt,name=length,proto3" json:"length,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
}

func (x *Notebook) Reset() {
	*x = Notebook{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Notebook) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notebook) ProtoMessage() {}

func (x *Notebook) ProtoReflect() protoreflect.Message {
	mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notebook.ProtoReflect.Descriptor instead.
func (*Notebook) Descriptor() ([]byte, []int) {
	return file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescGZIP(), []int{5}
}

func (x *Notebook) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Notebook) GetLength() int64 {
	if x != nil {
		return x.Length
	}
	return 0
}

func (x *Notebook) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type ListNotebooksResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Notebook `protobuf:"bytes,1,rep,name=Items,proto3" json:"Items,omitempty"`
}

func (x *ListNotebooksResponse) Reset() {
	*x = ListNotebooksResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListNotebooksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListNotebooksResponse) ProtoMessage() {}

func (x *ListNotebooksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListNotebooksResponse.ProtoReflect.Descriptor instead.
func (*ListNotebooksResponse) Descriptor() ([]byte, []int) {
	return file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescGZIP(), []int{6}
}

func (x *ListNotebooksResponse) GetItems() []*Notebook {
	if x != nil {
		return x.Items
	}
	return nil
}

type GetNotebookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkspaceID string `protobuf:"bytes,1,opt,name=workspaceID,proto3" json:"workspaceID,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetNotebookRequest) Reset() {
	*x = GetNotebookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotebookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotebookRequest) ProtoMessage() {}

func (x *GetNotebookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotebookRequest.ProtoReflect.Descriptor instead.
func (*GetNotebookRequest) Descriptor() ([]byte, []int) {
	return file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescGZIP(), []int{7}
}

func (x *GetNotebookRequest) GetWorkspaceID() string {
	if x != nil {
		return x.WorkspaceID
	}
	return ""
}

func (x *GetNotebookRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type GetNotebookResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *GetNotebookResponse) Reset() {
	*x = GetNotebookResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNotebookResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNotebookResponse) ProtoMessage() {}

func (x *GetNotebookResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNotebookResponse.ProtoReflect.Descriptor instead.
func (*GetNotebookResponse) Descriptor() ([]byte, []int) {
	return file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescGZIP(), []int{8}
}

func (x *GetNotebookResponse) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

var File_internal_context_workspace_interface_grpc_proto_notebook_proto protoreflect.FileDescriptor

var file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDesc = []byte{
	0x0a, 0x3e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x78, 0x74, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x6e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x67, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x20, 0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x22, 0x18, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x62,
	0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4d, 0x0a, 0x15, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x18, 0x0a, 0x16, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x38, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x65,
	0x62, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b,
	0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x44, 0x22, 0x70,
	0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x38, 0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x22, 0x3e, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x05, 0x49, 0x74, 0x65,
	0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52, 0x05, 0x49, 0x74, 0x65, 0x6d, 0x73,
	0x22, 0x4a, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x77, 0x6f, 0x72,
	0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2f, 0x0a, 0x13,
	0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x32, 0xc9, 0x02,
	0x0a, 0x0f, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x4f, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x62,
	0x6f, 0x6f, 0x6b, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x4f, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65,
	0x62, 0x6f, 0x6f, 0x6b, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f, 0x74, 0x65, 0x62,
	0x6f, 0x6f, 0x6b, 0x73, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4e, 0x6f,
	0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x46, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b,
	0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x65,
	0x62, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescOnce sync.Once
	file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescData = file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDesc
)

func file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescGZIP() []byte {
	file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescOnce.Do(func() {
		file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescData)
	})
	return file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDescData
}

var file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_internal_context_workspace_interface_grpc_proto_notebook_proto_goTypes = []interface{}{
	(*CreateNotebookRequest)(nil),  // 0: proto.CreateNotebookRequest
	(*CreateNotebookResponse)(nil), // 1: proto.CreateNotebookResponse
	(*DeleteNotebookRequest)(nil),  // 2: proto.DeleteNotebookRequest
	(*DeleteNotebookResponse)(nil), // 3: proto.DeleteNotebookResponse
	(*ListNotebooksRequest)(nil),   // 4: proto.ListNotebooksRequest
	(*Notebook)(nil),               // 5: proto.Notebook
	(*ListNotebooksResponse)(nil),  // 6: proto.ListNotebooksResponse
	(*GetNotebookRequest)(nil),     // 7: proto.GetNotebookRequest
	(*GetNotebookResponse)(nil),    // 8: proto.GetNotebookResponse
	(*timestamppb.Timestamp)(nil),  // 9: google.protobuf.Timestamp
}
var file_internal_context_workspace_interface_grpc_proto_notebook_proto_depIdxs = []int32{
	9, // 0: proto.Notebook.updatedAt:type_name -> google.protobuf.Timestamp
	5, // 1: proto.ListNotebooksResponse.Items:type_name -> proto.Notebook
	0, // 2: proto.NotebookService.CreateNotebook:input_type -> proto.CreateNotebookRequest
	2, // 3: proto.NotebookService.DeleteNotebook:input_type -> proto.DeleteNotebookRequest
	4, // 4: proto.NotebookService.ListNotebooks:input_type -> proto.ListNotebooksRequest
	7, // 5: proto.NotebookService.GetNotebook:input_type -> proto.GetNotebookRequest
	1, // 6: proto.NotebookService.CreateNotebook:output_type -> proto.CreateNotebookResponse
	3, // 7: proto.NotebookService.DeleteNotebook:output_type -> proto.DeleteNotebookResponse
	6, // 8: proto.NotebookService.ListNotebooks:output_type -> proto.ListNotebooksResponse
	8, // 9: proto.NotebookService.GetNotebook:output_type -> proto.GetNotebookResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_context_workspace_interface_grpc_proto_notebook_proto_init() }
func file_internal_context_workspace_interface_grpc_proto_notebook_proto_init() {
	if File_internal_context_workspace_interface_grpc_proto_notebook_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateNotebookRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateNotebookResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteNotebookRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteNotebookResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListNotebooksRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Notebook); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListNotebooksResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotebookRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNotebookResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_context_workspace_interface_grpc_proto_notebook_proto_goTypes,
		DependencyIndexes: file_internal_context_workspace_interface_grpc_proto_notebook_proto_depIdxs,
		MessageInfos:      file_internal_context_workspace_interface_grpc_proto_notebook_proto_msgTypes,
	}.Build()
	File_internal_context_workspace_interface_grpc_proto_notebook_proto = out.File
	file_internal_context_workspace_interface_grpc_proto_notebook_proto_rawDesc = nil
	file_internal_context_workspace_interface_grpc_proto_notebook_proto_goTypes = nil
	file_internal_context_workspace_interface_grpc_proto_notebook_proto_depIdxs = nil
}
