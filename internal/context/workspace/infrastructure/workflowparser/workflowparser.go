package workflowparser

import (
	"context"
	"fmt"
	"sync"
)

const (
	WDL       = "WDL"
	CWL       = "CWL"
	Snakemake = "SMK"
	Nextflow  = "NFL"
)

type ParserConfig interface{}

type WorkflowParser interface {
	ParseWorkflowVersion(ctx context.Context, mainWorkflowPath string) (string, error)
	ValidateWorkflowFiles(ctx context.Context, baseDir, mainWorkflowPath string) (string, error)
	GetWorkflowInputs(ctx context.Context, WorkflowFilePath string) (string, error)
	GetWorkflowOutputs(ctx context.Context, WorkflowFilePath string) (string, error)
	GetWorkflowGraph(ctx context.Context, WorkflowFilePath string) (string, error)
}

var parserCreateFunc = make(map[string]func(ParserConfig) WorkflowParser)

func RegisterParseCreateFunc(workflowType string, function func(ParserConfig) WorkflowParser) {
	parserCreateFunc[workflowType] = function
}

func init() {
	RegisterParseCreateFunc(WDL, func(config ParserConfig) WorkflowParser {
		configParam, ok := config.(WDLConfig)
		if !ok {
			panic("Invalid config type for WDL parser")
		}
		return NewWDLParser(configParam)
	})
	RegisterParseCreateFunc(CWL, func(config ParserConfig) WorkflowParser {
		configParam, ok := config.(CWLConfig)
		if !ok {
			panic("Invalid config type for WDL parser")
		}
		return NewCWLParser(configParam)
	})
}

var (
	instance *WorkflowParserFactory
	once     sync.Once
)

type WorkflowParserFactory struct {
	parsers map[string]WorkflowParser
}

func InitWorkflowParserFactory(configs map[string]ParserConfig) {
	once.Do(func() {
		instance = &WorkflowParserFactory{
			parsers: make(map[string]WorkflowParser),
		}

		for workflowType, config := range configs {
			createFunc, exists := parserCreateFunc[workflowType]
			if !exists {
				panic("no parser registered for workflow type")
			}
			instance.parsers[workflowType] = createFunc(config)
		}
	})
}

func GetFactory() *WorkflowParserFactory {
	return instance
}

func (f *WorkflowParserFactory) GetParser(workflowType string) (WorkflowParser, error) {
	parser, exists := f.parsers[workflowType]
	if !exists {
		return nil, fmt.Errorf("no parser found for workflow type: %s", workflowType)
	}
	return parser, nil
}
