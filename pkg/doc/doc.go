package doc

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

func NewDoc() *Document {
	return &Document{
		ID:    uuid.NewString(),
		Items: make([]*Item, 0),
	}
}

// Document is the top level struct for the SBOM document:
// https://www.ntia.doc.gov/files/ntia/publications/sbom_minimum_elements_report.pdf.
type Document struct {
	ID             string    `json:"id"`             // uuid
	Subject        string    `json:"subject"`        // memcached or gcr.io/image
	SubjectVersion string    `json:"subjectVersion"` // 1.6.9 or sha256:1234
	Format         string    `json:"format"`         // CycloneDX or SPDX
	FormatVersion  string    `json:"formatVersion"`  // 1.4 or 2.2 or 2.3
	Provider       string    `json:"provider"`       // syft, trivy
	Created        time.Time `json:"created"`        // 2021-07-01T00:00:00Z
	Items          []*Item   `json:"items"`
}

type Item struct {
	// ID is the unique identifier for the item (SPDXID or bom-ref)
	ID string `json:"id"`

	// Name is the name of the item (name or title)
	Name string `json:"name"`

	// Version is the version of the item (versionInfo or version)
	Version string `json:"version"`

	// components.externalReferences or components.properties
	Contexts []*Context `json:"pkgContexts"`

	mu sync.Mutex
}

func (p *Item) ToString() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("id: %s, name: %s, version: %s",
		p.ID, p.Name, p.Version))
	if p.Contexts != nil {
		for _, c := range p.Contexts {
			sb.WriteString(fmt.Sprintf(" - %+v", c))
		}
	}
	return sb.String()
}

func (p *Item) AddContext(t, k, v string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Contexts == nil {
		p.Contexts = make([]*Context, 0)
	}

	p.Contexts = append(p.Contexts, &Context{
		Type:  t,
		Key:   k,
		Value: v,
	})
}

type Context struct {
	Type  string `json:"ctxType"`
	Key   string `json:"ctxKey"`
	Value string `json:"ctxValue"`
}

func (c *Context) ToString() string {
	return fmt.Sprintf("type: %s, key: %s, value: %s", c.Type, c.Key, c.Value)
}
