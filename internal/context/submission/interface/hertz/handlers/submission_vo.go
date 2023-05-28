package handlers

type CreateSubmissionRequest struct {
	WorkspaceID    string         `path:"workspace_id"`
	Name           string         `json:"name"`
	WorkflowID     string         `json:"workflowID"`
	Description    *string        `json:"description"`
	Type           string         `json:"type"`
	Entity         *Entity        `json:"entity"`
	ExposedOptions ExposedOptions `json:"exposedOptions"`
	InOutMaterial  *InOutMaterial `json:"inOutMaterial"`
}

type Entity struct {
	DataModelID     string   `json:"dataModelID"`
	DataModelRowIDs []string `json:"dataModelRowIDs"`
	/** 输入配置，json 序列化后的 string
	  采用 json 序列化原因基于以下两点考虑：
	  - thrift/接口设计层面不允许 `Value` 类型不确定
	  - 在 inputs/outputs 层级进行序列化可使得 `bioos-server` 不处理 `Inputs`/`Outputs`(非 `this.xxx` 索引的输入) 就入库/提交给计算引擎，达到透传效果
	*/
	InputsTemplate string `json:"inputsTemplate"`
	/** 输出配置，json 序列化后的 string
	  采用 json 序列化原因基于以下两点考虑：
	  - thrift/接口设计层面不允许 `Value` 类型不确定
	  - 在 inputs/outputs 层级进行序列化可使得 `bioos-server` 不处理 `Inputs`/`Outputs`(非 `this.xxx` 索引的输入) 就入库/提交给计算引擎，达到透传效果
	*/
	OutputsTemplate string `json:"outputsTemplate"`
}

type InOutMaterial struct {
	InputsMaterial  string `json:"inputsMaterial"`
	OutputsMaterial string `json:"outputsMaterial"`
}

type ExposedOptions struct {
	ReadFromCache bool `json:"readFromCache"`
}

type CreateSubmissionResponse struct {
	ID string `json:"id"`
}

type CheckSubmissionRequest struct {
	WorkspaceID string `path:"workspace_id"`
	Name        string `json:"name"`
}

type CheckSubmissionResponse struct {
	IsNameExist bool `json:"isNameExist"`
}

type ListSubmissionsRequest struct {
	WorkspaceID string   `path:"workspace_id"`
	Page        int      `query:"page"`
	Size        int      `query:"size"`
	OrderBy     string   `query:"orderBy"`
	SearchWord  string   `query:"searchWord"`
	Exact       bool     `query:"exact"`
	WorkflowID  string   `query:"workflowID"`
	Status      []string `query:"status"`
	IDs         []string `query:"ids"`
}

type ListSubmissionsResponse struct {
	Page  int              `json:"page"`
	Size  int              `json:"size"`
	Total int              `json:"total"`
	Items []SubmissionItem `json:"items"`
}

type SubmissionItem struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Description     *string         `json:"description"`
	Type            string          `json:"type"`
	Status          string          `json:"status"`
	StartTime       int64           `json:"startTime"`
	FinishTime      *int64          `json:"finishTime"`
	Duration        int64           `json:"duration"`
	WorkflowVersion WorkflowVersion `json:"workflowVersion"`
	RunStatus       Status          `json:"runStatus"`
	Entity          *Entity         `json:"entity"`
	ExposedOptions  ExposedOptions  `json:"exposedOptions"`
	InOutMaterial   *InOutMaterial  `json:"inOutMaterial"`
}

type WorkflowVersion struct {
	ID        string `json:"id"`
	VersionID string `json:"versionID"`
}

type Status struct {
	Count        int64 `json:"count"`
	Pending      int64 `json:"pending"`
	Succeeded    int64 `json:"succeeded"`
	Failed       int64 `json:"failed"`
	Running      int64 `json:"running"`
	Cancelling   int64 `json:"cancelling"`
	Cancelled    int64 `json:"cancelled"`
	Queued       int64 `json:"queued"`
	Initializing int64 `json:"initializing"`
}

type DeleteSubmissionRequest struct {
	WorkspaceID string `path:"workspace_id"`
	ID          string `path:"id"`
}

type CancelSubmissionRequest struct {
	WorkspaceID string `path:"workspace_id"`
	ID          string `path:"id"`
}
