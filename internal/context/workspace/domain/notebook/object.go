package notebook

import (
	"path"
	"time"
)

// Notebook is domain object
type Notebook struct {
	Name       string
	Namespace  string
	Content    []byte
	UpdateTime time.Time
}

// Path return unique path of notebook, e.g. {workspace-id}/{name}
func (n *Notebook) Path() string {
	return Path(n.Namespace, n.Name)
}

// Path ...
func Path(namespace, name string) string {
	return path.Join(namespace, name)
}

// NamespaceOfPath ...
func NamespaceOfPath(name string) string {
	return path.Dir(name)
}

// NameOfPath ...
func NameOfPath(name string) string {
	return path.Base(name)
}
