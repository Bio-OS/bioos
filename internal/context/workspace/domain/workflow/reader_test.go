package workflow

import (
	"context"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/test/assert"

	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils/git"
)

const (
	goodNextflowContent = `
nextflow.enable.dsl = 2
`

	badNextflowContent = `
nextflow.enable.dsl = x2
`
)

func cloneRepo() (dir string, err error) {
	dir, err = os.MkdirTemp("", "test-nextflow-reader")
	if err != nil {
		return "", err
	}

	err = git.Clone(dir, "https://gitee.com/a2htray/rnaseq-nf", "", "v1.0.2")
	if err != nil {
		return "", err
	}
	return dir, err
}

func TestNextflowReader_ParseWorkflowVersion(t *testing.T) {
	log.RegisterLogger(nil)
	dir, err := os.MkdirTemp("", "test-nextflow-reader")
	assert.Nil(t, err)
	goodFile, err := os.Create(path.Join(dir, "main-good.nf"))
	assert.Nil(t, err)
	badFile, err := os.Create(path.Join(dir, "main-bad.nf"))
	assert.Nil(t, err)

	defer func() {
		goodFile.Close()
		badFile.Close()
		os.RemoveAll(dir)
	}()

	_, err = goodFile.Write([]byte(goodNextflowContent))
	assert.Nil(t, err)
	_, err = badFile.Write([]byte(badNextflowContent))
	assert.Nil(t, err)

	r := &NextflowReader{}
	version, err := r.ParseWorkflowVersion(context.Background(), path.Join(dir, "main-good.nf"))
	assert.Nil(t, err)
	assert.DeepEqual(t, "DSL2", version)

	version, err = r.ParseWorkflowVersion(context.Background(), path.Join(dir, "main-bad.nf"))
	assert.NotNil(t, err)

}

func TestNextflowReader_ValidateWorkflowFiles(t *testing.T) {
	log.RegisterLogger(nil)
	dir, err := cloneRepo()
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	r := &NextflowReader{}
	workflowVersion := &WorkflowVersion{
		Files: map[string]*WorkflowFile{},
	}
	err = r.ValidateWorkflowFiles(context.Background(), workflowVersion, dir, "main.nf")
	assert.Nil(t, err)
	assert.DeepEqual(t, 4, len(workflowVersion.Files))

	inputParams, err := r.GetWorkflowInputs(context.Background(), "")
	assert.Nil(t, err)
	for _, param := range inputParams {
		t.Log(param.Type, param.Name, param.Optional, param.Default)
	}

	outputParams, err := r.GetWorkflowOutputs(context.Background(), "")
	assert.Nil(t, err)
	for _, param := range outputParams {
		t.Log(param.Type, param.Name, param.Optional, param.Default)
	}
}

func TestNextflowReader_GetWorkflowGraph(t *testing.T) {
	r := NextflowReader{
		graphFilepath: "/path/to/your/dag.html",
	}
	graph, err := r.GetWorkflowGraph(context.Background(), "")
	assert.Nil(t, err)
	t.Log(graph)
}

func TestNextflowRunPreview(t *testing.T) {
	workdir := "/Users/a2htray/workspace/rnaseq-nf"
	cmd := exec.CommandContext(context.Background(), "nextflow", "run", "main.nf", "-preview", "-with-dag", "dag.html")
	cmd.Dir = workdir
	//cmd.Stdout = os.Stdout
	err := cmd.Start()
	assert.Nil(t, err)
	err = cmd.Wait()
	assert.Nil(t, err)
}
