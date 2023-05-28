package formatter

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type Formatter interface {
	Write(interface{})
}

func NewFormatter(format Format, output io.Writer) Formatter {
	switch format {
	case TextFormat:
		return NewText(output)
	case TableFormat:
		return NewTable(output)
	case JsonFormat:
		return NewJson(output)
	case YamlFormat:
		return NewYaml(output)
	default:
		return NewText(output)
	}
}

func NewText(output io.Writer) Formatter {
	return &Text{output}
}

func NewTable(output io.Writer) Formatter {
	return &Table{
		*tabwriter.NewWriter(output, 0, 8, 0, tableDelimiter[0], 0),
	}
}

func NewJson(output io.Writer) Formatter {
	return &Json{output}
}
func NewYaml(output io.Writer) Formatter {
	return &Yaml{output}
}

type Format string

const (
	TextFormat  Format = "text"
	TableFormat Format = "table"
	JsonFormat  Format = "json"
	YamlFormat  Format = "yaml"
)

func (f Format) Validate() error {
	switch f {
	case TextFormat, TableFormat, JsonFormat, YamlFormat:
		return nil
	}
	return fmt.Errorf("invalid formatter type. Valid formatter types are 'text', 'table', or 'json'")
}
