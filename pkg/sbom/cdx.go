package sbom

import (
	"errors"
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
			ID:         c.BOMRef,
			Originator: firstNonEmpty(c.Author, c.Publisher),
			Name:       c.Name,
			Version:    c.Version,
			License:    parseLicenses(c.Licenses),
			CPE:        c.CPE,
			SWID:       parseCDXSWID(c.SWID),
			PURL:       c.PackageURL,
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
		}
	}
	return sb.String()
}

func parseCDXContext(item *doc.Item, c cdx.Component) {
	if c.ExternalReferences != nil && len(*c.ExternalReferences) > 0 {
		for _, r := range *c.ExternalReferences {
			item.AddContext(ctxTypeReference, string(r.Type), r.URL)
		}
	}

	if c.Properties != nil && len(*c.Properties) > 0 {
		for _, r := range *c.Properties {
			item.AddContext(ctxTypeProperty, r.Name, r.Value)
		}
	}
}

func parseCDXSWID(val *cdx.SWID) string {
	if val == nil {
		return ""
	}

	return strings.TrimSpace(val.TagID + " " + val.Version)
}
