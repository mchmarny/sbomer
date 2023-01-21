package doc

import (
	"sync"
	"time"
)

func NewDoc() *Document {
	return &Document{
		Items: make([]*Item, 0),
	}
}

// Document is the top level struct for the SBOM document:
// https://www.ntia.doc.gov/files/ntia/publications/sbom_minimum_elements_report.pdf.
type Document struct {
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

	// Originator is the originator of the item (originator or author)
	Originator string `json:"originator"`

	// Name is the name of the item (name or title)
	Name string `json:"name"`

	// Version is the version of the item (versionInfo or version)
	Version string `json:"version"`

	License string `json:"license"` // licenseConcluded or licenses.license.id

	// CPE is the Common Platform Enumeration.
	CPE string `json:"cpe"` // Common Platform Enumeration

	// SWID is the Software Identification.
	SWID string `json:"swid"` // Software Identification

	// PURL is the Package Uniform Resource Locators.
	PURL string `json:"purl"` // components.purl

	Contexts []*Context `json:"pkgContexts"` // components.externalReferences or components.properties

	mu sync.Mutex
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
