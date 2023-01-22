package sbom

import (
	"errors"
	"strings"

	"github.com/mchmarny/sbomer/pkg/doc"
	"github.com/spdx/tools-golang/spdx/common"
	"github.com/spdx/tools-golang/spdx/v2_2"
	"github.com/spdx/tools-golang/spdx/v2_3"
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
			ID:      string(p.PackageSPDXIdentifier),
			Name:    p.PackageName,
			Version: firstNonEmpty(p.PackageVersion, missingValueDefault),
		}
		parseSPDXAttributeContext(item, p.PackageAttributionTexts)
		parseSPDX22Refs(item, p.PackageExternalReferences)
		parseSPDX22Ctx(item, p)
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
			ID:      string(p.PackageSPDXIdentifier),
			Name:    p.PackageName,
			Version: firstNonEmpty(p.PackageVersion, missingValueDefault),
		}
		parseSPDXAttributeContext(item, p.PackageAttributionTexts)
		parseSPDX23Refs(item, p.PackageExternalReferences)
		parseSPDX23Ctx(item, p)
		d.Items = append(d.Items, item)
	}

	return d, nil
}

func parseSPDX23Ctx(item *doc.Item, p *v2_3.Package) {
	if p.PackageFileName != "" {
		item.AddContext(ctxTypeComponent, "file", p.PackageFileName)
	}

	if p.PackageOriginator != nil {
		item.AddContext(ctxTypeComponent, p.PackageOriginator.OriginatorType,
			p.PackageOriginator.Originator)
	}

	if p.PackageSupplier != nil {
		item.AddContext(ctxTypeComponent, p.PackageSupplier.SupplierType,
			p.PackageSupplier.Supplier)
	}

	if p.PackageDownloadLocation != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindDownloadLocation), p.PackageDownloadLocation)
	}

	if len(p.PackageChecksums) > 0 {
		for _, c := range p.PackageChecksums {
			item.AddContext(ctxTypeChecksum, string(c.Algorithm), c.Value)
		}
	}

	if p.PackageHomePage != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindHomepage), p.PackageHomePage)
	}

	if p.PackageSourceInfo != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindSourceInfo), p.PackageSourceInfo)
	}

	if p.PackageSummary != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindSummary), p.PackageSummary)
	}

	if p.PrimaryPackagePurpose != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindPurpose), p.PrimaryPackagePurpose)
	}

	if p.ReleaseDate != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindReleaseDate), p.ReleaseDate)
	}

	if p.BuiltDate != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindBuildDate), p.BuiltDate)
	}

	if len(p.Files) > 0 {
		for _, f := range p.Files {
			item.AddContext(ctxTypeChecksum, f.FileName, string(f.FileSPDXIdentifier))
		}
	}
}

func parseSPDX22Ctx(item *doc.Item, p *v2_2.Package) {
	if p.PackageFileName != "" {
		item.AddContext(ctxTypeComponent, "file", p.PackageFileName)
	}

	if p.PackageOriginator != nil {
		item.AddContext(ctxTypeComponent, p.PackageOriginator.OriginatorType,
			p.PackageOriginator.Originator)
	}

	if p.PackageSupplier != nil {
		item.AddContext(ctxTypeComponent, p.PackageSupplier.SupplierType,
			p.PackageSupplier.Supplier)
	}

	if p.PackageDownloadLocation != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindDownloadLocation), p.PackageDownloadLocation)
	}

	if len(p.PackageChecksums) > 0 {
		for _, c := range p.PackageChecksums {
			item.AddContext(ctxTypeChecksum, string(c.Algorithm), c.Value)
		}
	}

	if p.PackageHomePage != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindHomepage), p.PackageHomePage)
	}

	if p.PackageSourceInfo != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindSourceInfo), p.PackageSourceInfo)
	}

	if p.PackageSummary != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindSummary), p.PackageSummary)
	}

	if len(p.Files) > 0 {
		for _, f := range p.Files {
			item.AddContext(ctxTypeChecksum, f.FileName, string(f.FileSPDXIdentifier))
		}
	}

	if p.PackageVerificationCode.Value != "" {
		item.AddContext(ctxTypeComponent, string(doc.ContextKindVerificationCode), p.PackageVerificationCode.Value)
	}
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
