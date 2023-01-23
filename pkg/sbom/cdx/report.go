package cdx

import (
	"bytes"
	"fmt"
	"strings"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/mchmarny/sbomer/pkg/sbom"
)

type CycloneDXReport struct {
	valid    bool
	docError error
	name     string

	creationToolName    int
	creationToolVersion int

	totalPackages  int
	hasLicense     int
	hasPackVersion int
	hasPackDigest  int
	hasPurl        int
	hasCPE         int
}

func (r *CycloneDXReport) GetName() string {
	return r.name
}

func (r *CycloneDXReport) Report() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d total packages\n", r.totalPackages))

	sb.WriteString(fmt.Sprintf("%d%% have versions.\n", sbom.PrettyPercent(r.hasPackVersion, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have licenses.\n", sbom.PrettyPercent(r.hasLicense, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have package digest.\n", sbom.PrettyPercent(r.hasPackDigest, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have purls.\n", sbom.PrettyPercent(r.hasPurl, r.totalPackages)))
	sb.WriteString(fmt.Sprintf("%d%% have CPEs.\n", sbom.PrettyPercent(r.hasCPE, r.totalPackages)))

	sb.WriteString(fmt.Sprintf("Has creation info? %v\n", r.hasCreationInfo()))
	sb.WriteString(fmt.Sprintf("Spec valid? %v\n", r.valid))
	return sb.String()
}

func (r *CycloneDXReport) hasCreationInfo() bool {
	return r.creationToolName > 0 &&
		r.creationToolVersion > 0 &&
		r.creationToolName == r.creationToolVersion
}

func (r *CycloneDXReport) IsSpecCompliant() sbom.ReportValue {
	if r.docError != nil {
		return sbom.ReportValue{
			Ratio:     0,
			Reasoning: r.docError.Error(),
		}
	}
	return sbom.ReportValue{Ratio: 1}
}

const rationPrecision = 2

func (r *CycloneDXReport) PackageIdentification() sbom.ReportValue {
	purlPercent := sbom.PrettyPercent(r.hasPurl, r.totalPackages)
	cpePercent := sbom.PrettyPercent(r.hasCPE, r.totalPackages)
	return sbom.ReportValue{
		// What percentage has both Purl & CPEs?
		Ratio:     float32(r.hasPurl+r.hasCPE) / float32(r.totalPackages*rationPrecision),
		Reasoning: fmt.Sprintf("%d%% have purls and %d%% have CPEs", purlPercent, cpePercent),
	}
}

func (r *CycloneDXReport) PackageVersions() sbom.ReportValue {
	return sbom.ReportValue{
		Ratio: float32(r.hasPackVersion) / float32(r.totalPackages),
	}
}

func (r *CycloneDXReport) PackageDigests() sbom.ReportValue {
	return sbom.ReportValue{
		Ratio: float32(r.hasPackDigest) / float32(r.totalPackages),
	}
}

func (r *CycloneDXReport) PackageLicenses() sbom.ReportValue {
	return sbom.ReportValue{
		Ratio: float32(r.hasLicense) / float32(r.totalPackages),
	}
}

func (r *CycloneDXReport) CreationInfo() sbom.ReportValue {
	// @@@
	return sbom.ReportValue{Ratio: 1}
}

func GetCycloneDXReport(b []byte) sbom.SbomReport {
	r := CycloneDXReport{}
	formats := []cdx.BOMFileFormat{cdx.BOMFileFormatJSON, cdx.BOMFileFormatXML}

	bom := new(cdx.BOM)
	for _, format := range formats {
		decoder := cdx.NewBOMDecoder(bytes.NewReader(b), format)
		if err := decoder.Decode(bom); err != nil {
			r.valid = false
			r.docError = err
		} else {
			r.valid = true
			r.docError = nil
			break
		}
	}

	if !r.valid {
		return &r
	}

	if bom.Metadata.Component != nil {
		r.name = bom.Metadata.Component.Name
	}

	if bom.Metadata.Tools != nil {
		for _, t := range *bom.Metadata.Tools {
			if t.Name != "" {
				r.creationToolName += 1
			}
			if t.Version != "" {
				r.creationToolVersion += 1
			}
		}
	}

	if bom.Components != nil {
		for _, p := range *bom.Components {
			r.totalPackages += 1
			if p.Licenses != nil && len(*p.Licenses) > 0 {
				r.hasLicense += 1
			}
			if p.Hashes != nil && len(*p.Hashes) > 0 {
				r.hasPackDigest += 1
			}
			if p.Version != "" {
				r.hasPackVersion += 1
			}
			if p.PackageURL != "" {
				r.hasPurl += 1
			}
			if p.CPE != "" {
				r.hasCPE += 1
			}
		}
	}

	return &r
}
