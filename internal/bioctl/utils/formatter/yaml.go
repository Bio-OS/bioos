package formatter

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type Yaml struct {
	writer io.Writer
}

func (f *Yaml) write(bytes []byte) {
	_, _ = fmt.Fprint(f.writer, string(bytes))
}

func (f *Yaml) newLine() {
	_, _ = fmt.Fprintln(f.writer)
}

func (f *Yaml) Write(o interface{}) {
	var b bytes.Buffer

	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	_ = yamlEncoder.Encode(o)

	f.write(b.Bytes())
	f.newLine()
}
