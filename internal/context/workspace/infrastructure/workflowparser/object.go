package workflowparser

import (
	"encoding/json"
	"time"
)

const (
	CommandExecuteTimeout = time.Minute * 3
	WDLVersionRegexpStr   = "^version\\s+([\\w-._]+)"
	CWLVersionRegexpStr   = "^cwlVersion\\s+([\\w-._]+)"
)

type WorkflowParam struct {
	// Name param name
	Name string `json:"name"`
	// Type param type
	Type string `json:"type"`
	// Optional param is optional
	Optional bool `json:"optional"`
	// Default param default value
	Default string `json:"default,omitempty"`
}

func workflowParamPOToDO(paramStr string) ([]WorkflowParam, error) {
	var params []WorkflowParam
	if len(paramStr) > 0 {
		if err := json.Unmarshal([]byte(paramStr), &params); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func workflowParamDoToPO(params []WorkflowParam) (string, error) {
	if len(params) == 0 {
		return "", nil
	}

	paramBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return string(paramBytes), nil
}

type FileParam struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func fileParamPOToDO(paramStr string) ([]FileParam, error) {
	var params []FileParam
	if len(paramStr) > 0 {
		if err := json.Unmarshal([]byte(paramStr), &params); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func fileParamDoToPO(params []FileParam) (string, error) {
	if len(params) == 0 {
		return "", nil
	}

	paramBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return string(paramBytes), nil
}
