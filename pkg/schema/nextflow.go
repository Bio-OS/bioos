package schema

const (
	TypeString  = "string"
	TypeNumber  = "number"
	TypeInteger = "integer"
	TypeBoolean = "boolean"
)

type NextflowSchema struct {
	Schema      string `json:"$schema"`
	ID          string `json:"$id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Definitions map[string]DefinitionSchema
}

type PropertySchema struct {
	Type        string      `json:"type"`
	Out         bool        `json:"out"`
	Format      string      `json:"format"`
	Description string      `json:"description"`
	Mimetype    string      `json:"mimetype"`
	Default     interface{} `json:"default"`
}

type DefinitionSchema struct {
	Title      string                    `json:"title"`
	Type       string                    `json:"type"`
	Required   []string                  `json:"required"`
	Properties map[string]PropertySchema `json:"properties"`
}
