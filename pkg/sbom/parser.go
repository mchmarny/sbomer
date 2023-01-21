package sbom

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"github.com/mchmarny/sbomer/pkg/doc"
	spdx_json "github.com/spdx/tools-golang/json"
	spdx_rdf "github.com/spdx/tools-golang/rdfloader"
	spdx_tv "github.com/spdx/tools-golang/tvloader"

	cdx "github.com/CycloneDX/cyclonedx-go"
)

const (
	ctxTypeReference = "reference"
	ctxTypeProperty  = "property"
)

func ParseDoc(b []byte) (*doc.Document, error) {
	// CycloneDX
	var bom cdx.BOM

	// CycloneDX JSON
	if err := cdx.NewBOMDecoder(bytes.NewReader(b), cdx.BOMFileFormatJSON).Decode(&bom); err == nil && bom.Metadata != nil {
		return cycloneToDoc(&bom)
	}

	// CycloneDX XML
	if err := cdx.NewBOMDecoder(bytes.NewReader(b), cdx.BOMFileFormatXML).Decode(&bom); err == nil && bom.Metadata != nil {
		return cycloneToDoc(&bom)
	}

	// SPDX v2.3
	if doc, err := spdx_json.Load2_3(bytes.NewReader(b)); err == nil && doc != nil {
		return spdx23ToDoc(doc)
	}

	if doc, err := spdx_rdf.Load2_3(bytes.NewReader(b)); err == nil && doc != nil {
		return spdx23ToDoc(doc)
	}

	if doc, err := spdx_tv.Load2_3(bytes.NewReader(b)); err == nil && doc != nil {
		return spdx23ToDoc(doc)
	}

	// SPDX v2.2
	if doc, err := spdx_json.Load2_2(bytes.NewReader(b)); err == nil && doc != nil {
		return spdx22ToDoc(doc)
	}

	if doc, err := spdx_rdf.Load2_2(bytes.NewReader(b)); err == nil && doc != nil {
		return spdx22ToDoc(doc)
	}

	if doc, err := spdx_tv.Load2_2(bytes.NewReader(b)); err == nil && doc != nil {
		return spdx22ToDoc(doc)
	}

	return nil, errors.New("invalid SBOM content")
}

const expectedImageParts = 2

func getImageParts(val string, d *doc.Document) {
	// digest
	if strings.Contains(val, "@") {
		parts := strings.Split(val, "@")
		if len(parts) == expectedImageParts {
			d.Subject = parts[0]
			d.SubjectVersion = parts[1]
			return
		}
		d.Subject = parts[0]
		return
	}

	// tag
	parts := strings.Split(val, ":")
	if len(parts) == expectedImageParts {
		d.Subject = parts[0]
		d.SubjectVersion = parts[1]
		return
	}

	d.Subject = parts[0]
}

func parseTime(val string) time.Time {
	if val == "" {
		return time.Now().UTC()
	}

	// 2021-07-01T00:00:00Z
	v, err := time.Parse("2006-01-02T15:04:05Z", val)
	if err != nil {
		return time.Now().UTC()
	}
	return v
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
