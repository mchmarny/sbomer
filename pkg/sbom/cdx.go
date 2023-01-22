package sbom

import (
	"errors"
	"fmt"
	"strings"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/mchmarny/sbomer/pkg/doc"
)

func cycloneToDoc(bom *cdx.BOM) (*doc.Document, error) {
	if bom == nil || bom.Metadata == nil || bom.Metadata.Component == nil || bom.Metadata.Component.Name == "" {
		return nil, errors.New("invalid CycloneDX doc")
	}

	d := doc.NewDoc()
	d.Subject = bom.Metadata.Component.Name
	d.SubjectVersion = bom.Metadata.Component.Version
	d.Format = firstNonEmpty(bom.BOMFormat, "CycloneDX")
	d.FormatVersion = bom.SpecVersion.String()
	d.Provider = parseGenerator(bom.Metadata)
	d.Created = parseTime(bom.Metadata.Timestamp)

	for _, c := range *bom.Components {
		item := &doc.Item{
			ID:      firstNonEmpty(c.BOMRef, fmt.Sprintf("%s-%s-%s", c.Type, c.Name, c.Version)),
			Name:    c.Name,
			Version: firstNonEmpty(c.Version, missingValueDefault),
		}

		parseCDXContext(item, c)
		d.Items = append(d.Items, item)
	}

	return d, nil
}

func parseGenerator(list *cdx.Metadata) string {
	if list == nil || list.Tools == nil || len(*list.Tools) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i, t := range *list.Tools {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(t.Name)
	}
	return sb.String()
}

func parseLicenses(list *cdx.Licenses) string {
	if list == nil || len(*list) == 0 {
		return ""
	}

	sb := strings.Builder{}
	for i, l := range *list {
		if i > 0 {
			sb.WriteString(", ")
		}
		if l.License != nil {
			sb.WriteString(l.License.ID)
			if l.License.Name != "" {
				sb.WriteString(" - ")
				sb.WriteString(l.License.Name)
			}
			if l.License.URL != "" {
				sb.WriteString(" - ")
				sb.WriteString(l.License.URL)
			}
		}
	}
	return sb.String()
}

func parseCDXContext(item *doc.Item, c cdx.Component) {
	if c.ExternalReferences != nil && len(*c.ExternalReferences) > 0 {
		for _, r := range *c.ExternalReferences {
			if r.URL != "" {
				item.AddContext(ctxTypeReference, string(r.Type), r.URL)
			}
		}
	}

	if c.Properties != nil && len(*c.Properties) > 0 {
		for _, r := range *c.Properties {
			if r.Value != "" {
				item.AddContext(ctxTypeProperty, r.Name, r.Value)
			}
		}
	}

	if lic := parseLicenses(c.Licenses); lic != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindLicense), lic)
	}

	if c.CPE != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindCPE), c.CPE)
	}

	if swid := parseCDXSWID(c.SWID); swid != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindSWID), swid)
	}

	if c.PackageURL != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindPURL), c.PackageURL)
	}

	if c.Type != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindType), string(c.Type))
	}

	if c.Supplier != nil && c.Supplier.Name != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindSupplier), c.Supplier.Name)
	}

	if c.Author != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindAuthor), c.Author)
	}

	if c.Publisher != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindPublisher), c.Publisher)
	}

	if c.Group != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindGroup), c.Group)
	}

	if c.Description != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindDescription), c.Description)
	}

	if c.Scope != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindScope), string(c.Scope))
	}

	if c.Hashes != nil && len(*c.Hashes) > 0 {
		for _, h := range *c.Hashes {
			if h.Value != "" {
				item.AddContext(ctxTypeProperty, string(h.Algorithm), h.Value)
			}
		}
	}
}

func parseCDXSWID(val *cdx.SWID) string {
	if val == nil {
		return ""
	}

	return strings.TrimSpace(val.TagID + " " + val.Version)
}
