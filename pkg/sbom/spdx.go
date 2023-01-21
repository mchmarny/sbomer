package sbom

import (
	"errors"
	"strings"

	"github.com/mchmarny/sbomer/pkg/doc"
	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_2"
	"github.com/spdx/tools-golang/spdx/v2_3"
)

const (
	pkgURLAbr    = "purl"
	spdxSpecName = "SPDX"
)

func spdx22ToDoc(bom *v2_2.Document) (*doc.Document, error) {
	if bom == nil {
		return nil, errors.New("invalid SPDX doc")
	}

	d := doc.NewDoc()
	getImageParts(bom.DocumentName, d)
	parseSPDXVersionParts(bom.SPDXVersion, d)
	d.Provider = parseSPDXCreator(bom.CreationInfo.Creators)
	d.Created = parseTime(bom.CreationInfo.Created)

	for _, p := range bom.Packages {
		item := &doc.Item{
			ID:         string(p.PackageSPDXIdentifier),
			Originator: parseSPDXOriginator(p.PackageOriginator, p.PackageSupplier),
			Name:       p.PackageName,
			Version:    p.PackageVersion,
			License:    firstNonEmpty(p.PackageLicenseConcluded, p.PackageLicenseDeclared),
			PURL:       getSPDX22RefValue(pkgURLAbr, p),
			// CPE:      ?,
			// SWID:     ?,
		}
		parseSPDXAttributeContext(item, p.PackageAttributionTexts)
		parseSPDX22Refs(item, p.PackageExternalReferences)
		d.Items = append(d.Items, item)
	}

	return d, nil
}

func spdx23ToDoc(bom *v2_3.Document) (*doc.Document, error) {
	if bom == nil {
		return nil, errors.New("invalid SPDX doc")
	}

	d := doc.NewDoc()
	getImageParts(bom.DocumentName, d)
	parseSPDXVersionParts(bom.SPDXVersion, d)
	d.Provider = parseSPDXCreator(bom.CreationInfo.Creators)
	d.Created = parseTime(bom.CreationInfo.Created)

	for _, p := range bom.Packages {
		item := &doc.Item{
			ID:         string(p.PackageSPDXIdentifier),
			Originator: parseSPDXOriginator(p.PackageOriginator, p.PackageSupplier),
			Name:       p.PackageName,
			Version:    p.PackageVersion,
			License:    firstNonEmpty(p.PackageLicenseConcluded, p.PackageLicenseDeclared),
			PURL:       getSPDX23RefValue(pkgURLAbr, p),
			// CPE:      ?,
			// SWID:     ?,
		}
		parseSPDXAttributeContext(item, p.PackageAttributionTexts)
		parseSPDX23Refs(item, p.PackageExternalReferences)
		d.Items = append(d.Items, item)
	}

	return d, nil
}

func parseSPDXOriginator(org *common.Originator, sup *common.Supplier) string {
	if org != nil {
		return org.Originator
	}

	if sup != nil {
		return sup.Supplier
	}

	return ""
}

func getSPDX23RefValue(k string, p *v2_3.Package) string {
	if p == nil || len(p.PackageExternalReferences) == 0 {
		return ""
	}
	for _, r := range p.PackageExternalReferences {
		if strings.EqualFold(r.RefType, k) {
			return r.Locator
		}
	}
	return ""
}

func getSPDX22RefValue(k string, p *v2_2.Package) string {
	if p == nil || len(p.PackageExternalReferences) == 0 {
		return ""
	}
	for _, r := range p.PackageExternalReferences {
		if strings.EqualFold(r.RefType, k) {
			return r.Locator
		}
	}
	return ""
}

const expectedAttributionParts = 2

func parseSPDX22Refs(item *doc.Item, list []*v2_2.PackageExternalReference) {
	if len(list) == 0 {
		return
	}

	for _, r := range list {
		item.AddContext(ctxTypeReference, r.RefType, r.Locator)
	}
}

func parseSPDX23Refs(item *doc.Item, list []*v2_3.PackageExternalReference) {
	if len(list) == 0 {
		return
	}

	for _, r := range list {
		item.AddContext(ctxTypeReference, r.RefType, r.Locator)
	}
}

func parseSPDXAttributeContext(item *doc.Item, list []string) {
	if len(list) == 0 {
		return
	}

	for _, a := range list {
		parts := strings.Split(a, ":")
		if len(parts) == expectedAttributionParts {
			item.AddContext(ctxTypeProperty,
				strings.TrimSpace(parts[0]),
				strings.TrimSpace(parts[1]))
		}
	}
}

func parseSPDXVersionParts(val string, d *doc.Document) {
	parts := strings.Split(val, "-")
	if len(parts) == expectedAttributionParts {
		d.Format = parts[0]
		d.FormatVersion = parts[1]
		return
	}
	d.Format = parts[0]
}

func parseSPDXCreator(list []common.Creator) string {
	if len(list) == 0 {
		return ""
	}

	sb := strings.Builder{}
	for _, c := range list {
		sb.WriteString(c.Creator)
	}

	return sb.String()
}
